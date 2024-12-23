/*
 *   Copyright (c) 2024 Andrey andrey4d.dev@gmail.com
 *   All rights reserved.
 */
package persistence

import (
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	kbov1alpha1 "github.com/andrey4d/kboperator/api/v1alpha1"
	"github.com/andrey4d/kboperator/internal/k8s/builder"
)

type Persistence struct {
	Namespace string
	Scheme    *runtime.Scheme
	builder   *builder.Builder
}

func NewPersistence(k *kbov1alpha1.KanikoBuild, scheme *runtime.Scheme) *Persistence {
	return &Persistence{
		Namespace: k.Namespace,
		Scheme:    scheme,
		builder:   builder.NewBuilder(k),
	}
}

func (p *Persistence) BuilderPvc() (*corev1.PersistentVolumeClaim, error) {
	accessMode := corev1.ReadWriteOnce

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.builder.BuilderName(),
			Namespace: p.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				accessMode,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(p.builder.VolumeSize()),
				},
			},
		},
	}

	if len(p.builder.Crd.Spec.Persistence.StorageClass) > 0 {
		pvc.Spec.StorageClassName = &p.builder.Crd.Spec.Persistence.StorageClass
	}

	if err := ctrl.SetControllerReference(p.builder.Crd, pvc, p.Scheme); err != nil {
		return nil, err
	}
	return pvc, nil
}

func (p *Persistence) ExtraPvcs() ([]*corev1.PersistentVolumeClaim, error) {
	return nil, nil
}
