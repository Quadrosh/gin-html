package repository

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
)

var ArticleMutex = &sync.Mutex{}

// ArticleStatus статус страницы
type ArticleStatus uint

const (
	// ArticleStatusDraft draft
	ArticleStatusDraft ArticleStatus = iota + 1
	// ArticleStatusPublished published
	ArticleStatusPublished
)

// ArticleStatusConstMap  map of PageStatus constants
var ArticleStatusConstMap = map[string]ArticleStatus{
	"ArticleStatusDraft":     ArticleStatusDraft,
	"ArticleStatusPublished": ArticleStatusPublished,
}

// ArticleStatusNameMap  map of PageStatus Names
var ArticleStatusNameMap = map[ArticleStatus]string{
	ArticleStatusDraft:     "Draft",
	ArticleStatusPublished: "Published",
}

// ArticleLayout article layout
type ArticleLayout uint

const (
	// ArticleLayoutHome home page
	ArticleLayoutHome ArticleLayout = iota + 1
	// ArticleLayoutPage public page
	ArticleLayoutPage
	// ArticleLayoutArticle article page
	ArticleLayoutArticle
)

// ArticleLayoutConstMap  map of ArticleLayout constants
var ArticleLayoutConstMap = map[string]ArticleLayout{
	"ArticleLayoutHome":    ArticleLayoutHome,
	"ArticleLayoutPage":    ArticleLayoutPage,
	"ArticleLayoutArticle": ArticleLayoutArticle,
}

// ArticleLayoutNameMap  map of ArticleLayout names
var ArticleLayoutNameMap = map[ArticleLayout]string{
	ArticleLayoutHome:    "Home",
	ArticleLayoutPage:    "Page",
	ArticleLayoutArticle: "Article",
}

// ArticleType тип страницы
type ArticleType uint

const (
	// ArticleTypePage page
	ArticleTypePage ArticleType = iota + 1
	// ArticleTypeArticle article
	ArticleTypeArticle
	// ArticleTypeNews news
	ArticleTypeNews
)

// ArticleTypeConstMap  map of PageType constants
var ArticleTypeConstMap = map[string]ArticleType{
	"ArticleTypePage":    ArticleTypePage,
	"ArticleTypeArticle": ArticleTypeArticle,
	"ArticleTypeNews":    ArticleTypeNews,
}

// ArticleTypeNameMap map of PageType names
var ArticleTypeNameMap = map[ArticleType]string{
	ArticleTypePage:    "Page",
	ArticleTypeArticle: "Article",
	ArticleTypeNews:    "News",
}

// Article страница
type Article struct {
	Model
	Type            ArticleType   `gorm:"null" json:"type" db:"type"`
	Hrurl           string        `gorm:"size:255;null;" json:"hrurl" db:"hrurl"`
	Title           string        `gorm:"size:255;null;" json:"title" db:"title"`
	Description     string        `gorm:"size:255;null;" json:"description" db:"description"`
	Keywords        string        `gorm:"size:1000;null;" json:"keywords" db:"keywords"`
	H1              string        `gorm:"size:255;null;" json:"h1" db:"h1"`
	PageDescription string        `gorm:"size:1000;null;" json:"page_description" db:"page_description"`
	Text            string        `gorm:"null;" json:"text" db:"text"`
	Status          ArticleStatus `gorm:"null" json:"status" db:"status"`
	Layout          ArticleLayout `gorm:"null;" json:"layout" db:"layout"`
}

// Articles are pages
type Articles []Article

// GetAll - get all articles. Fill 'as' variable
func (as *Articles) GetAll(db *gorm.DB) error {
	// err := u.preload(db).
	err := db.
		Order("id DESC").
		Find(as).Error
	if err != nil {
		return err
	}
	return nil
}

// GetAllPaged - get all pages with pagination
func (as *Articles) GetAllPaged(db *gorm.DB, page Pagination) (int, error) {
	_db := db.Order("id DESC")
	dbPaged := page.OffsetAndLimit(_db)
	err := dbPaged.Find(as).Error
	if err != nil {
		return 0, err
	}

	var count int64
	if err := db.Model(&Article{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// GetByID - get page by ID. Fill 'p' variable
func (a *Article) GetByID(db *gorm.DB, id uint32) error {
	// err := u.preload(db).
	err := db.
		// Preload("Article").
		Where("id = ?", id).
		Find(a).Error
	if err != nil {
		return err
	}
	return nil
}

// ByURL - get page by URL. Fill 'p' variable
func (a *Article) ByURL(db *gorm.DB, URL string) error {
	err := db.
		// Preload("Article").
		Where("hrurl = ?", URL).
		Find(a).Error
	if err != nil {
		return err
	}
	return nil
}

// Save - сохранить
func (a *Article) Save(db *gorm.DB) error {
	if err := db.Save(a).Error; err != nil {
		return err
	}

	return a.GetByID(db, a.ID)
}

// Delete - Удалить PlaylistItem
func (a *Article) Delete(db *gorm.DB, noTransaction bool) error {
	ArticleMutex.Lock()
	defer ArticleMutex.Unlock()

	var err error
	var action = func(tx *gorm.DB) error {
		if tErr := tx.Delete(a).Error; tErr != nil {
			return tErr
		}
		return nil
	}

	if noTransaction {
		err = action(db)
	} else {
		err = db.Transaction(action)
	}

	if err != nil {
		return fmt.Errorf("Delete() failed: %w", err)
	}

	return nil
}
