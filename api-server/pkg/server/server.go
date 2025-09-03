package server

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"kubeall.io/api-server/pkg/infra/apiserver"
	"kubeall.io/api-server/pkg/types"
	"log"
)

type Server interface {
	Start(context.Context, embed.FS) error
}

type serverImpl struct {
	restServer apiserver.RestServer
	config     *types.ServerConfig
}

func NewServer(cfg types.Config,
	restServer apiserver.RestServer) Server {
	return &serverImpl{
		config:     cfg.(*types.ServerConfig),
		restServer: restServer,
	}
}

// Start initializes and runs the server.
// It sets up a cancellable context, prints the configuration,
// initializes the cluster, and starts the REST server.
//
// Parameters:
//   - ctx: A context.Context for controlling the lifecycle of the server.
//   - fs: An embed.FS containing embedded files required for the server.
//
// Returns:
//   - error: An error if the server fails to start or encounters any issues during initialization.
//     Returns nil if the server starts successfully.
func (s *serverImpl) Start(ctx context.Context, fs embed.FS) error {
	s.printConfig()

	zap.L().Info("server started")
	err := s.restServer.GetEngine().Run(fmt.Sprintf("%s:%d", s.config.Http.Address, s.config.Http.Port))
	if err != nil {
		return err
	}
	return nil
}

func (s *serverImpl) printConfig() {
	if cfgString, err := json.Marshal(s.config); err != nil {
		log.Printf("failed to marshal config: %s", err)
	} else {
		zap.L().Info("the config takes effect", zap.String("config", string(cfgString)))
	}
}
