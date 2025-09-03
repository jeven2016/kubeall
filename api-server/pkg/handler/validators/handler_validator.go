package validators

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kubeall.io/api-server/pkg/infra/constants"
	"kubeall.io/api-server/pkg/infra/validator_resource"
	"kubeall.io/api-server/pkg/types"
	"slices"
)

// GetTranslator retrieves the Chinese translator from the provided ValidatorTranslator.
func GetTranslator(ctx *gin.Context, translator validator_resource.ValidatorTranslator) ut.Translator {
	//TODO: choose correct language
	return translator.Zh()
}

func GetValidate() *validator.Validate {
	validate := binding.Validator.Engine().(*validator.Validate)
	return validate
}

func ValidateResourceType(ctx *gin.Context, gvkResource *constants.GvkResource) (*schema.GroupVersionKind, types.ResourceType, error) {
	res := ctx.Param("resource")
	if res == "" {
		return nil, nil, types.FailWithErrorCode(ctx, constants.CodeRequired, map[string]string{"name": "resource"})
	}

	gvkRes, err := gvkResource.Get(res)
	if err != nil {
		return nil, nil, types.Fail(err)
	}

	resType, exists := ctx.Get(constants.ResourceType)
	if !exists {
		return nil, nil, types.FailWithErrorCode(ctx, constants.CodeRequired, map[string]string{"name": "resourceType"})
	}
	resourceType := resType.(types.ResourceType)
	return gvkRes, resourceType, nil
}

func ValidateListParams(ctx *gin.Context, translator validator_resource.ValidatorTranslator,
	query types.Query) error {
	result := ValidateParam(ctx, constants.PageQueryField, query.Pagination.Page, translator)
	if result != nil {
		return result
	}

	result = ValidateParam(ctx, constants.PageSizeQueryField, query.Pagination.PageSize, translator)
	if result != nil {
		return result
	}

	return nil
}

func ValidateParam(ctx *gin.Context, paramName string, paramValue uint,
	translator validator_resource.ValidatorTranslator) error {

	switch paramName {
	case constants.PageQueryField:
		if err, msg := ValidateNow(ctx, paramValue, "gte=1,lte=1000", translator); err != nil {
			return types.FailWithErrorCode(ctx, constants.CodeValidationFailed, map[string]string{"name": "page: " + msg})
		}
	case constants.PageSizeQueryField:
		if !slices.Contains(constants.AvailablePageSizes, int(paramValue)) {
			return types.FailWithErrorCode(ctx, constants.CodeInvalidParam, map[string]string{"name": "pageSize"})
		}
	}

	return nil
}

func ValidateNow(ctx *gin.Context, value any, tag string,
	translator validator_resource.ValidatorTranslator) (error, string) {
	err := GetValidate().Var(value, tag)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			msg := e.Translate(GetTranslator(ctx, translator))
			return e, msg
		}
	}
	return nil, ""
}

func ValidFilter(ctx *gin.Context) (map[string]string, error) {
	var filterMap map[string]string
	filter := ctx.Query(constants.FilterField)
	if filter != "" {
		if err := json.Unmarshal([]byte(filter), &filterMap); err != nil {
			return nil, err
		}
	}
	return filterMap, nil
}
