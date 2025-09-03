package apiserver

import (
	"context"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"kubeall.io/api-server/pkg/infra/clients"
	"kubeall.io/api-server/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
)

type ClusterResource interface {
	RestServer() RestServer
	RestConfig() *rest.Config
	ClusterCache() cache.Cache
	Cluster() cluster.Cluster
	Client() clients.ApiClient
	RuntimeClient() client.Client
	UpdateCluster(cls cluster.Cluster)
}

type clusterResourceImpl struct {
	config        *types.ServerConfig
	restServer    RestServer
	restConfig    *rest.Config
	clusterCache  cache.Cache
	cluster       cluster.Cluster
	client        clients.ApiClient
	runTimeClient client.Client
}

// NewClusterResource creates and returns a new Server instance.
func NewClusterResource(config types.Config, logger *zap.Logger,
	restServer RestServer, apiClient clients.ApiClient, schemeType types.SchemeType) ClusterResource {
	cfg := config.(*types.ServerConfig)
	cls := &clusterResourceImpl{
		config:     cfg,
		restServer: restServer,
		restConfig: apiClient.RestConfig(),
		client:     apiClient,
	}
	if err := cls.initCluster(context.Background(), schemeType); err != nil {
		zap.L().Error("failed to initialize cluster", zap.Error(err))
		panic(err)
	}
	return cls
}

// NewClusterResourceForCm the rest server is not required for controller mananger
func NewClusterResourceForCm(config types.Config, logger *zap.Logger,
	apiClient clients.ApiClient) ClusterResource {
	cfg := config.(*types.ServerConfig)

	//manager implements cluster interface, so no need to initialize an extra cache
	cls := &clusterResourceImpl{
		config:     cfg,
		restConfig: apiClient.RestConfig(),
		client:     apiClient,
	}
	return cls
}

// initCluster initializes the Kubernetes cluster client and starts the cluster.
// It sets up the necessary schemes, creates a new cluster client, and starts
// the cluster in a separate goroutine.
//
// Parameters:
//   - ctx: A context.Context that can be used to cancel the cluster initialization.
//
// Returns:
//   - error: An error if cluster initialization fails, nil otherwise.
func (s *clusterResourceImpl) initCluster(ctx context.Context, schemeType types.SchemeType) error {
	var sch *runtime.Scheme
	if schemeType == types.ApiServerScheme {
		sch = ServerScheme
	}
	if schemeType == types.CmScheme {
		sch = CmScheme
	}
	if sch == nil {
		panic("")
	}

	zap.L().Info("current scheme initialized is " + string(schemeType))
	c, err := cluster.New(s.restConfig, func(clusterOptions *cluster.Options) {
		clusterOptions.Scheme = sch
	})
	if err != nil {
		return err
	}
	s.cluster = c
	s.clusterCache = c.GetCache()
	s.runTimeClient = c.GetClient()

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		zap.S().Info("starts a background job for informers")
		err = s.cluster.Start(ctx)
		if err != nil {
			zap.L().Error("failed to start cluster", zap.Error(err))
			cancel()
		}
	}()
	return nil
}

// UpdateCluster update related resources by controller manager
func (s *clusterResourceImpl) UpdateCluster(cls cluster.Cluster) {
	s.cluster = cls
	s.clusterCache = cls.GetCache()
	s.runTimeClient = cls.GetClient()
}

func (s *clusterResourceImpl) RestServer() RestServer {
	return s.restServer
}
func (s *clusterResourceImpl) RestConfig() *rest.Config {
	return s.restConfig
}
func (s *clusterResourceImpl) ClusterCache() cache.Cache {
	return s.clusterCache
}
func (s *clusterResourceImpl) Cluster() cluster.Cluster {
	return s.cluster
}
func (s *clusterResourceImpl) Client() clients.ApiClient {
	return s.client
}

func (s *clusterResourceImpl) RuntimeClient() client.Client {
	return s.runTimeClient
}
