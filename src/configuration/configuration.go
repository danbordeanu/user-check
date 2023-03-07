package configuration

import (
	"user-check/utils"
)

// Configuration vars
type Configuration struct {
	Swagger CSwagger

	HttpPort          int32
	CleanupTimeoutSec int32
	Development       bool
	Tls               bool
	GinLogger         bool
	UseSwagger        bool
	Initialized       bool
	NpaUser           string
	NpaPassword       string
	LdapServerAddress string
	LdapGroup         string
	SearchPeople      string
	OncoGroup         string
	LdapCertFile      string
	ApiCertCrtFile    string
	ApiCertKeyFile    string
}

var appConfig Configuration

func AppConfig() *Configuration {
	if appConfig.Initialized == false {
		loadEnvironmentVariables()
		appConfig.Initialized = true
	}
	return &appConfig
}

// loadEnvironmentVariables load env variables
func loadEnvironmentVariables() {
	appConfig.CleanupTimeoutSec = utils.EnvOrDefaultInt32("SHUTDOWN_TIMEOUT", 300)
	// ldap server
	appConfig.LdapServerAddress = utils.EnvOrDefault("LDAP_ADDR", "ldaps://server.com:636")
	// NPA account info
	// user
	appConfig.NpaUser = utils.EnvOrDefault("NPA_USER", "npa@domain.com")
	// password
	// yes, i know how to use a vault and store password there, but this prj waay tooo simple
	appConfig.NpaPassword = utils.EnvOrDefault("NPA_PASSWORD", "passwd")
	// search people
	appConfig.SearchPeople = "OU=eCore Office,OU=People Accounts,DC=domain,DC=com"
	// onco group
	appConfig.OncoGroup = utils.EnvOrDefault("USER_GROUP", "group.users")
	// LDAP certificate file
	appConfig.LdapCertFile = utils.EnvOrDefault("LDAP_CERT_FILE", "cert.crt")
	// API SSL CRT file
	appConfig.ApiCertCrtFile = utils.EnvOrDefault("API_CERT_CRT_FILE", "server.crt")
	appConfig.ApiCertKeyFile = utils.EnvOrDefault("API_CERT_KEY_FILE", "private.key")
}
