package main

import (
	"context"
	"fmt"
	alchemistv1 "github.com/tapojit047/CRD-Controller/pkg/apis/fullmetal.com/v1"
	alchClientset "github.com/tapojit047/CRD-Controller/pkg/client/clientset/versioned"
	alchInformer "github.com/tapojit047/CRD-Controller/pkg/client/informers/externalversions/fullmetal.com/v1"
	alchLister "github.com/tapojit047/CRD-Controller/pkg/client/listers/fullmetal.com/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsInformer "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes" //clientset
	appsLister "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
	"time"
)

// Controller is the custom controller implementation for the Alchemist custom resource
type Controller struct {
	// clientset for custom resource Alchemist
	alchemistclient alchClientset.Interface
	// The cache of informer for Alchemist has Synced or not
	alchemistSynced cache.InformerSynced
	// Lister for the informer of the alchemist
	alchemistLister alchLister.AlchemistLister

	// clientset for kubernetes
	kubeclient kubernetes.Interface
	// Lister for the informer of the deployment
	deploymentLister appsLister.DeploymentLister
	// Cache of informer for deployment
	deploymentSynced cache.InformerSynced

	// And finally a queue which will hold all the works sequentially
	workQueue workqueue.RateLimitingInterface
}

// NewController returns a new sample controller
func NewController(alchemistClient alchClientset.Interface, alchemistInformer alchInformer.AlchemistInformer,
	kubeClient kubernetes.Interface, deploymentInformer appsInformer.DeploymentInformer) *Controller {

	controller := &Controller{
		alchemistclient: alchemistClient,
		alchemistSynced: alchemistInformer.Informer().HasSynced,
		alchemistLister: alchemistInformer.Lister(),

		kubeclient:       kubeClient,
		deploymentSynced: deploymentInformer.Informer().HasSynced,
		deploymentLister: deploymentInformer.Lister(),

		workQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Alchemist"),
	}

	alchemistInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueAlchemist,
		UpdateFunc: func(oldObj, newObj interface{}) {
			controller.enqueueAlchemist(newObj)
		},
	})
	return controller
}

func (controller *Controller) Run(ch chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer controller.workQueue.ShutDown()

	if ok := cache.WaitForCacheSync(ch, controller.alchemistSynced, controller.deploymentSynced); !ok {
		log.Println("cache was not synched")
	}
	log.Println("Alchemy started . . .")

	go wait.Until(controller.runWorker, time.Second, ch)
	<-ch
	return nil
}

func (controller *Controller) runWorker() {
	for controller.processNextWorkItem() {

	}
}

func (controller *Controller) processNextWorkItem() bool {
	// The Get() call gets the new item from the queue and deletes it
	item, shutdown := controller.workQueue.Get()
	if shutdown {
		return false
	}

	// We wrap this block in a func so, we can defer controller.workQueue.Done
	err := func(item interface{}) error {
		defer controller.workQueue.Done(item)

		var key string
		var ok bool
		if key, ok = item.(string); !ok {
			controller.workQueue.Forget(item)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", item))
			return nil
		}

		if err := controller.syncHandler(key); err != nil {
			controller.workQueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		controller.workQueue.Forget(item)
		return nil
	}(item)

	if err != nil {
		utilruntime.HandleError(err)
	}
	return true
}

func (controller *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invald resource key: %s", key))
		return nil
	}

	alchemist, err := controller.alchemistLister.Alchemists(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("alchemist '%s' in work queue, no longer exists", key))
			return nil
		}
		return err
	}

	deploymentName := alchemist.Spec.DeploymentName
	if deploymentName == "" {
		utilruntime.HandleError(fmt.Errorf("%s: deployment name can not be empty", key))
		return nil
	}

	deployment, err := controller.deploymentLister.Deployments(namespace).Get(deploymentName)
	if errors.IsNotFound(err) {
		deployment, err = controller.kubeclient.AppsV1().Deployments(namespace).Create(context.TODO(), newDeployment(alchemist), metav1.CreateOptions{})
	}
	if err != nil {
		return err
	}

	if alchemist.Spec.Replicas != nil && *alchemist.Spec.Replicas != *deployment.Spec.Replicas {
		log.Printf("Alchemist %s replicas: %d, deployment replicas: %d", name, *alchemist.Spec.Replicas, *deployment.Spec.Replicas)
		deployment, err = controller.kubeclient.AppsV1().Deployments(namespace).Update(context.TODO(), newDeployment(alchemist), metav1.UpdateOptions{})
	}

	if err != nil {
		return err
	}

	err = controller.updateAlchemistStatus(alchemist, deployment)
	if err != nil {
		return err
	}
	return nil
}

func (contoller *Controller) enqueueAlchemist(obj interface{}) {
	log.Println("Enqueueing Alchemist. . . ")
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	contoller.workQueue.AddRateLimited(key)
}

func (controller *Controller) updateAlchemistStatus(alchemist *alchemistv1.Alchemist, deployment *appsv1.Deployment) error {
	alchemistCopy := alchemist.DeepCopy()
	alchemistCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas

	_, err := controller.alchemistclient.FullmetalV1().Alchemists(alchemist.Namespace).Update(context.TODO(), alchemistCopy, metav1.UpdateOptions{})
	return err
}

func newDeployment(alchemist *alchemistv1.Alchemist) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      alchemist.Spec.DeploymentName,
			Namespace: alchemist.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: alchemist.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":        "api-server",
					"controller": alchemist.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":        "api-server",
						"controller": alchemist.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "api-server",
							Image: "tapojit047/api-server",
						},
					},
				},
			},
		},
	}
}
