[![Stories in Ready](https://badge.waffle.io/tleyden/zerocloud.png?label=ready&title=Ready)](https://waffle.io/tleyden/zerocloud)

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

### Create zerocloud stack with cloudformation template on BigDB's AWS account

```
aws cloudformation create-stack --stack-name "ZeroCloudStack" \
--template-body "file://path/to/zerocloud.template" \
--parameters ParameterKey=ZeroCloudAWSID,ParameterValue=123456789101 \
ParameterKey=ZCAccountID,ParameterValue=abcdefg1234der456ghijkl6789
```

For more options, refer to http://docs.aws.amazon.com/cli/latest/reference/cloudformation/create-stack.html

### Verify

* Spin up an EC2 instance in the BigDB Customer AWS account
* A new lease should be created in the zerocloud.db sqlite file (you can view with http://sqlitebrowser.org/)

NOTE: Rather than spinning up actual instances it's also possible to create CloudWatch Events directly via:

```
zerocloud-cli create cloudevent --help
```

## Development steps

## Regenerate code

After making changes to `design/design.go` or `design/models.go`, run:

```
./goa.sh
```

## Regenerating existing endpoints

This is a bit awkward.  If you have an existing file with endpoint code like `account.go` and you make changes to `design/design.go` to change any endpoints, do the following:

* Move the endpoint file(s) `account.go` to `account_old.go`
* Regenerate endpoint code by running `./goa.sh`
* Move the "implementation" calls from `account_old.go` to `account.go` (the actual implementation code lives in `account_impl.go`)
* Add any new implementation in `account_impl.go`
* Remove `account_old.go`

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

## Additional Documentation

### Swagger REST API spec

* Go to http://editor.swagger.io/
* Open url `https://raw.githubusercontent.com/tleyden/zerocloud/master/swagger/swagger.json?token=AASHrP6C1ju3bIx6xTr1mLKX4HKwP98Zks5Xyx8xwA%3D%3D`

## Project directories

* **app** generated by goa, contains non-model structs
* **cli** contains non-goa CLI tools, currently just the cloudevent_poller CLI
* **client** generated by goa, contains autogenerated REST client
* **cloudevent_poller** contains code for the cloudevent_poller
* **design** the main files that tell goa how to generate REST api (design.go) and Gorma models (models.go)
* **docs** additional project documentation / diagrams
* **models** generated by goa/gorma, contains gorm models.  Files with "extra" in name are custom written and not generated
* **swagger** generated swagger docs
* **tool** generated by goa, contains CLI wrappers around REST API for object CRUD


## Notes

### External Scheduler Providers

- https://github.com/gocraft/work  (background job runner w/ schedules)
- easycron.com (REST API)
- https://hook.io/cron
- http://dkron.io/




