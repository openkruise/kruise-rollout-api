package main

import (
	"context"
	"flag"
	"fmt"

	rolloutsv1beta1 "github.com/openkruise/kruise-rollout-api/client/clientset/versioned"
	rolloutapi "github.com/openkruise/kruise-rollout-api/rollouts/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := rolloutsv1beta1.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	client := clientset.RolloutsV1beta1().Rollouts("default")

	firststep := intstr.FromString("10%")

	rollout := rolloutapi.Rollout{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-rollout",
		},

		Spec: rolloutapi.RolloutSpec{
			WorkloadRef: rolloutapi.ObjectRef{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       "demo-deploy",
			},
			Strategy: rolloutapi.RolloutStrategy{
				Canary: &rolloutapi.CanaryStrategy{
					Steps: []rolloutapi.CanaryStep{
						{
							Replicas: &firststep,
						},
					},
				},
			},
		},
	}
	result, err := client.Create(context.TODO(), &rollout, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created rollout %q.\n", result.GetObjectMeta().GetName())

}
