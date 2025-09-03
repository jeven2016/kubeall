/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +enum
type ImageType string

const (
	ImageIso  ImageType = "iso"
	ImageDisk ImageType = "disk"
)

// +enum
type StorageBackend string

const (
	StorageBackendBackingImage StorageBackend = "backingimage"
	StorageBackendCDI          StorageBackend = "cdi"
)

// +enum
type ImageSourceType string

const (
	ImageSourceTypeUpload       ImageSourceType = "upload"
	ImageSourceTypeDownload     ImageSourceType = "download"
	ImageSourceTypeRestore      ImageSourceType = "restore"
	ImageSourceTypeClone        ImageSourceType = "clone"
	ImageSourceTypeExportVolume ImageSourceType = "export-from-volume"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ImageSpec defines the desired state of Image.
type ImageSpec struct {
	// +optional
	OsType string `json:"osType,omitempty"`

	// +optional
	OsVersion string `json:"osVersion,omitempty"`

	// +optional
	// +kubebuilder:default=disk
	// +kubebuilder:validation:Enum=iso;disk
	ImageType ImageType `json:"imageType"`

	// +optional
	StorageClassName string `json:"storageClassName,omitempty"`

	// +optional
	// +kubebuilder:validation:Optional
	ImageStorageClassName string `json:"imageStorageClassName,omitempty"`

	// +optional
	// +kubebuilder:default=backingimage
	// +kubebuilder:validation:Enum=backingimage;cdi
	StorageBackend StorageBackend `json:"backend"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=download;upload;restore;export-from-volume;clone
	ImageFrom ImageSourceType `json:"imageFrom"`

	// +optional
	SourceStorageClassName string `json:"sourceStorageClassName,omitempty"`

	// +optional
	StorageClassParameters map[string]string `json:"storageClassParameters"`
}

// ImageStatus defines the observed state of Image.
// status 的更新不会直接触发控制器的协调逻辑，因为 status 仅反映当前状态，而不代表用户意图。控制器通常只对 spec 的变化做出反应。
type ImageStatus struct {
	// +optional
	Progress int `json:"progress,omitempty"`

	// +optional
	Size int64 `json:"size,omitempty"`

	// +optional
	VirtualSize int64 `json:"virtualSize,omitempty"`

	// +optional
	State string `json:"state,omitempty"`

	// +optional
	Message string `json:"message,omitempty"`

	// +optional
	LastStateTransitionTime string `json:"lastStateTransitionTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Image is the Schema for the images API.
type Image struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageSpec   `json:"spec,omitempty"`
	Status ImageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ImageList contains a list of Image.
type ImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Image `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Image{}, &ImageList{})
}
