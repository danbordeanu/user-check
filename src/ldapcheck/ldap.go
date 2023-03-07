package ldapcheck

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"io/ioutil"
	"user-check/configuration"
	"user-check/utils/logger"
	"strings"
	"time"
)

type Provider struct {
	LdapServerAddress string
	NpaUser           string
	NpaPassword       string
	SearchPeople      string
	OncoGroup         string
	CertFile          string
}

var (
	ResultCount int
)

// New provide context, variables to be used for example
func New(ctx context.Context) (*Provider, error) {
	var (
		provider *Provider
	)

	log := logger.SugaredLogger().WithContextCorrelationId(ctx)
	log.Debug("Instantiating ldap check provider")

	provider = &Provider{}
	appConfig := configuration.AppConfig()

	// ldap server address
	if appConfig.LdapServerAddress == "" {
		return nil, fmt.Errorf("ldap address is not set")
	} else {
		provider.LdapServerAddress = appConfig.LdapServerAddress
		log.Debugf("LDAP_ADDRESS:%s", provider.LdapServerAddress)
	}

	// npa account
	if appConfig.NpaUser == "" {
		return nil, fmt.Errorf("npa account is not set")
	} else {
		provider.NpaUser = appConfig.NpaUser
		log.Debugf("NPA_USER:%s", provider.NpaUser)
	}

	// npa password
	if appConfig.NpaPassword == "" {
		return nil, fmt.Errorf("npa password is not set")
	} else {
		provider.NpaPassword = appConfig.NpaPassword
	}

	// search people
	if appConfig.SearchPeople == "" {
		return nil, fmt.Errorf("searchpeople is not set")
	} else {
		provider.SearchPeople = appConfig.SearchPeople
		log.Debugf("SEARCH_PEOPLE:%s", provider.SearchPeople)
	}

	// group group
	if appConfig.groupGroup == "" {
		return nil, fmt.Errorf("groupgroup is not set")
	} else {
		provider.groupGroup = appConfig.groupGroup
		log.Debugf("group_GROUP:%s", provider.groupGroup)
	}

	// cert file
	if appConfig.LdapCertFile == "" {
		return nil, fmt.Errorf("cert file not set")
	} else {
		provider.CertFile = appConfig.LdapCertFile
		log.Debugf("LDAP_CERT_FILE:%s", provider.CertFile)
	}

	return provider, nil
}

// CheckUserLdap check if user is in ldap group
func (p *Provider) CheckUserLdap(ctx context.Context, isidmap map[string]interface{}) (*ldap.SearchResult, error) {
	log := logger.SugaredLogger().WithContextCorrelationId(ctx).With("package", "group-license", "action", "check user ldap")

	l, err := p.FuncDialLdap(ctx)
	if err != nil {
		return &ldap.SearchResult{}, err
	}
	err = l.Bind(p.NpaUser, p.NpaPassword)
	if err != nil {
		log.Debugf("error binding:%s", err)
	}

	searchFilter := "(&(objectClass=user)(sAMAccountName=" + isidmap["isid"].(string) + "))"
	searchRequest := ldap.NewSearchRequest(
		p.SearchPeople, // The base dn to search
		2, 0, 0, 0, false,
		searchFilter, // The filter to apply
		[]string{"mail", "sn", "givenName", "memberOf"}, // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Debugf("Failed to search:%v", err)
		return sr, err
	}

	// any good developer is closing the connection after we done with it
	defer l.Close()

	// check how fast is this
	//defer TimeTaken(ctx, time.Now(), "checkuserldap")

	return sr, nil
}

// ReadCertFile read ldap certificate
func (p *Provider) ReadCertFile(ctx context.Context) []byte {
	log := logger.SugaredLogger().WithContextCorrelationId(ctx).With("package", "group-calculator", "action", "read cert file")
	res, err := ioutil.ReadFile(p.CertFile)
	if err != nil {
		log.Debugf("issue readint the cert file from disk:%s", err)
	}
	return res
}

