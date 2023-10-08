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

// AdminPagesResponse response for adminn pages index page
type AdminPagesResponse struct {
	responses.OkResponse
	Entries            []repository.Page
	PageTypeConstMap   map[string]repository.PageType
	PageStatusConstMap map[string]repository.PageStatus
	Pagination         repository.Pagination

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
	PageTypeConstMap   map[string]repository.PageType
	PageStatusConstMap map[string]repository.PageStatus
	Model              adminPageEntry
	Form               AdminPageForm

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
	Model              adminPageEntry
	PageTypeConstMap   map[string]repository.PageType
	PageStatusConstMap map[string]repository.PageStatus
	Pagination         repository.Pagination

	CSRFResponse
	PageMeta
}

type adminPageEntry struct {
	ID              uint32                `json:"id"`
	Type            repository.PageType   `json:"type"`
	ArticleID       uint                  `json:"article_id"`
	Hrurl           string                `json:"hrurl"`
	Title           string                `json:"title"`
	Description     string                `json:"description"`
	Keywords        string                `json:"keywords"`
	H1              string                `json:"h1"`
	PageDescription string                `json:"page_description"`
	Text            string                `json:"text" `
	Status          repository.PageStatus `json:"status" `
	CreatedAt       time.Time             `json:"created_at" format:"date-time"`
	UpdatedAt       time.Time             `json:"updated_at" format:"date-time"`
	DeletedAt       *time.Time            `json:"deleted_at" format:"date-time" `
}

// AdminPageForm  page entry form for admin
type AdminPageForm struct {
	URL             string                `form:"url" `
	Type            repository.PageType   `form:"type" binding:"required"`
	ArticleID       uint                  `form:"article_id" `
	Hrurl           string                `form:"hrurl" binding:"required"`
	Title           string                `form:"title"  binding:"required,max_length=120" `
	Description     string                `form:"description" binding:"max_length=250" `
	Keywords        string                `form:"keywords" `
	H1              string                `form:"h1" `
	PageDescription string                `form:"page_description" `
	Text            string                `form:"text" `
	Status          repository.PageStatus `form:"status"  binding:"required"`
	Errors          forms.Errors
}

func (to *adminPageEntry) convert(r *repository.Page) error {
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
	to.CreatedAt = r.CreatedAt
	to.UpdatedAt = r.UpdatedAt
	to.DeletedAt = r.DeletedAt

	return nil
}

