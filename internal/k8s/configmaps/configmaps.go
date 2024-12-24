/*
 *   Copyright (c) 2024 Andrey andrey4d.dev@gmail.com
 *   All rights reserved.
 */
package configmaps

import (
	kbov1alpha1 "github.com/andrey4d/kboperator/api/v1alpha1"
	"github.com/andrey4d/kboperator/internal/k8s/builder"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ConfigMap struct {
	Namespace string
	Scheme    *runtime.Scheme
	builder   *builder.Builder
}

func NewConfigMap(k *kbov1alpha1.KanikoBuild, scheme *runtime.Scheme) *ConfigMap {
	return &ConfigMap{
		Namespace: k.Namespace,
		Scheme:    scheme,
		builder:   builder.NewBuilder(k),
	}
}

func (c *ConfigMap) BuilderConfigMap() (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{ObjectMeta: c.builder.Metadata(),
		Data: map[string]string{
			"Dockerfile":  c.builder.Crd.Spec.Dockerfile,
			"config.json": dockerConfig(c.builder.Crd.Spec.DockerConfig.Registry, c.builder.Crd.Spec.DockerConfig.Auth),
		},
	}

	if err := ctrl.SetControllerReference(c.builder.Crd, cm, c.Scheme); err != nil {
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
