package apiserver

import (
	"embed"
	"encoding/json"
	ginI18n "github.com/gin-contrib/i18n"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"kubeall.io/api-server/pkg/infra/constants"
	"kubeall.io/api-server/pkg/types"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	// entrans "github.com/go-playground/validator/v10/translations/en"
)

type RestServer interface {
	Init(fs embed.FS)
	GetEngine() *gin.Engine
	RootGroup() *gin.RouterGroup
	NamespaceGroup() *gin.RouterGroup
	ClusterGroup() *gin.RouterGroup
}

type restServerImpl struct {
	config *types.ServerConfig
	engine *gin.Engine

	rootGroup      *gin.RouterGroup
	namespaceGroup *gin.RouterGroup
	clusterGroup   *gin.RouterGroup
}

func NewRestServer(config types.Config, fs embed.FS) RestServer {
	cfg := config.(*types.ServerConfig)
	restServer := &restServerImpl{
		config: cfg,
	}
	restServer.Init(fs)
	return restServer
}

func (r *restServerImpl) Init(fs embed.FS) {
	engine := r.setupRestServer(fs)
	r.engine = engine
	r.rootGroup = engine.Group(constants.RootUri)

	r.namespaceGroup = engine.Group(constants.NamespaceGroupUri, func(ctx *gin.Context) {
		// validate namespace is provided as expected
		namespace := ctx.Param("namespace")
		validate := binding.Validator.Engine().(*validator.Validate)
		err := validate.Var(namespace, "required")
		if err != nil {
			zap.L().Warn("namespace is required", zap.Error(err))
			ctx.AbortWithStatusJSON(http.StatusBadRequest,
				types.FailWithErrorCode(ctx, constants.CodeRequired, map[string]string{"name": "namespace"}))
		}
		// set resource type for current request
		ctx.Set(constants.ResourceType, types.NewResourceType(false, namespace))
	})

	r.clusterGroup = engine.Group(constants.ClusterGroupUri, func(ctx *gin.Context) {
		// set resource type for current request
		ctx.Set(constants.ResourceType, types.NewResourceType(true, ""))
	})
}

func (r *restServerImpl) setupRestServer(fs embed.FS) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	//Gin 默认的 MaxMultipartMemory 是 32MB（gin.Default() 中设置）。当上传的文件超过这个限制，Gin 会尝试将整个文件加载到内存中（multipart form 数据），导致内存使用量激增。
	//对于 100GB 的鏡像文件上传，内存占用可能达到数 GB 甚至更高（参考：上传 9MB 文件内存从 3MB 增到 30MB，）。这极有可能导致程序崩溃（OOM，Out of Memory）或服务器资源耗尽。
	engine.MaxMultipartMemory = 32 << 20

	// apply i18n middleware
	engine.Use(ginI18n.Localize(ginI18n.WithBundle(&ginI18n.BundleCfg{
		DefaultLanguage:  language.Chinese,
		FormatBundleFile: constants.ResourceBundleFormat,
		AcceptLanguage:   []language.Tag{language.Chinese},
		RootPath:         constants.ResourceRootDir,
		UnmarshalFunc:    json.Unmarshal,

		Loader: &ginI18n.EmbedLoader{
			FS: fs,
		},
	})))

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error logger.
	//   - Logs to stdout.
	//   - RFC3339 with local time format.
	engine.Use(ginzap.Ginzap(zap.L(), time.RFC3339, constants.UtcTimeUsed))

	// Logs all panic to error logger
	//   - stack means whether output the stack info.
	engine.Use(ginzap.RecoveryWithZap(zap.L(), r.config.LogSetting.PrintErrorStack))
	return engine
}

func (r *restServerImpl) GetEngine() *gin.Engine {
	return r.engine
}

func (r *restServerImpl) RootGroup() *gin.RouterGroup {
	return r.rootGroup
}

func (r *restServerImpl) NamespaceGroup() *gin.RouterGroup {
	return r.namespaceGroup
}

func (r *restServerImpl) ClusterGroup() *gin.RouterGroup {
	return r.clusterGroup
}
