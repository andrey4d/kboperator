/*
 *   Copyright (c) 2024 Andrey andrey4d.dev@gmail.com
 *   All rights reserved.
 */
package deployments

import (
	kbov1alpha1 "github.com/andrey4d/kboperator/api/v1alpha1"
	"github.com/andrey4d/kboperator/internal/k8s/builder"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Deployment struct {
	Name      string
	Namespace string
	Replicas  int32
	Scheme    *runtime.Scheme
}

func NewDeployment(scheme *runtime.Scheme, namespace string) *Deployment {
	return &Deployment{
		Scheme:    scheme,
		Namespace: namespace,
		Replicas:  1}
}

func (d *Deployment) BuilderDeployment(k *kbov1alpha1.KanikoBuild) (*appsv1.Deployment, error) {
	builder := builder.NewBuilder(k)
	containers := []corev1.Container{{
		Image:           k.Spec.Image,
		Name:            k.Spec.Name,
		ImagePullPolicy: corev1.PullIfNotPresent,
		// Ensure restrictive context for the container
		// More info: https://kubernetes.io/docs/concepts/security/pod-security-standards/#restricted
		SecurityContext: d.ContainerSecurityContext(),
		Command:         k.Spec.Command,
		VolumeMounts:    builder.VolumesMount(),
	}}

	template := corev1.PodTemplateSpec{
		ObjectMeta: d.PodMetadata(k),
		Spec: corev1.PodSpec{
			// Affinity: d.PodNodeAffinity(),
			SecurityContext: d.PodSecurityContext(),
			Containers:      containers,
			Volumes:         builder.PodVolumes(),
		},
	}

	spec := appsv1.DeploymentSpec{
		Replicas: &d.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: builder.LabelsForBuilder(),
		},
		Template: template,
	}

	deploy := &appsv1.Deployment{
		ObjectMeta: d.PodMetadata(k),
		Spec:       spec,
	}

	if err := ctrl.SetControllerReference(k, deploy, d.Scheme); err != nil {
		return nil, err
	}

	return deploy, nil
}

func (d *Deployment) PodMetadata(k *kbov1alpha1.KanikoBuild) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      k.Name,
		Namespace: d.Namespace,
		Labels:    builder.NewBuilder(k).LabelsForBuilder(),
	}
}

func (d *Deployment) PodSecurityContext() *corev1.PodSecurityContext {
	// IMPORTANT: seccomProfile was introduced with Kubernetes 1.19
	// If you are looking for to produce solutions to be supported
	// on lower versions you must remove this option.
	return &corev1.PodSecurityContext{
		RunAsNonRoot: &[]bool{true}[0],
		SeccompProfile: &corev1.SeccompProfile{
			Type: corev1.SeccompProfileTypeRuntimeDefault,
		},
	}
}

func (d *Deployment) PodNodeAffinity() *corev1.Affinity {
	// More info: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#node-affinity
	return &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{
					{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      "kubernetes.io/arch",
								Operator: "In",
								Values:   []string{"amd64", "arm64", "ppc64le", "s390x"},
							},
							{
								Key:      "kubernetes.io/os",
								Operator: "In",
								Values:   []string{"linux"},
							},
						},
					},
				},
			},
		},
	}
}

func (d *Deployment) PodNodeSelector() map[string]string {
	// More info: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#node-affinity
	return map[string]string{
		"kubernetes.io/arch": "amd64",
		"kubernetes.io/os":   "linux",
	}
}

func (d *Deployment) ContainerSecurityContext() *corev1.SecurityContext {
	// WARNING: Ensure that the image used defines an UserID in the Dockerfile
	// otherwise the Pod will not run and will fail with "container has runAsNonRoot and image has non-numeric user"".
	// If you want your workloads admitted in namespaces enforced with the restricted mode in OpenShift/OKD vendors
	// then, you MUST ensure that the Dockerfile defines a User ID OR you MUST leave the "RunAsNonRoot" and
	// "RunAsUser" fields empty.
	return &corev1.SecurityContext{
		RunAsNonRoot:             &[]bool{true}[0],
		RunAsUser:                &[]int64{1001}[0],
		AllowPrivilegeEscalation: &[]bool{false}[0],
		Capabilities: &corev1.Capabilities{
			Drop: []corev1.Capability{
				"ALL",
			},
		},
	}
}
