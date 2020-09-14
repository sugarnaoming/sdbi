package kubernetes

import (
	"context"
	"strings"

	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Kube struct {
	ns     string
	client *kubernetes.Clientset
}

func New() (*Kube, error) {
	kube := new(Kube)

	kubeConfigPath := clientcmd.NewDefaultClientConfigLoadingRules()
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(kubeConfigPath, nil)
	ns, _, err := config.Namespace()
	if err != nil {
		//TODO error handling
		return kube, err
	}
	kube.ns = ns

	clientConfig, err := config.ClientConfig()
	if err != nil {
		// TODO error handling
		return kube, err
	}

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		// TODO error handling
		return kube, err
	}
	kube.client = clientset

	return kube, nil
}

func (k *Kube) BuildPods() (*v1.PodList, error) {
	return k.client.CoreV1().Pods(k.ns).List(context.TODO(), metav1.ListOptions{LabelSelector: "tier=builds"})
}

func (k *Kube) NodeFilter(pods *v1.PodList, nodeNames ...string) []v1.Pod {
	var filtered []v1.Pod

	for _, pod := range pods.Items {
		for _, nodeName := range nodeNames {
			if strings.Contains(pod.Spec.NodeName, nodeName) {
				filtered = append(filtered, pod)
			}
		}
	}

	return filtered
}

func (k *Kube) ExtractId(pods []v1.Pod) []string {
	var ids []string

	for _, pod := range pods {
		ids = append(ids, pod.Labels["sdbuild"])
	}

	return ids
}
