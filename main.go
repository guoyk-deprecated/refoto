package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
)

var (
	envPort  = 4000
	envTitle = "Refoto"
	envDebug = false
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

	if err = envInt("REFOTO_PORT", &envPort); err != nil {
		return
	}
	if err = envStr("REFOTO_TITLE", &envTitle); err != nil {
		return
	}
	if err = envBool("REFOTO_DEBUG", &envDebug); err != nil {
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
	e.HTTPErrorHandler = r.ErrorHandler("error")
	e.Use(middleware.Recover())

	err = e.Start(fmt.Sprintf(":%d", envPort))
}
