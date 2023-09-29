package repository

import "gorm.io/gorm"

// Pagination - paging structure
type Pagination struct {
	ItemFirst   int
	ItemLast    int
	Total       int
	PageSize    int
	CurrentPage int
}

// OffsetAndLimit - применить смещения постраничное
func (p Pagination) OffsetAndLimit(db *gorm.DB) *gorm.DB {
	if p.CurrentPage == 0 {
		db = db.Offset(int(p.CurrentPage * p.PageSize)).Limit(int(p.PageSize))
	} else {
		db = db.Offset(int((p.CurrentPage - 1) * p.PageSize)).Limit(int(p.PageSize))
	}
	return db
}

func (p *Pagination) SetTotal(total int) {
	p.Total = total
	if p.CurrentPage > 1 {
		p.ItemFirst = (p.PageSize * (p.CurrentPage - 1)) + 1
	} else {
		p.ItemFirst = 1
	}
	if p.Total-p.ItemFirst > p.PageSize {
		p.ItemLast = p.ItemFirst + p.PageSize - 1
	} else {
		p.ItemLast = p.ItemFirst + (p.Total - p.ItemFirst)
	}
}
