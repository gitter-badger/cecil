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
- **common.go** -- COntains common utility functions.
- **core.go** -- Contains the all the initialization code for the core package.
- **core_test.go** -- core package test.
- **db-models.go** -- Contains the database models.
- **email-action-handler.go** -- Contains the handler function for lease approval|extension|termination link endpoints.
- **email-templates.go** -- Will contain the templates of the emails sent out for specific scenarios (new lease, lease expired, instance terminated, etc.).
- **list-regions-handler.go** -- Contains the handler function for listing the regions for a particular cloudaccount along with their status (active or not).
- **mock_ec2.go** -- Contains a mock of the EC2 API.
- **mock_mailgun.go** -- Contains a mock of the Mailgun API.
- **mock_sqs.go** -- Contains a mock of the SQS API.
- **new-lease-queue-consumer.go** -- Contains the consumer function for the NewLeaseQueue.
- **periodic-jobs.go** -- Contains the periodic job functions
- **sync-regions-handler.go** -- Contains the handler function for changing the status of a region (active: true|false).
- **task-consumers.go** -- Contains some of the functions that consume tasks from queues; some got their own file because are big.
- **task-structs.go** -- Contains the structs of the tasks passed in-out of queues.
- **temporary** -- Is the folder that contains a temporary server that runs the core package.
- **transmission.go** -- Contains the `Transmission` and its methods; `Transmission` is what an SQS message is parsed to.


## Endpoint usage examples

### GET /email_action/leases/:lease_uuid/:instance_id/:action?t=token_once&s=signature

This endpoint is used to allow lease manipulation without login. Each link of this kind is signed to not allow arbitrary command execution.

### POST /accounts/:account_id/cloudaccounts/:cloudaccount_id/owners

This endpoint is used to add an owner to a cloudaccount's owner whitelist. An owner in the whitelist can own leases, for which they will be responsible.

An example of this endpoint's usage is:

```
POST /accounts/1/cloudaccounts/1/owners

Body:
{
  "email" : "example@example.com"
}

Success response:
{"message":"owner added successfully"}
```

CURL example:

```
curl \
-H "Content-Type: application/json" \
-X POST \
-d '{"email":"example@example.com"}' \
http://localhost:8080/accounts/1/cloudaccounts/1/owners
```
