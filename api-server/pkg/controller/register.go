package controller

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func AsReconciler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(ReconcileHandler)),
		fx.ResultTags(`group:"reconcilers"`),
	)
}

var CtrlModule = fx.Module("controllers",
	fx.Provide(
		AsReconciler(NewImageReconciler),
		AsReconciler(NewBackingImageReconciler),
		AsReconciler(NewVmReconciler),

		//receive group resources
		fx.Annotate(
			NewManager,
			fx.ParamTags(`group:"reconcilers"`),
		),
	),
	// Invoke the provided function after the module is fully initialized to register routes
	fx.Invoke(func(mgr Manager, logger *zap.Logger) {
		mgr.Start()
	}),
)
