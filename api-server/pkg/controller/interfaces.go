package controller

import (
	"context"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ReconcileHandler interface {
	reconcile.Reconciler
	SetupWithManager(mgr ctrl.Manager) error
}

type ReconcileHook[T client.Object] interface {
	GetResource(ctx context.Context, req ctrl.Request) (T, error)
	GetClient() client.Client
	Finalizer() string
	OnAddFinalizer(resource T)
	OnRemove(ctx context.Context, req ctrl.Request, obj T) (ctrl.Result, error)
	OnChange(ctx context.Context, obj T) error
	DeepCopy(T) T
}
