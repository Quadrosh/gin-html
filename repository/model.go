package repository

import "time"

// Model - базовая модель для хранения в базе данных
type Model struct {
	ID        uint32     `gorm:"primary_key;auto_increment" json:"id" goqu:"pk"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at" format:"date-time" db:"created_at"` // Здесь тег db нужен только для генератора goqu
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at" format:"date-time" db:"updated_at"` // Здесь тег db нужен только для генератора goqu
	DeletedAt *time.Time `sql:"index" gorm:"null" json:"deleted_at" format:"date-time" db:"deleted_at"`          // Здесь тег db нужен только для генератора goqu
}

// ActionModel - базовая модель
type ActionModel struct {
	CreatedByID uint32 `json:"-" db:"created_by_id"`                                                                                                                               // Здесь тег db нужен только для генератора goqu
	CreatedBy   *User  `json:"created_by" gorm:"foreignkey:CreatedByID;PRELOAD:true;association_save_reference:false;save_associations:false;association_autoupdate:false" db:"-"` // Здесь тег db нужен только для генератора goqu

	UpdatedByID uint32 `json:"-" db:"updated_by_id"`                                                                                                                               // Здесь тег db нужен только для генератора goqu
	UpdatedBy   *User  `json:"updated_by" gorm:"foreignkey:UpdatedByID;PRELOAD:true;association_save_reference:false;save_associations:false;association_autoupdate:false" db:"-"` // Здесь тег db нужен только для генератора goqu

}

// AllModels - automigration models
var AllModels = []interface{}{
	&User{},
	&Page{},
}
