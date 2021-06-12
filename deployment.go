package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	//"time"
	//"k8s.io/apimachinery/pkg/api/errors"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	//pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceAll)
	dep, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	// for _, d := range list.Items {
	// 	//fmt.Printf(" * %s (%d replicas)\n", d.Spec.Template.Spec, *d.Spec.Replicas)
	// 	//fmt.Println(d.version)
	// }
	//if err != nil {
	//panic(err.Error())
	//}
	//fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	for _, pod := range dep.Items {
		i := pod.GetObjectMeta().GetAnnotations()["kubectl.kubernetes.io/last-applied-configuration"]
		var result map[string]interface{}
		json.Unmarshal([]byte(i), &result)
		fmt.Println(pod.GetName(), pod.GetNamespace(), result["apiVersion"], result["kind"])
	}
}
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}
