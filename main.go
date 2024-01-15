package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	projectCrd "github.com/martin-helmich/kubernetes-crd-example/api/types/v1alpha1"
	projectClient "github.com/martin-helmich/kubernetes-crd-example/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
}

func main() {

	var config *rest.Config
	var err error

	if kubeconfig == "" {
		log.Printf("using in-cluster configuration")
		// return a *rest.Config
		config, err = rest.InClusterConfig()
	} else {
		log.Printf("using configuration from '%s'", kubeconfig)
		// return *rest.Config
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		panic(err)
	}

	// register the Project and ProjectList schema structures
	projectCrd.AddToScheme(scheme.Scheme)

	// get a customresource client from the given config
	projectClientSet, err := projectClient.NewForConfig(config)

	if err != nil {
		panic(err)
	}

	projectList, err := projectClientSet.ProjectsCrdV1Alpha1().Projects("default").List(metav1.ListOptions{})

	if err != nil {
		panic(err)
	}

	fmt.Printf("projects found: %+v\n", projectList)

	store := WatchResources(*projectClientSet.ProjectsCrdV1Alpha1(), "default")

	for {
		projectsFromStore := store.List()
		fmt.Printf("project in store: %d\n", len(projectsFromStore))

		time.Sleep(2 * time.Second)
	}
}
