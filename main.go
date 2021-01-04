package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"os"
)

var (
	envPort     = 4000
	envTitle    = "Refoto"
	envDebug    = false
	envMySQLDSN = ""
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
	if err = envStr("REFOTO_MYSQL_DSN", &envMySQLDSN); err != nil {
		return
	}

	var db *gorm.DB
	if db, err = gorm.Open(mysql.Open(envMySQLDSN), &gorm.Config{}); err != nil {
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

	routes(e, db)

	err = e.Start(fmt.Sprintf(":%d", envPort))
}

func routes(e *echo.Echo, db *gorm.DB) {
	e.GET("/", func(ctx echo.Context) error {
		type DataGirl struct {
			ID        string
			AvatarURL string
		}
		type DataEvent struct {
			ID    string
			Name  string
			Girls []DataGirl
		}
		type Data struct {
			Title  string
			Events []DataEvent
		}
		var data Data
		data.Title = envTitle
		// TODO: remove dummy data
		mi := 5 + rand.Intn(10)
		for i := 0; i < mi; i++ {
			event := DataEvent{
				ID:   fmt.Sprintf("dummy-%02d", i),
				Name: fmt.Sprintf("测试数据-%02d", i),
			}
			mj := 7 + rand.Intn(10)
			for j := 0; j < mj; j++ {
				event.Girls = append(event.Girls, DataGirl{
					ID:        fmt.Sprintf("dummy-%02d-%02d", i, j),
					AvatarURL: "https://via.placeholder.com/150",
				})
			}
			data.Events = append(data.Events, event)
		}
		return ctx.Render(http.StatusOK, "index", data)
	})
}
