package auth

import (
	"errors"
	"net/http"

	"github.com/quadrosh/gin-html/repository"
	resources "github.com/quadrosh/gin-html/resources/ru"
	"gorm.io/gorm"
)

func CheckUser(
	r *http.Request,
	role repository.UserRoleType,
	db *gorm.DB,
	apiSecret string,
) (*repository.User, error) {
	var (
		user = &repository.User{}
		err  error
	)

	user, err = GetUserByToken(r, db, apiSecret)

	if err != nil {
		return nil, err
	}

	if user.RoleType == repository.UserRoleTypeNone ||
		user.BlockedTime != nil ||
		user.FiredDate != nil {
		return nil, err
	}

	switch role {
	case repository.UserRoleTypeAdmin:
		if user.RoleType != repository.UserRoleTypeAdmin {
			return nil, errors.New(resources.Forbidden())
		}
	case repository.UserRoleTypeUser:
		if user.RoleType != repository.UserRoleTypeUser &&
			user.RoleType != repository.UserRoleTypeAdmin {
			return nil, errors.New(resources.Forbidden())
		}
	default:
		return nil, errors.New(resources.Forbidden())
	}

	return user, nil
}

// GetUserByToken - получить пользователя по токену
func GetUserByToken(
	r *http.Request,
	db *gorm.DB,
	apiSecret string,
) (*repository.User, error) {
	id, err := ExtractTokenID(r, apiSecret)
	if err != nil {
		return nil, err
	}
	var user = &repository.User{}
	err = user.GetByID(db, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
