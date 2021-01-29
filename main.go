package main

import (
	"context"
	"flag"
	"time"

	"k8s.io/klog/v2"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	appsv1mf "k8s.io/client-go/applyconfigurations/apps/v1"
	metav1mf "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags("", "/home/soltysh/kubeconfig_b")
	if err != nil {
		klog.Fatal(err)
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatal(err)
	}

	deploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "hello-node",
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "hello-node"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "hello-node"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "hello-node",
							Image: "k8s.gcr.io/serve_hostname",
						},
					},
				},
			},
		},
	}

	unstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(deploy)
	if err != nil {
		klog.Fatal(err)
	}
	deploymf := appsv1mf.Deployment()
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, deploymf)
	if err != nil {
		klog.Fatal(err)
	}
	actual, err := kubeClient.AppsV1().Deployments("soltysh").Apply(context.TODO(), deploymf, "test-client", metav1.ApplyOptions{})
	if err != nil {
		klog.Fatal(err)
	}
	klog.Infof("Created:\n%#v\n", actual)

	time.Sleep(5 * time.Second)

	deploymf.
		SetObjectMeta(metav1mf.ObjectMeta().
			SetName("hello-node").
			SetAnnotations(map[string]string{"myannotation": "myvalue"}),
		).
		SetSpec(appsv1mf.DeploymentSpec().
			SetReplicas(int32(2)),
		)

	actual, err = kubeClient.AppsV1().Deployments("soltysh").Apply(context.TODO(), deploymf, "test-client", metav1.ApplyOptions{})
	if err != nil {
		klog.Fatal(err)
	}
	klog.Infof("Updated:\n%#v\n", actual)
}
