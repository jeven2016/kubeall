package server

import (
	"embed"
	"go.uber.org/fx"
	"kubeall.io/api-server/pkg/controller"
	"kubeall.io/api-server/pkg/infra"
	"kubeall.io/api-server/pkg/service"
	"kubeall.io/api-server/pkg/types"
)

func RegisterControllerManagerModules(params *types.StartupParams, localeFs embed.FS) {
	fx.New(
		infra.NewInfraModuleForCm(params, localeFs, types.CmScheme),
		service.Module,
		controller.CtrlModule,
		//fx.WithLogger(
		//	func() fxevent.Logger {
		//		return fxevent.NopLogger
		//	},
		//),
	).Run()
}
