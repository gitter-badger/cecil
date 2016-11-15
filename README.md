[![CircleCI](https://circleci.com/gh/tleyden/cecil.svg?style=svg&circle-token=95a33d3c7729a0423eb4acdf306a8ebf398647d3)](https://circleci.com/gh/tleyden/cecil)

# Cecil - The [C]ustodian for your [CL]oud

Cecil minimizes cost waste from **forgotten EC2 instances** on AWS by imposing a **strict leasing mechanism** on all EC2 instances that are started under it's watch.

Whenever a new EC2 instance is started in a Cecil-monitored AWS account, a lease will be created and assigned to the user that is declared in the `CecilOwner` tag, or assigned to the admin user if no owner is specified.  The owner will be notified by email before the lease expires, and the instance will be automatically shut down when the lease expires unless it is renewed via the REST API.  

Cecil has been designed specifically in order to meet the working requirements of [Couchbase](http://www.couchbase.com) for testing large scale distributed database functionality/performance on AWS clusters.

It was created to allow maximum developer/tester agility when it came to cloud resources -- in particular to avoid having to file IT tickets whenever you need an EC2 instance.  See the [backstory](docs/backstory.md) for more details on why it was created.

# Example Deployment

![](docs/architecture-flowcharts/system-overview-diagram.png)

* Acme.co represents **you** or **your project**.  It's assumed you already have an AWS account, possibly multiple.
* The Acme Cecil Service is expected to be run by **your IT department** using a separate AWS account dedicated for Cecil, and must be hosted somewhere that the REST endpoint will be publicly accessible.  It's not run by a 3rd party, because there is no third party.  Cecil is software, not a service, but it is packaged as a service for maximum decoupling.
* Although not shown, there can be more tenants than just the Acme.co tenant.  For example if Acme.co had a subsidiary called SubAcme.co, a new tenant could be created which had it's own AWS accounts for each of it's departments.

# User Interaction Example

![](docs/architecture-flowcharts/interaction-diagram.png)

1. A developer who has direct access to your AWS account (and they should!  see the [backstory](docs/backstory.md)) spins up one or more EC2 instances
1. This information propagates **across AWS account boundaries** into the Cecil process.  It starts out in the Acme.co AWS account, then gets pushed to the Cecil SQS queue which is running in the dedicated AWS account for Cecil.
1. Cecil emails the developer via the Mailgun API and informs them of the new instance and the lease that has been opened against it.  At this point, the developer can terminate the instance directly by clicking through a link in the email.
1. On Wednesday it will send the developer another email informing them that unless action is taken, their instance will be terminated in 24 hours.
1. On Thursday the developer clicks a link in the email and renews the lease.
1. On Saturday Cecil informs the developer they need to renew their lease or their instance will be terminated.
1. On Sunday, since the developer has taken no action, Cecil terminates the instance and informs the developer.

# Features

* Minimize forgotten EC2 instances by forcing users to refresh leases on resources still in use
* Configurable lease expiration times, number of renewals allowed, maximum number of leases per user
* Optionally require users to add owner tags to their EC2 instances
* Supports cross-account usage via STS role assumption.
* 100% Open Source (Apache2)

# Getting started: Installation and setup

The installation and configuration process has been broken up into separate documents:

1. [Install and configure the Cecil Service](docs/InstallCecilService.md)
   * Create a dedicated AWS account for Cecil, or re-use an existing separate AWS account
   * `go get` the code and build it -- or use a docker container
   * Set your `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables
   * Expose the REST API
   * Run the Cecil processs `./cecil`
1. [Configure one or more of your AWS accounts to be monitored by Cecil](docs/ConfigureAWSAccount.md)
   * Create a new account using the Cecil REST API
   * Create a new cloud account under that account, again with the Cecil REST API, which will return generated cloudformation templates 
   * Run `aws cloudformation create` on the cloudformation templates to setup CloudWatch Events and SNS in your AWS account 

## Related projects

* [Netflix Janitor Monkey](https://github.com/Netflix/SimianArmy/wiki/Janitor-Home)
* [Capital One Cloud Custodian](https://github.com/capitalone/cloud-custodian)

## Additional docs

* [Cecil backstory](docs/backstory.md)
* [Internal Developer Docs](docs/Dev.md) - useful if you want to contribute to Cecil development



