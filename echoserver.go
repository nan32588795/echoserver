package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"echoserver/config"
	"echoserver/handler"
)

func main() {

	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal("設定ファイルの読み込み失敗:", err)
	}

	db, err := sql.Open("postgres", cfg.DB.ConnString())
	// connStr := "user=postgres password=postgres dbname=postgres sslmode=disable"
	// db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	boil.SetDB(db)
	boil.DebugMode = true
	boil.DebugWriter = os.Stdout

	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	e.GET("/hello", func(c echo.Context) error {
		name := c.QueryParam("name")
		if name == "" {
			name = "ゲスト"
		}
		return c.String(http.StatusOK, "こんにちは、"+name+" さん！")
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

	e.Logger.Fatal(e.Start(":8080"))
}
