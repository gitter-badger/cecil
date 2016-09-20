package core

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// @@@@@@@@@@@@@@@ router handles @@@@@@@@@@@@@@@

func (s *Service) EmailActionHandler(c *gin.Context) {

	err := s.verifySignature(c)
	if err != nil {
		fmt.Println("verification error:", err)

		c.JSON(404, gin.H{
			"error": "url not found",
		})
		return
	}

	switch c.Param("action") {
	case "approve":
		fmt.Printf("approval of lease for %v initiated", c.Param("instance_id"))

		s.ExtenderQueue.TaskQueue <- ExtenderTask{
			TokenOnce:  c.Query("t"),
			UUID:       c.Param("lease_uuid"),
			InstanceID: c.Param("instance_id"),
			ExtendBy:   time.Duration(ZCDefaultLeaseDuration),
			Approving:  true,
		}

		// TODO: give immediately a response, from here
		c.JSON(202, gin.H{
			"instanceId": c.Param("instance_id"),
			"message":    "Approval request received",
		})
		return

	case "extend":
		fmt.Printf("extension of lease for %v initiated", c.Param("instance_id"))

		s.ExtenderQueue.TaskQueue <- ExtenderTask{
			TokenOnce:  c.Query("t"),
			UUID:       c.Param("lease_uuid"),
			InstanceID: c.Param("instance_id"),
			ExtendBy:   time.Duration(ZCDefaultLeaseDuration),
			Approving:  false,
		}

		// TODO: give immediately a response, from here
		c.JSON(202, gin.H{
			"instanceId": c.Param("instance_id"),
			"message":    "Extension initiated",
		})
		return

	case "terminate":
		fmt.Printf("termination of lease and instance %v initiated", c.Param("instance_id"))

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
