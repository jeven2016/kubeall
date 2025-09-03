package controller

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"kubeall.io/api-server/pkg/infra/constants"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type Reconciler interface {
	Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error)
}

type DefaultReconciler[T client.Object] struct {
	client.Client
	hook ReconcileHook[T]
}

func (d DefaultReconciler[T]) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	resource, err := d.hook.GetResource(ctx, req)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// delete resources
	if !resource.GetDeletionTimestamp().IsZero() {
		return d.deleteResource(ctx, req, resource)
	}

	// create or update
	if !controllerutil.ContainsFinalizer(resource, constants.DefaultFinalizer) {
		if result, err := d.addFinalizer(ctx, resource); err != nil {
			return result, err
		}
	} else {
		if err = d.hook.OnChange(ctx, resource); err != nil {
			msg := fmt.Sprintf("failed to ensure related resource created for %s %s ",
				resource.GetObjectKind(), resource.GetNamespace())
			return ctrl.Result{}, errors.Wrap(err, msg)
		}
	}
	return ctrl.Result{}, nil
}

func (d DefaultReconciler[T]) addFinalizer(ctx context.Context, resource T) (ctrl.Result, error) {
	newResource := d.hook.DeepCopy(resource)
	controllerutil.AddFinalizer(newResource, constants.DefaultFinalizer)
	if err := d.Patch(ctx, newResource, client.MergeFrom(resource)); err != nil {
		msg := fmt.Sprintf("failed to add finalizer for %s %s ",
			newResource.GetObjectKind(), newResource.GetNamespace())
		return ctrl.Result{}, errors.Wrap(err, msg)
	}
	return ctrl.Result{}, nil
}

func (d DefaultReconciler[T]) deleteResource(ctx context.Context, req ctrl.Request, resource T) (ctrl.Result, error) {
	if controllerutil.ContainsFinalizer(resource, constants.DefaultFinalizer) {
		// firstly, remove finalizer
		if result, err := d.hook.OnRemove(ctx, req, resource); err != nil {
			return result, err
		}

		// remove finalizer
		newResource := d.hook.DeepCopy(resource)
		controllerutil.RemoveFinalizer(newResource, constants.DefaultFinalizer)
		if err := d.Patch(ctx, newResource, client.MergeFrom(resource)); err != nil {
			msg := fmt.Sprintf("failed to remove finalizer for %s %s ",
				newResource.GetObjectKind(), newResource.GetNamespace())
			return ctrl.Result{}, errors.Wrap(err, msg)
		}
	}
	return ctrl.Result{}, nil
}
