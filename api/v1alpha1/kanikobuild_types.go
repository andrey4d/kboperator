/*
Copyright 2024 andrey4d.dev@gmail.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KanikoBuildSpec defines the desired state of KanikoBuild.
type KanikoBuildSpec struct {
	Name         string            `json:"name,omitempty" default:"builder"`
	Image        string            `json:"image,omitempty" default:"gcr.io/kaniko-project/executor:latest"`
	Context      string            `json:"context,omitempty" default:"/workspace"`
	Dockerfile   string            `json:"dockerfile,omitempty"`
	Destination  string            `json:"destination,omitempty"`
	Certificate  string            `json:"certificate,omitempty"`
	DockerConfig DockerConfig      `json:"docker_config,omitempty"`
	Command      []string          `json:"command,omitempty"`
	Args         []string          `json:"args,omitempty"`
	Persistence  PersistenceVolume `json:"persistence,omitempty"`
}

type DockerConfig struct {
	Registry string `json:"registry,omitempty"`
	Auth     string `json:"auth,omitempty"`
}

type PersistenceVolume struct {
	Enabled      bool     `json:"enabled,omitempty"`
	VolumeSize   string   `json:"volumeSize,omitempty" Defaults:"10Gi"`
	StorageClass string   `json:"storageClass,omitempty"`
	ExtraVolumes []Volume `json:"extraVolumes,omitempty"`
}

type Volume struct {
	corev1.Volume `json:",inline"`
	// MountPath is the path where this volume should be mounted
	MountPath string `json:"mountPath"`
}

// KanikoBuildStatus defines the observed state of KanikoBuild.
type KanikoBuildStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KanikoBuild is the Schema for the kanikobuilds API.
type KanikoBuild struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KanikoBuildSpec   `json:"spec,omitempty"`
	Status KanikoBuildStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KanikoBuildList contains a list of KanikoBuild.
type KanikoBuildList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KanikoBuild `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KanikoBuild{}, &KanikoBuildList{})
}
