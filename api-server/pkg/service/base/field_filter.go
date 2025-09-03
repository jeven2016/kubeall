package service

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"kubeall.io/api-server/pkg/infra/constants"
	"kubeall.io/api-server/pkg/types"
	"slices"
	"strings"
)

type FieldFilter interface {
	FilterBy(context.Context, []runtime.Object, map[string]string) ([]runtime.Object, error)
}

type fieldFilterImpl struct {
}

func NewFieldFilter() FieldFilter {
	return &fieldFilterImpl{}
}

func (f fieldFilterImpl) FilterBy(ctx context.Context, list []runtime.Object,
	filters map[string]string) ([]runtime.Object, error) {
	if len(filters) == 0 {
		return list, nil
	}
	finalList := list
	for k, v := range filters {
		if !slices.Contains(constants.SupportedFilterFields, constants.Field(k)) {
			return nil, types.FailWithErrorCode(ctx, constants.CodeInvalidParam, map[string]string{"name": k})
		}
		finalList = f.filterByField(finalList, k, v)
	}
	return finalList, nil
}

func (f fieldFilterImpl) filterByField(list []runtime.Object, field string, value string) []runtime.Object {
	finalList := make([]runtime.Object, 0)
	switch constants.Field(field) {
	case constants.FieldName:
		for _, v := range list {
			if strings.Contains(v.(metav1.Object).GetName(), value) {
				finalList = append(finalList, v)
			}
		}
		return finalList
	default:
		return list
	}
}
