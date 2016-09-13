package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// @@@@@@@@@@@@@@@ router handles @@@@@@@@@@@@@@@

func (s *Service) ApproverHandle(c *gin.Context) {
	// s.ExtenderQueue.TaskQueue <- ExtenderTask{}

	err := s.verifySignature(c)
	if err != nil {
		fmt.Println("verification error:", err)

		c.JSON(200, gin.H{
			"message": "not_valid",
			"error":   err,
		})
		return
	}

	fmt.Printf("approval of %v initiated", c.Param("leaseID"))

	c.JSON(200, gin.H{
		"message": "valid",
		"error":   err,
	})

}

func (s *Service) ExtenderHandle(c *gin.Context) {
	s.ExtenderQueue.TaskQueue <- ExtenderTask{}

	fmt.Printf("renewal of %v initiated", c.Param("leaseID"))

	c.JSON(200, gin.H{
		"message": "hello",
	})
}

func (s *Service) TerminatorHandle(c *gin.Context) {
	s.TerminatorQueue.TaskQueue <- TerminatorTask{}

	fmt.Printf("termination of %v initiated", c.Param("leaseID"))
	// /welcome?firstname=Jane&lastname=Doe
	// lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

	c.JSON(200, gin.H{
		"message": "hello",
	})
}
