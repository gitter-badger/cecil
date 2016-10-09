package core

import (
	"github.com/gin-gonic/gin"
)

func (s *Service) CreateAccount(c *gin.Context) {

	/*
				POST /accounts

		REQUEST:
		{
			"email":"example@example.com",
			"name":"Example",
			"surname":"example"
		}

		// validate email

		RESPONSE:
		   {
				"id":1,
				"email":"example@example.com",
				"verified":false
		   }
	*/

	/*
				    Email with Verification token +
		           instructions to create API token
	*/

}

func (s *Service) ValidateAccount(c *gin.Context) {
	/*
				   POST /account/:account_id/api_token

		REQUEST:
				   {
						"verification_token":"98wtyw4t8h3nc94t34t3gtgc643n7t347gtc396tbgb36"
				   }

		RESPONSE:
				   {
						"id":1,
						"email":"example@example.com",
						"verified":true
						"api_token":"key-giowg9w9g49tgh439hy9384hy943hy934hy4u39t8439y"
				   }
	*/
}
