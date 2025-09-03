package apiserver

import (
	metallbv1beta1 "go.universe.tf/metallb/api/v1beta1"
	metallbv1beta2 "go.universe.tf/metallb/api/v1beta2"
	corev1 "k8s.io/api/core/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	k8scheme "k8s.io/client-go/kubernetes/scheme"
	kav1 "kubeall.io/api-server/pkg/generated/kubeall.io/v1"
	lhscheme "kubeall.io/api-server/pkg/generated/longhorn/clientset/versioned/scheme"
	kubevirtscheme "kubevirt.io/client-go/kubevirt/scheme"
)

// CmScheme for controller manager
var CmScheme = runtime.NewScheme()

var ServerScheme = runtime.NewScheme()

// for api server
func init() {
	metav1.AddToGroupVersion(ServerScheme, metav1.SchemeGroupVersion)
	utilruntime.Must(k8scheme.AddToScheme(ServerScheme))
	utilruntime.Must(corev1.AddToScheme(ServerScheme))
	utilruntime.Must(apiextv1.AddToScheme(ServerScheme))
	utilruntime.Must(lhscheme.AddToScheme(ServerScheme))
	utilruntime.Must(kav1.AddToScheme(ServerScheme))
	utilruntime.Must(kubevirtscheme.AddToScheme(ServerScheme))
	utilruntime.Must(metallbv1beta1.AddToScheme(ServerScheme))
	utilruntime.Must(metallbv1beta2.AddToScheme(ServerScheme))
}

// for controller manager
func init() {
	utilruntime.Must(k8scheme.AddToScheme(CmScheme))
	utilruntime.Must(corev1.AddToScheme(CmScheme))
	utilruntime.Must(lhscheme.AddToScheme(CmScheme))
	utilruntime.Must(kav1.AddToScheme(CmScheme))
	utilruntime.Must(kubevirtscheme.AddToScheme(CmScheme))
}
