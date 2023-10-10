package controllers

import "github.com/quadrosh/gin-html/repository"

type publicArticleEntry struct {
	ID              uint32                   `json:"id"`
	Hrurl           string                   `json:"hrurl"`
	Title           string                   `json:"title"`
	Description     string                   `json:"description"`
	Keywords        string                   `json:"keywords"`
	ArticleID       uint                     `json:"article_id"`
	H1              string                   `json:"h1"`
	PageDescription string                   `json:"page_description"`
	Text            string                   `json:"text" `
	Status          repository.ArticleStatus `json:"status" `
	Layout          repository.ArticleLayout `json:"layout" `
}

func (to *publicArticleEntry) convert(r *repository.Article) error {
	to.ID = r.ID
	to.Hrurl = r.Hrurl
	to.Title = r.Title
	to.Description = r.Description
	to.Keywords = r.Keywords
	to.H1 = r.H1
	to.PageDescription = r.PageDescription
	to.Text = r.Text
	to.Status = r.Status
	to.Layout = r.Layout

	return nil
}
