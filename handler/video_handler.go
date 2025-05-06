package handler

import (
	"github.com/labstack/echo/v4"
)

// インターフェースの定義
type VideoService interface {
	UploadFile(c echo.Context) error
	DownloadFile(c echo.Context) error
	GetVideos(c echo.Context) error
}
