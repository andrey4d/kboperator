/*
 *   Copyright (c) 2024 Andrey andrey4d.dev@gmail.com
 *   All rights reserved.
 */
package jobs

import (
	kbov1alpha1 "github.com/andrey4d/kboperator/api/v1alpha1"
	"github.com/andrey4d/kboperator/internal/k8s/builder"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Job struct {
	Namespace string
	Scheme    *runtime.Scheme
	builder   *builder.Builder
}

func NewJob(k *kbov1alpha1.KanikoBuild, scheme *runtime.Scheme) *Job {
	return &Job{
		Namespace: k.Namespace,
		Scheme:    scheme,
		builder:   builder.NewBuilder(k),
	}
}

func (j *Job) BuilderJob() (*batchv1.Job, error) {

	containers := []corev1.Container{{
		Image:           j.builder.BuilderImage(j.builder.Crd),
		Name:            j.builder.BuilderName(),
		ImagePullPolicy: corev1.PullIfNotPresent,
		Command:         j.builder.Command(),
		Args:            j.builder.Args(),
		VolumeMounts:    j.builder.VolumesMount(),
	}}

	template := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: j.builder.LabelsForBuilder(),
		},
		Spec: corev1.PodSpec{
			Containers:    containers,
			Volumes:       j.builder.PodVolumes(),
			RestartPolicy: "OnFailure",
		},
	}

	spec := batchv1.JobSpec{
		Template:     template,
		BackoffLimit: &[]int32{3}[0],
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      j.builder.BuilderName(),
			Namespace: j.Namespace,
		},
		Spec: spec,
	}

	if err := ctrl.SetControllerReference(j.builder.Crd, job, j.Scheme); err != nil {
		return nil, err
	}

	return job, nil

}
