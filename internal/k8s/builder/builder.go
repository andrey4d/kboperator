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
	BuilderCrd *kbov1alpha1.KanikoBuild
}

func NewBuilder(k *kbov1alpha1.KanikoBuild) *Builder {
	return &Builder{BuilderCrd: k}
}

func (b *Builder) LabelsForBuilder() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":     b.BuilderName(),
		"app.kubernetes.io/instance": b.BuilderName(),
	}
}

func (b *Builder) PodVolumes() []corev1.Volume {
	return []corev1.Volume{
		{Name: "dockerfile", VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				DefaultMode:          &[]int32{0644}[0],
				LocalObjectReference: corev1.LocalObjectReference{Name: b.BuilderName()},
			},
		}},
	}
}

func (b *Builder) VolumesMount() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{Name: "dockerfile", MountPath: b.GetContext() + "/Dockerfile", SubPath: "Dockerfile"},
		{Name: "dockerfile", MountPath: "/kaniko/.docker/config.json", SubPath: "config.json"},
	}
}

func (b *Builder) Command() []string {
	if len(b.BuilderCrd.Spec.Command) > 0 {
		return b.BuilderCrd.Spec.Command
	}
	return []string{
		"/kaniko/executor",
		"--context=" + b.GetContext(),
		"--dockerfile=" + b.GetContext() + "/Dockerfile",
		"--destination=" + b.BuilderCrd.Spec.Destination,
	}
}

func (b *Builder) Args() []string {
	if len(b.BuilderCrd.Spec.Args) > 0 {
		return b.BuilderCrd.Spec.Args
	}
	return []string{}
}

func (b *Builder) GetContext() string {
	if b.BuilderCrd.Spec.Context != "" {
		return b.BuilderCrd.Spec.Context
	}
	return "/workspace"
}

func (b *Builder) BuilderImage(k *kbov1alpha1.KanikoBuild) string {
	if b.BuilderCrd.Spec.Image != "" {
		return b.BuilderCrd.Spec.Image
	}
	return "gcr.io/kaniko-project/executor:latest"
}

func (b *Builder) BuilderName() string {
	if b.BuilderCrd.Spec.Name != "" {
		return b.BuilderCrd.Spec.Name
	}
	return "builder"
}
