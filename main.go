package main

import (
	"context"
	"flag"
	"fmt"
	alchemistclientset "github.com/tapojit047/CRD/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	_ "k8s.io/code-generator"
	"log"
	"path/filepath"
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
	client, err := alchemistclientset.NewForConfig(config)
	if err != nil {
		log.Printf("Building config from flags, %s", err.Error())
	}
	alchemists, err := client.FullmetalV1().Alchemists("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {

	}
	fmt.Printf("length of alchemists is %d and name is %s\n", len(alchemists.Items), alchemists.Items[0].Name)
}
