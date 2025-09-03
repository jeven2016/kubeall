package server

import (
	"context"
	"embed"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"kubeall.io/api-server/pkg/handler"
	"kubeall.io/api-server/pkg/infra"
	"kubeall.io/api-server/pkg/service"
	"kubeall.io/api-server/pkg/types"
)

func NewServerModule() fx.Option {
	return fx.Options(
		// 暂时没有传递的参数
		fx.Provide(
			NewServer,
		),
	)
}

func RegisterModules(params *types.StartupParams, localeFs embed.FS) {
	fx.New(
		infra.NewInfraModule(params, localeFs, types.ApiServerScheme),
		NewServerModule(),
		service.Module,
		handler.Module,
		fx.Invoke(func(lifecycle fx.Lifecycle, logger *zap.Logger, server Server) error {
			return server.Start(context.Background(), localeFs)
		}),
		fx.WithLogger(
			func() fxevent.Logger {
				return fxevent.NopLogger
			},
		),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()
}
