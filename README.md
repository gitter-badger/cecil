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

### Setup

- I login to zerocloud via github
     - Confirm email, so can send messages to admin user
- Add AWS account to be monitored:
     - Name the account (eg, Couchbase Mobile)
     - Setup AssumeRole stuff for ZeroCloud
     - v2: I choose a default lease period / policy (1 day / 1 week)
     - v2: I customize the required tag (default: owner)
     - v2: I customized the approve email domains (eg, *.couchbase.com)
     - v2: I assign owner tags to any existing resources (note: this requires permission to write tags to ec2 instances)

### Ongoing

- After resources spun up, contact owner and negotiate lease, up to one week by default, otherwise shut down
- Before lease is up, email owner and renegotiate lease, otherwise shut down

## Technical Pieces

- A Golang webserver that does github login and stores a user object in Sync Gateway  (could be a separate blog post)
- Documentation on how a Couchbase Admin can setup AssumeRole stuff for ZeroCloud (will be a script at some point — also blog post material)
- Ability for Couchbase User to respond to a lease expiration notification (API)
     - Renew lease
     - Shutdown immediately
- Event loop
     - Detector
          - Poll AWS and look for new resources that the system doesn’t know about yet
     - Notifier
          - View query on expired leases that haven’t been notified yet, send notification via Amazon SES
     - Reaper
          - If a lease has expired and has not been renewed by the deadline, then shutdown the resource

Note: the event loop should probably be inverted so that it’s turned into an event stream, and then there is code to react to all incoming events.  So the type of events the reactor would deal with are:

- ResourceChange
     - New resources added
     - Existing resources updated or deleted
- LeaseExpiringSoon
- LeaseExpired

This will help with testing to isolate components and events.  Also, when converting to CloudTrail Events, will make life easier.  In fact, for the ResourceChange event (most complicated by far!) it should be as close to CloudTrail Events as possilble.

Then one component will be in charge of polling AWS and generating ResourceChange events.

The LeaseExpiringSoon and LeaseExpired events could be driven by a scheduler, possibly even 3rd party externally managed service to outsource all of that logic/code.  It would just need a webhook callback with an opaque token.

## Prove out architecture

- Create private github repo called zerocloud
- Create another aws account called customer
- Setup the external access for another account (personal tleyden AWS) to delete ec2 instances on customer account, add to zerocloud README docs
- Play with CloudWatch Events on customer aws account
- Allow tleyden AWS to tap into customer aws account CloudWatch events, add to zerocloud README docs

## External Scheduler Providers

- https://github.com/gocraft/work  (background job runner w/ schedules)
- easycron.com (REST API)
- https://hook.io/cron
- http://dkron.io/

## Where should code live?

- Lambda
- Amazon ECS

## What should the DB be?

- SyncGw + CBS
  - Advantage: mobile app could rock!!
  - Advantage: dogfooding
- Amazon DynamoDB
  - Scalable / Standard / Hosted / Always Just Works
- Amazon RDS / Postgres

## Data Flow

- Stuff happens in EC2, instances are launched, etc

## ZeroCloud Business Objects

### Account

Each customer (eg, BigDB) will need to have an Customer object associated with it

### AWSAccount

A customer (eg, BigDB), might have 12 different AWS accounts.  They will need to do some one-time manual setup for each account they add.

### Users (eg, Cloud Admin or Test/Dev at BigDB)

- Role (admin | dev)

### AccountPolicy

This is something the Admin users maintain.  Applies to account as a whole.

- Default Lease Time
- Max Lease Time
- Allowed Email Address Pattern (eg, *.bigdb.com)
- Max Instances Allowed (eg, 10)
- Etc

### User Policy

Certain policies might need to be per-user and override account policy

- User ID
- Max Instances Allowed (eg, 50)

### Lease

When a new instance is detected by the injestor, if no lease is associated with it, then a new lease is created and assigned.  If the Owner tag is present on the instance, and maps to a known user, then the lease is assigned to that user.  Otherwise, assigned to an Admin user.

- CloudWatch Event -- the upstream cloudwatch event that triggered this lease to be created
- Owner User ID
- Expires At -- when it expires


## Open Source User Story

* Write a a config.yml with EC2 account credentials (expects to be different AWS acct, v2 could allow same acct)
* Run `docker compose up`
* Go to ip address of website
* Admin signup
* Any new EC2 instances will get default leases