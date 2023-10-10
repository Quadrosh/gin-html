package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/quadrosh/gin-html/forms"
	"github.com/quadrosh/gin-html/internal/constants"
	"github.com/quadrosh/gin-html/render"
	"github.com/quadrosh/gin-html/repository"
	resources "github.com/quadrosh/gin-html/resources/ru"
	"github.com/quadrosh/gin-html/responses"
)

// AdminArticlesResponse response for adminn pages index page
type AdminArticlesResponse struct {
	responses.OkResponse
	Entries               []repository.Article
	ArticleTypeConstMap   map[string]repository.ArticleType
	ArticleStatusConstMap map[string]repository.ArticleStatus
	ArticleLayoutConstMap map[string]repository.ArticleLayout
	Pagination            repository.Pagination

	CSRFResponse
	PageMeta
}

type PageMeta struct {
	Title       string
	Description string
}

// AdminPageFormPageResponse response for admin edit page entry page
type AdminPageFormPageResponse struct {
	responses.OkResponse
	ArticleTypeConstMap   map[string]repository.ArticleType
	ArticleStatusConstMap map[string]repository.ArticleStatus
	ArticleLayoutConstMap map[string]repository.ArticleLayout
	Model                 adminArticleEntry
	Form                  AdminArticleForm

	CSRFResponse
	PageMeta
}

type SuccessJsonResponse struct {
	responses.OkResponse
	Redirect string `json:"redirect"`

	CSRFResponse
}

// AdminPageViewPageResponse ответ страницы страницы сайта в админке
type AdminPageViewPageResponse struct {
	responses.OkResponse
	Model                 adminArticleEntry
	ArticleTypeConstMap   map[string]repository.ArticleType
	ArticleStatusConstMap map[string]repository.ArticleStatus
	ArticleLayoutConstMap map[string]repository.ArticleLayout
	Pagination            repository.Pagination

	CSRFResponse
	PageMeta
}

type adminArticleEntry struct {
	ID              uint32                   `json:"id"`
	Type            repository.ArticleType   `json:"type"`
	Hrurl           string                   `json:"hrurl"`
	Title           string                   `json:"title"`
	Description     string                   `json:"description"`
	Keywords        string                   `json:"keywords"`
	H1              string                   `json:"h1"`
	PageDescription string                   `json:"page_description"`
	Text            string                   `json:"text" `
	Layout          repository.ArticleLayout `json:"layout"`
	Status          repository.ArticleStatus `json:"status" `
	CreatedAt       time.Time                `json:"created_at" format:"date-time"`
	UpdatedAt       time.Time                `json:"updated_at" format:"date-time"`
	DeletedAt       *time.Time               `json:"deleted_at" format:"date-time" `
}

// AdminArticleForm  page entry form for admin
type AdminArticleForm struct {
	URL             string                   `form:"url" `
	Type            repository.ArticleType   `form:"type" binding:"required"`
	Hrurl           string                   `form:"hrurl" binding:"required"`
	Title           string                   `form:"title"  binding:"required,max_length=120" `
	Description     string                   `form:"description" binding:"max_length=250" `
	Keywords        string                   `form:"keywords" `
	H1              string                   `form:"h1" `
	PageDescription string                   `form:"page_description" `
	Text            string                   `form:"text" `
	Status          repository.ArticleStatus `form:"status"  binding:"required"`
	Layout          repository.ArticleLayout `form:"layout" `

	Errors forms.Errors
}

func (to *adminArticleEntry) convert(r *repository.Article) error {
	to.ID = r.ID
	to.Type = r.Type
	to.Hrurl = r.Hrurl
	to.Title = r.Title
	to.Description = r.Description
	to.Keywords = r.Keywords
	to.H1 = r.H1
	to.PageDescription = r.PageDescription
	to.Text = r.Text
	to.Status = r.Status
	to.Layout = r.Layout
	to.CreatedAt = r.CreatedAt
	to.UpdatedAt = r.UpdatedAt
	to.DeletedAt = r.DeletedAt

	return nil
}

func (to *adminArticleEntry) convertForm(r *repository.Article, f *AdminArticleForm) error {

	if f != nil {
		to.Type = f.Type
		to.Layout = f.Layout
		to.Hrurl = f.Hrurl
		to.Title = f.Title
		to.Description = f.Description
		to.Keywords = f.Keywords
		to.H1 = f.H1
		to.PageDescription = f.PageDescription
		to.Text = f.Text
		to.Status = f.Status
	} else if r != nil {
		to.ID = r.ID
		to.Type = r.Type
		to.Layout = r.Layout
		to.Hrurl = r.Hrurl
		to.Title = r.Title
		to.Description = r.Description
		to.Keywords = r.Keywords
		to.H1 = r.H1
		to.PageDescription = r.PageDescription
		to.Text = r.Text
		to.Status = r.Status

	}

	return nil
}

