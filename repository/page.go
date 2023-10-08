package repository

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
)

// PageStatus статус страницы
type PageStatus uint

var PageMutex = &sync.Mutex{}

const (
	// PageStatusDraft черновик
	PageStatusDraft PageStatus = iota + 1
	// PageStatusPublished опубликовано
	PageStatusPublished
)

// PageStatusConstMap  мапа констант PageStatus
var PageStatusConstMap = map[string]PageStatus{
	"PageStatusDraft":     PageStatusDraft,
	"PageStatusPublished": PageStatusPublished,
}

// PageType тип страницы
type PageType uint

const (
	// PageTypeCommon страница сайта
	PageTypeCommon PageType = iota + 1
	// PageTypeArticle страница статьи
	PageTypeArticle
)

// PageTypeConstMap  мапа констант PageType
var PageTypeConstMap = map[string]PageType{
	"PageTypeCommon":  PageTypeCommon,
	"PageTypeArticle": PageTypeArticle,
}

// Page страница
type Page struct {
	Model
	Type            PageType   `gorm:"null" json:"type" db:"type"`
	Hrurl           string     `gorm:"size:255;null;" json:"hrurl" db:"hrurl"`
	Title           string     `gorm:"size:255;null;" json:"title" db:"title"`
	Description     string     `gorm:"size:255;null;" json:"description" db:"description"`
	Keywords        string     `gorm:"size:1000;null;" json:"keywords" db:"keywords"`
	H1              string     `gorm:"size:255;null;" json:"h1" db:"h1"`
	PageDescription string     `gorm:"size:1000;null;" json:"page_description" db:"page_description"`
	Text            string     `gorm:"null;" json:"text" db:"text"`
	Status          PageStatus `gorm:"null" json:"status" db:"status"`
	ArticleID       *uint      `gorm:"null;" json:"article_id" db:"article_id"`
}

// Pages are pages
type Pages []Page

// GetAll - get all pages. Fill 'as' variable
func (ps *Pages) GetAll(db *gorm.DB) error {
	// err := u.preload(db).
	err := db.
		Order("id DESC").
		Find(ps).Error
	if err != nil {
		return err
	}
	return nil
}

// GetAllPaged - get all pages with pagination
func (ps *Pages) GetAllPaged(db *gorm.DB, page Pagination) (int, error) {
	_db := db.Order("id DESC")
	dbPaged := page.OffsetAndLimit(_db)
	err := dbPaged.Find(ps).Error
	if err != nil {
		return 0, err
	}

	var count int64
	if err := db.Model(&Page{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// GetByID - get page by ID. Fill 'p' variable
func (p *Page) GetByID(db *gorm.DB, id uint32) error {
	// err := u.preload(db).
	err := db.
		// Preload("Article").
		Where("id = ?", id).
		Find(p).Error
	if err != nil {
		return err
	}
	return nil
}

// ByURL - get page by URL. Fill 'p' variable
func (p *Page) ByURL(db *gorm.DB, URL string) error {
	err := db.
		// Preload("Article").
		Where("hrurl = ?", URL).
		Find(p).Error
	if err != nil {
		return err
	}
	return nil
}

// Save - сохранить
func (p *Page) Save(db *gorm.DB) error {
	if err := db.Save(p).Error; err != nil {
		return err
	}

	return p.GetByID(db, p.ID)
}

// Delete - Удалить PlaylistItem
func (p *Page) Delete(db *gorm.DB, noTransaction bool) error {
	PageMutex.Lock()
	defer PageMutex.Unlock()

	var err error
	var action = func(tx *gorm.DB) error {
		if tErr := tx.Delete(p).Error; tErr != nil {
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
