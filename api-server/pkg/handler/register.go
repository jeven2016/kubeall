package handler

import (
	"go.uber.org/fx"
	basehandler "kubeall.io/api-server/pkg/handler/base"
	"kubeall.io/api-server/pkg/handler/image"
	"kubeall.io/api-server/pkg/handler/route"
	"kubeall.io/api-server/pkg/handler/vm"
)

var Module = fx.Module("md_handler",
	fx.Provide(
		route.AsRoute(basehandler.NewBaseHandler),
		route.AsRoute(image.NewImageHandler),
		route.AsRoute(vm.NewVmHandler),

		// Register routes to the route manager
		//进行注解，表明接收包含“routes”组内容的切片
		fx.Annotate(
			route.NewRouteManager,
			fx.ParamTags(`group:"routes"`),
		),
	),
	// Invoke the provided function after the module is fully initialized to register routes
	fx.Invoke(func(mgr route.RouteManager) {}),
)
