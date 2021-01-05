package main

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	"math/rand"
	"mime/multipart"
	"net/http"
	"sort"
	"strconv"
)

func sessionIsAdmin(ctx echo.Context) bool {
	sess, _ := session.Get("session", ctx)
	if yes, ok := sess.Values["is_admin"].(bool); ok {
		return yes
	}
	return false
}

func sessionSetAdmin(ctx echo.Context, isAdmin bool) {
	sess, _ := session.Get("session", ctx)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["is_admin"] = isAdmin
	_ = sess.Save(ctx.Request(), ctx.Response())
}

func requireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if !sessionIsAdmin(ctx) {
			return errors.New("没有管理员权限")
		}
		return next(ctx)
	}
}

func setupRoutes(e *echo.Echo, db *gorm.DB, bucket *oss.Bucket) {
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(envSecret))))
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{TokenLookup: "form:_csrf"}))
	e.GET("/", func(ctx echo.Context) error {
		type Data struct {
			Title   string
			IsAdmin bool
			CSRF    string
			Events  []Event
		}
		var data Data
		data.CSRF, _ = ctx.Get("csrf").(string)
		data.IsAdmin = sessionIsAdmin(ctx)
		data.Title = envTitle
		if err := db.Preload("Girls").Order("id DESC").Find(&data.Events).Error; err != nil {
			return err
		}
		for i := range data.Events {
			e := &data.Events[i]
			sort.Slice(e.Girls, func(i, j int) bool {
				return e.Girls[i].ID < e.Girls[j].ID
			})
		}
		//AvatarPath: "https://via.placeholder.com/150",
		return ctx.Render(http.StatusOK, "index", data)
	})
	e.GET("/girls/:girl_id", func(ctx echo.Context) error {
		type Data struct {
			Title            string
			IsAdmin          bool
			TokenExisted     bool
			TokenMatched     bool
			Contact          string
			CSRF             string
			Girl             Girl
			Event            Event
			PhotosOriginal   []Photo
			PhotosRoughTuned []Photo
			PhotosFineTuned  []Photo
		}
		var data Data
		data.CSRF, _ = ctx.Get("csrf").(string)
		data.Contact = envContact
		data.IsAdmin = sessionIsAdmin(ctx)
		data.Title = envTitle
		if err := db.Preload("Photos").Find(&data.Girl, ctx.Param("girl_id")).Error; err != nil {
			return err
		}
		token := ctx.QueryParam("token")
		data.TokenExisted = token != ""
		data.TokenMatched = token == data.Girl.Token
		data.PhotosFineTuned = data.Girl.PhotosWithKind(PhotoKindFineTuned)
		data.PhotosRoughTuned = data.Girl.PhotosWithKind(PhotoKindRoughTuned)
		data.PhotosOriginal = data.Girl.PhotosWithKind(PhotoKindOriginal)
		if err := db.Find(&data.Event, data.Girl.EventID).Error; err != nil {
			return err
		}
		return ctx.Render(http.StatusOK, "girl", data)
	})
	e.GET("/admin/sign_in/:admin_token", func(ctx echo.Context) error {
		if ctx.Param("admin_token") == envAdminToken {
			sessionSetAdmin(ctx, true)
		}
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	})
	e.GET("/admin/sign_out", func(ctx echo.Context) error {
		sessionSetAdmin(ctx, false)
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	})
	e.POST("/events", func(ctx echo.Context) error {
		if err := db.Create(&Event{Name: ctx.FormValue("name")}).Error; err != nil {
			return err
		}
		return ctx.Redirect(http.StatusSeeOther, "/")
	}, requireAdmin)
	e.POST("/girls", func(ctx echo.Context) error {
		var err error
		var header *multipart.FileHeader
		if header, err = ctx.FormFile("avatar"); err != nil {
			return err
		}
		var file multipart.File
		if file, err = header.Open(); err != nil {
			return err
		}
		defer file.Close()
		var eventID int
		if eventID, err = strconv.Atoi(ctx.FormValue("event_id")); err != nil {
			return err
		}
		var relPath string
		if relPath, err = ossUploadFile(bucket, header.Filename, file); err != nil {
			return err
		}
		if err = db.Create(&Girl{
			EventID:    uint(eventID),
			AvatarPath: relPath,
			Token:      fmt.Sprintf("%08d", 1+rand.Intn(99999999)),
		}).Error; err != nil {
			return err
		}
		return ctx.Redirect(http.StatusSeeOther, "/")
	})
}
