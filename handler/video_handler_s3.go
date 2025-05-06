package handler

import (
	"bytes"
	"net/http"

	"github.com/labstack/echo/v4"
)

type VideoHandlerS3 struct {
}

func NewVideoHandlerS3() *VideoHandlerS3 {
	return &VideoHandlerS3{}
}

func (v *VideoHandlerS3) UploadFile(c echo.Context) error {
	// S3にファイルをアップロードする処理を実装
	return c.String(http.StatusOK, "ファイルアップロード成功: "+"sample.txt")
}

func (v *VideoHandlerS3) DownloadFile(c echo.Context) error {
	// S3からファイルをダウンロードする処理を実装
	// オンメモリのダミーファイルを生成する
	var buf bytes.Buffer
	buf.WriteString("これはオンメモリのファイルです\n")

	// クライアントに「ファイルとしてダウンロードさせる」ためにヘッダを設定
	c.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="sample.txt"`)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)

	// io.Reader を使って Stream 返却
	return c.Stream(http.StatusOK, echo.MIMETextPlain, &buf)
}

func (v *VideoHandlerS3) GetVideos(c echo.Context) error {
	// S3から動画ファイルのリストを取得する処理を実装
	var fileList []string = []string{
		"sample1.mp4",
		"sample2.mp4",
		"sample3.mp4",
	}
	return c.JSON(http.StatusOK, fileList)
}
