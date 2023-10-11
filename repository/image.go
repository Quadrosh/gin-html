package repository

// ImageFile image file
type ImageFile struct {
	Model
	Name string `gorm:"size:510;null;unique;" json:"name" db:"name"`
}
