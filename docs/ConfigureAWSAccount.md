## Endpoint usage examples

## Create account

```bash
curl -X POST \
-H "Cache-Control: no-cache" \
-d '{
	"email":"example@example.com",
	"name":"John",
	"surname":"Smith"
}' \
"http://127.0.0.1:8080/accounts"
```

Response:

```json
{
  "email": "example@example.com",
  "account_id": 1,
  "response": "An email has been sent to the specified address with a verification token and instructions.",
  "verified": false
}
```
You will receive an email with a vefication code.


## Verify account and get api token

```bash
curl -X POST \
-H "Cache-Control: no-cache" \
-d '{"verification_token":"0d78a4e0-9922-4b55-93d7-5adfd0f589be7b9a0fa6-c5bc-4991-9f8e-b8bdbc429343322e200c-ab6c-4189-9e81-453ab0b34d56"}' \
"http://0.0.0.0:8080/accounts/1/api_token"
```

Response:

```json
{
  "account_id": 1,
  "api_token": "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoxLCJpYXQiOjE0Nzc0MDg1MzJ9.tr5Ark32AIQyYfM4AnQuC4I6ROQsP7PUSuz6hMR5EOMjDEHQ74A6JKxxR08OkdIgA8NCLw7a8oUyKqDc4XalrQKIq--FCZzf47dswMsJNjtwZPPFTX1hLjhsvuuQiVvtm39jjJL_t4l-ICa0oKX8nrJNGmB5epVR3KMPySlXXShUx-vc77P6My4WOpLIZV8lyeVlobRvLxfCKyXtqxKSRiu0-oJ1rXxCDkcGVvGFMk8vVjYeXDHM4dITuoweb_1TVHxRelePKtpuw5BEyakYXJmLI7m3eQYk8Pv9sBpviS2KhGjq9qPG6kweopGNCuYsrF0L1x5YZ3jWcBL0-KpK2g",
  "email": "example@example.com",
  "verified": true
}
```

Use the api token to manage your account.

## Add CloudAccount

```bash
curl -X POST \
-H "Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoxLCJpYXQiOjE0Nzc0MDg1MzJ9.tr5Ark32AIQyYfM4AnQuC4I6ROQsP7PUSuz6hMR5EOMjDEHQ74A6JKxxR08OkdIgA8NCLw7a8oUyKqDc4XalrQKIq--FCZzf47dswMsJNjtwZPPFTX1hLjhsvuuQiVvtm39jjJL_t4l-ICa0oKX8nrJNGmB5epVR3KMPySlXXShUx-vc77P6My4WOpLIZV8lyeVlobRvLxfCKyXtqxKSRiu0-oJ1rXxCDkcGVvGFMk8vVjYeXDHM4dITuoweb_1TVHxRelePKtpuw5BEyakYXJmLI7m3eQYk8Pv9sBpviS2KhGjq9qPG6kweopGNCuYsrF0L1x5YZ3jWcBL0-KpK2g" \
-H "Cache-Control: no-cache" \
-d '{
	"aws_id":"0123456789"
}' \
"http://0.0.0.0:8080/accounts/1/cloudaccounts"
```

Response:

```json
{
  "aws_id": "0123456789",
  "cloudaccount_id": 1,
  "initial_setup_cloudformation_url": "/accounts/1/cloudaccounts/1/cecil-aws-initial-setup.template",
  "region_setup_cloudformation_url": "/accounts/1/cloudaccounts/1/cecil-aws-region-setup.template"
}
```

Before this cloudaccount is active, you need to setup the Cecil stacks on your AWS account:

1. The first stack is the **initial stack**. It's a one-time only setup, and will be valid for the whole AWS account.
2.  The second stack is the **region stack**. This stack is to be created on each region you want to monitor with Cecil.

To setup the stacks, download them from the urls provided in this response. And then use AWS cli or AWS web gui to set them up.


## Cloudformation template for initial setup

