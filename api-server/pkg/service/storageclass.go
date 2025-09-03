package service

import (
	"context"
	storagev1 "k8s.io/api/storage/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeall.io/api-server/pkg/infra/apiserver"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type StorageClass interface {
	Get(ctx context.Context, name string) (*storagev1.StorageClass, error)
	Create(ctx context.Context, storageClass *storagev1.StorageClass) (*storagev1.StorageClass, error)
}

type storageClassImpl struct {
	clusterResource apiserver.ClusterResource
}

func NewStorageClass(clusterResource apiserver.ClusterResource) StorageClass {
	return &storageClassImpl{clusterResource: clusterResource}
}

func (s storageClassImpl) Get(ctx context.Context, name string) (*storagev1.StorageClass, error) {
	var sc = &storagev1.StorageClass{}
	err := s.clusterResource.ClusterCache().Get(ctx, client.ObjectKey{
		Name: name,
	}, sc)
	return sc, err
}

func (s storageClassImpl) Create(ctx context.Context, storageClass *storagev1.StorageClass) (*storagev1.StorageClass, error) {
	sc, err := s.clusterResource.Client().K8sClient().StorageV1().StorageClasses().
		Create(ctx, storageClass, v1.CreateOptions{})
	return sc, err
}
