package controller

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	kav1 "kubeall.io/api-server/pkg/generated/kubeall.io/v1"
	"kubeall.io/api-server/pkg/infra/apiserver"
	"kubeall.io/api-server/pkg/infra/constants"
	"kubeall.io/api-server/pkg/service"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	contrl "sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const imageControllerName = "image"

// ImageReconciler reconciles a Image object
type ImageReconciler struct {
	ReconcileHook[*kav1.Image]
	client.Client
	manager         Manager
	clusterResource apiserver.ClusterResource
	imageService    service.ImageService
	reconciler      Reconciler
}

func NewImageReconciler(clusterResource apiserver.ClusterResource, imageService service.ImageService) ReconcileHandler {
	return &ImageReconciler{clusterResource: clusterResource, imageService: imageService}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ImageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Client = mgr.GetClient()
	r.reconciler = DefaultReconciler[*kav1.Image]{hook: r, Client: mgr.GetClient()}
	return ctrl.NewControllerManagedBy(mgr).
		For(&kav1.Image{}).
		Named(imageControllerName).
		WithEventFilter(predicate.Funcs{
			UpdateFunc: func(e event.UpdateEvent) bool {
				// if status is changed due to the backing image's status changed.
				// In this case no need to reconcile the image itself, just ignored the event.
				oldImage := e.ObjectOld.(*kav1.Image)
				newImage := e.ObjectNew.(*kav1.Image)
				statusNotChanged := reflect.DeepEqual(oldImage.Status, newImage.Status)

				//if status not changed that means the image needs to update
				if statusNotChanged {
					zap.L().Info("image will be reconciled", zap.String("old", oldImage.Name))
				}
				return statusNotChanged
			},
			CreateFunc: func(e event.CreateEvent) bool {
				return true
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				return true
			},
			GenericFunc: func(e event.GenericEvent) bool {
				return false
			},
		}).
		WithOptions(contrl.Options{MaxConcurrentReconciles: constants.MaxConcurrentReconciles}).
		Complete(r)
}

func (r *ImageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return r.reconciler.Reconcile(ctx, req)
}

func (r *ImageReconciler) GetResource(ctx context.Context, req ctrl.Request) (*kav1.Image, error) {
	var image = &kav1.Image{}
	err := r.Get(ctx, req.NamespacedName, image)
	return image, err
}

func (r *ImageReconciler) GetClient() client.Client {
	return r.Client
}

func (r *ImageReconciler) Finalizer() string {
	return constants.DefaultFinalizer
}

func (r *ImageReconciler) OnRemove(ctx context.Context, req ctrl.Request, obj *kav1.Image) (ctrl.Result, error) {
	if err := r.imageService.DeleteImageResources(ctx, obj); err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed clear related resources for image"+obj.Name)
	}
	zap.L().Info("image is removed", zap.String("name", obj.Name))
	return ctrl.Result{}, nil
}

func (r *ImageReconciler) OnChange(ctx context.Context, obj *kav1.Image) error {
	if err := r.imageService.EnsureImageResources(ctx, obj); err != nil {
		return errors.Wrap(err, "failed to ensure related resource created for image "+obj.Name)
	}
	zap.L().Info("image is reconciled", zap.String("name", obj.Name))
	return nil
}

func (r *ImageReconciler) DeepCopy(obj *kav1.Image) *kav1.Image {
	return obj.DeepCopy()
}

func (r *ImageReconciler) OnAddFinalizer(obj *kav1.Image) {

}
