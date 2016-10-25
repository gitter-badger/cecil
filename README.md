[![CircleCI](https://circleci.com/gh/tleyden/zerocloud.svg?style=svg&circle-token=0b966949f6517187f0a2cece8aac8be59e0182a3)](https://circleci.com/gh/tleyden/zerocloud)

# Mission

Allow your devs and testers unfettered access to create AWS instances yet make it impossible for them to forget about unused resources and let your AWS bill spin out of control.

You define your policies, we enforce them.

# Get code

```
go get -t github.com/tleyden/zerocloud/...
```

# ZEROCLOUD AWS Setup (AWS web GUI)

- Go to https://console.aws.amazon.com/billing/home?#/account and save your Account ID
- Go to https://console.aws.amazon.com/cloudformation/home and click "Create stack"
	- (make sure ZeroCloudQueue does not exists)
	- use docs/cloudformation-templates/zerocloud-root.template
	- give whatever unique name to stack (e.g. "ZeroCloudRootStack")
	- allow and create
	- wait stack creation
- Go to https://console.aws.amazon.com/iam/home
	- click on "Users"
	- click on "ZeroCloudRootUser"
	- click on "Security Credentials" tab
	- click on "Create Access Key" and save/download the credentials
	- don't close this window with this user logged in (will need later)

# ZEROCLOUD AWS Setup (AWS CLI)


```
aws cloudformation create-stack --stack-name "ZeroCloudRootStack" \
--template-body "file://./docs/cloudformation-templates/zerocloud-root.template"
```

# BigDB AWS Setup (AWS web GUI)

Start the ZeroCloud service with `go run main.go`

- Login to another AWS account (this will be the customer) in another web browser (e.g. Firefox if you used Chrome), or incognito window.
- Go to https://console.aws.amazon.com/cloudformation/home and click "Create stack"
	- use docs/cloudformation-templates/zerocloud-aws-initial-setup.template
	- give whatever name to stack (e.g. "ZeroCloudInitialStack")
	- in parameters, for "IAMRoleExternalID" write "hithere"
	- in parameters, for "ZeroCloudAWSID" write the AWS id you saved before
	- allow and create
	- wait stack creation

To setup a region:
- Go to https://console.aws.amazon.com/cloudformation/home and click "Create stack"
	- use docs/cloudformation-templates/zerocloud-aws-region-setup.template
	- give whatever name to stack (e.g. "ZeroCloudRegionStack")
	- in parameters, for "ZeroCloudAWSID" write the AWS id you saved before
	- allow and create
	- wait stack creation

After a couple minutes you (BigDB) should receive an email confirm the region has been setup.


# BigDB AWS Setup (AWS CLI)


```
aws cloudformation create-stack --stack-name "ZeroCloudInitialStack" \
--template-body "file://./docs/cloudformation-templates/zerocloud-aws-initial-setup.template" \
--parameters ParameterKey=ZeroCloudAWSID,ParameterValue=123456789101 \
ParameterKey=IAMRoleExternalID,ParameterValue=hithere
```

Wait for the creation to complete and run

```
aws cloudformation create-stack --stack-name "ZeroCloudRegionStack" \
--template-body "file://./docs/cloudformation-templates/zerocloud-aws-region-setup.template" \
--parameters ParameterKey=ZeroCloudAWSID,ParameterValue=123456789101
```

# Service setup

- Open a terminal tab/window
- cd to `github.com/tleyden/zerocloud/`
- Enter each of the following commands in the terminal (with a leading space), completing with the proper values

```
 export AWS_REGION=<fill in here>

 export AWS_ACCESS_KEY_ID=<fill in here>
 export AWS_SECRET_ACCESS_KEY=<fill in here>
 export AWS_ACCOUNT_ID=<fill in here>

 export MAILERDOMAIN=mg.zerocloud.co
 export MAILERAPIKEY=<fill in here>
 export MAILERPUBLICAPIKEY=<fill in here>
```

# Run

Run `go run main.go` or use the [docker container](docs/docker/README.md)

Now, create an instance on BigDB's account without a ZeroCloudOwner tag.

BigDB's admin will receive an email (might not arrive immediately; using a sandbox mailgun account).

Try also with

`ZeroCloudOwner = nope`

`ZeroCloudOwner = someone@unauthorized.site`

`ZeroCloudOwner = dev@bigbd.io` (replacing `dev@bigbd.io` with what you wrote on line 198 in `core.go`)

`ZeroCloudOwner = admin@bigdb.io` (replacing `admin@bigdb.io` with BigDB's admin email you used on line 175 in `core.go`)


## core package contents

- **add-owner-handler.go** -- Contains the handler function for adding a new owner to owner's whitelist for a cloudaccount.
- **aws.go** -- Contains SQS structs and DefaultEc2ServiceFactory.
- **common.go** -- Contains common utility functions.
- **core.go** -- Contains the all the initialization code for the core package.
- **core_test.go** -- core package test.
- **db-models.go** -- Contains the database models.
- **email-action-handler.go** -- Contains the handler function for lease approval|extension|termination link endpoints.
- **email-templates.go** -- Will contain the templates of the emails sent out for specific scenarios (new lease, lease expired, instance terminated, etc.).
- **mock_ec2.go** -- Contains a mock of the EC2 API.
- **mock_mailgun.go** -- Contains a mock of the Mailgun API.
- **mock_sqs.go** -- Contains a mock of the SQS API.
- **new-lease-queue-consumer.go** -- Contains the consumer function for the NewLeaseQueue.
- **periodic-jobs.go** -- Contains the periodic job functions
- **service.go** -- Contains the Service struct and the initialization methods (to setup queues, db, external services, etc.)
- **task-consumers.go** -- Contains some of the functions that consume tasks from queues; some got their own file because are big.
- **task-structs.go** -- Contains the structs of the tasks passed in-out of queues.
- **transmission.go** -- Contains the `Transmission` and its methods; `Transmission` is what an SQS message is parsed to.

## Endpoint usage examples

## Create account

```
curl -X POST \
-H "Cache-Control: no-cache" \
-d '{
	"email":"example@example.com",
	"name":"Example",
	"surname":"Example"
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

```
curl -X POST \
-H "Cache-Control: no-cache" \
-d '{"verification_token":"0d78a4e0-9922-4b55-93d7-5adfd0f589be7b9a0fa6-c5bc-4991-9f8e-b8bdbc429343322e200c-ab6c-4189-9e81-453ab0b34d56"}' "http://0.0.0.0:8080/accounts/1/api_token"
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

```curl
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
  "initial_setup_cloudformation_url": "/accounts/1/cloudaccounts/1/zerocloud-aws-initial-setup.template",
  "region_setup_cloudformation_url": "/accounts/1/cloudaccounts/1/zerocloud-aws-region-setup.template"
}
```

Before this cloudaccount is active, you need to setup the Cecil stacks on your AWS account:

1. The first stack is the **initial stack**. It's a one-time only setup, and will be valid for thr whole AWS account.
2.  The second stack is the **region stack**. This stack is to be created on each region you want to monitor with Cecil.

To setup the stacks, download them from the urls provided in this response. ANd then use AWS cli or AWS web gui to set them up.

## Add email to owner tag whitelist

```curl
curl -X POST \
-H "Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoxLCJpYXQiOjE0Nzc0MDg1MzJ9.tr5Ark32AIQyYfM4AnQuC4I6ROQsP7PUSuz6hMR5EOMjDEHQ74A6JKxxR08OkdIgA8NCLw7a8oUyKqDc4XalrQKIq--FCZzf47dswMsJNjtwZPPFTX1hLjhsvuuQiVvtm39jjJL_t4l-ICa0oKX8nrJNGmB5epVR3KMPySlXXShUx-vc77P6My4WOpLIZV8lyeVlobRvLxfCKyXtqxKSRiu0-oJ1rXxCDkcGVvGFMk8vVjYeXDHM4dITuoweb_1TVHxRelePKtpuw5BEyakYXJmLI7m3eQYk8Pv9sBpviS2KhGjq9qPG6kweopGNCuYsrF0L1x5YZ3jWcBL0-KpK2g" \
-H "Cache-Control: no-cache" \
-d '{"email":"slavomir.balsan@gmail.com"}' \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/owners"
```

Response:

```json
{
  "message": "owner added successfully to whitelist"
}
```

## Download cloudformation template for initial setup

```
curl -X GET \
-H "Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoxLCJpYXQiOjE0Nzc0MDg1MzJ9.tr5Ark32AIQyYfM4AnQuC4I6ROQsP7PUSuz6hMR5EOMjDEHQ74A6JKxxR08OkdIgA8NCLw7a8oUyKqDc4XalrQKIq--FCZzf47dswMsJNjtwZPPFTX1hLjhsvuuQiVvtm39jjJL_t4l-ICa0oKX8nrJNGmB5epVR3KMPySlXXShUx-vc77P6My4WOpLIZV8lyeVlobRvLxfCKyXtqxKSRiu0-oJ1rXxCDkcGVvGFMk8vVjYeXDHM4dITuoweb_1TVHxRelePKtpuw5BEyakYXJmLI7m3eQYk8Pv9sBpviS2KhGjq9qPG6kweopGNCuYsrF0L1x5YZ3jWcBL0-KpK2g" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/zerocloud-aws-initial-setup.template"
```

Response:

```
{
    "AWSTemplateFormatVersion": "2010-09-09",
    "Description": "Initial global zerocloud setup (to do once)",
			....
	[TRUNCATED]
```

## Download cloudformation template for REGION setup

```
curl -X GET \
-H "Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoxLCJpYXQiOjE0Nzc0MDg1MzJ9.tr5Ark32AIQyYfM4AnQuC4I6ROQsP7PUSuz6hMR5EOMjDEHQ74A6JKxxR08OkdIgA8NCLw7a8oUyKqDc4XalrQKIq--FCZzf47dswMsJNjtwZPPFTX1hLjhsvuuQiVvtm39jjJL_t4l-ICa0oKX8nrJNGmB5epVR3KMPySlXXShUx-vc77P6My4WOpLIZV8lyeVlobRvLxfCKyXtqxKSRiu0-oJ1rXxCDkcGVvGFMk8vVjYeXDHM4dITuoweb_1TVHxRelePKtpuw5BEyakYXJmLI7m3eQYk8Pv9sBpviS2KhGjq9qPG6kweopGNCuYsrF0L1x5YZ3jWcBL0-KpK2g" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/zerocloud-aws-region-setup.template"
```

Response:

```
{
    "AWSTemplateFormatVersion": "2010-09-09",
    "Description": "Setup zerocloud on your aws account for a region",
		....
  [TRUNCATED]
```
