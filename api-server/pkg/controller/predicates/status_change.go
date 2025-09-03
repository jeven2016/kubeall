package predicates

import (
	"go.uber.org/zap"
	lhv1beta2 "kubeall.io/api-server/pkg/generated/longhorn/apis/longhorn/v1beta2"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type StatusChangePredicate struct {
	predicate.Funcs
}

func (StatusChangePredicate) Update(e event.UpdateEvent) bool {
	oldBi, oldOk := e.ObjectOld.(*lhv1beta2.BackingImage)

	newBi, newOk := e.ObjectNew.(*lhv1beta2.BackingImage)
	if !oldOk || !newOk {
		return false
	}

	if !newBi.GetDeletionTimestamp().IsZero() {
		zap.L().Info("no need to update while it's being deleted", zap.String("backingImage", newBi.GetName()))
		return false
	}

	if len(newBi.Status.DiskFileStatusMap) == 0 {
		zap.L().Info("no need to update while no progress reported", zap.String("backingImage", newBi.GetName()))
		return false
	}

	// check if status changed
	if !reflect.DeepEqual(oldBi.Status, newBi.Status) {
		zap.L().Info("the image's progress updated", zap.String("backingImage", newBi.GetName()))
		return true
	}
	return false
}

func (StatusChangePredicate) Create(_ event.CreateEvent) bool {
	return false
}

func (StatusChangePredicate) Delete(_ event.DeleteEvent) bool {
	return false
}

func (StatusChangePredicate) Generic(_ event.GenericEvent) bool {
	return false
}
