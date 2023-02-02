package main

import (
	"flag"
	alchClientset "github.com/tapojit047/CRD-Controller/pkg/client/clientset/versioned"
	alchInformer "github.com/tapojit047/CRD-Controller/pkg/client/informers/externalversions"
	kubeInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	_ "k8s.io/code-generator"
	"log"
	"path/filepath"
	"time"
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
		log.Printf("Building config from flags, %s\n", err.Error())
		panic(err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Building kubernetes clientset, %s\n", err.Error())
		panic(err)
	}

	alchemistClient, err := alchClientset.NewForConfig(config)
	if err != nil {
		log.Printf("Building alchemist clientset, %s\n", err.Error())
		panic(err)
	}

	kubeInfoFactory := kubeInformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	alchemistInfoFactory := alchInformer.NewSharedInformerFactory(alchemistClient, time.Second*30)
	controller := NewController(alchemistClient, alchemistInfoFactory.Fullmetal().V1().Alchemists(), kubeClient, kubeInfoFactory.Apps().V1().Deployments())

	ch := make(chan struct{})
	alchemistInfoFactory.Start(ch)
	kubeInfoFactory.Start(ch)

	if err := controller.Run(ch); err != nil {
		log.Printf("error running controller %s\n", err.Error())
	}
}
