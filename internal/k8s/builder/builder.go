/*
 *   Copyright (c) 2024 Andrey andrey4d.dev@gmail.com
 *   All rights reserved.
 */
package builder

import (
	kbov1alpha1 "github.com/andrey4d/kboperator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

type Builder struct {
	Crd *kbov1alpha1.KanikoBuild
}

func NewBuilder(k *kbov1alpha1.KanikoBuild) *Builder {
	return &Builder{Crd: k}
}

func (b *Builder) LabelsForBuilder() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":     b.BuilderName(),
		"app.kubernetes.io/instance": b.BuilderName(),
	}
}

func (b *Builder) PodVolumes() []corev1.Volume {

	volumes := []corev1.Volume{{
		Name: "dockerfile",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				DefaultMode:          &[]int32{0644}[0],
				LocalObjectReference: corev1.LocalObjectReference{Name: b.BuilderName()},
			},
		}},
	}

	if b.Crd.Spec.Persistence.Enabled {
		volumes = append(volumes, corev1.Volume{Name: "workspace", VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: b.BuilderName(),
				ReadOnly:  false,
			},
		}})
	}
	return volumes
}

func (b *Builder) VolumesMount() []corev1.VolumeMount {
	volumesMount := []corev1.VolumeMount{
		{Name: "dockerfile", MountPath: "/kaniko/.docker/config.json", SubPath: "config.json"},
	}
	if b.Crd.Spec.Dockerfile != "" {
		volumesMount = append(volumesMount,
			corev1.VolumeMount{Name: "dockerfile", MountPath: "/kaniko/Dockerfile", SubPath: "Dockerfile"})
	}
	if b.Crd.Spec.Persistence.Enabled {
		volumesMount = append(volumesMount,
			corev1.VolumeMount{Name: "workspace", MountPath: b.GetContext(), ReadOnly: false})
	}
	return volumesMount
}

func (b *Builder) Command() []string {
	if len(b.Crd.Spec.Command) > 0 {
		return b.Crd.Spec.Command
	}
	return []string{
		"/kaniko/executor",
		"--context=" + b.GetContext(),
		"--dockerfile=/kaniko/Dockerfile",
		"--destination=" + b.Crd.Spec.Destination,
	}
}

func (b *Builder) Args() []string {
	if len(b.Crd.Spec.Args) > 0 {
		return b.Crd.Spec.Args
	}
	return []string{}
}

func (b *Builder) GetContext() string {
	if b.Crd.Spec.Context != "" {
		return b.Crd.Spec.Context
	}
	return "/workspace"
}

func (b *Builder) BuilderImage(k *kbov1alpha1.KanikoBuild) string {
	if b.Crd.Spec.Image != "" {
		return b.Crd.Spec.Image
	}
	return "gcr.io/kaniko-project/executor:latest"
}

func (b *Builder) BuilderName() string {
	if b.Crd.Spec.Name != "" {
		return b.Crd.Spec.Name
	}
	return "builder"
}

func (b *Builder) VolumeSize() string {
	if b.Crd.Spec.Persistence.VolumeSize != "" {
		return b.Crd.Spec.Persistence.VolumeSize
	}
	return "1Gi"
}
