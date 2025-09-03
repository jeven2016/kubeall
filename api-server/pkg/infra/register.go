package infra

import (
	"embed"
	"go.uber.org/fx"
	"kubeall.io/api-server/pkg/infra/apiserver"
	"kubeall.io/api-server/pkg/infra/clients"
	"kubeall.io/api-server/pkg/infra/config"
	"kubeall.io/api-server/pkg/infra/constants"
	"kubeall.io/api-server/pkg/infra/validator_resource"
	"kubeall.io/api-server/pkg/types"

	"kubeall.io/api-server/pkg/infra/logger"
)

func NewInfraModule(params *types.StartupParams, localeFs embed.FS, schemeType types.SchemeType) fx.Option {
	return fx.Options(
		// 提供构造函数所需的参数
		fx.Supply(params, localeFs, schemeType),
		fx.Provide(
			config.NewServerConfig,
			logger.NewLogger,
			clients.NewClients,
			apiserver.NewClusterResource,
			apiserver.NewRestServer,
			constants.NewGvkResource,
			validator_resource.NewValidatorTranslator,
		),
	)
}

func NewInfraModuleForCm(params *types.StartupParams, localeFs embed.FS, schemeType types.SchemeType) fx.Option {
	return fx.Options(
		// 提供构造函数所需的参数
		fx.Supply(params, localeFs, schemeType),
		fx.Provide(
			config.NewServerConfig,
			logger.NewLogger,
			clients.NewClients,
			apiserver.NewClusterResourceForCm,
			constants.NewGvkResource,
			validator_resource.NewValidatorTranslator,
		),
	)
}
