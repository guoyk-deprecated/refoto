package main

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	var err error
	defer func(err *error) {
		if *err != nil {
			log.Println("exited with error:", (*err).Error())
			os.Exit(1)
		} else {
			log.Println("exited")
		}
	}(&err)

	rand.Seed(time.Now().Unix())

	log.SetOutput(os.Stdout)

	if err = setupEnv(); err != nil {
		return
	}

	var db *gorm.DB
	if db, err = setupDB(); err != nil {
		return
	}

	var bucket *oss.Bucket
	if bucket, err = setupOSS(); err != nil {
		return
	}

	var r *TemplateEngine
	if r, err = NewTemplateEngine(TemplateEngineOptions{
		Debug: envDebug,
		Dir:   "views",
		Ext:   "gohtml",
	}); err != nil {
		return
	}

	e := echo.New()
	e.Debug = envDebug
	e.HideBanner = true
	e.HidePort = true
	e.Renderer = r
	e.HTTPErrorHandler = r.ErrorHandler("error", envTitle)
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	setupRoutes(e, db, bucket)

	err = e.Start(":4000")
}
