package config

import (
	"html/template"
	"log"
)

// AppConfig is application config
type AppConfig struct {
	UseCache            bool
	CommonTemplateCache map[string]*template.Template
	AdminTemplateCache  map[string]*template.Template
	ArtistTemplateCache map[string]*template.Template
	InfoLog             *log.Logger
	ErrorLog            *log.Logger
	InProduction        bool
	// Session             *scs.SessionManager
	// MailChan            chan models.MailData
	MailSMTP            string
	MailPort            int
	MailLogin           string
	MailPassword        string
	MailEncrypted       bool
	MailFrom            string
	MailTo              string
	SmsTo               string
	SmsKey              string
	AllowOrigin         string
	IsDebug             bool
	IsDbDebug           bool
	ApiTokenExpireSec   uint
	ApiSecret           string
	DefaultEntriesCount int
	CWD                 string
}
