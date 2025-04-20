package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"echoserver/config"
	"echoserver/handler"
	"echoserver/validator"
)

func main() {

	db, err := sql.Open(config.GlobalConfig.DriverName(), config.GlobalConfig.ConnString())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	boil.SetDB(db)
	boil.DebugMode = true
	boil.DebugWriter = os.Stdout

	e := echo.New()
	e.Validator = validator.NewCustomValidator()

	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	// タイムアウトミドルウェア（3秒）
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		// Skipper: func(c echo.Context) bool {
		// 	// 特定のパスのみこのミドルウェアを適用する
		// 	return c.Path() != "/timeout/test2"
		// },
		// ErrorMessage: "忙しいから話かけないで",
		Timeout: 3 * time.Second,
	}))

	e.GET("/hello", func(c echo.Context) error {
		name := c.QueryParam("name")
		if name == "" {
			name = "ゲスト"
		}
		return c.String(http.StatusOK, "こんにちは、"+name+" さん！")
	})

	e.GET("/timeout/test1", func(c echo.Context) error {
		timeoutStr := c.QueryParam("timeout")
		sleepStr := c.QueryParam("sleep")
		var timeoutSec int
		var sleepSec int
		if t, err := strconv.Atoi(timeoutStr); err == nil {
			timeoutSec = t
		} else {
			timeoutSec = 5
		}
		if s, err := strconv.Atoi(sleepStr); err == nil {
			sleepSec = s
		} else {
			sleepSec = 100
		}

		ctx, cancel := context.WithTimeout(c.Request().Context(), time.Duration(timeoutSec)*time.Second)
		defer cancel()

		_, err := db.ExecContext(ctx, fmt.Sprintf("SELECT SLEEP(%d)", sleepSec))
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return c.String(http.StatusGatewayTimeout, "タイムアウト")
			} else {
				return c.String(http.StatusInternalServerError, "DBエラー")
			}
		}
		return c.String(http.StatusOK, "DBクエリ成功")
	})

	e.GET("/timeout/test2", func(c echo.Context) error {
		sleepStr := c.QueryParam("sleep")
		var sleepSec int
		if s, err := strconv.Atoi(sleepStr); err == nil {
			sleepSec = s
		} else {
			sleepSec = 100
		}

		db.ExecContext(c.Request().Context(), fmt.Sprintf("SELECT SLEEP(%d)", sleepSec))
		return c.String(http.StatusOK, "DBクエリ成功")
	})

	videoHandler := handler.NewVideoHandler()
	e.POST("/videos/upload", videoHandler.UploadFile)
	e.GET("/videos/download/:filename", videoHandler.DownloadFile)
	e.GET("/videos/list", videoHandler.GetVideos)

	userHandler := handler.NewUserHandler()
	e.POST("/users", userHandler.CreateUser)
	e.GET("/users", userHandler.GetUsers)
	e.GET("/users/:id", userHandler.GetUserByID)
	e.PUT("/users/:id", userHandler.UpdateUser)
	e.DELETE("/users/:id", userHandler.DeleteUser)

	// server := &http.Server{
	// 	Addr:         ":8080",
	// 	Handler:      e,
	// 	ReadTimeout:  3 * time.Second,
	// 	WriteTimeout: 3 * time.Second,
	// }
	// e.Logger.Fatal(server.ListenAndServe())

	e.Logger.Fatal(e.Start(":8080"))
}
