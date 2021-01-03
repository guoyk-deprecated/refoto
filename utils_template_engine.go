package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type TemplateEngineOptions struct {
	Dir   string
	Ext   string
	Debug bool
}

type TemplateEngine struct {
	opts TemplateEngineOptions
	t    *template.Template
}

func (t *TemplateEngine) ErrorHandler(name string) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := err.Error()

		if he, ok := err.(*echo.HTTPError); ok {
		again:
			if he.Internal != nil {
				if he2, ok := he.Internal.(*echo.HTTPError); ok {
					he = he2
					goto again
				}
			}

			code = he.Code
			message = fmt.Sprintf("%v", he.Message)
		}

		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				err = c.NoContent(code)
			} else {
				err = c.Render(code, name, echo.Map{"Message": message})
			}
			if err != nil {
				c.Echo().Logger.Error(err)
			}
		}
	}
}

func (t *TemplateEngine) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if t.opts.Debug {
		if err := t.Reload(); err != nil {
			return err
		}
	}
	return t.t.ExecuteTemplate(w, name, data)
}

func (t *TemplateEngine) Reload() (err error) {
	newT := template.New("__root__")
	if err = filepath.Walk(t.opts.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != t.opts.Ext {
			return nil
		}
		var rel string
		if rel, err = filepath.Rel(t.opts.Dir, path); err != nil {
			return err
		}
		name := strings.TrimSuffix(filepath.ToSlash(rel), t.opts.Ext)
		var buf []byte
		if buf, err = ioutil.ReadFile(path); err != nil {
			return err
		}
		if newT, err = newT.New(name).Parse(string(buf)); err != nil {
			return err
		}
		log.Println("template loaded:", name)
		return nil
	}); err != nil {
		return
	}
	t.t = newT
	return
}

func NewTemplateEngine(opts TemplateEngineOptions) (te *TemplateEngine, err error) {
	if opts.Dir == "" {
		opts.Dir = "views"
	}
	if opts.Ext == "" {
		opts.Ext = ".gohtml"
	}
	if !strings.HasPrefix(opts.Ext, ".") {
		opts.Ext = "." + opts.Ext
	}
	te = &TemplateEngine{
		opts: opts,
	}
	if err = te.Reload(); err != nil {
		return
	}
	return
}
