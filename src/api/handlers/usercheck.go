package handlers

import (
	"github.com/gin-gonic/gin"
	"user-check/api/response"
	"user-check/ldapcheck"
	"user-check/utils"
	"user-check/utils/go-stats/concurrency"
	"user-check/utils/logger"
)

// UserCheck godoc
// @Summary UserCheck
// @Description This will validate if user is part of the group
// @Produce json
// @Param isid path string true "User isid"
// @Success 200 {string} string "true or false"
// @Router /v1/usercheck/{isid} [get]
func UserCheck(c *gin.Context) {
	concurrency.GlobalWaitGroup.Add(1)
	defer concurrency.GlobalWaitGroup.Done()

	log := logger.SugaredLogger().WithContextCorrelationId(c)

	// let's do a map, we love maps :D
	isidmap := map[string]interface{}{
		"isid": c.Param("isid"),
	}

	var (
		userLdapProvider *ldapcheck.Provider
		err              error
		status           string
	)

	log.Debugf("Payload: user isid:%v", isidmap["isid"].(string))

	ctx := c.Request.Context()

	if userLdapProvider, err = ldapcheck.New(ctx); err != nil {
		log.Errorf("Error while initializing Ldap Provider: %s", err)
		response.FailureResponse(c, nil, utils.HttpError{Code: 500, Err: err})
		return // return here because we don't want to continue if we failed to initialize ldap provider
	}

	newusersearch, err := userLdapProvider.CheckUserLdap(ctx, isidmap)
	if err != nil {
		log.Errorf("check user exists failed: %v", err)
		response.FailureResponse(c, nil, utils.HttpError{Code: 400, Err: err})
		return
	} else {
		log.Debugf("check user exists returned:%s", *newusersearch)
	}

	if len(newusersearch.Entries) == 0 {
		log.Infof("no info in Ldap found for isid:%s", isidmap["isid"].(string))
		status = "false"
	} else {
		log.Infof("there is info in Ldap for isid:%s", isidmap["isid"].(string))

		for _, entry := range newusersearch.Entries {
			// print just email
			log.Debugf("%s: %v\n", entry.DN, entry.GetAttributeValue("mail"))
			//log.Debugf("%s: %v\n", entry.DN, entry.GetAttributeValue("mail"), entry.GetAttributeValues("memberOf"))

			if userLdapProvider.IsUserInGroup(ctx, entry.GetAttributeValues("memberOf")) {
				status = "true"
			} else {
				status = "false"
			}
		}
	}

	response.SuccessResponse(c, c.MustGet("correlation_id").(string), status)

}
