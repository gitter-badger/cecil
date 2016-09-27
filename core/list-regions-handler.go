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
}
