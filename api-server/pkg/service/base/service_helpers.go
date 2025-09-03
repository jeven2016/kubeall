package service

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kubeall.io/api-server/pkg/types"
	"math"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateObjectList(scheme *runtime.Scheme, listGvk schema.GroupVersionKind) (client.ObjectList, error) {
	var objList client.ObjectList
	if scheme.Recognizes(listGvk) {
		gvkObject, err := scheme.New(listGvk)
		if err != nil {
			return nil, err
		}
		objList = gvkObject.(client.ObjectList)
	} else {
		ul := &unstructured.UnstructuredList{}
		ul.SetGroupVersionKind(listGvk)
		objList = ul
	}
	return objList, nil
}

func CreateObject(scheme *runtime.Scheme, objGvk schema.GroupVersionKind) (client.Object, error) {
	var obj client.Object
	if scheme.Recognizes(objGvk) {
		gvkObject, err := scheme.New(objGvk)
		if err != nil {
			return nil, err
		}
		obj = gvkObject.(client.Object)
	} else {
		ul := &unstructured.Unstructured{}
		ul.SetGroupVersionKind(objGvk)
		obj = ul
	}
	return obj, nil
}

func Paginate[T runtime.Object](list []T, query types.Query, pageResult *types.PageResult) {
	page := int(query.Pagination.Page)
	pageSize := int(query.Pagination.PageSize)

	pageResult.Page = page
	pageResult.PageSize = pageSize

	if len(list) != 0 {
		pageResult.TotalPages = int(math.Ceil(float64(len(list)) / float64(pageSize)))
	}

	startIndex := (page - 1) * pageSize
	if startIndex >= len(list) {
		return
	}

	endIndex := startIndex + pageSize
	if endIndex > len(list) {
		endIndex = len(list)
	}

	pageList := list[startIndex:endIndex]
	for _, item := range pageList {
		pageResult.Items = append(pageResult.Items, item)
	}
}
