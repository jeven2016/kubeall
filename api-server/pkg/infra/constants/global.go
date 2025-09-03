package constants

const (
	NamespaceAll         = "all"
	ResourceRootDir      = "./resources/locales/"
	ResourceBundleFormat = "json"
	UtcTimeUsed          = false
	ResourceType         = "resourceType"
	PageQueryField       = "page"
	PageSizeQueryField   = "pageSize"
	SortByQueryField     = "sortBy"
	SortOrderField       = "sortOrder"
	FilterField          = "filter"
	DefaultPage          = "1"
	DefaultPageSize      = "10"

	RootUri                = "/api/v1"
	ClusterGroupUri        = RootUri + "/clusters"
	NamespaceGroupUri      = RootUri + "/namespaces/:namespace"
	ResourceUri            = "/:resource"
	ResourceNameUri        = ResourceUri + "/:name"
	ResourceImageUri       = "/images"
	ResourceVmUri          = "/vms"
	ResourceImageUploadUri = ResourceImageUri + "/:imageName/upload"
	ResourceParam          = "resource"
	ImageResourceParam     = "images"

	JsonFormat                   = "json"
	YamlFormat                   = "yaml"
	DefaultBackingImageNamespace = "longhorn-system"
	BackingImagePrefix           = "bi-"
	LonghornDriver               = "driver.longhorn.io"
	ParamBiImageName             = "backingImage"

	LabelImage          = "kubeall.io/image"
	LabelImageNamespace = "kubeall.io/imageNamespace"
	DefaultFinalizer    = "kubeall.io/finalizer"

	VarLonghornUploadUiPrefix = "LONGHORN_UPLOAD_URL_PREFIX"

	AnnotationPvcTemplates = "kubeall.io/pvcTemplates"

	MaxConcurrentReconciles = 2

	ValidateImageType = "required,oneof=iso disk"
)

var (
	AvailablePageSizes = []int{10, 20, 50, 100}
	//BackingImageUploadUri = "http://longhorn-backend.longhorn-system:9500/v1/backingimages"
	BackingImageUploadUri = "http://192.168.1.66:31492/v1/backingimages"
)