First download it:

```bash
curl -X GET \
-H "Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoxLCJpYXQiOjE0Nzc0MDg1MzJ9.tr5Ark32AIQyYfM4AnQuC4I6ROQsP7PUSuz6hMR5EOMjDEHQ74A6JKxxR08OkdIgA8NCLw7a8oUyKqDc4XalrQKIq--FCZzf47dswMsJNjtwZPPFTX1hLjhsvuuQiVvtm39jjJL_t4l-ICa0oKX8nrJNGmB5epVR3KMPySlXXShUx-vc77P6My4WOpLIZV8lyeVlobRvLxfCKyXtqxKSRiu0-oJ1rXxCDkcGVvGFMk8vVjYeXDHM4dITuoweb_1TVHxRelePKtpuw5BEyakYXJmLI7m3eQYk8Pv9sBpviS2KhGjq9qPG6kweopGNCuYsrF0L1x5YZ3jWcBL0-KpK2g" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/cecil-aws-initial-setup.template" > cecil-aws-initial-setup.template
```

Then install it:

```bash
export AWS_ACCESS_KEY_ID=<access key id of cecil user aws account>
export AWS_SECRET_ACCESS_KEY=<access key id of cecil user aws account>
aws cloudformation create-stack --stack-name "CecilStack" --template-body "file://cecil-aws-initial-setup.template" --region us-east-1 --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM
```

Or alternatively you can upload this in the Cloudformation section of the AWS web UI.

## Cloudformation template for REGION setup

```bash
curl -X GET \
-H "Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoxLCJpYXQiOjE0Nzc0MDg1MzJ9.tr5Ark32AIQyYfM4AnQuC4I6ROQsP7PUSuz6hMR5EOMjDEHQ74A6JKxxR08OkdIgA8NCLw7a8oUyKqDc4XalrQKIq--FCZzf47dswMsJNjtwZPPFTX1hLjhsvuuQiVvtm39jjJL_t4l-ICa0oKX8nrJNGmB5epVR3KMPySlXXShUx-vc77P6My4WOpLIZV8lyeVlobRvLxfCKyXtqxKSRiu0-oJ1rXxCDkcGVvGFMk8vVjYeXDHM4dITuoweb_1TVHxRelePKtpuw5BEyakYXJmLI7m3eQYk8Pv9sBpviS2KhGjq9qPG6kweopGNCuYsrF0L1x5YZ3jWcBL0-KpK2g" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/cecil-aws-region-setup.template" > cecil-aws-region-setup.template
```

Then install it:

```bash
aws cloudformation create-stack --stack-name "CecilUSEastStack" --template-body "file://cecil-aws-region-setup.template" --region us-east-1
```

After this has been successfully setup by AWS, you will receive an email from Cecil.

## Add email to owner tag whitelist

```bash
curl -X POST \
-H "Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoxLCJpYXQiOjE0Nzc0MDg1MzJ9.tr5Ark32AIQyYfM4AnQuC4I6ROQsP7PUSuz6hMR5EOMjDEHQ74A6JKxxR08OkdIgA8NCLw7a8oUyKqDc4XalrQKIq--FCZzf47dswMsJNjtwZPPFTX1hLjhsvuuQiVvtm39jjJL_t4l-ICa0oKX8nrJNGmB5epVR3KMPySlXXShUx-vc77P6My4WOpLIZV8lyeVlobRvLxfCKyXtqxKSRiu0-oJ1rXxCDkcGVvGFMk8vVjYeXDHM4dITuoweb_1TVHxRelePKtpuw5BEyakYXJmLI7m3eQYk8Pv9sBpviS2KhGjq9qPG6kweopGNCuYsrF0L1x5YZ3jWcBL0-KpK2g" \
-H "Cache-Control: no-cache" \
-d '{"email":"someone.legit@example.com"}' \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/owners"
```

Response:

```json
{
  "message": "owner added successfully to whitelist"
}
```
