package qiniu

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/viper"
	"mime/multipart"
	"time"
)

// UploadFile 上传文件
func (*QiNiu) UploadFile(file *multipart.FileHeader) (string, string, error) {
	bucket := viper.GetString("qiniu.bucket")
	ak := viper.GetString("qiniu.access-key")
	sk := viper.GetString("qiniu.secret-key")
	imgPath := viper.GetString("qiniu.img-path")
	putPolicy := storage.PutPolicy{Scope: bucket}
	mac := qbox.NewMac(ak, sk)

	// TODO token 存 redis
	upToken := putPolicy.UploadToken(mac)
	cfg := Config()
	formUploader := storage.NewFormUploader(cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{Params: map[string]string{"x:name": "github logo"}}

	f, openError := file.Open()
	if openError != nil {
		return "", "", openError
	}
	fileKey := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename) // 文件名格式 自己可以改 建议保证唯一性
	putErr := formUploader.Put(context.Background(), &ret, upToken, fileKey, f, file.Size, &putExtra)
	if putErr != nil {
		return "", "", putErr
	}

	return imgPath + "/" + ret.Key, ret.Key, nil
}
