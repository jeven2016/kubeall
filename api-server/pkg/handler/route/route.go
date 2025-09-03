package route

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"kubeall.io/api-server/pkg/infra/apiserver"
)

// Route represents a route
type Route interface {
	RegisterRoutes(rootGroup *gin.RouterGroup, namespaceGroup *gin.RouterGroup, clusterGroup *gin.RouterGroup)
}

// RouteManager represents a route manager
type RouteManager interface{}

// routeManagerImpl represents a route manager implementation
type routeManagerImpl struct{}

// NewRouteManager creates and initializes a new RouteManager.
// It registers all the provided routes and returns a RouteManager instance.
//
// Parameters:
//   - routes: A slice of Route interfaces. Each Route in the slice will have its RegisterRoutes method called.
//
// Returns:
//   - RouteManager: A new instance of RouteManager after registering all provided routes.
func NewRouteManager(routes []Route, restServer apiserver.RestServer) RouteManager {
	for _, r := range routes {
		r.RegisterRoutes(restServer.RootGroup(), restServer.NamespaceGroup(), restServer.ClusterGroup())
	}
	return &routeManagerImpl{}
}

// AsRoute wraps a function to be used as a Route in the dependency injection framework.
// It annotates the function to be recognized as a Route and groups it under the "routes" tag.
//
// Parameters:
//   - f: any - The function to be wrapped as a Route. It should implement the Route interface.
//
// Returns:
//   - any - The annotated function, ready to be used in the dependency injection container.
func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}
