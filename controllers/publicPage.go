package controllers

import "github.com/quadrosh/gin-html/repository"

type publicPageEntry struct {
	ID              uint32                `json:"id"`
	Hrurl           string                `json:"hrurl"`
	Title           string                `json:"title"`
	Description     string                `json:"description"`
	Keywords        string                `json:"keywords"`
	ArticleID       uint                  `json:"article_id"`
	H1              string                `json:"h1"`
	PageDescription string                `json:"page_description"`
	Text            string                `json:"text" `
	Status          repository.PageStatus `json:"status" `
}

func (to *publicPageEntry) convert(r *repository.Page) error {
	to.ID = r.ID
	to.Hrurl = r.Hrurl
	to.Title = r.Title
	to.Description = r.Description
	to.Keywords = r.Keywords
	to.H1 = r.H1
	to.PageDescription = r.PageDescription
	to.Text = r.Text
	to.Status = r.Status
	if r.ArticleID != nil {
		to.ArticleID = *r.ArticleID
	}

	return nil
}
