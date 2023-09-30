package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	Form               AdminEditPageForm

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

// AdminEditPageForm  page entry form for admin
type AdminEditPageForm struct {
	URL             string `form:"url" `
	Type            string `form:"type" binding:"required"`
	ArticleID       string `form:"article_id" `
	Hrurl           string `form:"hrurl" binding:"required"`
	Title           string `form:"title" `
	Description     string `form:"description" `
	Keywords        string `form:"keywords" `
	H1              string `form:"h1" `
	PageDescription string `form:"page_description" `
	Text            string `form:"text" `
	Status          string `form:"status" `
	Errors          map[string][]string
}

// AdminPageCreatePage - Админка -> создание страницы
// @Summary Страница создания страницы
// @Description Админка - Страница создания страницы
// @ID AdminPageCreatePage
// @Tags admin page
// @Produce  html
// @Success 200 {object} AdminArtistFormPageResponse "Успех"
// @Router /admin/page/create [GET]
func (ctl *Controller) AdminPageCreatePage(ctx *gin.Context) {

	// form, _ := ctx.App.Session.Get(r.Context(), "form").(forms.Form)
	var form AdminEditPageForm

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
	})
	if err != nil {
		log.Fatal(err)
	}
}
