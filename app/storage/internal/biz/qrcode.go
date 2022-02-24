package biz

import (
	"stb-library/lib/qrcode"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type QrcodeUseCase struct {
	defaultFileDir DefaultFileDir
	log            *log.Helper
}

func NewQrcodeCase(defaultDir DefaultFileDir, logger log.Logger) *QrcodeUseCase {
	return &QrcodeUseCase{defaultFileDir: defaultDir, log: log.NewHelper(logger)}
}

func (q *QrcodeUseCase) qrcodeEncoder(mes string) (string, error) {
	return qrcode.GenerateQR(mes)
}

func (q *QrcodeUseCase) qrcodeDecoder(imageURL string) (string, error) {
	code, err := qrcode.Decode(imageURL)
	if err != nil {
		return "", err
	}
	return code.Content, nil
}

// QrcodeDecoder 对纯数字的二维码内容无法解析，反馈为乱码
func (q *QrcodeUseCase) QrcodeDecoder(ctx *gin.Context) ([]string, error) {
	filePaths, err := getAllFormFile(ctx, q.defaultFileDir.DefaultAssetsPath)
	if err != nil {
		return nil, err
	}
	contents := []string{}
	for _, filePath := range filePaths {
		content, err := q.qrcodeDecoder(filePath)
		if err != nil {
			return contents, err
		}

		contents = append(contents, content)
	}
	return contents, nil
}
