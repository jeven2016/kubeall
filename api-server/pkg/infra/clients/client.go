package clients

import (
	"go.uber.org/zap"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	lhclient "kubeall.io/api-server/pkg/generated/longhorn/clientset/versioned"
	"kubeall.io/api-server/pkg/infra/utils"
	"kubeall.io/api-server/pkg/types"
	kvclient "kubevirt.io/client-go/kubevirt"
)

type ApiClient interface {
	K8sClient() *kubernetes.Clientset
	RestConfig() *rest.Config
	LonghornClient() *lhclient.Clientset
	KubevirtClient() *kvclient.Clientset
}

type apiClientsImpl struct {
	config         types.Config
	restConfig     *rest.Config
	k8sClient      *kubernetes.Clientset
	longhornClient *lhclient.Clientset
	kvClient       *kvclient.Clientset
}

func NewClients(cfg types.Config) ApiClient {
	apiClient := &apiClientsImpl{
		config: cfg,
	}
	apiClient.createClients()
	return apiClient
}

// CreateClients 根据kubernetes配置文件创建各类客户端
func (a *apiClientsImpl) createClients() {
	kubeConfigFile := a.config.(*types.ServerConfig).KubeConfig
	if kubeConfigFile != "" {
		if exist, _ := utils.IsFileExists(kubeConfigFile); !exist {
			zap.L().Error("the kube config file doesn't exist, ", zap.String("configFile", kubeConfigFile))
			return
		}
	}

	var err error

	a.restConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigFile)
	utilruntime.Must(err)

	// k8s
	a.k8sClient, err = kubernetes.NewForConfig(a.restConfig)
	utilruntime.Must(err)

	a.longhornClient, err = lhclient.NewForConfig(a.restConfig)
	utilruntime.Must(err)

	// kubevirt client
	a.kvClient = kvclient.NewForConfigOrDie(a.restConfig)

}

func (a *apiClientsImpl) K8sClient() *kubernetes.Clientset {
	return a.k8sClient
}

func (a *apiClientsImpl) RestConfig() *rest.Config {
	return a.restConfig
}

func (a *apiClientsImpl) LonghornClient() *lhclient.Clientset {
	return a.longhornClient
}

func (a *apiClientsImpl) KubevirtClient() *kvclient.Clientset {
	return a.kvClient
}
