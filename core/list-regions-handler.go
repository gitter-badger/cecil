// THIS WILL BE DELETED
// THIS WILL BE DELETED
// THIS WILL BE DELETED
// THIS WILL BE DELETED
// THIS WILL BE DELETED
// THIS WILL BE DELETED
// THIS WILL BE DELETED
// THIS WILL BE DELETED
// THIS WILL BE DELETED
// THIS WILL BE DELETED
// THIS WILL BE DELETED
// THIS WILL BE DELETED
// THIS WILL BE DELETED
// THIS WILL BE DELETED
// THIS WILL BE DELETED
// THIS WILL BE DELETED
package core

import (
	"github.com/gin-gonic/gin"
)

func (s *Service) ListRegionsHandler(c *gin.Context) {
	// validate parameters
	// check whether account exists
	// check whether cloudaccount exists
	// return regions object entry from db
	// if entry does not exist, return an object anyway with all .active=false

	// SELECT * FROM Regions WHERE cloud_account_id = cloudaccount.id
	// if len(regions) == 0 ...
	// ... then return an empty RegionsList{}
	// Else, for each region in regions, set RegionsList[region.name].active = region.active
}
