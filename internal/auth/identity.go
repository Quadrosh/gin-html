package auth

import "github.com/quadrosh/gin-html/repository"

// Identity holds identification data
type Identity struct {
	User         *repository.User
	IsAuthorized bool
}