// FuncDialLdap dial ldap
func (p *Provider) FuncDialLdap(ctx context.Context) (*ldap.Conn, error) {

	log := logger.SugaredLogger().WithContextCorrelationId(ctx).With("package", "group-calculator", "action", "dial ldap")

	combinedCerts := p.ReadCertFile(ctx)

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(combinedCerts)

	l, err := ldap.DialURL(p.LdapServerAddress, ldap.DialWithTLSConfig(&tls.Config{
		Rand:                        nil,
		Time:                        nil,
		Certificates:                nil,
		GetCertificate:              nil,
		GetClientCertificate:        nil,
		GetConfigForClient:          nil,
		VerifyPeerCertificate:       nil,
		VerifyConnection:            nil,
		RootCAs:                     nil,
		NextProtos:                  nil,
		ServerName:                  "",
		ClientAuth:                  tls.RequireAndVerifyClientCert,
		ClientCAs:                   caCertPool,
		InsecureSkipVerify:          true,
		CipherSuites:                nil,
		SessionTicketsDisabled:      false,
		ClientSessionCache:          nil,
		MinVersion:                  0,
		MaxVersion:                  0,
		CurvePreferences:            nil,
		DynamicRecordSizingDisabled: false,
		Renegotiation:               0,
		KeyLogWriter:                nil,
	}))

	if err != nil {
		log.Debugf("error dialling up ldap server:%s", err)
		return l, err
	}

	return l, err
}

// QueryUserGroupLdap count total number of users in ldap group
func (p *Provider) QueryUserGroupLdap(ctx context.Context) (*ldap.SearchResult, error) {
	log := logger.SugaredLogger().WithContextCorrelationId(ctx).With("package", "group-license", "action", "get users from  security group")

	l, err := p.FuncDialLdap(ctx)
	if err != nil {
		return &ldap.SearchResult{}, err
	}

	err = l.Bind(p.NpaUser, p.NpaPassword)
	if err != nil {
		log.Debugf("error binding:%s", err)
	}

	searchGroupFilter := "(&(objectCategory=group)(cn=" + p.groupGroup + "))"

	searchRequestGroups := ldap.NewSearchRequest(
		"CN=Groups,DC=domain,DC=com", // The base dn to search
		2, 0, 0, 0, false,
		searchGroupFilter,  // The filter to apply
		[]string{"member"}, // A list attributes to retrieve
		nil,
	)
	srg, err := l.Search(searchRequestGroups)
	if err != nil {
		log.Debugf("Failed to search group:%v", err)
		return srg, err
	}

	// any good developer is closing the connection after we done with it
	defer l.Close()
	//defer TimeTaken(ctx, time.Now(), "ldapquerygroup")
	log.Debugf("ldap query user group returned:%s", srg)
	return srg, err
}

// LdapCountObjects count members in ldap group
func LdapCountObjects(ctx context.Context, newusercount *ldap.SearchResult) int {
	log := logger.SugaredLogger().WithContextCorrelationId(ctx).With("package", "go-user-check", "action", "count members in ldap security group")
	for _, entry := range newusercount.Entries {
		ResultCount = strings.Count(fmt.Sprint(entry.GetAttributeValues("member")), "CN=")
	}
	log.Debugf("total numbers of members in security group is: %d", ResultCount)
	//defer TimeTaken(ctx, time.Now(), "ldapcountobjects")
	return ResultCount
}

// IsUserInGroup check if user is in group
func (p *Provider) IsUserInGroup(ctx context.Context, list []string) bool {
	log := logger.SugaredLogger().WithContextCorrelationId(ctx).With("package", "go-user-check", "action", "check if user is in the security group")
	//log.Debug(list)
	groupgropupbase := "CN=" + p.groupGroup + ",CN=Security,CN=Groups,DC=domain,DC=com"
	for _, v := range list {
		if v == groupgropupbase {
			log.Infof("user is in the security group:%s", groupgropupbase)
			return true
		}
	}
	log.Infof("user is not in the security group:%s", groupgropupbase)
	return false
}

// TimeTaken simple function returning how long it takes to execute a function
// be sure deffer is disabled in the functions in prod env
// use this just for internal debugging
func TimeTaken(ctx context.Context, t time.Time, name string) {
	log := logger.SugaredLogger().WithContextCorrelationId(ctx).With("package", "go-user-check", "action", "time check function")
	elapsed := time.Since(t)
	//elapsed := float64(time.Since(t) / time.Millisecond)
	log.Debugf("Time to execute function:%s was:%s", name, elapsed)
}
