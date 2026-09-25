package main

import (
	"fmt"

	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		cfg := config.New(ctx, "")
		kubeContext := cfg.Require("kubeContext")
		namespaceName := cfg.Require("namespace")
		deploymentName := cfg.Require("deploymentName")
		image := cfg.Require("image")
		tag := cfg.Require("tag")
		replicas := cfg.RequireInt("replicas")
		serviceName := cfg.Require("serviceName")
		servicePort := cfg.RequireInt("servicePort")
		targetPort := cfg.RequireInt("targetPort")

		provider, err := kubernetes.NewProvider(ctx, "k8s-provider", &kubernetes.ProviderArgs{
			Context: pulumi.String(kubeContext),
		})

		if err != nil {
			return err
		}

		appLabels := pulumi.StringMap{
			"app": pulumi.String(deploymentName),
		}

		namespace, err := corev1.NewNamespace(ctx, namespaceName, &corev1.NamespaceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String(namespaceName),
			},
		}, pulumi.Provider(provider))
		if err != nil {
			return err
		}

		deployment, err := appsv1.NewDeployment(ctx, deploymentName, &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String(deploymentName),
				Namespace: namespace.Metadata.Name(),
			},
			Spec: appsv1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: appLabels,
				},
				Replicas: pulumi.Int(replicas),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Name:   pulumi.String(deploymentName),
						Labels: appLabels,
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							corev1.ContainerArgs{
								Name:  pulumi.String(deploymentName),
								Image: pulumi.String(image + ":" + tag),
								Env: corev1.EnvVarArray{
									corev1.EnvVarArgs{
										Name:  pulumi.String("PORT"),
										Value: pulumi.String(fmt.Sprint(targetPort)),
									},
								},
							}},
					},
				},
			},
		}, pulumi.Provider(provider))
		if err != nil {
			return err
		}

		service, err := corev1.NewService(ctx, serviceName, &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String(serviceName),
				Namespace: namespace.Metadata.Name(),
			},
			Spec: &corev1.ServiceSpecArgs{
				Selector: appLabels,
				Ports: corev1.ServicePortArray{
					corev1.ServicePortArgs{
						Name:       pulumi.String("http"),
						Port:       pulumi.Int(servicePort),
						TargetPort: pulumi.Int(targetPort),
					},
				},
			},
		}, pulumi.Provider(provider))
		if err != nil {
			return err
		}

		ctx.Export("name", deployment.Metadata.Name())
		ctx.Export("namespace", namespace.Metadata.Name())
		ctx.Export("service_name", service.Metadata.Name())

		return nil
	})
}
