package basehandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kubeall.io/api-server/pkg/handler/route"
	"kubeall.io/api-server/pkg/infra/constants"
	"kubeall.io/api-server/pkg/infra/validator_resource"
	service "kubeall.io/api-server/pkg/service/base"
	"kubeall.io/api-server/pkg/types"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type BaseHandler interface {
	route.Route
	List(*gin.Context)
	Get(*gin.Context)
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

type baseHandlerImpl struct {
	gvkResource *constants.GvkResource
	baseService service.BaseService
	translator  validator_resource.ValidatorTranslator
}

func NewBaseHandler(baseService service.BaseService,
	gvkResource *constants.GvkResource, translator validator_resource.ValidatorTranslator) BaseHandler {
	return &baseHandlerImpl{
		gvkResource: gvkResource,
		baseService: baseService,
		translator:  translator,
	}
}

func (b baseHandlerImpl) RegisterRoutes(_ *gin.RouterGroup, namespaceGroup *gin.RouterGroup, clusterGroup *gin.RouterGroup) {
	namespaceGroup.GET(constants.ResourceUri, b.List)
	namespaceGroup.POST(constants.ResourceUri, b.Create)
	namespaceGroup.GET(constants.ResourceNameUri, b.Get)
	namespaceGroup.PUT(constants.ResourceNameUri, b.Update)
	namespaceGroup.DELETE(constants.ResourceNameUri, b.Delete)

	clusterGroup.GET(constants.ResourceUri, b.List)
	clusterGroup.POST(constants.ResourceUri, b.Create)
	clusterGroup.GET(constants.ResourceNameUri, b.Get)
	clusterGroup.PUT(constants.ResourceNameUri, b.Update)
	clusterGroup.DELETE(constants.ResourceNameUri, b.Delete)
}

func (b baseHandlerImpl) List(ctx *gin.Context) {
	HandleList(ctx, b.gvkResource, b.translator, b.baseService)
}

func (b baseHandlerImpl) Get(ctx *gin.Context) {
	var gvkRes *schema.GroupVersionKind
	var resourceType types.ResourceType
	var err error
	var name = ctx.Param("name")

	if gvkRes, resourceType, err = CheckResourceType(ctx, b.gvkResource); err != nil {
		zap.L().Warn("failed to get resources", zap.Error(err))
		return
	}

	if err = CheckName(ctx, b.translator); err != nil {
		zap.L().Warn("failed to delete resources", zap.Error(err), zap.String("name", name))
		return
	}

	obj, err := b.baseService.Get(ctx, *gvkRes, resourceType, name)
	if err != nil {
		zap.L().Warn("failed to get resource", zap.String("resource", gvkRes.Kind), zap.String("name", name),
			zap.Error(err))
		AbortRequest(ctx, types.Fail(err), 0)
		return
	}
	if obj == nil {
		zap.L().Warn("failed to get target resource", zap.String("resource", gvkRes.Kind),
			zap.String("name", name))
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, obj)
}

func (b baseHandlerImpl) Create(ctx *gin.Context) {
	var obj client.Object
	if obj = b.parseBody(ctx); obj == nil {
		return
	}

	if err := b.baseService.Create(ctx, obj); err != nil {
		zap.L().Warn("failed to create resource", zap.Any("error", err))
		AbortRequest(ctx, types.Fail(err), 0)
		return
	}
	ctx.Status(http.StatusOK)
}

func (b baseHandlerImpl) Update(ctx *gin.Context) {
	var obj client.Object

	if obj = b.parseBody(ctx); obj == nil {
		return
	}

	if err := b.baseService.Update(ctx, obj); err != nil {
		zap.L().Warn("failed to update resource", zap.Any("error", err))
		AbortRequest(ctx, types.Fail(err), 0)
		return
	}
	ctx.Status(http.StatusOK)
}

func (b baseHandlerImpl) Delete(ctx *gin.Context) {
	var gvkRes *schema.GroupVersionKind
	var resourceType types.ResourceType
	var err error
	var name = ctx.Param("name")

	if gvkRes, resourceType, err = CheckResourceType(ctx, b.gvkResource); err != nil {
		zap.L().Warn("failed to delete resources", zap.Error(err), zap.String("name", name))
		return
	}

	if err = CheckName(ctx, b.translator); err != nil {
		zap.L().Warn("failed to delete resources", zap.Error(err), zap.String("name", name))
		return
	}

	err = b.baseService.Delete(ctx, *gvkRes, resourceType, name)
	if err != nil {
		zap.L().Warn("failed to get resource", zap.String("resource", gvkRes.Kind), zap.String("name", name),
			zap.Error(err))
		AbortRequest(ctx, err, 0)
		return
	}

	ctx.Status(http.StatusOK)
}

func (b baseHandlerImpl) parseBody(ctx *gin.Context) client.Object {
	var obj client.Object
	var gvkRes *schema.GroupVersionKind
	var err error

	format := ctx.DefaultQuery("format", "json")
	if format != constants.JsonFormat && format != constants.YamlFormat {
		zap.L().Warn("invalid format to unmarshall resource", zap.Any("format", format))
		AbortRequest(ctx, types.FailWithErrorCode(ctx, constants.CodeInvalidParam, map[string]string{"name": "format"}), http.StatusBadRequest)
		return nil
	}

	if gvkRes, _, err = CheckResourceType(ctx, b.gvkResource); err != nil {
		zap.L().Warn("failed to check resource type", zap.Error(err))
		return nil
	}

	if obj, err = b.baseService.CreateObject(*gvkRes); err != nil {
		zap.L().Warn("failed to check resource type", zap.Error(err))
		AbortRequest(ctx, types.FailWithErrorCode(ctx, constants.CodeInvalidData, nil), http.StatusBadRequest)
		return nil
	}

	if format == constants.JsonFormat {
		err = ctx.ShouldBindJSON(obj)
	}

	if err != nil {
		zap.L().Warn("failed to unmarshall json", zap.String("format", format), zap.Any("error", err))
		AbortRequest(ctx, types.Fail(err), http.StatusBadRequest)
	}
	return obj
}
