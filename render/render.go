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
	CSRF(ctx *gin.Context) // adds token to response
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

// PublicTemplate render main template
func PublicTemplate(
	app *config.AppConfig,
	engine *gin.Engine,
	ctx *gin.Context,
	tmplName string,
	obj iCSRF) error {
	var templatesPath = filepath.Join(app.CWD, "templates/public")

	var functions = template.FuncMap{}

	var tCache map[string]*template.Template
	if app.UseCache && app.PublicTemplateCache != nil {
		tCache = app.PublicTemplateCache
	} else {
		tCache, _ = CreateTemplateCache(templatesPath, functions)
		if app.PublicTemplateCache == nil {
			app.PublicTemplateCache = tCache
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
	var (
		templatesPath = filepath.Join(app.CWD, "templates/admin")
		tCache        map[string]*template.Template
		err           error
		functions     = template.FuncMap{}
	)

	if app.UseCache && app.PublicTemplateCache != nil {
		tCache = app.PublicTemplateCache
	} else {
		tCache, err = CreateTemplateCache(templatesPath, functions)
		if err != nil {
			return fmt.Errorf("CreateTemplateCache error %s", err)
		}
		if app.PublicTemplateCache == nil {
			app.PublicTemplateCache = tCache
		}
	}
	t, ok := tCache[tmplName]
	if !ok {
		return errors.New(fmt.Sprintf(`can't find template "%s" in tCache map`, tmplName))
	}
	engine.SetHTMLTemplate(t)
	obj.CSRF(ctx)
	ctx.HTML(http.StatusOK, tmplName, obj)
	return nil
}

// UserTemplate render user template
func UserTemplate(
	app *config.AppConfig,
	engine *gin.Engine,
	ctx *gin.Context,
	tmplName string,
	obj iCSRF) error {
	var templatesPath = filepath.Join(app.CWD, "templates/user")

	var functions = template.FuncMap{}

	var tCache map[string]*template.Template
	if app.UseCache && app.UserTemplateCache != nil {
		tCache = app.UserTemplateCache
	} else {
		tCache, _ = CreateTemplateCache(templatesPath, functions)
		if app.UserTemplateCache == nil {
			app.UserTemplateCache = tCache
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
