package e2e

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	rolloutsv1beta1 "github.com/openkruise/kruise-rollout-api/client/clientset/versioned"
	"github.com/openkruise/kruise-rollout-api/rollouts/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var _ = Describe("Rollout E2E Tests", func() {
	var (
		clientset *rolloutsv1beta1.Clientset
		namespace = "default"
	)

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	Expect(err).NotTo(HaveOccurred())

	clientset, err = rolloutsv1beta1.NewForConfig(config)
	Expect(err).NotTo(HaveOccurred())

	Context("Rollout Operations", func() {
		var rolloutDemo *v1beta1.Rollout

		BeforeEach(func() {
			firststep := intstr.FromString("10%")

			rolloutDemo = &v1beta1.Rollout{
				ObjectMeta: metav1.ObjectMeta{
					Name: "demo-rollout",
				},
				Spec: v1beta1.RolloutSpec{
					WorkloadRef: v1beta1.ObjectRef{
						APIVersion: "apps/v1",
						Kind:       "Deployment",
						Name:       "demo-deploy",
					},
					Strategy: v1beta1.RolloutStrategy{
						Canary: &v1beta1.CanaryStrategy{
							Steps: []v1beta1.CanaryStep{
								{
									Replicas: &firststep,
								},
							},
						},
					},
				},
			}
		})

		It("should create a Rollout", func() {
			rollout := rolloutDemo.DeepCopy()
			rollout.Name = "test-create"
			result, err := clientset.RolloutsV1beta1().Rollouts(namespace).Create(context.TODO(), rollout, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Name).To(Equal(rollout.Name))
		})

		It("should update a Rollout", func() {
			rollout := rolloutDemo.DeepCopy()
			rollout.Name = "test-update"
			result, err := clientset.RolloutsV1beta1().Rollouts(namespace).Create(context.TODO(), rollout, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			result.Spec.Strategy.Canary.Steps[0].Replicas = &intstr.IntOrString{
				Type:   intstr.String,
				StrVal: "20%",
			}
			updatedResult, err := clientset.RolloutsV1beta1().Rollouts(namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedResult.Spec.Strategy.Canary.Steps[0].Replicas.StrVal).To(Equal("20%"))
		})

		It("should delete a Rollout", func() {
			rollout := rolloutDemo.DeepCopy()
			rollout.Name = "test-delete"
			result, err := clientset.RolloutsV1beta1().Rollouts(namespace).Create(context.TODO(), rollout, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			err = clientset.RolloutsV1beta1().Rollouts(namespace).Delete(context.TODO(), result.Name, metav1.DeleteOptions{})
			Expect(err).NotTo(HaveOccurred())

			_, err = clientset.RolloutsV1beta1().Rollouts(namespace).Get(context.TODO(), result.Name, metav1.GetOptions{})
			Expect(err).To(HaveOccurred())
		})
	})
})

func TestRolloutE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rollout E2E Tests")
}
