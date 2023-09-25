package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/quadrosh/gin-html/helpers"
	"github.com/quadrosh/gin-html/repository"
	"gorm.io/gorm"
)

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func (cli *CommandLine) createAdmin(db *gorm.DB, firstName, lastName, email string) error {
	fmt.Printf("Creating user:  %s, %s, %s\n", firstName, lastName, email)

	user := repository.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,

		Access:   1,
		RoleType: repository.UserRoleTypeAdmin,
		Status:   repository.UserStatusActive,

		AuthKey:            helpers.GenerateSecureToken(15),
		PasswordResetToken: helpers.GenerateSecureToken(15),
		PasswordHash:       helpers.GenerateSecureToken(15),
	}

	// fmt.Printf("admin to create  %v\n", prettyPrint(user))

	if err := db.Create(&user).Error; err != nil {
		log.Printf("error during creating user %v\n", err)
		log.Fatalln(err)
		return err
	}

	var users repository.Users
	if err := db.Model(&repository.Users{}).
		Find(&users).Error; err != nil {
		log.Fatalln(err)
		return err
	}

	found := false
	for _, user := range users {
		if user.Email == email {
			log.Printf("successfully created user (admin) ID: %d, email: %s, firstName: %s, lastName: %s \n", user.ID, user.Email, user.FirstName, user.LastName)
			// port := os.Getenv("PORT")
			log.Printf("password reset link: http://localhost:%s/password-reset/%s \n", os.Getenv("PORT"), user.PasswordResetToken)

			found = true
		}
	}

	if !found {
		log.Printf("it seems that user not created \n")

	}

	return nil

}
func (cli *CommandLine) createUser(db *gorm.DB, firstName, lastName, email string) error {
	fmt.Printf("Creating user:  %s, %s, %s\n", firstName, lastName, email)

	user := repository.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,

		Access:   1,
		RoleType: repository.UserRoleTypeUser,
		Status:   repository.UserStatusActive,

		AuthKey:            helpers.GenerateSecureToken(15),
		PasswordResetToken: helpers.GenerateSecureToken(15),
		PasswordHash:       helpers.GenerateSecureToken(15),
	}

	// fmt.Printf("user to create  %v\n", prettyPrint(user))

	if err := db.Create(&user).Error; err != nil {
		log.Printf("error during creating user %v\n", err)
		log.Fatalln(err)
		return err
	}

	var users repository.Users
	if err := db.Model(&repository.Users{}).
		Find(&users).Error; err != nil {
		log.Fatalln(err)
		return err
	}

	found := false
	for _, user := range users {
		if user.Email == email {
			log.Printf("successfully created user (user) ID: %d, email: %s, firstName: %s, lastName: %s \n", user.ID, user.Email, user.FirstName, user.LastName)
			log.Printf("password reset link: http://localhost:%s/password-reset/%s \n", os.Getenv("PORT"), user.PasswordResetToken)

			found = true
		}
	}

	if !found {
		log.Printf("it seems that user not created \n")

	}

	return nil

}