// AdminArticleIndexPage - Администратор -> страницы
// @Summary Список страниц
// @Description Админка - страницы
// @ID AdminArticleIndexPageм
// @Tags admin pages
// @Produce  html
// @Success 200 {object} AdminArticlesResponse "Успех"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /admin/articles [GET]
func (ctl *Controller) AdminArticleIndexPage(ctx *gin.Context) {

	pagination := ctl.GetPagination(ctx.Request)
	db := ctl.Db

	articles := repository.Articles{}

	total, err := articles.GetAllPaged(db, pagination)
	if err != nil {
		ctl.ErrorPage(ctx, http.StatusBadRequest, err)
		return
	}
	pagination.SetTotal(total)

	err = render.AdminTemplate(ctl.App, ctl.Engine, ctx, "admin-article-index.page.tmpl", &AdminArticlesResponse{
		OkResponse: responses.OkResponse{
			Success: true,
			Error:   ctl.GetSessionString(ctx, constants.SessionKeyError, true),
			Info:    ctl.GetSessionString(ctx, constants.SessionKeyInfo, true),
		},
		Entries:               articles,
		ArticleTypeConstMap:   repository.ArticleTypeConstMap,
		ArticleStatusConstMap: repository.ArticleStatusConstMap,
		ArticleLayoutConstMap: repository.ArticleLayoutConstMap,
		Pagination:            pagination,
		PageMeta: PageMeta{
			Title:       "Articles",
			Description: "Article list",
		},
	})
	if err != nil {
		log.Panic(err)
	}
}

