package service

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"kubeall.io/api-server/pkg/infra/constants"
)

type FieldComparator interface {
	Compare(first, second runtime.Object, order constants.SortOrder) bool
}

type CreationTimestampComparator struct{}

func (c CreationTimestampComparator) Compare(first, second runtime.Object, order constants.SortOrder) bool {
	firstTime := first.(metav1.Object).GetCreationTimestamp()
	secondTime := second.(metav1.Object).GetCreationTimestamp()

	if order == constants.Desc {
		return firstTime.After(secondTime.Time)
	}
	return firstTime.Before(&secondTime)
}

type NameComparator struct{}

func (n NameComparator) Compare(first, second runtime.Object, order constants.SortOrder) bool {
	firstName := first.(metav1.Object).GetName()
	secondName := second.(metav1.Object).GetName()
	if order == constants.Desc {
		return firstName > secondName
	}
	return firstName < secondName
}
