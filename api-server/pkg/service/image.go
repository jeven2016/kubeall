package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kav1 "kubeall.io/api-server/pkg/generated/kubeall.io/v1"
	lhv1beta2 "kubeall.io/api-server/pkg/generated/longhorn/apis/longhorn/v1beta2"
	"kubeall.io/api-server/pkg/infra/apiserver"
	"kubeall.io/api-server/pkg/infra/constants"
	"kubeall.io/api-server/pkg/infra/utils"
	baseservice "kubeall.io/api-server/pkg/service/base"
	"kubeall.io/api-server/pkg/types"
	"mime/multipart"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const sleepTime = 3 * time.Second
const biImagePrefix = "bi"
const bufferSize = 32 * 1024 * 1024 // 32MB buffer cache
const NameMaximumLength = 40        // the max length of backing image's name

var (
	reclaimPolicy        = corev1.PersistentVolumeReclaimDelete
	allowVolumeExpansion = true
	volumeBindingMode    = storagev1.VolumeBindingImmediate
)

type ImageService interface {
	Upload(ctx context.Context, imageName string, req multipart.File, fileSize int64, request *http.Request) error
	EnsureImageResources(ctx context.Context, image *kav1.Image) error
	DeleteImageResources(ctx context.Context, image *kav1.Image) error
	UpdateStatus(ctx context.Context, imgStatus *kav1.ImageStatus, biImage *lhv1beta2.BackingImage) error
	ListImagesByType(ctx context.Context, namespace, imageType string) ([]kav1.Image, error)
}

type imageServiceImpl struct {
	clusterResource apiserver.ClusterResource
	storageClass    StorageClass
	baseService     baseservice.BaseService
	imageGvk        *schema.GroupVersionKind
}

func NewImageService(clusterResource apiserver.ClusterResource, sc StorageClass, baseService baseservice.BaseService,
	gvkResource *constants.GvkResource) (ImageService, error) {
	imageGvk, err := gvkResource.Get(constants.ImageResourceParam)
	if err != nil {
		zap.L().Fatal("no gvk resource defined", zap.Error(err))
		return nil, err
	}
	return &imageServiceImpl{
		clusterResource: clusterResource,
		storageClass:    sc,
		baseService:     baseService,
		imageGvk:        imageGvk,
	}, err
}

func (i imageServiceImpl) Upload(ctx context.Context, imageName string, file multipart.File, fileSize int64, request *http.Request) error {
	err := i.waitImage(ctx, imageName)
	if err != nil {
		return err
	}

	//upload the image
	err = i.uploadImageContent(imageName, file, fileSize, request)
	if err != nil {
		return err
	}

	zap.L().Info("image upload successfully", zap.String("imageName", imageName))
	return nil
}

func (i imageServiceImpl) EnsureImageResources(ctx context.Context, image *kav1.Image) error {
	biImage, err := i.ensureBackingImage(ctx, image)
	if err != nil {
		zap.L().Warn("failed to ensure backingimage to be created", zap.String("imageName", image.Name), zap.Error(err))
		return err
	}
	zap.L().Info("backing image existed", zap.String("biImage", biImage.Name))
	err = i.ensureStorageClass(ctx, image, biImage)
	if err != nil {
		zap.L().Warn("failed to ensure storageClass to be created", zap.String("imageName", image.Name), zap.Error(err))
		return err
	}
	zap.L().Info("backing image existed", zap.String("biImage", biImage.Name))
	return nil
}

