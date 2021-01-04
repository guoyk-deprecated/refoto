package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
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

func setupRoutes(e *echo.Echo, db *gorm.DB) {
	e.GET("/", func(ctx echo.Context) error {
		type Data struct {
			Title   string
			IsAdmin bool
			Events  []Event
		}
		var data Data
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
}
