CURRENT_DIR=$(shell pwd)
# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

KIND_VERSION ?= v0.18.0
CLUSTER_NAME ?= rollout-ci-testing

all: generate gen-openapi-schema

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./rollouts/..."
	@hack/generate_client.sh

CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
CONTROLLER_GEN_VERSION = v0.16.0
controller-gen: ## Download controller-gen locally if necessary.
ifeq ("$(shell $(CONTROLLER_GEN) --version)", "Version: ${CONTROLLER_GEN_VERSION}")
else
	rm -rf $(CONTROLLER_GEN)
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@${CONTROLLER_GEN_VERSION})
endif

OPENAPI_GEN = $(shell pwd)/bin/openapi-gen
module=$(shell go list -f '{{.Module}}' k8s.io/kube-openapi/cmd/openapi-gen | awk '{print $$1}')
module_version=$(shell go list -m $(module) | awk '{print $$NF}' | head -1)
openapi-gen: ## Download openapi-gen locally if necessary.
ifeq ("$(shell command -v $(OPENAPI_GEN) 2> /dev/null)", "")
	$(call go-get-tool,$(OPENAPI_GEN),k8s.io/kube-openapi/cmd/openapi-gen@$(module_version))
else
	@echo "openapi-gen is already installed."
endif

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

ensure-kind:
ifeq ("$(shell command -v $(PROJECT_DIR)/bin/kind 2> /dev/null)", "")
	@echo "Downloading kind version $(KIND_VERSION)"
	GOBIN=$(PROJECT_DIR)/bin go install sigs.k8s.io/kind@$(KIND_VERSION)
else
	@echo "kind is already installed."
endif

delete-cluster: ensure-kind
	@echo "Deleting kind cluster $(CLUSTER_NAME) if it exists"
	bin/kind delete cluster --name $(CLUSTER_NAME) || true

create-cluster: ensure-kind
	bin/kind create cluster --name $(CLUSTER_NAME)

.PHONY: run-e2e-test
run-e2e-test:
	@echo "Installing Rollouts CRDs"
	kubectl apply -f https://raw.githubusercontent.com/openkruise/rollouts/refs/heads/master/config/crd/bases/rollouts.kruise.io_rollouts.yaml
	@echo "Running E2E tests"
	go test -v ./tests/e2e/...

.PHONY: gen-schema-only
gen-schema-only:
	go run cmd/gen-schema/main.go

.PHONY: gen-openapi-schema
gen-openapi-schema: gen-rollouts-openapi
	go run cmd/gen-schema/main.go

.PHONY: gen-rollouts-openapi
gen-rollouts-openapi: openapi-gen
	$(OPENAPI_GEN) \
	  	--go-header-file hack/boilerplate.go.txt \
		--input-dirs github.com/openkruise/kruise-rollout-api/rollouts/v1beta1,github.com/openkruise/kruise-rollout-api/rollouts/v1alpha1 \
		--output-package pkg/rollouts/ \
  		--report-filename pkg/rollouts/violation_exceptions.list \
  		-o $(CURRENT_DIR)
