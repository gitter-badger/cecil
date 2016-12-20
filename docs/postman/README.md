## How to import to Postman:

#### Instructions:

1. Open Postman
2. Click on "Import"
3. Import `cecil.postman_collection.json`
4. Make sure to run it with a "cecil_environment"

Run the first API request with your name and email address.

After you receive the email with `verification_token`, paste it as payload in the second API request.

Now you can run the other endpoints as the JWT token from the second response has been added to the environment.
