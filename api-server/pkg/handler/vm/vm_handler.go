package vm

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	basehandler "kubeall.io/api-server/pkg/handler/base"
	"kubeall.io/api-server/pkg/handler/route"
	"kubeall.io/api-server/pkg/infra/constants"
	"kubeall.io/api-server/pkg/infra/validator_resource"
	"kubeall.io/api-server/pkg/service"
	"kubeall.io/api-server/pkg/types"
	"net/http"
)

type VmHandler interface {
	route.Route
	Create(*gin.Context)
}

type vmHandlerImpl struct {
	vmService  service.VmService
	translator validator_resource.ValidatorTranslator
}

func NewVmHandler(vmService service.VmService, translator validator_resource.ValidatorTranslator) VmHandler {
	return &vmHandlerImpl{
		vmService, translator,
	}
}

func (v vmHandlerImpl) Create(ctx *gin.Context) {
	var vmRequest types.VmRequest
	if err := ctx.ShouldBindJSON(&vmRequest); err != nil {
		zap.L().Warn("failed to unmarshall vm request", zap.Any("error", err))
		basehandler.AbortRequest(ctx, types.Fail(err), http.StatusBadRequest)
		return
	}
	pvcsJson, err := json.Marshal(vmRequest.Pvs)
	if err != nil {
		zap.L().Warn("failed to marshall pvcs", zap.Any("error", err))
		basehandler.AbortRequest(ctx, types.Fail(err), http.StatusBadRequest)
		return
	}
	vm := vmRequest.Vm
	vm.Annotations[constants.AnnotationPvcTemplates] = string(pvcsJson)
	if err = v.vmService.Create(ctx, &vm); err != nil {
		zap.L().Warn("failed to create vm ", zap.String("vmName", vm.Name),
			zap.String("namespace", vm.Namespace), zap.Error(err))
		basehandler.AbortRequest(ctx, types.Fail(err), 0)
		return
	}
	zap.S().Info(fmt.Sprintf("vm(%s/%s) is created", vm.Namespace, vm.Name))
	ctx.Status(http.StatusCreated)
}

func (v vmHandlerImpl) RegisterRoutes(rootGroup *gin.RouterGroup, namespaceGroup *gin.RouterGroup, clusterGroup *gin.RouterGroup) {
	namespaceGroup.POST(constants.ResourceVmUri, v.Create)
}
