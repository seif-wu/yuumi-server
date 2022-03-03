package qiniu

import (
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/viper"
)

// DeleteFile 删除文件
func (*QiNiu) DeleteFile(key string) error {
	bucket := viper.GetString("qiniu.bucket")
	ak := viper.GetString("qiniu.access-key")
	sk := viper.GetString("qiniu.secret-key")
	mac := qbox.NewMac(ak, sk)
	cfg := Config()
	bucketManager := storage.NewBucketManager(mac, cfg)
	if err := bucketManager.Delete(bucket, key); err != nil {
		return err
	}
	return nil
}
