package core

import (
	"github.com/gin-gonic/gin"
)

func (s *Service) SyncRegionsHandler(c *gin.Context) {
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
			"processing":false,
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
