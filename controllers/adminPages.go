package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/quadrosh/gin-html/render"
	"github.com/quadrosh/gin-html/repository"
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
	Page               adminPageEntry
	PageTypeConstMap   map[string]repository.PageType
	PageStatusConstMap map[string]repository.PageStatus
	Form               AdminPageForm

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

	// form, _ := ctx.App.Session.Get(r.Context(), "form").(forms.Form)
	var form AdminPageForm

	// if err := ctx.ShouldBind(&form); err != nil {
	// 	ctl.ErrorJSON(ctx, err)
	// 	return
	// }

	form.URL = "/admin/page/create"

	var entry = adminPageEntry{}

	// err := entry.convertForm(&repository.Page{}, &form)

	err := render.AdminTemplate(ctl.App, ctl.Engine, ctx, "admin-page-create.page.tmpl", &AdminPageFormPageResponse{
		OkResponse: responses.OkResponse{
			Success: true,
		},
		PageTypeConstMap:   repository.PageTypeConstMap,
		PageStatusConstMap: repository.PageStatusConstMap,
		Page:               entry,
		Form:               form,
		PageMeta: PageMeta{
			Title: "Edit page form",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

// AdminPageForm  page entry form for admin
type AdminPageForm struct {
	URL             string `form:"url" `
	Type            string `form:"type" binding:"required"`
	ArticleID       string `form:"article_id" `
	Hrurl           string `form:"hrurl" binding:"required"`
	Title           string `form:"title"  binding:"required,max_length=130" `
	Description     string `form:"description" binding:"max_length=250" `
	Keywords        string `form:"keywords" `
	H1              string `form:"h1" `
	PageDescription string `form:"page_description" `
	Text            string `form:"text" `
	Status          string `form:"status"  binding:"required"`
	Errors          map[string][]string
}

// AdminPageCreatePost - create Page post request
// @Summary Create Page model post
// @Description Create Page model by admin, post request
// @ID AdminPageCreatePost
// @Tags admin page post
// @Produce  json
// @Success 200 {object} AdminArtistFormPageResponse "Success"
// @Router /admin/page/create [POST]
func (ctl *Controller) AdminPageCreatePost(ctx *gin.Context) {

	// var pageURL = "/admin/page/create"

	var form AdminPageForm

	if err := ctx.ShouldBindWith(&form, binding.Form); err != nil {
		ctl.ErrorJSON(ctx, err)
		// отобразить страницу с ошибками
		// ctx.Redirect(303, pageURL)

		return
	}

	// if !form.Valid() {
	// 	// отобразить страницу с ошибками
	// 	ctx.App.Session.Put(r.Context(), "form", form)
	// 	http.Redirect(w, r, pageURL, http.StatusSeeOther)
	// 	return
	// }

	// var (
	// 	ctxDatabase, _ = r.Context().Value(middleware.ContextKeyDatabase).(*gorm.DB)
	// 	page           = repository.Page{}
	// )

	// page.Artist = repository.Artist{}

	// artistIDStr := r.Form.Get("artist_id")
	// if artistIDStr != "" && artistIDStr != "0" {
	// 	_artistID, err := strconv.ParseUint(artistIDStr, 10, 32)
	// 	if err != nil {
	// 		ctx.App.Session.Put(r.Context(), "error", err.Error())
	// 		ctx.App.Session.Put(r.Context(), "form", form)
	// 		http.Redirect(w, r, pageURL, http.StatusSeeOther)
	// 		return
	// 	}
	// 	// _uintArtistID := uint(_artistID)
	// 	var artist = repository.Artist{}
	// 	err = artist.GetByID(ctxDatabase, uint32(_artistID))
	// 	if err != nil || artist.ID == 0 {
	// 		ctx.App.Session.Put(r.Context(), "error", resources.ArtistNotFound(uint32(_artistID)))
	// 		ctx.App.Session.Put(r.Context(), "form", form)
	// 		http.Redirect(w, r, pageURL, http.StatusSeeOther)
	// 		return
	// 	}
	// 	page.Artist = artist
	// 	// page.ArtistID = &_uintArtistID
	// } else {
	// 	page.ArtistID = nil
	// }

	// page.Hrurl = r.Form.Get("hrurl")
	// page.Title = r.Form.Get("title")
	// page.Description = r.Form.Get("description")
	// page.Keywords = r.Form.Get("keywords")
	// page.H1 = r.Form.Get("h1")
	// page.PageDescription = r.Form.Get("page_description")
	// page.Text = r.Form.Get("text")

	// statusInt, err := strconv.Atoi(r.Form.Get("status"))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// page.Status = repository.PageStatus(statusInt)

	// typeInt, err := strconv.Atoi(r.Form.Get("type"))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// page.Type = repository.PageType(typeInt)

	// if err := page.Save(ctxDatabase); err != nil {
	// 	ctx.App.Session.Put(r.Context(), "error", err.Error())
	// 	ctx.App.Session.Put(r.Context(), "form", form)
	// 	http.Redirect(w, r, pageURL, http.StatusSeeOther)
	// 	return
	// }

	// ctx.App.Session.Remove(r.Context(), "form")

	// http.Redirect(w, r, fmt.Sprintf("/admin/page/%s", strconv.FormatUint(uint64(page.ID), 10)), http.StatusSeeOther)
	// return

	ctx.Redirect(303, "/admin/page/newPageID")
	return
}