func (i imageServiceImpl) ensureBackingImage(ctx context.Context, image *kav1.Image) (*lhv1beta2.BackingImage, error) {
	imageName := image.Name
	biName := i.buildBackingImageName(imageName)

	var biImage = &lhv1beta2.BackingImage{}
	err := i.clusterResource.ClusterCache().Get(ctx, client.ObjectKey{
		Namespace: constants.DefaultBackingImageNamespace,
		Name:      biName,
	}, biImage)

	if err != nil {
		if k8serrors.IsNotFound(err) {
			// create if not exist
			biImage = &lhv1beta2.BackingImage{
				ObjectMeta: v1.ObjectMeta{
					Name:      biName,
					Namespace: constants.DefaultBackingImageNamespace,
					Labels: map[string]string{
						constants.LabelImage:          imageName,
						constants.LabelImageNamespace: image.Namespace,
					},
				},
				Spec: lhv1beta2.BackingImageSpec{
					SourceType: lhv1beta2.BackingImageDataSourceType(image.Spec.ImageFrom),
				},
			}
			lhClient := i.clusterResource.Client().LonghornClient().LonghornV1beta2()
			return lhClient.BackingImages(constants.DefaultBackingImageNamespace).
				Create(ctx, biImage, v1.CreateOptions{})
		} else {
			return nil, err
		}
	}
	return biImage, nil
}

func (i imageServiceImpl) buildBackingImageName(imageName string) string {
	name := fmt.Sprintf("%s-%s", biImagePrefix, imageName)
	if len(name) > NameMaximumLength {
		name = name[:NameMaximumLength]
	}
	return name
}

func (i imageServiceImpl) uploadImageContent(imageName string, file multipart.File, fileSize int64, request *http.Request) error {
	// 4. 创建管道
	pr, pw := io.Pipe()
	defer pr.Close()

	bodyWriter := multipart.NewWriter(pw)

	// 5. 启动goroutine将上传文件写入管道
	go func() {
		defer func() { _ = pw.Close() }()
		defer func() { _ = bodyWriter.Close() }()
		part, err := bodyWriter.CreateFormFile("chunk", "blob")
		if err != nil {
			return
		}
		if _, err = io.Copy(part, file); err != nil {
			return
		}
	}()

	//bi image name is invlalid todo
	imageName = i.buildBackingImageName(imageName)
	uploadUrl := fmt.Sprintf("%s/%s?action=upload&size=%d",
		utils.GetEnv(constants.VarLonghornUploadUiPrefix, &constants.BackingImageUploadUri), imageName, fileSize)

	httpClient := &http.Client{Timeout: time.Minute * 30}
	resp, err := httpClient.Post(uploadUrl, bodyWriter.FormDataContentType(), pr)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	// get response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to get response body: %w", err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("failed to upload image's content: %d %s", resp.StatusCode, string(body))
	}
	return nil
}

func (i imageServiceImpl) waitImage(ctx context.Context, imageName string) error {
	retries := 20
	for j := 0; j < retries; j++ {
		biImage, err := i.GetBackingImageDataSource(ctx, imageName)
		if err != nil {
			if k8serrors.IsNotFound(err) {
				time.Sleep(sleepTime)
				continue
			}
			return err
		}
		if biImage != nil {

			// ready to upload while the status is pending
			if biImage.Status.CurrentState == lhv1beta2.BackingImageStatePending {
				zap.L().Info("the backing image's state is pending, upload later", zap.String("imageName", imageName))
				break
			}
			if biImage.Status.CurrentState == lhv1beta2.BackingImageStateFailed {
				zap.L().Warn("the backing image's state is failed", zap.String("imageName", imageName),
					zap.Any("state", biImage.Status.CurrentState))
				return types.FailWithErrorCode(ctx, constants.CodeBackingImageCreatedError, nil)
			}
		}
		time.Sleep(sleepTime)
	}
	return nil
}

func (i imageServiceImpl) GetBackingImage(ctx context.Context, imageName string) (*lhv1beta2.BackingImage, error) {
	var bi lhv1beta2.BackingImage
	err := i.clusterResource.ClusterCache().Get(ctx, client.ObjectKey{
		Namespace: constants.DefaultBackingImageNamespace,
		Name:      constants.BackingImagePrefix + imageName,
	}, &bi)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			zap.L().Warn("no backing image found", zap.String("name", imageName), zap.Error(err))
			return nil, nil
		}
		return nil, err
	}
	return &bi, nil
}

