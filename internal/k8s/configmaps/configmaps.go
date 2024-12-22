/*
 *   Copyright (c) 2024 Andrey andrey4d.dev@gmail.com
 *   All rights reserved.
 */
package configmaps

import (
	kbov1alpha1 "github.com/andrey4d/kboperator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ConfigMap struct {
	Scheme    *runtime.Scheme
	Namespace string
}

func NewConfigMap(scheme *runtime.Scheme, namespace string) *ConfigMap {
	return &ConfigMap{Scheme: scheme, Namespace: namespace}
}

func (c *ConfigMap) BuilderConfigMap(k *kbov1alpha1.KanikoBuild) (*corev1.ConfigMap, error) {

	metadata := metav1.ObjectMeta{Name: k.Name, Namespace: c.Namespace}

	cm := &corev1.ConfigMap{ObjectMeta: metadata,
		Data: map[string]string{
			"Dockerfile":  k.Spec.Dockerfile,
			"config.json": dockerConfig(k.Spec.DockerConfig.Registry, k.Spec.DockerConfig.Auth),
		},
	}

	if err := ctrl.SetControllerReference(k, cm, c.Scheme); err != nil {
		return nil, err
	}

	return cm, nil
}

func dockerConfig(registry, auth string) string {
	return `{
	"auths": {
		"` + registry + `": {
		  "auth": "` + auth + `"
		}
	  }
	}`
}
