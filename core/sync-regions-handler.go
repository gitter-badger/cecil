package core

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Service) SyncRegionsHandler(c *gin.Context) {

	account, err := s.FetchAccountByID(c.Param("account_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}

	cloudAccount, err := s.FetchCloudAccountByID(c.Param("cloudaccount_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}

	// check whether everything is consistent
	if !account.IsOwnerOf(cloudAccount) {
		c.JSON(404, gin.H{
			"error": "error",
		})
		return
	}

	/* inputPayload is
	{
		"region-name-1":
		{
			"active":true,
		}
	}

	One region at a time.

	@@@@

	{
		"region-name-1":
		{
			"active":false,
			"processing_started_at":"timestamp"
		}
	}
	*/

	// validate parameters
	// check whether account exists
	// check whether cloudaccount exists
	// unmarshal json
	// validate struct
	// get current regions entry
	// compare the inputPayload and the entry from the db
	// if db.regions[region].status != inputPayload[region].status, then change status;
	// call SQS to doublecheck the current status
	// call SQS to subscribe/unsubscribe
	// if success, update in the db
	// return new status
}

func IsValidRegion(r string) bool {
	for regionIndex := range ValidRegions {
		if strings.EqualFold(r, ValidRegions[regionIndex]) {
			return true
		}
	}
	return false
}

var ValidRegions = map[string]bool{
	"us-east-1":      false,
	"us-west-2":      false,
	"us-west-1":      false,
	"eu-west-1":      false,
	"eu-central-1":   false,
	"ap-southeast-1": false,
	"ap-northeast-1": false,
	"ap-southeast-2": false,
	"ap-northeast-2": false,
	"ap-south-1":     false,
	"sa-east-1":      false,
}

type CurrentRegions struct {
	US_east_1      bool `json:"us-east-1"`
	US_west_2      bool `json:"us-west-2"`
	US_west_1      bool `json:"us-west-1"`
	EU_west_1      bool `json:"eu-west-1"`
	EU_central_1   bool `json:"eu-central-1"`
	AP_southeast_1 bool `json:"ap-southeast-1"`
	AP_northeast_1 bool `json:"ap-northeast-1"`
	AP_southeast_2 bool `json:"ap-southeast-2"`
	AP_northeast_2 bool `json:"ap-northeast-2"`
	AP_south_1     bool `json:"ap-south-1"`
	SA_east_1      bool `json:"sa-east-1"`
}
