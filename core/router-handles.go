package core

import (
	"time"

	"github.com/gin-gonic/gin"
)

// @@@@@@@@@@@@@@@ router handles @@@@@@@@@@@@@@@

func (s *Service) EmailActionHandler(c *gin.Context) {

	err := s.verifySignature(c)
	if err != nil {
		logger.Warn("Signature verification error", "error", err)

		c.JSON(404, gin.H{
			"error": "url not found",
		})
		return
	}

	switch c.Param("action") {
	case "approve":
		logger.Info("Approval of lease initiated", "instance_id", c.Param("instance_id"))

		var lease Lease
		var leasesFound int64
		s.DB.Table("leases").Where(&Lease{
			InstanceID: c.Param("instance_id"),
			UUID:       c.Param("lease_uuid"),
			Terminated: false,
		}).First(&lease).Count(&leasesFound)

		if leasesFound == 0 {
			logger.Warn("No lease found for extension", "count", leasesFound)
			c.JSON(410, gin.H{
				"message": "error",
			})
			return
		}
		if leasesFound > 1 {
			logger.Warn("Multiple leases found for extension", "count", leasesFound)
			c.JSON(410, gin.H{
				"message": "error",
			})
			return
		}

		if lease.TokenOnce != c.Query("t") {
			// TODO: return this info to the http request
			logger.Warn("lease.TokenOnce != c.Query(\"t\")")
			c.JSON(410, gin.H{
				"message": "link expired",
			})
			return
		}

		s.ExtenderQueue.TaskQueue <- ExtenderTask{
			//TokenOnce:  c.Query("t"),
			//UUID:       c.Param("lease_uuid"),
			//InstanceID: c.Param("instance_id"),
			Lease:     lease,
			ExtendBy:  time.Duration(ZCDefaultLeaseDuration),
			Approving: true,
		}

		// TODO: give immediately a response, from here
		c.JSON(202, gin.H{
			"instanceId": c.Param("instance_id"),
			"message":    "Approval request received",
		})
		return

	case "extend":
		logger.Info("Extension of lease initiated", "instance_id", c.Param("instance_id"))

		var lease Lease
		var leasesFound int64
		s.DB.Table("leases").Where(&Lease{
			InstanceID: c.Param("instance_id"),
			UUID:       c.Param("lease_uuid"),
			Terminated: false,
		}).First(&lease).Count(&leasesFound)

		if leasesFound == 0 {
			logger.Warn("No lease found for extension", "count", leasesFound)
			c.JSON(410, gin.H{
				"message": "error",
			})
			return
		}
		if leasesFound > 1 {
			logger.Warn("Multiple leases found for extension", "count", leasesFound)
			c.JSON(410, gin.H{
				"message": "error",
			})
			return
		}

		if lease.TokenOnce != c.Query("t") {
			// TODO: return this info to the http request
			logger.Warn("lease.TokenOnce != c.Query(\"t\")")
			c.JSON(410, gin.H{
				"message": "link expired",
			})
			return
		}

		s.ExtenderQueue.TaskQueue <- ExtenderTask{
			//TokenOnce:  c.Query("t"),
			//UUID:       c.Param("lease_uuid"),
			//InstanceID: c.Param("instance_id"),
			Lease:     lease,
			ExtendBy:  time.Duration(ZCDefaultLeaseDuration),
			Approving: false,
		}

		// TODO: give immediately a response, from here
		c.JSON(202, gin.H{
			"instanceId": c.Param("instance_id"),
			"message":    "Extension initiated",
		})
		return

	case "terminate":
		logger.Info("Termination of lease initiated", "instance_id", c.Param("instance_id"))

		var leaseCount int64
		var leaseToBeTerminated Lease
		s.DB.Table("leases").Where(&Lease{
			InstanceID: c.Param("instance_id"),
			UUID:       c.Param("lease_uuid"),
			Terminated: false,
		}).First(&leaseToBeTerminated).Count(&leaseCount)

		if leaseCount != 1 {
			c.JSON(410, gin.H{
				"message": "error",
			})
			return
		}

		if leaseToBeTerminated.TokenOnce != c.Query("t") {
			c.JSON(406, gin.H{
				"message": "link not usable anymore",
			})
			return
		}

		s.TerminatorQueue.TaskQueue <- TerminatorTask{Lease: leaseToBeTerminated}

		c.JSON(202, gin.H{
			"instanceId": c.Param("instance_id"),
			"message":    "Termination initiated",
		})
	}

}
