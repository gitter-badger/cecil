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


# BigDB AWS Setup

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

Run `go run temp-service.go`.

Now, create an instance on BigDB's account without a ZeroCloudOwner tag.

BigDB's admin will receive an email (might not arrive immediately; using a sandbox mailgun account).

Try also with

`ZeroCloudOwner = nope`

`ZeroCloudOwner = someone@unauthorized.site`

`ZeroCloudOwner = dev@bigbd.io` (replacing `dev@bigbd.io` with what you wrote on line 198 in `core.go`)

`ZeroCloudOwner = admin@bigdb.io` (replacing `admin@bigdb.io` with BigDB's admin email you used on line 175 in `core.go`)

The relevant values are on line 49, 50 and 51 in `core.go`

```
ZCMaxLeasesPerOwner    = 2
ZCDefaultLeaseDuration = time.Minute * 1 // might be: time.Hour * 24 * 3 (i.e. 3 days)
ZCDefaultTruceDuration = time.Minute * 1 // the period before terminating non-approved instances
```