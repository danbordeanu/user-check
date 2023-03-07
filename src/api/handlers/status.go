package handlers

import (
	"github.com/gin-gonic/gin"
	"user-check/api/response"
	"user-check/configuration"
	"user-check/ldapcheck"
	"user-check/utils"
	"user-check/utils/go-stats/concurrency"
	"user-check/utils/logger"
	"os"
)

// Status godoc
// @Summary HealthCheck Endpoint
// @Description This return API status
// @Produce json
// @Success 200 {string} string "api pid and ldap status"
// @Router /v1/status [get]
func Status(c *gin.Context) {
	concurrency.GlobalWaitGroup.Add(1)
	defer concurrency.GlobalWaitGroup.Done()

	log := logger.SugaredLogger().WithContextCorrelationId(c)

	log.Info("Return API status")

	var (
		userLdapProvider *ldapcheck.Provider
		err              error
		ldapstatus       string
	)

	type status struct {
		LdapStatus string
		ProcessPid int64
	}

	ctx := c.Request.Context()

	if userLdapProvider, err = ldapcheck.New(ctx); err != nil {
		log.Errorf("Error while initializing Ldap Provider: %s", err)
		response.FailureResponse(c, nil, utils.HttpError{Code: 500, Err: err})
		return // return here because we don't want to continue if we failed to initialize ldap provider
	}

	_, err = userLdapProvider.FuncDialLdap(ctx)
	if err != nil {
		log.Errorf("seems ldap dialing not working: %v", err)
		ldapstatus = configuration.LdapDown
	} else {
		log.Info("ldap is up&running, we can dial")
		ldapstatus = configuration.LdapUp
	}

	response.SuccessResponse(c, c.MustGet("correlation_id").(string), status{
		LdapStatus: ldapstatus,
		ProcessPid: int64(os.Getpid()),
	})

}
