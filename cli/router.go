package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"gorm.io/gorm"
)

// CommandLine is an interface for manage it
type CommandLine struct {
}

func (cli *CommandLine) printUsage() {
	fmt.Println("something went wrong, command line typical usage:")
	fmt.Println(" createuser -firstname NAME -lastname NAME -email EMAIL  // creates new user")

}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

// Run command line application
func (cli *CommandLine) Run(db *gorm.DB) {

	const (
		createAdminCmd = "create_admin"
		createUserCmd  = "create_user"
		pwResetLinkCmd = "password_reset_link"
		usersCmd       = "users"
	)

	cli.validateArgs()

	createAdminFlags := flag.NewFlagSet(createAdminCmd, flag.ExitOnError)
	createAdminFirstName := createAdminFlags.String("firstname", "", "The first name")
	createAdminLastName := createAdminFlags.String("lastname", "", "The last name")
	createAdminEmail := createAdminFlags.String("email", "", "Email")

	createUserFlags := flag.NewFlagSet(createUserCmd, flag.ExitOnError)
	createUserFirstName := createUserFlags.String("firstname", "", "The first name")
	createUserLastName := createUserFlags.String("lastname", "", "The last name")
	createUserEmail := createUserFlags.String("email", "", "Email")

	pwResetLinkFlags := flag.NewFlagSet(pwResetLinkCmd, flag.ExitOnError)
	pwResetLinkrEmail := pwResetLinkFlags.String("email", "", "Email")

	switch os.Args[1] {

	case createAdminCmd:
		err := createAdminFlags.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case createUserCmd:
		err := createUserFlags.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case pwResetLinkCmd:
		err := pwResetLinkFlags.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case usersCmd: // no arguments - no flags
		err := cli.PrintUsers(db)
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if createAdminFlags.Parsed() {
		if *createAdminFirstName == "" {
			createAdminFlags.Usage()
			runtime.Goexit()
		}
		if *createAdminLastName == "" {
			createAdminFlags.Usage()
			runtime.Goexit()
		}
		if *createAdminEmail == "" {
			createAdminFlags.Usage()
			runtime.Goexit()
		}

		err := cli.createAdmin(db, *createAdminFirstName, *createAdminLastName, *createAdminEmail)
		if err != nil {
			log.Panic(err)
		}
	}

	if createUserFlags.Parsed() {
		if *createUserFirstName == "" {
			createUserFlags.Usage()
			runtime.Goexit()
		}
		if *createUserLastName == "" {
			createUserFlags.Usage()
			runtime.Goexit()
		}
		if *createUserEmail == "" {
			createUserFlags.Usage()
			runtime.Goexit()
		}

		err := cli.createUser(db, *createUserFirstName, *createUserLastName, *createUserEmail)
		if err != nil {
			log.Panic(err)
		}
	}

	if pwResetLinkFlags.Parsed() {
		if *pwResetLinkrEmail == "" {
			pwResetLinkFlags.Usage()
			runtime.Goexit()
		}

		err := cli.PasswordResetLink(db, *pwResetLinkrEmail)
		if err != nil {
			log.Panic(err)
		}
	}

}
