package core

import (
	"time"

	"github.com/gin-gonic/gin"
)

// EmailActionHandler accepts and validates requests from links in emails;
// actions are: approve|extend|terminate
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

		var leaseToBeApproved Lease
		var leaseCount int64
		s.DB.Table("leases").Where(&Lease{
			InstanceID: c.Param("instance_id"),
			UUID:       c.Param("lease_uuid"),
			Terminated: false,
		}).Count(&leaseCount).First(&leaseToBeApproved)

		if leaseCount == 0 {
			logger.Warn("No lease found for approval", "count", leaseCount)
			c.JSON(410, gin.H{
				"message": "error",
			})
			return
		}
		if leaseCount > 1 {
			logger.Warn("Multiple leases found for approval", "count", leaseCount)
			c.JSON(410, gin.H{
				"message": "error",
			})
			return
		}

		if leaseToBeApproved.TokenOnce != c.Query("t") {
			logger.Warn("leaseToBeApproved.TokenOnce != c.Query(\"t\")")
			c.JSON(410, gin.H{
				"message": "link expired",
			})
			return
		}

		s.ExtenderQueue.TaskQueue <- ExtenderTask{
			Lease:     leaseToBeApproved,
			ExtendBy:  time.Duration(ZCDefaultLeaseDuration),
			Approving: true,
		}

		c.JSON(202, gin.H{
			"instanceId": c.Param("instance_id"),
			"message":    "Approval request received",
		})
		return

	case "extend":
		logger.Info("Extension of lease initiated", "instance_id", c.Param("instance_id"))

		var leaseToBeExtended Lease
		var leaseCount int64
		s.DB.Table("leases").Where(&Lease{
			InstanceID: c.Param("instance_id"),
			UUID:       c.Param("lease_uuid"),
			Terminated: false,
		}).Count(&leaseCount).First(&leaseToBeExtended)

		if leaseCount == 0 {
			logger.Warn("No lease found for extension", "count", leaseCount)
			c.JSON(410, gin.H{
				"message": "error",
			})
			return
		}
		if leaseCount > 1 {
			logger.Warn("Multiple leases found for extension", "count", leaseCount)
			c.JSON(410, gin.H{
				"message": "error",
			})
			return
		}

		if leaseToBeExtended.TokenOnce != c.Query("t") {
			logger.Warn("leaseToBeExtended.TokenOnce != c.Query(\"t\")")
			c.JSON(410, gin.H{
				"message": "link expired",
			})
			return
		}

		s.ExtenderQueue.TaskQueue <- ExtenderTask{
			Lease:     leaseToBeExtended,
			ExtendBy:  time.Duration(ZCDefaultLeaseDuration),
			Approving: false,
		}

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
		}).Count(&leaseCount).First(&leaseToBeTerminated)

		if leaseCount == 0 {
			logger.Warn("No lease found for approval", "count", leaseCount)
			c.JSON(410, gin.H{
				"message": "error",
			})
			return
		}
		if leaseCount > 1 {
			logger.Warn("Multiple leases found for approval", "count", leaseCount)
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
