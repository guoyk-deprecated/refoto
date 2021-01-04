package main

import (
	"errors"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	"net/http"
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

func setupRoutes(e *echo.Echo, db *gorm.DB) {
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
		//AvatarURL: "https://via.placeholder.com/150",
		return ctx.Render(http.StatusOK, "index", data)
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
		//TODO: create girl with avatar
		return ctx.Redirect(http.StatusSeeOther, "/")
	})
}
