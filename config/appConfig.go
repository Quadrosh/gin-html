package config

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

// AppConfig is application config
type AppConfig struct {
	UseCache           bool
	MainTemplateCache  map[string]*template.Template
	AdminTemplateCache map[string]*template.Template
	InfoLog            *log.Logger
	ErrorLog           *log.Logger
	InProduction       bool
	MailSMTP           string
	MailPort           int
	MailLogin          string
	MailPassword       string
	MailEncrypted      bool
	MailFrom           string
	MailTo             string
	SmsTo              string
	SmsKey             string
	AllowOrigin        string
	IsDebug            bool
	IsDbDebug          bool
	ApiTokenExpireSec  uint
	ApiSecret          string
	CWD                string
	LocalConfig
}

func (app *AppConfig) LoadConfig() {

	const fileNameConfigEnv = "config.env"

	// try to init config from executing file
	var err error
	exe, _ := os.Executable()
	var env = filepath.Join(filepath.Dir(exe), fileNameConfigEnv)
	if _, err = os.Stat(env); err == nil { // file not found
		err = godotenv.Load(env)
		if err != nil {
			log.Printf("Error getting Env file: %s, reason: %+v", env, err)
		} else {
			log.Printf("Env file: %s loaded successfully", env)
		}
	}

	var (
		_, b, _, _  = runtime.Caller(0)
		projectRoot = filepath.Join(filepath.Dir(b), "../..")
	)
	app.CWD = projectRoot

	cwd, _ := os.Getwd()
	var cwdenv = filepath.Join(cwd, fileNameConfigEnv)
	if err != nil &&
		len(cwdenv) != 0 {
		// from CWD
		if _, err = os.Stat(cwdenv); err == nil {
			err = godotenv.Load(cwdenv)
			if err != nil {
				log.Printf("Error getting Env file: %s, reason: %+v", cwdenv, err)
			} else {
				log.Printf("Env file: %s loaded successfully", cwdenv)
			}
		}
	}

	if err != nil {
		// Last attempt - godotenv default
		err = godotenv.Load()
		log.Printf("final check Env file: .env")
	}

	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	inProduction := os.Getenv("PRODUCTION") // flag.Bool("production", true, "Application is in production")
	ssl := os.Getenv("SSL")                 //flag.Bool("ssl", false, "is SSL")
	sslSert := os.Getenv("SSL_CERT")        // flag.String("sslsert", "", "SSL sertificat")
	sslKey := os.Getenv("SSL_KEY")          // flag.String("sslkey", "", "SSL key")

	useCache := os.Getenv("TEMPLATE_CACHE") // flag.Bool("cache", true, "Use template cache")
	dbHost := os.Getenv("DB_HOST")          // flag.String("dbhost", "localhost", "Database host")
	dbName := os.Getenv("DB_NAME")          // flag.String("dbname", "", "Database name")
	dbUser := os.Getenv("DB_USER")          // flag.String("dbuser", "", "Database user")
	dbPass := os.Getenv("DB_PASS")          // flag.String("dbpass", "", "Database password")
	dbPort := os.Getenv("DB_PORT")          // flag.String("dbport", "5432", "Database port")
	port := os.Getenv("PORT")               // flag.String("port", "8080", "Application serve port")
	dbSSL := os.Getenv("DB_SSL")            // flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")

	mailSMTP := os.Getenv("MAIL_SMTP")           // flag.String("mailsmtp", "localhost", "Mailer smtp port")
	mailPort := os.Getenv("MAIL_PORT")           //flag.Int("mailport", 25, "Mailer port")
	mailLogin := os.Getenv("MAIL_LOGIN")         //flag.String("maillogin", "", "Mailer login")
	mailPass := os.Getenv("MAIL_PASS")           //flag.String("mailpass", "", "Mailer password")
	mailEncrypted := os.Getenv("MAIL_ENCRYPTED") //flag.Bool("mailencrypted", false, "Mailer encryption")
	mailFrom := os.Getenv("MAIL_FROM")           //flag.String("mailfrom", "", "Mailer sender address")
	// mailTo := os.Getenv("MAIL_SMTP")             //flag.String("mailto", "", "Mailer admin address")
	smsTo := os.Getenv("SMS_TO")   //flag.String("smsto", "", "SMS order recepient")
	smsKey := os.Getenv("SMS_KEY") //flag.String("smskey", "", "SMS api key")

	flag.Parse()

	if dbName == "" || dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	if ssl != "" {
		if sslSert == "" {
			fmt.Println("Missing SSL sert flag")
			os.Exit(1)
		}
		if sslKey == "" {
			fmt.Println("Missing SSL key flag")
			os.Exit(1)
		}
	}

	app.MailSMTP = mailSMTP
	app.MailPort, _ = strconv.Atoi(mailPort)
	app.MailLogin = mailLogin
	app.MailPassword = mailPass
	app.MailEncrypted = mailEncrypted == "true"
	app.MailFrom = mailFrom

	app.SmsTo = smsTo
	app.SmsKey = smsKey
	app.AllowOrigin = os.Getenv("ALLOW_ORIGIN")
	app.IsDebug, _ = strconv.ParseBool(os.Getenv("IS_DEBUG"))
	app.IsDbDebug, _ = strconv.ParseBool(os.Getenv("IS_DB_DEBUG"))

	ApiTokenExpireSecU64, _ := strconv.ParseUint(os.Getenv("API_TOKEN_EXPIRE_SEC"), 10, 32)
	app.ApiTokenExpireSec = uint(ApiTokenExpireSecU64)
	app.ApiSecret = os.Getenv("API_SECRET")

	lConfig := LocalConfig{
		InProduction: inProduction == "true",
		UseCache:     useCache == "true",
		DbHost:       dbHost,
		DbName:       dbName,
		DbUser:       dbUser,
		DbPass:       dbPass,
		DbPort:       dbPort,
		AppPort:      port,
		SSL:          ssl == "true",
		SSLKey:       sslKey,
		SSLCert:      sslSert,
		DbSSL:        dbSSL,
	}

	app.LocalConfig = lConfig

}
