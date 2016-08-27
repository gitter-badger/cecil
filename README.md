## Mission

Allow your devs and testers unfettered access to create AWS instances yet make it impossible for them to forget about unused resources and let your AWS bill spin out of control.

You define your policies, we enforce them.

## Clone repo

```
go get -t github.com/tleyden/zerocloud/...
```

## Setup steps

### AWS Account configuration

Currently by hand, but ideally should be a CloudFormation script

* Add IAM Role to BigDB Customer which has
   * Attached policy of AmazonEC2FullAccess (for now -- later need to minimize access)
   * Trust relationship with the AWS account number of the ZeroCloud AWS account and the externalid (eg, bigdb)
* Add CloudWatch Event Rule which pushes all CloudWatch Events to an SNS topic
* Add a subscription to the SNS topic which pushes to an SQS queue in the ZeroCloud AWS account
   * In my case, I think I did this from the ZeroCloud AWS SQS UI because it was easiest
   * Ideally this could be done when the customer runs the Cloudformation
   * If not, the SNS topic ARN could be given to ZeroCloud and it could add the subscription on it's end
* Needs outputs:
  * IAM Role ARN to give to ZeroCloud
  * (Maybe) SNS topic ARN to give to ZeroCloud

### Run REST server

```
./zerocloud
```

If that doesn't work, you might need to add `$GOPATH/bin` to your `PATH`

You can also run it via:

```
cd $GOPATH/src/github.com/tleyden/zerocloud
go run main.go
```

### Run CloudWatch Event SQS poller

This process polls the ZeroCloud SQS for new CloudWatch Events from any customers, adds the instnace tags, and then pushes them to the ZeroCloud REST API

```
cd $GOPATH/src/github.com/tleyden/zerocloud/cli
go run main.go poll_cloudevent_queue --help
```

See the `--help` for parameter instructions

### Add ZeroCloud Account for BigDB Customer via REST/CLI

```
zerocloud-cli create account --payload '{"lease_expires_in": 3, "lease_expires_in_units": "days", "name": "BigDB"}'
```

### Add AWS Account for BigDB Customer via REST/CLI

```
zerocloud-cli create cloudaccount --accountID 1 --payload '{"assume_role_arn": "arn:aws:iam::788612350743:role/ZeroCloud", "assume_role_external_id": "bigdb", "cloudprovider": "AWS", "name": "BigDB.cos perf testing AWS account", "upstream_account_id": "98798079879"}'
```

### Verify

* Spin up an EC2 instance in the BigDB Customer AWS account
* A new lease should be created in the zerocloud.db sqlite file (you can view with http://sqlitebrowser.org/)

NOTE: Rather than spinning up actual instances it's also possible to create CloudWatch Events directly via:

```
zerocloud-cli create cloudevent --help
```

## MVP

### Userstory 1: Auto-terminate instances after three days

* Start ZeroCloud REST + CloudWatch Event SQS poller
* Add ZeroCloud Account for BigDB Customer via REST/CLI
* Add AWS Account for BigDB Customer via REST/CLI
* Spin up an EC2 instance in the BigDB Customer AWS account
* Three days later, ZeroCloud will shutdown the instance

### Userstory 2: Add email confirmation

Same as Userstory 1, but:

* Two days after instance is spun up, the BigDB admin receives an email with the following:
    * Your instance (i-dafsaf) will but shutdown in 24 hours unless you extend your lease
    * Link to REST API URL which will extend the lease another three days
    * Link to REST API URL which will shutdown the instance immmediately
* By default, one day later, ZeroCloud will shutdown the instance

### Userstory 3: Allow owner tags

Same as Userstory 2, but if the instance has a ZeroCloudOwner tag with an email address, then that user will be contacted regarding the lease expiry rather than the BigDB admin

## External Scheduler Providers

- https://github.com/gocraft/work  (background job runner w/ schedules)
- easycron.com (REST API)
- https://hook.io/cron
- http://dkron.io/