func (i imageServiceImpl) GetBackingImageDataSource(ctx context.Context, imageName string) (*lhv1beta2.BackingImageDataSource, error) {
	var bi lhv1beta2.BackingImageDataSource
	err := i.clusterResource.ClusterCache().Get(ctx, client.ObjectKey{
		Namespace: constants.DefaultBackingImageNamespace,
		Name:      i.buildBackingImageName(imageName),
	}, &bi)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			zap.L().Warn("no backing image datasource found", zap.String("name", imageName), zap.Error(err))
			return nil, nil
		}
		return nil, err
	}
	return &bi, nil
}

func (i imageServiceImpl) ensureStorageClass(ctx context.Context, image *kav1.Image, biImage *lhv1beta2.BackingImage) error {
	scName := image.Spec.StorageClassName
	if scName == "" {
		scName = image.Name
	}

	sc, err := i.storageClass.Get(ctx, image.Name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// create if not exist
			params := map[string]string{constants.ParamBiImageName: biImage.Name}
			imageParams := image.Spec.StorageClassParameters
			if imageParams != nil {
				for k, v := range imageParams {
					params[k] = v
				}
			}

			sc = &storagev1.StorageClass{
				ObjectMeta: v1.ObjectMeta{
					Name: scName,
				},
				Provisioner:          constants.LonghornDriver,
				Parameters:           params,
				ReclaimPolicy:        &reclaimPolicy,
				AllowVolumeExpansion: &allowVolumeExpansion,
				VolumeBindingMode:    &volumeBindingMode,
			}
			_, err = i.storageClass.Create(ctx, sc)
			return err
		}
		return err
	}
	return nil
}

func (i imageServiceImpl) DeleteImageResources(ctx context.Context, image *kav1.Image) error {
	// delete backing image
	lhClient := i.clusterResource.Client().LonghornClient()
	biName := i.buildBackingImageName(image.Name)
	err := lhClient.LonghornV1beta2().BackingImages(constants.DefaultBackingImageNamespace).Delete(ctx, biName, v1.DeleteOptions{})
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return err
		}
	}

	// delete storage class
	err = i.clusterResource.Client().K8sClient().StorageV1().StorageClasses().Delete(ctx, image.Name, v1.DeleteOptions{})
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (i imageServiceImpl) UpdateStatus(ctx context.Context, imgStatus *kav1.ImageStatus, biImage *lhv1beta2.BackingImage) error {
	// if backing image associated to an image
	imageName, imgOk := biImage.Labels[constants.LabelImage]
	imageNamespace, nsOk := biImage.Labels[constants.LabelImageNamespace]
	if imgOk && nsOk {
		image := &kav1.Image{}
		err := i.clusterResource.ClusterCache().Get(ctx, client.ObjectKey{
			Name:      imageName,
			Namespace: imageNamespace,
		}, image)

		if err != nil {
			return err
		}

		//update its status to track the uploading progress of backing image
		newImage := image.DeepCopy()
		newImage.Status = *imgStatus
		patch, _ := json.Marshal(newImage)
		if err = i.clusterResource.RuntimeClient().Status().
			Patch(ctx, newImage, client.RawPatch(client.Merge.Type(), patch)); err != nil {
			zap.L().Warn("failed to patch image's status",
				zap.String("name", newImage.Name), zap.Error(err))
			return nil
		}
	}
	zap.L().Info("image's status is updated", zap.String("image", imageName), zap.Any("status", imgStatus))
	return nil
}

func (i imageServiceImpl) ListImagesByType(ctx context.Context, namespace, imageType string) ([]kav1.Image, error) {
	resType := types.NewResourceType(false, namespace)
	query := types.Query{}
	result, err := i.baseService.List(ctx, *i.imageGvk, resType, query, false)
	if err != nil {
		return []kav1.Image{}, err
	}
	var images []kav1.Image
	for _, item := range result.Items {
		img := item.(*kav1.Image)
		if img.Spec.ImageType == kav1.ImageType(imageType) {
			images = append(images, *img)
		}
	}
	return images, nil
}
