package core

import (
	"github.com/gin-gonic/gin"
)

func (s *Service) CreateAccountHandler(c *gin.Context) {

	/*
				POST /accounts

		REQUEST:
		{
			"email":"example@example.com",
			"name":"Example",
			"surname":"example"
		}

		// validate email
		// check whether there is already an account with that same email address
		// create a new account in db: verified:false, verification_token:78w3t823gt32tg4gt674gt74g..., etc.
		// send email with verification token and instructions
		// return response

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

func (s *Service) ValidateAccountHandler(c *gin.Context) {
	/*
				   POST /account/:account_id/api_token

		REQUEST:
				   {
						"verification_token":"98wtyw4t8h3nc94t34t3gtgc643n7t347gtc396tbgb36"
				   }

		// check verification_token length
		// find in db a non-verifed account with that verification_token
		// check whether they match
		// generate api_token

		RESPONSE:
				   {
						"id":1,
						"email":"example@example.com",
						"verified":true
						"api_token":"key-giowg9w9g49tgh439hy9384hy943hy934hy4u39t8439y"
				   }
	*/
}
