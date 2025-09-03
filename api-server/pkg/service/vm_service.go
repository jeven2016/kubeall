package service

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeall.io/api-server/pkg/infra/apiserver"
	"kubeall.io/api-server/pkg/infra/constants"
	kv1 "kubevirt.io/api/core/v1"
)

type VmService interface {
	Create(ctx context.Context, vm *kv1.VirtualMachine) error
	DeleteDisks(ctx context.Context, vm *kv1.VirtualMachine) error
	CreateDisks(ctx context.Context, vm *kv1.VirtualMachine) error
}

type vmServiceImpl struct {
	clusterResource apiserver.ClusterResource
}

func NewVmService(clusterResource apiserver.ClusterResource) VmService {
	return &vmServiceImpl{
		clusterResource: clusterResource,
	}
}

func (v vmServiceImpl) Create(ctx context.Context, vm *kv1.VirtualMachine) error {
	kvClient := v.clusterResource.Client().KubevirtClient().KubevirtV1()

	_, err := kvClient.VirtualMachines(vm.Namespace).Create(ctx, vm, metav1.CreateOptions{})
	if err != nil {

		return err
	}
	return nil
}

func (v vmServiceImpl) DeleteDisks(ctx context.Context, vm *kv1.VirtualMachine) error {
	pvcs, err := v.marshallPvcs(vm)
	if err != nil || pvcs == nil {
		return err
	}

	for _, pvcName := range pvcs {
		err = v.clusterResource.Client().K8sClient().CoreV1().PersistentVolumeClaims(vm.Namespace).
			Delete(ctx, pvcName.Name, metav1.DeleteOptions{})

		if err != nil {
			if k8serrors.IsNotFound(err) {
				continue
			}

			zap.L().Warn("failed to delete pvc",
				zap.String("namespace", vm.Namespace),
				zap.String("name", vm.Name), zap.Error(err))
		}
	}
	zap.L().Info("vm's pvcs are ensured to be deleted",
		zap.String("namespace", vm.Namespace),
		zap.String("name", vm.Name))
	return nil
}

func (v vmServiceImpl) CreateDisks(ctx context.Context, vm *kv1.VirtualMachine) error {
	pvcs, err := v.marshallPvcs(vm)
	if err != nil || pvcs == nil {
		return err
	}

	for _, pvc := range pvcs {
		_, err = v.clusterResource.Client().K8sClient().CoreV1().PersistentVolumeClaims(vm.Namespace).
			Create(ctx, &pvc, metav1.CreateOptions{})

		if err != nil {
			if k8serrors.IsAlreadyExists(err) {
				continue
			}

			zap.L().Warn("failed to create pvc",
				zap.String("namespace", vm.Namespace),
				zap.String("name", vm.Name), zap.Error(err))
			return err
		}
	}
	zap.L().Info("vm's pvcs are ensured to be created",
		zap.String("namespace", vm.Namespace),
		zap.String("name", vm.Name))
	return nil
}

func (v vmServiceImpl) marshallPvcs(vm *kv1.VirtualMachine) ([]corev1.PersistentVolumeClaim, error) {
	pvcTemplate, ok := vm.Annotations[constants.AnnotationPvcTemplates]
	if ok {
		var pvcs []corev1.PersistentVolumeClaim
		if err := json.Unmarshal([]byte(pvcTemplate), &pvcs); err != nil {
			zap.L().Warn("failed to Unmarshal pvcs from vm's annotation",
				zap.String("namespace", vm.Namespace),
				zap.String("name", vm.Name), zap.Error(err))
			return nil, err
		}
		return pvcs, nil
	}
	zap.L().Info("no pvcs found form vm's annotations, ignored", zap.String("namespace", vm.Namespace),
		zap.String("name", vm.Name))
	return nil, nil
}
