package constants

import (
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type GvkResource struct {
	namespaceResource map[string]*schema.GroupVersionKind
	clusterResource   map[string]*schema.GroupVersionKind
}

func NewGvkResource() *GvkResource {
	gvk := &GvkResource{
		namespaceResource: make(map[string]*schema.GroupVersionKind),
		clusterResource:   make(map[string]*schema.GroupVersionKind),
	}
	gvk.initGvkResource()
	return gvk
}

func (g GvkResource) initGvkResource() {
	g.namespaceResource["pods"] = &schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"}
	g.namespaceResource["services"] = &schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"}
	g.namespaceResource["secrets"] = &schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"}
	g.namespaceResource["configmaps"] = &schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ConfigMap"}
	g.namespaceResource["pvs"] = &schema.GroupVersionKind{Group: "", Version: "v1", Kind: "PersistentVolume"}
	g.namespaceResource["pvcs"] = &schema.GroupVersionKind{Group: "", Version: "v1", Kind: "PersistentVolumeClaim"}
	g.namespaceResource["nodes"] = &schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Node"}
	g.namespaceResource["namespaces"] = &schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"}
	g.namespaceResource["deployments"] = &schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	g.namespaceResource["daemonsets"] = &schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "DaemonSet"}
	g.namespaceResource["statefulsets"] = &schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"}
	g.namespaceResource["ingresses"] = &schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"}
	g.namespaceResource["jobs"] = &schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "Job"}
	g.namespaceResource["cronjobs"] = &schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "Cronjob"}
	g.namespaceResource["hpas"] = &schema.GroupVersionKind{Group: "autoscaling", Version: "v1", Kind: "HorizontalPodAutoScaler"}
	g.namespaceResource["crds"] = &schema.GroupVersionKind{Group: "apiextensions.k8s.io", Version: "v1", Kind: "CustomResourceDefinition"}
	g.namespaceResource["sc"] = &schema.GroupVersionKind{Group: "storage.k8s.io", Version: "v1", Kind: "StorageClass"}

	g.namespaceResource["pools"] = &schema.GroupVersionKind{Group: "metallb.io", Version: "v1beta1", Kind: "IPAddressPool"}

	g.namespaceResource["images"] = &schema.GroupVersionKind{Group: "api.kubeall.io", Version: "v1", Kind: "Image"}
	g.namespaceResource["globalsettings"] = &schema.GroupVersionKind{Group: "api.kubeall.io", Version: "v1", Kind: "GlobalSettings"}

	g.namespaceResource["vms"] = &schema.GroupVersionKind{Group: "kubevirt.io", Version: "v1", Kind: "VirtualMachine"}
	g.namespaceResource["vminstances"] = &schema.GroupVersionKind{Group: "kubevirt.io", Version: "v1", Kind: "VirtualMachineInstance"}
	g.namespaceResource["vmiReplicaSets"] = &schema.GroupVersionKind{Group: "kubevirt.io", Version: "v1", Kind: "VirtualMachineInstanceReplicaSet"}
	g.namespaceResource["vmiPresents"] = &schema.GroupVersionKind{Group: "kubevirt.io", Version: "v1", Kind: "VirtualMachineInstancePreset"}
	g.namespaceResource["vmiMigration"] = &schema.GroupVersionKind{Group: "kubevirt.io", Version: "v1", Kind: "VirtualMachineInstanceMigration"}
	g.namespaceResource["kubevirt"] = &schema.GroupVersionKind{Group: "kubevirt.io", Version: "v1", Kind: "KubeVirt"}

	crs := rbacv1.SchemeGroupVersion.WithKind("ClusterRole")
	g.clusterResource["clusterroles"] = &crs

	crbs := rbacv1.SchemeGroupVersion.WithKind("ClusterRoleBinding")
	g.clusterResource["clusterrolebindings"] = &crbs
}

// Get retrieves the GroupVersionKind for a given resource key.
// It searches first in namespace resources, then in cluster resources.
//
// Parameters:
//   - key: A string representing the resource key to look up.
//
// Returns:
//   - A pointer to a schema.GroupVersionKind if found, or nil if not found.
func (g GvkResource) Get(key string) (*schema.GroupVersionKind, error) {
	if val, ok := g.namespaceResource[key]; ok {
		return val, nil
	}
	if val, ok := g.clusterResource[key]; ok {
		return val, nil
	}
	return nil, InvalidKind
}
