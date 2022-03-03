package oss

import (
	"mime/multipart"
	"yuumi/internal/pkg/oss/qiniu"
)

type Local struct {
	Path string `mapstructure:"path" json:"path" yaml:"path" `
}

type OSS interface {
	UploadFile(file *multipart.FileHeader) (string, string, error)
	DeleteFile(key string) error
}

var Oss OSS

const (
	QiNiu = iota
)

func NewOss(ossType int) OSS {
	switch ossType {
	case QiNiu:
		Oss = &qiniu.QiNiu{}
	}

	return Oss
}
