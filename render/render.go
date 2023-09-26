package render

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/quadrosh/gin-html/config"
)

// iCSRF interface handling csrf
type iCSRF interface {
	CSRF(ctx *gin.Context)
}

// CreateTemplateCache creates template cache
func CreateTemplateCache(path string, funcs template.FuncMap) (map[string]*template.Template, error) {
	tMap := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/page/*.tmpl", path))
	if err != nil {
		log.Println("error: ", err)
		return tMap, err
	}

	for _, page := range pages {

		name := filepath.Base(page)
		tmpl, err := template.Must(template.New(name).Funcs(funcs), err).ParseFiles(page)
		if err != nil {
			return tMap, err
		}

		layoutMatches, err := filepath.Glob(fmt.Sprintf("%s/layout/*.tmpl", path))
		if err != nil {
			return tMap, err
		}

		if len(layoutMatches) > 0 {
			tmpl, err = tmpl.ParseGlob(fmt.Sprintf("%s/layout/*.tmpl", path))
			if err != nil {
				return tMap, err
			}
		}

		blockMatches, err := filepath.Glob(fmt.Sprintf("%s/block/*.tmpl", path))
		if err != nil {
			return tMap, err
		}

		if len(blockMatches) > 0 {
			tmpl, err = tmpl.ParseGlob(fmt.Sprintf("%s/block/*.tmpl", path))
			if err != nil {
				return tMap, err
			}
		}

		tMap[name] = tmpl
	}

	return tMap, nil
}

// MainTemplate render main template
func MainTemplate(
	app *config.AppConfig,
	engine *gin.Engine,
	ctx *gin.Context,
	tmplName string,
	obj iCSRF) error {
	var pathToMainTemplates = filepath.Join(app.CWD, "templates/main")

	var functions = template.FuncMap{}

	var tCache map[string]*template.Template
	if app.UseCache && app.MainTemplateCache != nil {
		tCache = app.MainTemplateCache
	} else {
		tCache, _ = CreateTemplateCache(pathToMainTemplates, functions)
		if app.MainTemplateCache == nil {
			app.MainTemplateCache = tCache
		}
	}

	t, ok := tCache[tmplName]
	if !ok {
		return errors.New("can't get template from cache")
	}

	engine.SetHTMLTemplate(t)

	obj.CSRF(ctx)

	ctx.HTML(http.StatusOK, tmplName, obj)

	return nil
}

// AdminTemplate render admin template
func AdminTemplate(
	app *config.AppConfig,
	engine *gin.Engine,
	ctx *gin.Context,
	tmplName string,
	obj iCSRF) error {
	var pathToMainTemplates = filepath.Join(app.CWD, "templates/admin")

	var functions = template.FuncMap{}

	var tCache map[string]*template.Template
	if app.UseCache && app.MainTemplateCache != nil {
		tCache = app.MainTemplateCache
	} else {
		tCache, _ = CreateTemplateCache(pathToMainTemplates, functions)
		if app.MainTemplateCache == nil {
			app.MainTemplateCache = tCache
		}
	}

	t, ok := tCache[tmplName]
	if !ok {
		return errors.New("can't get template from cache")
	}

	engine.SetHTMLTemplate(t)

	obj.CSRF(ctx)

	ctx.HTML(http.StatusOK, tmplName, obj)

	return nil
}
