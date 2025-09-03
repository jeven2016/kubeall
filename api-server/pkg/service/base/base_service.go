package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kubeall.io/api-server/pkg/infra/apiserver"
	"kubeall.io/api-server/pkg/infra/constants"
	"kubeall.io/api-server/pkg/types"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BaseService defines the interface for basic service operations.
type BaseService interface {
	List(ctx context.Context, gvk schema.GroupVersionKind, resType types.ResourceType, query types.Query, isPaginated bool) (*types.PageResult, error)
	Get(ctx context.Context, gvk schema.GroupVersionKind, resType types.ResourceType, name string) (client.Object, error)
	Delete(ctx context.Context, gvk schema.GroupVersionKind, resType types.ResourceType, name string) error
	Create(ctx context.Context, obj client.Object) error
	Update(ctx context.Context, obj client.Object) error
	CreateObject(gvk schema.GroupVersionKind) (client.Object, error)
}

// baseServiceImpl is an implementation of the BaseService interface.
type baseServiceImpl struct {
	clusterCache  cache.Cache
	runtimeClient client.Client
	fieldSorter   FieldSorter
	fieldFilter   FieldFilter
}

func NewBaseService(clusterRes apiserver.ClusterResource, sorter FieldSorter, fieldFilter FieldFilter) BaseService {
	baseSvcImpl := &baseServiceImpl{
		clusterCache:  clusterRes.ClusterCache(),
		runtimeClient: clusterRes.RuntimeClient(),
		fieldSorter:   sorter,
		fieldFilter:   fieldFilter,
	}
	return baseSvcImpl
}

func (b baseServiceImpl) List(ctx context.Context, gvk schema.GroupVersionKind,
	resType types.ResourceType, query types.Query, isPaginated bool) (*types.PageResult, error) {
	var err error
	var listOpts = &client.ListOptions{}

	listGvk := gvk
	listGvk.Kind = listGvk.Kind + "List"
	objList, err := CreateObjectList(b.runtimeClient.Scheme(), listGvk)
	if err != nil {
		return nil, err
	}

	if !resType.ClusterResource() {
		ns := resType.Namespace()
		if ns == constants.NamespaceAll {
			ns = v1.NamespaceAll
		}
		listOpts.Namespace = ns
	}
	err = b.clusterCache.List(ctx, objList, listOpts)
	if err != nil {
		zap.L().Warn("failed to list objects", zap.Any("gvk", gvk),
			zap.Any("resourceType", resType),
			zap.Error(err))
		return nil, types.FailWithErrorCode(ctx.(*gin.Context),
			constants.CodeInternalError, map[string]string{"error": err.Error()})
	}

	list, err := meta.ExtractList(objList)
	if err != nil {
		zap.L().Warn("failed to extract list objects", zap.Any("gvk", gvk),
			zap.Any("resourceType", resType),
			zap.Error(err))
		return nil, types.FailWithErrorCode(ctx.(*gin.Context),
			constants.CodeInternalError, map[string]string{"error": err.Error()})
	}

	//filter
	if list, err = b.fieldFilter.FilterBy(ctx, list, query.Filters); err != nil {
		zap.L().Warn("failed to filter by field(s)", zap.Any("filters", query.Filters),
			zap.Error(err))
		return nil, err
	}

	//sort
	if err = b.fieldSorter.SortByField(ctx, list, query.SortBy, string(query.SortOrder)); err != nil {
		zap.L().Warn("failed to sort list by this field", zap.String("sortBy", query.SortBy),
			zap.String("sortOrder", string(query.SortOrder)), zap.Error(err))
		return nil, err
	}

	//pagination
	pageResult := &types.PageResult{Items: make([]any, 0), TotalItems: len(list)}
	if isPaginated {
		Paginate(list, query, pageResult)
	} else {
		for _, item := range list {
			pageResult.Items = append(pageResult.Items, item)
		}
	}

	return pageResult, nil
}

func (b baseServiceImpl) Get(ctx context.Context, gvk schema.GroupVersionKind,
	resType types.ResourceType, name string) (client.Object, error) {
	objKey := b.createObjectKey(resType, name)
	obj, err := CreateObject(b.runtimeClient.Scheme(), gvk)
	if err != nil {
		return nil, err
	}
	if err = b.clusterCache.Get(ctx, objKey, obj); err != nil {
		if k8serrors.IsNotFound(err) {
			zap.L().Warn("no resource hit for deleting", zap.String("name", name),
				zap.Any("resourceType", resType), zap.Error(err))
			return nil, types.FailWithStatusCode(http.StatusNotFound)
		}
		zap.L().Warn("failed to get resource", zap.String("name", name),
			zap.Any("resourceType", resType), zap.Error(err))
		return nil, err
	}
	return obj, nil
}

func (b baseServiceImpl) Delete(ctx context.Context, gvk schema.GroupVersionKind,
	resType types.ResourceType, name string) error {
	obj, err := CreateObject(b.runtimeClient.Scheme(), gvk)
	if err != nil {
		return err
	}
	obj.SetNamespace(resType.Namespace())
	obj.SetName(name)
	if err = b.runtimeClient.Delete(ctx, obj); err != nil {
		if k8serrors.IsNotFound(err) {
			zap.L().Warn("no resource found for deleting", zap.String("name", name),
				zap.Any("resourceType", resType), zap.Error(err))
			return types.FailWithStatusCode(http.StatusNotFound)
		}
	}
	zap.L().Warn("failed to delete resource", zap.String("name", name),
		zap.Any("resourceType", resType), zap.Error(err))
	return err
}

func (b baseServiceImpl) Create(ctx context.Context, obj client.Object) error {
	if err := b.runtimeClient.Create(ctx, obj); err != nil {
		zap.L().Warn("failed to create resource", zap.Any("error", err))
		return err
	}
	return nil
}

func (b baseServiceImpl) Update(ctx context.Context, obj client.Object) error {
	if err := b.runtimeClient.Update(ctx, obj); err != nil {
		zap.L().Warn("failed to update resource", zap.Any("error", err))
		return err
	}
	return nil
}

func (b baseServiceImpl) CreateObject(gvk schema.GroupVersionKind) (client.Object, error) {
	return CreateObject(b.runtimeClient.Scheme(), gvk)
}

// createObjectKey constructs an ObjectKey for the given resource type and name
func (b baseServiceImpl) createObjectKey(resType types.ResourceType, name string) client.ObjectKey {
	objKey := client.ObjectKey{Name: name}
	if !resType.ClusterResource() {
		objKey.Namespace = resType.Namespace()
	}
	return objKey
}
