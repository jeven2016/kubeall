package image

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	basehandler "kubeall.io/api-server/pkg/handler/base"
	"kubeall.io/api-server/pkg/handler/route"
	"kubeall.io/api-server/pkg/handler/validators"
	"kubeall.io/api-server/pkg/infra/constants"
	"kubeall.io/api-server/pkg/infra/validator_resource"
	"kubeall.io/api-server/pkg/service"
	baseservice "kubeall.io/api-server/pkg/service/base"
	"kubeall.io/api-server/pkg/types"
	"net/http"
)

type ImageHandler interface {
	route.Route
	Upload(ctx *gin.Context)
	ListImages(ctx *gin.Context)
}

type imageHandlerImpl struct {
	baseService  baseservice.BaseService
	imageService service.ImageService
	gvkResource  *constants.GvkResource
	translator   validator_resource.ValidatorTranslator
}

func (i imageHandlerImpl) ListImages(ctx *gin.Context) {
	imageType := ctx.Query("type")
	if imageType == "" {
		ctx.Params = append(ctx.Params, gin.Param{Key: constants.ResourceParam, Value: constants.ImageResourceParam})
		basehandler.HandleList(ctx, i.gvkResource, i.translator, i.baseService)
		return
	}
	err, _ := validators.ValidateNow(ctx, imageType, constants.ValidateImageType, i.translator)
	if err != nil {
		zap.L().Warn("invalid imageType", zap.String("imageType", imageType), zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, types.FailWithErrorCode(ctx,
			constants.CodeInvalidParam, map[string]string{"name": "type"}))
		return
	}
	images, err := i.imageService.ListImagesByType(ctx, ctx.Param("namespace"), imageType)
	if err != nil {
		zap.L().Warn("failed to get images", zap.String("type", imageType),
			zap.Error(err))
		basehandler.AbortRequest(ctx, err, 0)
		return
	}
	ctx.JSON(http.StatusOK, images)
}

func NewImageHandler(baseService baseservice.BaseService, imageService service.ImageService, gvkResource *constants.GvkResource,
	translator validator_resource.ValidatorTranslator) ImageHandler {
	return &imageHandlerImpl{
		baseService, imageService, gvkResource, translator,
	}
}

// Upload a vm image. meanwhile engine.MaxMultipartMemory is set to 32MB
func (i imageHandlerImpl) Upload(ctx *gin.Context) {
	imageName := ctx.Param("imageName")
	err, errMsg := validators.ValidateNow(ctx, imageName, "required", i.translator)
	if err != nil {
		zap.L().Warn("invalid image name", zap.String("name", imageName), zap.Error(err),
			zap.String("errMsg", errMsg))
		basehandler.AbortRequestWithMessage(ctx, "name"+errMsg, http.StatusBadRequest)
		return
	}
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		zap.L().Warn("failed get file from from", zap.String("name", imageName), zap.Error(err))
		basehandler.AbortRequest(ctx, types.Fail(err), http.StatusBadRequest)
		return
	}
	fileSize := header.Size
	zap.L().Info("the image to upload", zap.String("name", imageName),
		zap.Int64("size", fileSize))

	if err = i.imageService.Upload(ctx, imageName, file, fileSize, ctx.Request); err != nil {
		zap.L().Warn("failed to upload image", zap.Error(err))
		basehandler.AbortRequest(ctx, err, http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func (b imageHandlerImpl) RegisterRoutes(_ *gin.RouterGroup, namespaceGroup *gin.RouterGroup, _ *gin.RouterGroup) {
	namespaceGroup.POST(constants.ResourceImageUploadUri, b.Upload)
	namespaceGroup.GET(constants.ResourceImageUri, b.ListImages)
}
