package basehandler

type ResourceType interface {
	ClusterResource() bool
	Namespace() string
}

type clusterResource struct {
	isClusterResource bool
	namespace         string
}

func (c clusterResource) ClusterResource() bool {
	return c.isClusterResource
}

func (c clusterResource) Namespace() string {
	return c.namespace
}

func NewResourceType(isClusterResource bool, namespace string) ResourceType {
	return &clusterResource{
		isClusterResource: isClusterResource,
		namespace:         namespace,
	}
}
