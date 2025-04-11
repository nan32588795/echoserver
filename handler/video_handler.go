package handler

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

var VIDEO_PATH = "storage/videos"

type VideoHandler struct {
}

func NewVideoHandler() *VideoHandler {
	os.MkdirAll(VIDEO_PATH, os.ModePerm)
	return &VideoHandler{}
}

func (h *VideoHandler) UploadFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ファイルが含まれていません")
	}

	src, err := file.Open()
	if err != nil {
		c.Logger().Error("Openに失敗", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer src.Close()

	dstPath := filepath.Join(VIDEO_PATH, file.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		c.Logger().Error("Createに失敗", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		c.Logger().Error("Copyに失敗", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.String(http.StatusOK, "ファイルアップロード成功: "+file.Filename)
}

func (h *VideoHandler) DownloadFile(c echo.Context) error {
	filename := c.Param("filename")
	filePath := filepath.Join(VIDEO_PATH, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.String(http.StatusNotFound, "ファイルが見つかりません")
	}

	return c.Attachment(filePath, filename)
}

func (h *VideoHandler) GetVideos(c echo.Context) error {
	files, err := os.ReadDir(VIDEO_PATH)
	if err != nil {
		c.Logger().Error("ReadDirに失敗: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "ファイル一覧の取得に失敗しました")
	}

	var fileList []string
	for _, f := range files {
		if !f.IsDir() {
			fileList = append(fileList, f.Name())
		}
	}

	return c.JSON(http.StatusOK, fileList)
}
