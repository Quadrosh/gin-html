package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quadrosh/gin-html/render"
	"github.com/quadrosh/gin-html/repository"
	"github.com/quadrosh/gin-html/responses"
)

// AdminPagesResponse ответ страницы страниц сайта
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
