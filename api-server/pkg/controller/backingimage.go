package controller

import (
	"context"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"
	"kubeall.io/api-server/pkg/controller/predicates"
	kav1 "kubeall.io/api-server/pkg/generated/kubeall.io/v1"
	lhv1beta2 "kubeall.io/api-server/pkg/generated/longhorn/apis/longhorn/v1beta2"
	"kubeall.io/api-server/pkg/infra/apiserver"
	"kubeall.io/api-server/pkg/service"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	contrl "sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type BackingImageReconciler struct {
	client.Client
	manager         Manager
	clusterResource apiserver.ClusterResource
	imageService    service.ImageService
}

func NewBackingImageReconciler(clusterResource apiserver.ClusterResource, imageService service.ImageService) ReconcileHandler {
	return &BackingImageReconciler{clusterResource: clusterResource, imageService: imageService}
}

// SetupWithManager sets up the controller with the Manager.
func (r *BackingImageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Client = mgr.GetClient()

	// watch the change event of backing image's status
	return ctrl.NewControllerManagedBy(mgr).
		Named("backingImageController").
		Watches(
			&lhv1beta2.BackingImage{},
			handler.EnqueueRequestsFromMapFunc(
				func(ctx context.Context, h client.Object) []reconcile.Request {
					biImage := h.(*lhv1beta2.BackingImage)
					return []reconcile.Request{
						{
							NamespacedName: types.NamespacedName{
								Name:      biImage.Name,
								Namespace: biImage.Namespace,
							},
						},
					}
				}),
			builder.WithPredicates(predicates.StatusChangePredicate{})).
		WithOptions(contrl.Options{MaxConcurrentReconciles: 2}).
		Complete(r)
}

func (r *BackingImageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	zap.L().Info("backing image reconciler triggered", zap.String("backingImage", req.Name),
		zap.String("namespace", req.Namespace))
	biImage, err := r.GetResource(ctx, req)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if len(biImage.Status.DiskFileStatusMap) == 0 {
		return ctrl.Result{}, nil
	}

	imageStatus := &kav1.ImageStatus{}
	for _, v := range biImage.Status.DiskFileStatusMap {
		imageStatus.Size = biImage.Status.Size

		if v.State == lhv1beta2.BackingImageStateReady {
			imageStatus.VirtualSize = biImage.Status.VirtualSize
		}
		//the progress's value 100 doesn't mean the image is finished to upload meanwhile the state is not ready
		p := v.Progress
		if v.State != "ready" && v.Progress == 100 {
			p = 99
		}
		imageStatus.State = string(v.State)
		imageStatus.Progress = p
		imageStatus.LastStateTransitionTime = v.LastStateTransitionTime
		imageStatus.Message = v.Message
		break
	}
	zap.L().Info("staus map", zap.Any("statusMap", biImage.Status.DiskFileStatusMap))
	err = r.imageService.UpdateStatus(ctx, imageStatus, biImage)
	return ctrl.Result{}, err
}

func (r *BackingImageReconciler) GetResource(ctx context.Context, req ctrl.Request) (*lhv1beta2.BackingImage, error) {
	biImage := &lhv1beta2.BackingImage{}
	err := r.Get(ctx, req.NamespacedName, biImage)
	return biImage, err
}
