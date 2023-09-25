package repository

import (
	"sync"
	"time"

	"github.com/quadrosh/gin-html/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Rule - правило проверки доступа
type UserRule uint

var UserMutex = &sync.Mutex{}

const (
	// UserRuleRoleUser - check access User
	UserRuleRoleUser UserRule = iota
	// UserRuleRoleAdmin  - check access Admin
	UserRuleRoleAdmin
)

// UserRoleType - Тип роли
type UserRoleType uint

const (
	// UserRoleTypeNone - no role
	UserRoleTypeNone UserRoleType = iota
	// UserRoleTypeAdmin - admin
	UserRoleTypeAdmin
	// UserRoleTypeUser - user
	UserRoleTypeUser
)

// UserRoleTypeConstMap  мапа констант UserRoleType
var UserRoleTypeConstMap = map[string]UserRoleType{
	"UserRoleTypeNone":  UserRoleTypeNone,
	"UserRoleTypeAdmin": UserRoleTypeAdmin,
	"UserRoleTypeUser":  UserRoleTypeUser,
}

// UserRole - роль пользователя
type UserRole struct {
	Type UserRoleType `gorm:"null;" json:"type"`
}

// CanSignIn - определение, позволяет ли роль авторизоваться в систему
func (urt UserRoleType) CanSignIn() bool {
	switch urt {
	case UserRoleTypeAdmin,
		UserRoleTypeUser:
		return true
	}
	return false
}

// UserStatus статус пользователя
type UserStatus uint

const (
	//UserStatusNew соискатель
	UserStatusNew UserStatus = iota + 1
	//UserStatusActive активный допуск
	UserStatusActive
	//UserStatusRefused доступ отклонен
	UserStatusRefused
	//UserStatusSuspended доступ приостановлен
	UserStatusSuspended
	//UserStatusFired уволен
	UserStatusFired
	// UserStatusValidationWaiting ожидает проверочных мероприятий
	UserStatusValidationWaiting
	// UserStatusValidationFailed Не прошел проверочные мероприятия
	UserStatusValidationFailed
	// UserStatusValidationDone Прошёл проверочные мероприятия
	UserStatusValidationDone
)

// UserStatusConstMap  мапа констант UserStatus
var UserStatusConstMap = map[string]UserStatus{
	"UserStatusNew":               UserStatusNew,
	"UserStatusActive":            UserStatusActive,
	"UserStatusRefused":           UserStatusRefused,
	"UserStatusSuspended":         UserStatusSuspended,
	"UserStatusFired":             UserStatusFired,
	"UserStatusValidationWaiting": UserStatusValidationWaiting,
	"UserStatusValidationFailed":  UserStatusValidationFailed,
	"UserStatusValidationDone":    UserStatusValidationDone,
}

// Users - users collection
type Users []User

// User - site user
type User struct {
	Model
	FirstName          string       `gorm:"size:255;null;" json:"first_name" db:"first_name"`
	LastName           string       `gorm:"size:255;null;" json:"last_name" db:"last_name"`
	MiddleName         string       `gorm:"size:255;null;" json:"patronymic" db:"middle_name"`
	Access             uint         `gorm:"not null;" json:"access"`
	Email              string       `gorm:"size:100;uniqueIndex" json:"email" db:"email"` // TODO not null
	Phone              string       `gorm:"size:100;null;" json:"phone" db:"phone"`
	Telegram           string       `gorm:"size:100;null;" json:"telegram" db:"telegram"`
	Skype              string       `gorm:"size:100;null;" json:"skype" db:"skype"`
	WhatsApp           string       `gorm:"size:100;null;" json:"whatsapp" db:"whatsapp"`
	PasswordHash       string       `gorm:"size:1024;not null;" json:"-" db:"password_hash"`
	PasswordResetToken string       `gorm:"size:1024" json:"-" db:"password_reset_token"`
	AuthKey            string       `gorm:"size:1024" json:"-" db:"auth_key"`
	Address            string       `gorm:"size:11024;null" json:"address"`
	BlockedTime        *time.Time   `json:"blocked_time" format:"date-time" db:"blocked_time"`
	BlockedReason      string       `json:"blocked_reason" db:"blocked_reason"`
	FiredDate          *time.Time   `json:"fired_date" format:"date-time" db:"fired_date"`
	RejectDate         *time.Time   `json:"reject_date" format:"date-time" db:"reject_date"`
	RejectReason       string       `json:"reject_reason" db:"reject_reason"`
	Status             UserStatus   `gorm:"null;" json:"status" db:"status"`
	RoleType           UserRoleType `gorm:"null;" json:"role_type" db:"role_type"`

	DisplayName string  `gorm:"-"  json:"display_name" db:"-"`
	StatusName  string  `gorm:"-" json:"status_name" db:"-"`
	Password    *string `gorm:"-" json:"password" db:"-"`
}

// GetByPasswordResetToken - get user by password_reset_token
func (u *User) GetByPasswordResetToken(db *gorm.DB, token string) error {
	err := db.
		Where("password_reset_token = ?", token).
		Find(u).Error
	if err != nil {
		return err
	}
	return nil
}

/*HashPassword - Hash string 'psw' and set u.PasswordHash*/
func (u *User) HashPassword(psw string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(psw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) verifyPassword(psw string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(psw))
}

/*SignIn user. If success, u* instance will be wrapped*/
func (u *User) SignIn(db *gorm.DB, email, psw string) error {
	email = utils.Normalize(email)
	tmp := User{}
	err := db.Model(tmp).
		Where(&User{Email: email}).
		Take(&tmp).Error
	if err != nil {
		return err
	}

	err = tmp.verifyPassword(psw)
	if err != nil {
		return err
	}

	err = db.Model(tmp).Where("id = ?", tmp.ID).Take(u).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateAuthKey обновляет AuthKey для пользователя
func (u *User) UpdateAuthKey(db *gorm.DB, token string) error {
	err := db.Model(&User{}).
		Where("id = ?", u.ID).
		Update("auth_key", token).Error
	if err != nil {
		return err
	}

	return nil
}

// UserCanSettings - params for check access
type UserCanSettings struct {
	Rule       UserRule
	ObjectType uint
	ObjectID   interface{}
}

// Can - check access
func (u *User) Can(db *gorm.DB, s UserCanSettings) bool {
	var count int64
	switch s.Rule {
	case UserRuleRoleUser:
		db.
			Model(&User{}).
			Where("id = ?", u.ID).
			Where("role_type = ?", UserRoleTypeUser).
			Count(&count)
		return count > 0
	case UserRuleRoleAdmin:
		db.
			Model(&User{}).
			Where("id = ?", u.ID).
			Where("role_type = ?", UserRoleTypeAdmin).
			Count(&count)
		return count > 0
	default:
		return false
	}

}

// GetByID - get user by ID. Fill 'u' variable
func (u *User) GetByID(db *gorm.DB, id uint32) error {
	err := db.
		Where("id = ?", id).
		Find(u).Error
	if err != nil {
		return err
	}
	return nil
}
