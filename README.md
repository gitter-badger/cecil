[![CircleCI](https://circleci.com/gh/tleyden/cecil.svg?style=svg&circle-token=95a33d3c7729a0423eb4acdf306a8ebf398647d3)](https://circleci.com/gh/tleyden/cecil)

# Cecil - The [C]ustodian for your [CL]oud 

Cecil is an EC2 instance garbage collector, similar to [Netflix Janitor Monkey](https://github.com/Netflix/SimianArmy/wiki/Janitor-Home), geared towards **development and testing** use cases of AWS.  It works by imposing a **leasing mechanism** on all instances started under it's watch, and requires users to continually extend leases on instances in order to prevent them from being garbage collected.

Cecil was developed and is in use at [Couchbase](http://www.couchbase.com) to facilitate performance testing of the Couchbase distributed NoSQL database on AWS. See the [backstory](docs/backstory.md) for more details on why it was created.

# Features

* Monitor multiple AWS accounts
* Cross-account usage via STS role assumption [System Diagram](docs/architecture-flowcharts/system-overview-diagram.png)
* Tag-based instance grouping mechanism
* Recognizes Cloudformation and AutoScalingGroup instance grouping mechanisms
* Configurable lease expiration times, number of renewals allowed, maximum number of leases per user
* Slack integration

# How it works

![](docs/architecture-flowcharts/interaction-diagram.png)

1. A developer who has direct access to your AWS account spins up one or more EC2 instances
1. This information propagates **across AWS account boundaries** into the Cecil process.  It starts out in the Acme.co AWS account, then gets pushed to the Cecil SQS queue which is running in the dedicated AWS account for Cecil.
1. Cecil emails the developer via the Mailgun API and informs them of the new instance and the lease that has been opened against it.  At this point, the developer can terminate the instance directly by clicking through a link in the email.
1. On Wednesday it will send the developer another email informing them that unless action is taken, their instance will be terminated in 24 hours.
1. On Thursday the developer clicks a link in the email and renews the lease.
1. On Saturday Cecil informs the developer they need to renew their lease or their instance will be terminated.
1. On Sunday, since the developer has taken no action, Cecil terminates the instance and informs the developer.

# Single Account: Installation and setup

If you only have a single AWS account you want to monitor, you can run Cecil in the same account you want to monitor.

TODO

# Cross Account: Installation and setup

If you want to monitor multiple AWS accounts, you probably want to dedicate one of the accounts to run Cecil in (or create a new one), and run Cecil in that account and configure it to monitor the rest of the AWS accounts.

The installation and configuration process has been broken up into separate documents:

1. [Install and configure the Cecil Service](docs/InstallCecilService.md)
   * Create a dedicated AWS account for Cecil, or re-use an existing separate AWS account
   * Build the code
   * Set environment variables + config file
   * Run the Cecil processs
1. [Configure one or more of your AWS accounts to be monitored by Cecil](docs/ConfigureAWSAccount.md)
   * Create a new account using the Cecil REST API
   * Create a new cloud account, which will return generated cloudformation templates 
   * Run `aws cloudformation create` on the cloudformation templates to setup CloudWatch Events and SNS in your AWS account 


# Documentation Index

| Name  | Category | Description | 
| ------------- | ------------- | ------------- |
| [InstallCecilService.md](docs/InstallCecilService.md)  | Setup  | Install and configure Cecil service |
| [ConfigureAWSAccount.md](docs/ConfigureAWSAccount.md)  | Setup  | Configure Tenant(s) + AWS account(s) |
| [DeployToCloud.md](docs/DeployToCloud.md)  | Setup  | Deploy to various IaaS/PaaS/CaaS Providers |
| [StaticConfig.md](docs/StaticConfig.md)  | Setup  | Optional Static Config options for MailGun and JWT keypair |
| [Api.md](docs/Api.md)  | API  | Ad-hoc REST API docs |
| [swagger.json](goa/swagger/swagger.json) + [.yaml](goa/swagger/swagger.yaml)  | API  | Swagger / OpenAPI REST docs (auto-generated) |
| [postman/README.md](docs/postman/README.md) | API  | Using the GUI Postman REST API client |







## Related projects

* [Netflix Janitor Monkey](https://github.com/Netflix/SimianArmy/wiki/Janitor-Home)
    * [NJM vs Cecil](docs/backstory.md)
* [Capital One Cloud Custodian](https://github.com/capitalone/cloud-custodian)

## Additional docs

* [Cecil backstory](docs/backstory.md)
* [Internal Developer Docs](docs/Dev.md) - useful if you want to contribute to Cecil development


