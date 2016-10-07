[![CircleCI](https://circleci.com/gh/tleyden/zerocloud.svg?style=svg&circle-token=0b966949f6517187f0a2cece8aac8be59e0182a3)](https://circleci.com/gh/tleyden/zerocloud) 

# Mission

Allow your devs and testers unfettered access to create AWS instances yet make it impossible for them to forget about unused resources and let your AWS bill spin out of control.

You define your policies, we enforce them.

# Get code

```
go get -t github.com/tleyden/zerocloud/...
```

# ZEROCLOUD AWS Setup

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

```
aws cloudformation create-stack --stack-name "ZeroCloudRegionStack" \
--template-body "file://./docs/cloudformation-templates/zerocloud-aws-region-setup.template" \
--parameters ParameterKey=ZeroCloudAWSID,ParameterValue=123456789101
```

# (DEPRECATED) BigDB AWS Setup

- Login to another AWS account (this wiil be the customer) in another web browser (e.g. Firefox if you used Chrome), or incognito window.
- Go to https://console.aws.amazon.com/cloudformation/home and click "Create stack"
	- use docs/cloudformation-templates/zerocloud-user.template
	- give whatever name to stack (e.g. "ZeroCloudStack")
	- in parameters, for "IAMRoleExternalID" write "hithere"
	- in parameters, for "ZeroCloudAWSID" write the AWS id you saved before
	- allow and create
	- wait stack creation
- Go to https://console.aws.amazon.com/sns/v2/home
	- click on "Topics" on the left, and the click on the ARN of "ZeroCloudTopic"
	- copy the ARN of ZeroCLoudTopic
	- go to the the ZeroCloud AWS window you left open
		- go to https://console.aws.amazon.com/sqs/home
		- click on ZeroCloudQueue
		- the click on "Queue Actions" and "Subscribe queue to SNS topic"
		- enter the ARN you copied before and click on "Subscribe"
		- the "Topic Subscription Result" should say "Successfully subscribed the following queue to the SNS topic ..."

# Service setup

- Open a terminal tab/window
- cd to `github.com/tleyden/zerocloud/core/temporary`
- Enter each of the following commands in the terminal (with a leading space), completing with the proper values

```
 export AWS_REGION=us-east-1

 export AWS_ACCESS_KEY_ID=<key here>
 export AWS_SECRET_ACCESS_KEY=<secret here>
 export AWS_ACCOUNT_ID=<ZEROCLOUD AWS ID here>
```

- Now open `github.com/tleyden/zerocloud/core/core.go` in a text editor
- Go to line 175 (and subsequent)
- Change `Email: "slv.balsan@gmail.com",` to BigDB's admin email.
- On line 179, change `AWSID:      859795398601,` to the BigDB's aws account.
- On line 180, change `ExternalID: "slavomir",` with "hithere"
- On line 192, change `Email:          "slv.balsan@gmail.com",` to BigDB's admin email.
- On line 198, change `Email:          "slavomir.balsan@gmail.com",` to an email address of a developer (who will be whitelisted to create leases).

# Run

Run `go run main.go` or use the [docker container](docs/docker/README.md)

Now, create an instance on BigDB's account without a ZeroCloudOwner tag.

BigDB's admin will receive an email (might not arrive immediately; using a sandbox mailgun account).

Try also with

`ZeroCloudOwner = nope`

`ZeroCloudOwner = someone@unauthorized.site`

`ZeroCloudOwner = dev@bigbd.io` (replacing `dev@bigbd.io` with what you wrote on line 198 in `core.go`)

`ZeroCloudOwner = admin@bigdb.io` (replacing `admin@bigdb.io` with BigDB's admin email you used on line 175 in `core.go`)

The relevant values are on line 49, 50 and 51 in `core.go`

```
ZCMaxPerOwner    = 2
s.Config.Lease.Duration = time.Minute * 1 // might be: time.Hour * 24 * 3 (i.e. 3 days)
ZCDefaultTruceDuration = time.Minute * 1 // the period before terminating non-approved instances
```

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

### GET /accounts/:account_id/cloudaccounts/:cloudaccount_id/regions

### PATCH /accounts/:account_id/cloudaccounts/:cloudaccount_id/regions

