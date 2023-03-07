package handlers

import (
	"github.com/gin-gonic/gin"
	"user-check/api/response"
	"user-check/ldapcheck"
	"user-check/utils"
	"user-check/utils/go-stats/concurrency"
	"user-check/utils/logger"
)

// UserGroupCount godoc
// @Summary UserGroupCount
// @Description This will return number of users in the group
// @Produce json
// @Success 200 {string} string "success or failure"
// @Router /v1/usercount [get]
func UserGroupCount(c *gin.Context) {
	concurrency.GlobalWaitGroup.Add(1)
	defer concurrency.GlobalWaitGroup.Done()

	log := logger.SugaredLogger().WithContextCorrelationId(c)

	var (
		userLdapProvider *ldapcheck.Provider
		err              error
		result           int
	)

	log.Infof("Count users in specific group")
	ctx := c.Request.Context()

	if userLdapProvider, err = ldapcheck.New(ctx); err != nil {
		log.Errorf("Error while initializing Ldap Provider: %s", err)
		response.FailureResponse(c, nil, utils.HttpError{Code: 500, Err: err})
		return
	}

	secgroupusercount, err := userLdapProvider.QueryUserGroupLdap(ctx)
	if err != nil {
		log.Errorf("ldap query user group failed: %v", err)
		response.FailureResponse(c, nil, utils.HttpError{Code: 400, Err: err})
		return
	} else {
		result = ldapcheck.LdapCountObjects(ctx, secgroupusercount)
		response.SuccessResponse(c, c.MustGet("correlation_id").(string), result)
	}
}

