package basehandler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kubeall.io/api-server/pkg/handler/validators"
	"kubeall.io/api-server/pkg/infra/constants"
	"kubeall.io/api-server/pkg/infra/validator_resource"
	service "kubeall.io/api-server/pkg/service/base"
	"kubeall.io/api-server/pkg/types"
	"net/http"
	"strconv"
)

func GetHttpCode(err error, code int) int {
	if code > 0 {
		return code
	}
	isResult, result := types.IsResult(err)
	if isResult {
		if result != nil && result.StatusCode > 0 {
			return result.StatusCode
		}
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func AbortRequest(ctx *gin.Context, err error, customStatusCode int) {
	code := GetHttpCode(err, customStatusCode)
	if code == http.StatusNotFound {
		ctx.AbortWithStatus(code)
	} else {
		ctx.AbortWithStatusJSON(code, types.Fail(err))
	}
}

func AbortRequestWithMessage(ctx *gin.Context, message string, customStatusCode int) {
	if customStatusCode == http.StatusNotFound {
		ctx.AbortWithStatus(customStatusCode)
	} else {
		ctx.AbortWithStatusJSON(customStatusCode, &types.Result{
			Message:    message,
			ErrorCode:  constants.CodeInvalidParam,
			StatusCode: http.StatusBadRequest,
		})
	}
}

func CheckResourceType(ctx *gin.Context, gvkResource *constants.GvkResource) (*schema.GroupVersionKind,
	types.ResourceType, error) {
	if gvkRes, resourceType, err := validators.ValidateResourceType(ctx, gvkResource); err != nil {
		AbortRequest(ctx, types.Fail(err), http.StatusBadRequest)
		return nil, nil, err
	} else {
		return gvkRes, resourceType, err
	}
}

func CheckName(ctx *gin.Context, translator validator_resource.ValidatorTranslator) error {
	if err, msg := validators.ValidateNow(ctx, ctx.Param("name"),
		"required", translator); err != nil {
		AbortRequest(ctx, types.Fail(errors.New(msg)), http.StatusBadRequest)
		return err
	}
	return nil
}

func ConvertIntValue(ctx *gin.Context, field, defaultValue string) (int, error) {
	page := ctx.DefaultQuery(field, defaultValue)
	intValue, err := strconv.Atoi(page)
	if err != nil {
		zap.L().Warn("invalid"+field, zap.String(field, page))
		AbortRequest(ctx, types.FailWithErrorCode(ctx, constants.CodeInvalidParam,
			map[string]string{"name": field}), http.StatusBadRequest)
		return 0, err
	}
	return intValue, nil
}

func HandleList(ctx *gin.Context, gvkResource *constants.GvkResource,
	translator validator_resource.ValidatorTranslator, baseService service.BaseService) {
	var gvkRes *schema.GroupVersionKind
	var resourceType types.ResourceType
	var err error
	if gvkRes, resourceType, err = CheckResourceType(ctx, gvkResource); err != nil {
		zap.L().Warn("failed to list the resources", zap.Error(err))
		return
	}

	pageInt, err := ConvertIntValue(ctx, constants.PageQueryField, constants.DefaultPage)
	if err != nil {
		return
	}
	pageSizeInt, err := ConvertIntValue(ctx, constants.PageSizeQueryField, constants.DefaultPageSize)
	if err != nil {
		return
	}

	sortOrder := ctx.Query(constants.SortOrderField)
	sortBy := ctx.Query(constants.SortByQueryField)

	filterMap, err := validators.ValidFilter(ctx)
	if err != nil {
		zap.L().Warn("failed to unmarshal filter", zap.String("filter",
			ctx.Query(constants.FilterField)), zap.Error(err))
		AbortRequest(ctx, types.Fail(err), http.StatusBadRequest)
		return
	}

	query := types.Query{
		Pagination: types.Pagination{
			Page:     uint(pageInt),
			PageSize: uint(pageSizeInt),
		},
		SortOrder: constants.SortOrder(sortOrder),
		SortBy:    sortBy,
		Filters:   filterMap,
	}
	result := validators.ValidateListParams(ctx, translator, query)
	if result != nil {
		zap.L().Warn("failed to list resources", zap.Any("result", result))
		AbortRequest(ctx, result, http.StatusBadRequest)
		return
	}

	objList, err := baseService.List(ctx, *gvkRes, resourceType.(types.ResourceType), query, true)
	if err != nil {
		zap.L().Warn("failed to list resources", zap.String("resource", gvkRes.Kind), zap.Error(err))
		AbortRequest(ctx, types.Fail(err), 0)
		return
	}
	ctx.JSON(http.StatusOK, objList)
}
