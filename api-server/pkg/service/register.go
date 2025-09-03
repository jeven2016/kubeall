package service

import (
	"go.uber.org/fx"
	baseservice "kubeall.io/api-server/pkg/service/base"
)

var Module = fx.Module("service",
	fx.Provide(
		baseservice.NewBaseService,
		baseservice.NewFieldSorter,
		baseservice.NewFieldFilter,
		NewImageService,
		NewStorageClass,
		NewVmService,
	),
)
