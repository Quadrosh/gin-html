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

	cli.validateArgs()

	createAdminCmd := flag.NewFlagSet("createadmin", flag.ExitOnError)
	createAdminFirstName := createAdminCmd.String("firstname", "", "The first name")
	createAdminLastName := createAdminCmd.String("lastname", "", "The last name")
	createAdminEmail := createAdminCmd.String("email", "", "Email")

	createUserCmd := flag.NewFlagSet("createuser", flag.ExitOnError)
	createUserFirstName := createUserCmd.String("firstname", "", "The first name")
	createUserLastName := createUserCmd.String("lastname", "", "The last name")
	createUserEmail := createUserCmd.String("email", "", "Email")

	switch os.Args[1] {

	case "createadmin":
		err := createAdminCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createuser":
		err := createUserCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if createAdminCmd.Parsed() {
		if *createAdminFirstName == "" {
			createAdminCmd.Usage()
			runtime.Goexit()
		}
		if *createAdminLastName == "" {
			createAdminCmd.Usage()
			runtime.Goexit()
		}
		if *createAdminEmail == "" {
			createAdminCmd.Usage()
			runtime.Goexit()
		}

		err := cli.createAdmin(db, *createAdminFirstName, *createAdminLastName, *createAdminEmail)
		if err != nil {
			log.Fatal(err)
		}
	}

	if createUserCmd.Parsed() {
		if *createUserFirstName == "" {
			createUserCmd.Usage()
			runtime.Goexit()
		}
		if *createUserLastName == "" {
			createUserCmd.Usage()
			runtime.Goexit()
		}
		if *createUserEmail == "" {
			createUserCmd.Usage()
			runtime.Goexit()
		}

		err := cli.createUser(db, *createUserFirstName, *createUserLastName, *createUserEmail)
		if err != nil {
			log.Fatal(err)
		}
	}

}
