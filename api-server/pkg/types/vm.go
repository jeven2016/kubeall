package types

import (
	corev1 "k8s.io/api/core/v1"
	kv1 "kubevirt.io/api/core/v1"
)

type VmRequest struct {
	Vm  kv1.VirtualMachine
	Pvs []corev1.PersistentVolumeClaim
}