// AdminArticleCreatePage - create Page page
// @Summary Page for create Page model
// @Description Page for create Page model by admin
// @ID AdminArticleCreatePage
// @Tags admin page
// @Produce  html
// @Success 200 {object} AdminArtistFormPageResponse "Успех"
// @Router /admin/article/create [GET]
func (ctl *Controller) AdminArticleCreatePage(ctx *gin.Context) {
	var (
		pageUrl = "/admin/article/create"
		form    AdminArticleForm
		entry   = adminArticleEntry{}
		session = sessions.Default(ctx)
	)

	sessionForm, ok := session.Get(constants.SessionKeyForm).(AdminArticleForm)
	if ok && &sessionForm != nil && sessionForm.URL == pageUrl { // form exists in session
		ctl.DeleteFromSession(ctx, constants.SessionKeyForm)
		form = sessionForm // Нужен флаг в сессии что бы грузить оттуда сейчас, который стирается после загрузки
	} else {
		if err := ctx.ShouldBind(&form); err != nil {
			form.Errors = ctl.FormErrors(&form, err.(validator.ValidationErrors))
			// TODO отделить пользовательские ошибки от системных
		}
		form.URL = pageUrl
	}

	entry.convertForm(nil, &form)

	err := render.AdminTemplate(ctl.App, ctl.Engine, ctx, "admin-article-create.page.tmpl", &AdminPageFormPageResponse{
		OkResponse: responses.OkResponse{
			Success: true,
			Error:   ctl.GetSessionString(ctx, constants.SessionKeyError, true),
			Info:    ctl.GetSessionString(ctx, constants.SessionKeyInfo, true),
		},
		ArticleTypeConstMap:   repository.ArticleTypeConstMap,
		ArticleStatusConstMap: repository.ArticleStatusConstMap,
		ArticleLayoutConstMap: repository.ArticleLayoutConstMap,
		Model:                 entry,
		Form:                  form,
		PageMeta: PageMeta{
			Title: "Edit page form",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

// AdminArticleCreatePost - create Page post request
// @Summary Create Page model post
// @Description Create Page model by admin, post request
// @ID AdminArticleCreatePost
// @Tags admin page post
// @Produce  json
// @Success 200 {object}  "Success"
// @Router /admin/article/create [POST]
func (ctl *Controller) AdminArticleCreatePost(ctx *gin.Context) {
	var pageURL = "/admin/article/create"

	var form AdminArticleForm
	if err := ctx.ShouldBindWith(&form, binding.Form); err != nil {
		form.Errors = ctl.FormErrors(&form, err.(validator.ValidationErrors))
		form.URL = pageURL
		if err := ctl.SetToSession(ctx, constants.SessionKeyForm, form); err != nil {
			ctl.ErrorJSON(ctx, err, true)
			return
		}
		ctx.Redirect(http.StatusFound, pageURL) // validation errors show
		return
	}

	var (
		db      = ctl.Db
		article = repository.Article{}
	)

	article.Hrurl = form.Hrurl
	article.Title = form.Title
	article.Description = form.Description
	article.Keywords = form.Keywords
	article.H1 = form.H1
	article.PageDescription = form.PageDescription
	article.Text = form.Text
	article.Layout = form.Layout
	article.Status = form.Status
	article.Type = form.Type

	if err := article.Save(db); err != nil {
		if err := ctl.SetToSession(ctx, constants.SessionKeyError, err.Error()); err != nil {
			ctl.ErrorJSON(ctx, err, true)
			return
		}
		form.URL = pageURL
		if err := ctl.SetToSession(ctx, constants.SessionKeyForm, form); err != nil {
			ctl.ErrorJSON(ctx, err, true)
			return
		}
		ctx.Redirect(http.StatusFound, pageURL)
		return
	}
	// success
	ctl.DeleteFromSession(ctx, constants.SessionKeyForm)
	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/article/%s", strconv.FormatUint(uint64(article.ID), 10)))
	return
}

// AdminArticleViewPage - Admin page view page
// @Summary page view page
// @Description Admin page view page
// @ID AdminArticleViewPage
// @Tags admin page
// @Produce  html
// @Success 200 {object} AdminPageViewPageResponse "Success"
// @Failure 200 json {object} OkResponse
// @Router /admin/article/:id [GET]
func (ctl *Controller) AdminArticleViewPage(ctx *gin.Context) {

	var strID = ctx.Param("id")
	if strID == "" {
		ctl.ErrorJSON(ctx, errors.New(resources.InvalidID()), false)
		return
	}
	ID, err := strconv.Atoi(strID)
	if err != nil {
		ctl.ErrorJSON(ctx, errors.New(err.Error()), false)
		return
	}

	var (
		db      = ctl.Db
		article = repository.Article{}
	)

	err = article.GetByID(db, uint32(ID))
	if err != nil {
		ctl.ErrorJSON(ctx, errors.New(err.Error()), false)
		return
	}

	var entry = adminArticleEntry{}
	entry.convert(&article)

	err = render.AdminTemplate(ctl.App, ctl.Engine, ctx, "admin-article-view.page.tmpl", &AdminPageViewPageResponse{

		OkResponse: responses.OkResponse{
			Success: true,
		},

		ArticleTypeConstMap:   repository.ArticleTypeConstMap,
		ArticleStatusConstMap: repository.ArticleStatusConstMap,
		ArticleLayoutConstMap: repository.ArticleLayoutConstMap,
		Model:                 entry,
		PageMeta: PageMeta{
			Title: fmt.Sprintf("Article  #%d ", article.ID),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

// AdminArticleEditPage - edit page entry page be admin
// @Summary edit page entry
// @Description  edit page form page for admin
// @ID AdminArticleEditPage
// @Param id path int false "Page ID"
// @Tags admin page
// @Produce  html
// @Success 200 {object} AdminPageFormPageResponse "Success"
// @Failure 200 json {object} OkResponse
// @Router /admin/article/:id/edit [GET]
func (ctl *Controller) AdminArticleEditPage(ctx *gin.Context) {

	var strID = ctx.Param("id")
	if strID == "" {
		ctl.ErrorJSON(ctx, errors.New(resources.InvalidID()), false)
		return
	}
	ID, err := strconv.Atoi(strID)
	if err != nil {
		ctl.ErrorJSON(ctx, errors.New(err.Error()), false)
		return
	}

	var (
		pageURL = fmt.Sprintf("/admin/article/%s/edit", strID)
		db      = ctl.Db
		article = repository.Article{}
		form    AdminArticleForm
	)

	if err = article.GetByID(db, uint32(ID)); err != nil {
		ctl.ErrorJSON(ctx, errors.New(err.Error()), false)
		return
	}
	var entry = adminArticleEntry{}
	var session = sessions.Default(ctx)
	sessionForm, ok := session.Get(constants.SessionKeyForm).(AdminArticleForm)
	if ok && &sessionForm != nil && sessionForm.URL == pageURL { // form exists in session
		// session.Delete(constants.SessionKeyForm)
		// session.Save()
		ctl.DeleteFromSession(ctx, constants.SessionKeyForm)
		form = sessionForm
		err = entry.convertForm(nil, &form)
	} else {
		err = entry.convertForm(&article, nil)
	}
	if err != nil {
		ctl.ErrorJSON(ctx, fmt.Errorf("Error while convert model: %w", err), false)
		return
	}

	err = render.AdminTemplate(ctl.App, ctl.Engine, ctx, "admin-article-edit.page.tmpl", &AdminPageFormPageResponse{
		OkResponse: responses.OkResponse{
			Success: true,
			Error:   ctl.GetSessionString(ctx, constants.SessionKeyError, true),
			Info:    ctl.GetSessionString(ctx, constants.SessionKeyInfo, true),
		},
		ArticleTypeConstMap:   repository.ArticleTypeConstMap,
		ArticleStatusConstMap: repository.ArticleStatusConstMap,
		ArticleLayoutConstMap: repository.ArticleLayoutConstMap,
		Model:                 entry,
		Form:                  form,
		PageMeta: PageMeta{
			Title: "Edit page form",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

// AdminArticleEditPost - edit Article post request
// @Summary Edit Article model post
// @Description Edit Article model by admin, post request
// @ID AdminArticleEditPost
// @Tags admin Article edit post
// @Produce  json
// @Success 200 redirect "Success"
// @Router /admin/article/:id/edit [POST]
func (ctl *Controller) AdminArticleEditPost(ctx *gin.Context) {

	var strID = ctx.Param("id")
	if strID == "" {
		ctl.ErrorJSON(ctx, errors.New(resources.InvalidID()), false)
		return
	}
	ID, err := strconv.Atoi(strID)
	if err != nil {
		ctl.ErrorJSON(ctx, errors.New(err.Error()), false)
		return
	}

	var (
		pageURL = fmt.Sprintf("/admin/article/%s/edit", strID)
		db      = ctl.Db
		article = repository.Article{}
		form    AdminArticleForm
	)

	if err = article.GetByID(db, uint32(ID)); err != nil {
		ctl.ErrorJSON(ctx, errors.New(err.Error()), false)
		return
	}

	if err := ctx.ShouldBindWith(&form, binding.Form); err != nil {
		form.Errors = ctl.FormErrors(&form, err.(validator.ValidationErrors))
		form.URL = pageURL
		if err := ctl.SetToSession(ctx, constants.SessionKeyForm, form); err != nil {
			ctl.ErrorJSON(ctx, err, true)
			return
		}
		ctx.Redirect(http.StatusFound, pageURL) // validation errors show
		return
	}

	article.Hrurl = form.Hrurl
	article.Title = form.Title
	article.Description = form.Description
	article.Keywords = form.Keywords
	article.H1 = form.H1
	article.PageDescription = form.PageDescription
	article.Text = form.Text
	article.Layout = form.Layout
	article.Status = form.Status
	article.Type = form.Type

	if err := article.Save(db); err != nil {
		if err := ctl.SetToSession(ctx, constants.SessionKeyError, err.Error()); err != nil {
			ctl.ErrorJSON(ctx, err, true)
			return
		}
		form.URL = pageURL
		if err := ctl.SetToSession(ctx, constants.SessionKeyForm, form); err != nil {
			ctl.ErrorJSON(ctx, err, true)
			return
		}
		ctx.Redirect(http.StatusFound, pageURL)
		return
	}

	// success
	ctl.DeleteFromSession(ctx, constants.SessionKeyForm)
	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/article/%s", strconv.FormatUint(uint64(article.ID), 10)))
	return
}

// AdminArticleDelete - delete Article request
// @Summary Delete Article model
// @Description Delete Article model by admin, request
// @ID AdminArticleDelete
// @Tags admin Article delete
// @Produce  json
// @Success 200 redirect "/admin/articles"  "Success"
// @Router /admin/article/:id/delete [GET]
func (ctl *Controller) AdminArticleDelete(ctx *gin.Context) {

	var strID = ctx.Param("id")
	if strID == "" {
		ctl.ErrorJSON(ctx, errors.New(resources.InvalidID()), false)
		return
	}

	ID, err := strconv.Atoi(strID)
	if err != nil {
		ctl.ErrorJSON(ctx, errors.New(err.Error()), false)
		return
	}

	var (
		pageURL = fmt.Sprintf("/admin/article/%s", strID)
		db      = ctl.Db
		article = repository.Article{}
	)

	if err = article.GetByID(db, uint32(ID)); err != nil {
		ctl.ErrorJSON(ctx, errors.New(err.Error()), false)
		return
	}

	if err := article.Delete(db, false); err != nil {
		if err := ctl.SetToSession(ctx, constants.SessionKeyError, err.Error()); err != nil {
			ctl.ErrorJSON(ctx, err, true)
			return
		}
		ctx.Redirect(http.StatusFound, pageURL)
		return
	}

	if err := ctl.SetToSession(ctx, constants.SessionKeyInfo, resources.DeleteSuccessful()); err != nil {
		ctl.ErrorJSON(ctx, err, true)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/admin/articles")
	return
}
