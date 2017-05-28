[![CircleCI](https://circleci.com/gh/tleyden/cecil.svg?style=svg&circle-token=95a33d3c7729a0423eb4acdf306a8ebf398647d3)](https://circleci.com/gh/tleyden/cecil) [![Golang](https://img.shields.io/badge/Go-1.7-brightgreen.svg)](https://golang.org/) 

[![Launch Cecil](https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png)](https://console.aws.amazon.com/cloudformation/home?region=us-east-1#/stacks/new?stackName=CecilRootStack&templateURL=http://tleyden-misc.s3.amazonaws.com/cecil/cecil-root.template) [![Deploy to Docker Cloud](https://files.cloud.docker.com/images/deploy-to-dockercloud.svg)](https://cloud.docker.com/stack/deploy/?repo=https://github.com/tleyden/cecil) 


# 🤖 Cecil - an AWS EC2 instance garbage collector

Have you ever launched an EC2 instance and then forgotten to shut it down, and then it sat there for weeks racking up pointless AWS charges? 💰

Using AWS for **development and testing** is great, but it's all too easy to accumulate costly unused cloud resources.  In a larger organization, it takes effort to manually track down who owns these resources, and whether they are still in use or not.

Cecil was created to solve this problem using automation and a self-serve approach.

1. Whenever you start a new EC2 instance, Cecil assigns a lease to you for that instance and holds you accountable.
1. When the lease is about to expire, Cecil will notify you by email and give you a chance to renew it if you're still actually using it.
1. Unless you renew the lease, Cecil will automatically shut it down.

Cecil was developed at [Couchbase](http://www.couchbase.com) to reduce costs of development and testing use of AWS.  Couchbase developers have the freedom to spin up cloud resources without having to wait for approval by an IT department, which leads to high productivity, but at the risk of cost waste if resources are not cleaned up when they are no longer needed.  Cecil was created to minimize the cost waste without interfering with developer productivity.

Why another [Netflix Janitor Monkey](https://github.com/Netflix/SimianArmy/wiki/Janitor-Home)? 🙈 Janitor Monkey seemed a little tied to the Netflix production use case, rather than a developer sandbox use case.

# How it works

🛠 **One-time setup**

1. Install Cecil and configure it to monitor Cloudwatch Event streams of one or more AWS accounts
1. Configure Cecil with your {AWS key pair -> email address} mappings, so that new instance leases will get assigned to the right person based on the AWS key pair
1. Alternatively, inform your users that they need to add a special `CecilOwner` tag with their email address to all instances they launch, so the leases will be assigned to them

🚀 **Each time an EC2 instance is launched**

1. When a new instance is detected on the CloudWatch Event stream, a lease will be created and assigned to the person who launched it, or the admin user if the owner can't be identified.
1. When the lease is about to expire (3 days later by default), the owner is notified by email and given a chance to extend the lease.
1. Once the lease expires and is not extended, then the instance associated with the lease will get shutdown (terminated).

# Features

* ✅ Monitor multiple AWS accounts from a single Cecil instance via STS role assumption
* ✅ Stream based approach via Cloudwatch Events
* ✅ Treats Cloudformations and AutoScalingGroups as individual units
* ✅ Explicitly group instances into a single lease via a custom tag
* ✅ Assign leases based on SSH key or an owner tag
* ✅ Configurable lease expiration times, number of renewals allowed, maximum number of leases per user

# Typical workflow 

![](docs/architecture-flowcharts/interaction-diagram.png)

# Documentation + Resources

1. 📓 [Cecil Manual](http://tleyden-misc.s3.amazonaws.com/cecil/index.html) -- primary documentation, start here.  ([up-to-date-version](docs/index.asciidoc))
1. 📺 [Screencast: up and running (20 mins)](http://tleyden-misc.s3.amazonaws.com/cecil/CecilScreencastHD.mp4)
1. ⚙ [REST API reference](http://petstore.swagger.io/?url=https://gist.githubusercontent.com/tleyden/274e0605cb530deaf0c2c97f55644b00/raw/bdff0dccefee214f3ba588b0d49f8c70b52e9ada/cecil-api.yaml)
1. 📰 [Gitter Community](https://gitter.im/tleyden/cecil) - coming soon






