package service

import (
	"context"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"kubeall.io/api-server/pkg/infra/constants"
	"kubeall.io/api-server/pkg/types"
	"slices"
	"sort"
)

type FieldSorter interface {
	SortByField(ctx context.Context, list []runtime.Object, vmField, order string) error
}

type fieldSorterImpl struct {
	fieldMap map[constants.Field]FieldComparator
}

func NewFieldSorter() FieldSorter {
	fieldMap := map[constants.Field]FieldComparator{
		constants.FieldCreationTimeStamp: &CreationTimestampComparator{},
		constants.FieldName:              &NameComparator{},
	}
	return &fieldSorterImpl{fieldMap}
}

func (f fieldSorterImpl) SortByField(ctx context.Context, list []runtime.Object, field string, order string) error {
	realField := constants.Field(field)
	realSortOrder := constants.SortOrder(order)

	if field == "" {
		realField = constants.FieldCreationTimeStamp
	}
	if order == "" {
		realSortOrder = constants.Desc
	}
	if !slices.Contains(constants.SupportedSortFields, realField) {
		zap.L().Warn("unsupported sort field", zap.String("field", field))
		return types.FailWithErrorCode(ctx, constants.CodeInvalidParam, map[string]string{"name": "field"})
	}
	if !slices.Contains(constants.SupportedSortOrders, realSortOrder) {
		zap.L().Warn("unsupported sort order", zap.String("order", order))
		return types.FailWithErrorCode(ctx, constants.CodeInvalidParam, map[string]string{"name": "order"})
	}

	if comparator, ok := f.fieldMap[realField]; ok {
		sort.Slice(list, func(i, j int) bool {
			return comparator.Compare(list[i], list[j], realSortOrder)
		})
	} else {
		zap.L().Warn("unknown field to sort the list", zap.String("field", field),
			zap.String("order", order))
	}
	return nil
}
