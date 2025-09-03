package controller

import (
	"fmt"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"kubeall.io/api-server/pkg/infra/apiserver"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"time"
)

type Manager interface {
	Start()
	GetManager() manager.Manager
	ClusterResource() apiserver.ClusterResource
}

type managerImpl struct {
	mgr             manager.Manager
	clusterResource apiserver.ClusterResource
	logger          *zap.Logger
}

func (m *managerImpl) GetManager() manager.Manager {
	return m.mgr
}

func (m *managerImpl) ClusterResource() apiserver.ClusterResource {
	return m.clusterResource
}

func NewManager(reconcilers []ReconcileHandler, clusterResource apiserver.ClusterResource, logger *zap.Logger) Manager {
	m := &managerImpl{clusterResource: clusterResource, logger: logger}
	m.Initialize()
	m.injectManager(reconcilers)
	return m
}

func (m *managerImpl) Initialize() {
	certDir := ""
	port := 8085
	leaderElect := false
	LeaseDuration := 30 * time.Second
	RetryPeriod := 10 * time.Second
	RenewDeadline := 5 * time.Second

	webhookServer := webhook.NewServer(webhook.Options{
		CertDir: certDir, //TODO
		Port:    port,
	})
	opts := ctrl.Options{
		Scheme:                 apiserver.CmScheme,
		WebhookServer:          webhookServer,
		HealthProbeBindAddress: fmt.Sprintf(":%d", port),
		Metrics: server.Options{
			BindAddress: "0", // 禁用 metrics 服务器
		},
	}

	if leaderElect {
		opts.LeaderElection = true
		opts.LeaderElectionNamespace = "kubeall-system"
		opts.LeaderElectionID = "kubeall-leader-election"
		opts.LeaseDuration = &LeaseDuration
		opts.RetryPeriod = &RetryPeriod
		opts.RenewDeadline = &RenewDeadline
	}

	restConfig := m.clusterResource.RestConfig()
	mgr, err := manager.New(restConfig, opts)
	if err != nil {
		zap.L().Error("failed to start controller manager")
		os.Exit(1)
	}

	// update manager as global cluster object
	m.clusterResource.UpdateCluster(mgr)
	m.mgr = mgr

	logger := zapr.NewLogger(m.logger)

	// inject to controller-runtime's logger
	ctrl.SetLogger(logger)
}

func (m *managerImpl) Start() {
	zap.L().Info("controller manager started")
	if err := m.mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		zap.L().Error("failed to run manager", zap.Error(err))
		os.Exit(1)
	}
}

func (m *managerImpl) injectManager(reconcilers []ReconcileHandler) {
	for _, h := range reconcilers {
		err := h.SetupWithManager(m.mgr)
		if err != nil {
			zap.L().Error("failed to setup manager", zap.Error(err))
			utilruntime.Must(err)
		}
	}
	zap.L().Info("controllers are set up")
}
