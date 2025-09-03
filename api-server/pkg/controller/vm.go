package controller

import (
	"context"
	"kubeall.io/api-server/pkg/infra/apiserver"
	"kubeall.io/api-server/pkg/infra/constants"
	"kubeall.io/api-server/pkg/service"
	kv1 "kubevirt.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	contrl "sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"
)

const vmControllerName = "virtualMachine"

// VmReconciler reconciles a virtual machine object
type VmReconciler struct {
	ReconcileHook[*kv1.VirtualMachine]
	client.Client
	manager         Manager
	clusterResource apiserver.ClusterResource
	vmService       service.VmService
	reconciler      Reconciler
}

func NewVmReconciler(clusterResource apiserver.ClusterResource, vmService service.VmService) ReconcileHandler {
	return &VmReconciler{clusterResource: clusterResource, vmService: vmService}
}

func (v *VmReconciler) SetupWithManager(mgr ctrl.Manager) error {
	v.Client = mgr.GetClient()
	v.reconciler = DefaultReconciler[*kv1.VirtualMachine]{hook: v, Client: mgr.GetClient()}
	return ctrl.NewControllerManagedBy(mgr).
		For(&kv1.VirtualMachine{}).
		Named(vmControllerName).
		WithEventFilter(predicate.Funcs{
			//only vm's creation/deletion triggers the reconciler
			CreateFunc:  func(e event.CreateEvent) bool { return true },
			UpdateFunc:  func(e event.UpdateEvent) bool { return false },
			DeleteFunc:  func(e event.DeleteEvent) bool { return true },
			GenericFunc: func(e event.GenericEvent) bool { return false },
		}).
		WithOptions(contrl.Options{MaxConcurrentReconciles: constants.MaxConcurrentReconciles}).
		Complete(v)
}

func (v *VmReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return v.reconciler.Reconcile(ctx, req)
}

func (v *VmReconciler) GetResource(ctx context.Context, req ctrl.Request) (*kv1.VirtualMachine, error) {
	var vm = &kv1.VirtualMachine{}
	err := v.Get(ctx, req.NamespacedName, vm)
	return vm, err
}

func (v *VmReconciler) GetClient() client.Client {
	return v.Client
}

func (v *VmReconciler) Finalizer() string {
	return constants.DefaultFinalizer
}

func (v *VmReconciler) OnRemove(ctx context.Context, req ctrl.Request, obj *kv1.VirtualMachine) (ctrl.Result, error) {
	if err := v.vmService.DeleteDisks(ctx, obj); err != nil {
		return ctrl.Result{RequeueAfter: 2 * time.Second}, err
	}

	return ctrl.Result{}, nil
}

func (v *VmReconciler) OnChange(ctx context.Context, obj *kv1.VirtualMachine) error {
	//ensure pvc created
	if err := v.vmService.CreateDisks(ctx, obj); err != nil {
		return err
	}
	return nil
}

func (v *VmReconciler) DeepCopy(obj *kv1.VirtualMachine) *kv1.VirtualMachine {
	return obj.DeepCopy()
}

func (v *VmReconciler) OnAddFinalizer(obj *kv1.VirtualMachine) {
	// this method invoked during vm's creatio
}