func (to *adminPageEntry) convertForm(r *repository.Page, f *AdminPageForm) error {

	if f != nil {
		to.Type = f.Type
		to.ArticleID = f.ArticleID
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
		if r.ArticleID != nil {
			to.ArticleID = *r.ArticleID
		}
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

// AdminPageIndexPage - Администратор -> страницы
// @Summary Список страниц
// @Description Админка - страницы
// @ID AdminPageIndexPageм
// @Tags admin pages
// @Produce  html
// @Success 200 {object} AdminPagesResponse "Успех"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /admin/pages [GET]
func (ctl *Controller) AdminPageIndexPage(ctx *gin.Context) {

	pagination := ctl.GetPagination(ctx.Request)
	db := ctl.Db

	pages := repository.Pages{}

	total, err := pages.GetAllPaged(db, pagination)
	if err != nil {
		ctl.ErrorPage(ctx, http.StatusBadRequest, err)
		return
	}
	pagination.SetTotal(total)

	err = render.AdminTemplate(ctl.App, ctl.Engine, ctx, "admin-page-index.page.tmpl", &AdminPagesResponse{
		OkResponse: responses.OkResponse{
			Success: true,
			Error:   ctl.GetSessionString(ctx, constants.SessionKeyError, true),
			Info:    ctl.GetSessionString(ctx, constants.SessionKeyInfo, true),
		},
		Entries:            pages,
		PageTypeConstMap:   repository.PageTypeConstMap,
		PageStatusConstMap: repository.PageStatusConstMap,
		Pagination:         pagination,
		PageMeta: PageMeta{
			Title:       "Индекс страниц сайта",
			Description: "Список страниц сайта",
		},
	})
	if err != nil {
		log.Panic(err)
	}
}

// AdminPageCreatePage - create Page page
// @Summary Page for create Page model
// @Description Page for create Page model by admin
// @ID AdminPageCreatePage
// @Tags admin page
// @Produce  html
// @Success 200 {object} AdminArtistFormPageResponse "Успех"
// @Router /admin/page/create [GET]
func (ctl *Controller) AdminPageCreatePage(ctx *gin.Context) {
	var (
		pageUrl = "/admin/page/create"
		form    AdminPageForm
		entry   = adminPageEntry{}
		session = sessions.Default(ctx)
	)

	sessionForm, ok := session.Get(constants.SessionKeyForm).(AdminPageForm)
	if ok && &sessionForm != nil && sessionForm.URL == pageUrl { // form exists in session
		session.Delete(constants.SessionKeyForm)
		form = sessionForm // Нужен флаг в сессии что бы грузить оттуда сейчас, который стирается после загрузки
	} else {
		if err := ctx.ShouldBind(&form); err != nil {
			form.Errors = ctl.FormErrors(&form, err.(validator.ValidationErrors))
			// TODO отделить пользовательские ошибки от системных
		}
		form.URL = pageUrl
	}

	entry.convertForm(nil, &form)

	err := render.AdminTemplate(ctl.App, ctl.Engine, ctx, "admin-page-create.page.tmpl", &AdminPageFormPageResponse{
		OkResponse: responses.OkResponse{
			Success: true,
			Error:   ctl.GetSessionString(ctx, constants.SessionKeyError, true),
			Info:    ctl.GetSessionString(ctx, constants.SessionKeyInfo, true),
		},
		PageTypeConstMap:   repository.PageTypeConstMap,
		PageStatusConstMap: repository.PageStatusConstMap,
		Model:              entry,
		Form:               form,
		PageMeta: PageMeta{
			Title: "Edit page form",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

// AdminPageCreatePost - create Page post request
// @Summary Create Page model post
// @Description Create Page model by admin, post request
// @ID AdminPageCreatePost
// @Tags admin page post
// @Produce  json
// @Success 200 {object}  "Success"
// @Router /admin/page/create [POST]
func (ctl *Controller) AdminPageCreatePost(ctx *gin.Context) {
	var pageURL = "/admin/page/create"
	var session = sessions.Default(ctx)

	var form AdminPageForm
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
		db   = ctl.Db
		page = repository.Page{}
	)

	page.Hrurl = form.Hrurl
	page.Title = form.Title
	page.Description = form.Description
	page.Keywords = form.Keywords
	page.H1 = form.H1
	page.PageDescription = form.PageDescription
	page.Text = form.Text
	page.Status = form.Status
	page.Type = form.Type

	if err := page.Save(db); err != nil {
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
	session.Delete(constants.SessionKeyForm) // redirect to page view page
	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/page/%s", strconv.FormatUint(uint64(page.ID), 10)))
	return
}

// AdminPageViewPage - Admin page view page
// @Summary page view page
// @Description Admin page view page
// @ID AdminPageViewPage
// @Tags admin page
// @Produce  html
// @Success 200 {object} AdminPageViewPageResponse "Success"
// @Failure 200 json {object} OkResponse
// @Router /admin/page/:id [GET]
func (ctl *Controller) AdminPageViewPage(ctx *gin.Context) {

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
		db   = ctl.Db
		page = repository.Page{}
	)

	err = page.GetByID(db, uint32(ID))
	if err != nil {
		ctl.ErrorJSON(ctx, errors.New(err.Error()), false)
		return
	}

	var entry = adminPageEntry{}
	entry.convert(&page)

	err = render.AdminTemplate(ctl.App, ctl.Engine, ctx, "admin-page-view.page.tmpl", &AdminPageViewPageResponse{

		OkResponse: responses.OkResponse{
			Success: true,
		},

		PageTypeConstMap:   repository.PageTypeConstMap,
		PageStatusConstMap: repository.PageStatusConstMap,
		Model:              entry,
		PageMeta: PageMeta{
			Title: fmt.Sprintf("Page  #%d ", page.ID),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

// AdminPageEditPage - edit page entry page be admin
// @Summary edit page entry
// @Description  edit page form page for admin
// @ID AdminPageEditPage
// @Param id path int false "Page ID"
// @Tags admin page
// @Produce  html
// @Success 200 {object} AdminPageFormPageResponse "Success"
// @Failure 200 json {object} OkResponse
// @Router /admin/page/:id/edit [GET]
func (ctl *Controller) AdminPageEditPage(ctx *gin.Context) {

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
		pageURL = fmt.Sprintf("/admin/page/%s/edit", strID)
		db      = ctl.Db
		page    = repository.Page{}
		form    AdminPageForm
	)

	if err = page.GetByID(db, uint32(ID)); err != nil {
		ctl.ErrorJSON(ctx, errors.New(err.Error()), false)
		return
	}
	var entry = adminPageEntry{}
	var session = sessions.Default(ctx)
	sessionForm, ok := session.Get(constants.SessionKeyForm).(AdminPageForm)
	if ok && &sessionForm != nil && sessionForm.URL == pageURL { // form exists in session
		session.Delete(constants.SessionKeyForm)
		form = sessionForm
		err = entry.convertForm(nil, &form)
	} else {
		err = entry.convertForm(&page, nil)
	}
	if err != nil {
		ctl.ErrorJSON(ctx, fmt.Errorf("Error while convert model: %w", err), false)
		return
	}

	err = render.AdminTemplate(ctl.App, ctl.Engine, ctx, "admin-page-edit.page.tmpl", &AdminPageFormPageResponse{
		OkResponse: responses.OkResponse{
			Success: true,
			Error:   ctl.GetSessionString(ctx, constants.SessionKeyError, true),
			Info:    ctl.GetSessionString(ctx, constants.SessionKeyInfo, true),
		},
		PageTypeConstMap:   repository.PageTypeConstMap,
		PageStatusConstMap: repository.PageStatusConstMap,
		Model:              entry,
		Form:               form,
		PageMeta: PageMeta{
			Title: "Edit page form",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

// AdminPageEditPost - edit Page post request
// @Summary Edit Page model post
// @Description Edit Page model by admin, post request
// @ID AdminPageEditPost
// @Tags admin page edit post
// @Produce  json
// @Success 200 redirect "Success"
// @Router /admin/page/:id/edit [POST]
func (ctl *Controller) AdminPageEditPost(ctx *gin.Context) {

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
		pageURL = fmt.Sprintf("/admin/page/%s/edit", strID)
		db      = ctl.Db
		page    = repository.Page{}
		form    AdminPageForm
		session = sessions.Default(ctx)
	)

	if err = page.GetByID(db, uint32(ID)); err != nil {
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

	page.Hrurl = form.Hrurl
	page.Title = form.Title
	page.Description = form.Description
	page.Keywords = form.Keywords
	page.H1 = form.H1
	page.PageDescription = form.PageDescription
	page.Text = form.Text
	page.Status = form.Status
	page.Type = form.Type

	if err := page.Save(db); err != nil {
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
	session.Delete(constants.SessionKeyForm) // redirect to page view page

	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/page/%s", strconv.FormatUint(uint64(page.ID), 10)))
	return
}

// AdminPageDelete - delete Page request
// @Summary Delete Page model
// @Description Delete Page model by admin, request
// @ID AdminPageDelete
// @Tags admin page delete
// @Produce  json
// @Success 200 redirect "/admin/pages"  "Success"
// @Router /admin/page/:id/delete [GET]
func (ctl *Controller) AdminPageDelete(ctx *gin.Context) {

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
		pageURL = fmt.Sprintf("/admin/page/%s", strID)
		db      = ctl.Db
		page    = repository.Page{}
	)

	if err = page.GetByID(db, uint32(ID)); err != nil {
		ctl.ErrorJSON(ctx, errors.New(err.Error()), false)
		return
	}

	if err := page.Delete(db, false); err != nil {
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

	ctx.Redirect(http.StatusSeeOther, "/admin/pages")
	return
}
