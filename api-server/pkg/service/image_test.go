package service

import (
	"fmt"
	"kubeall.io/api-server/pkg/infra/constants"
	"kubeall.io/api-server/pkg/infra/utils"
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {
	uploadUrl := fmt.Sprintf("%s/%s?action=upload&",
		utils.GetEnv(constants.VarLonghornUploadUiPrefix, &constants.BackingImageUploadUri), "bi-win")
	reqUpload, err := http.NewRequest(uploadUrl, http.MethodPost, nil)
	if err != nil {
		t.Fatal(err)
	}
	println(reqUpload)
}
