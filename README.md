[![CircleCI](https://circleci.com/gh/tleyden/cecil.svg?style=svg&circle-token=95a33d3c7729a0423eb4acdf306a8ebf398647d3)](https://circleci.com/gh/tleyden/cecil) [![ViewTheAPI](https://img.shields.io/badge/REST%20API-latest-brightgreen.svg)](http://cecil.viewtheapi.io)  [![Golang](https://img.shields.io/badge/Go-1.8-brightgreen.svg)](https://golang.org/) [![Apache 2](https://img.shields.io/badge/license-Apache%202-blue.svg )](https://www.apache.org/licenses/LICENSE-2.0) 

[![Launch Cecil](https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png)](https://console.aws.amazon.com/cloudformation/home?region=us-east-1#/stacks/new?stackName=CecilRootStack&templateURL=http://tleyden-misc.s3.amazonaws.com/cecil/cecil-root.template) [![Deploy to Docker Cloud](https://files.cloud.docker.com/images/deploy-to-dockercloud.svg)](https://cloud.docker.com/stack/deploy/?repo=https://github.com/tleyden/cecil) 


# 🤖 Cecil - an AWS EC2 instance garbage collector

Cecil is a monitoring and clean-up tool designed to make it as hard as possible to let EC2 instance garbage accumulate and rack up pointless AWS charges. 💰 It's geared towards doing **development and testing** in AWS.

It uses automation and a self-serve approach to minimizing EC2 waste:

1. Whenever you start a new EC2 instance, Cecil assigns a lease to you for that instance and notifies you via email.
1. When the lease is about to expire, Cecil will notify you by email and give you a chance to renew it if you're still actually using it.
1. Unless you renew the lease, Cecil will automatically terminate the EC2 instance.

Cecil was developed at [Couchbase](http://www.couchbase.com) [![Couchbase](docs/images/couchbase.png)](http://www.couchbase.com) to help control AWS costs related to performance testing of it's [open source distributed NoSQL database](https://developer.couchbase.com/documentation/server/current/architecture/architecture-intro.html).


# Features

* ✅ Monitor multiple AWS accounts from a single Cecil instance via STS role assumption
* ✅ Stream based approach via Cloudwatch Events
* ✅ Treats Cloudformations and AutoScalingGroups as individual units
* ✅ Explicitly group instances into a single lease via a custom tag
* ✅ Assign leases based on SSH key or an owner tag
* ✅ Configurable lease expiration times, number of renewals allowed, maximum number of leases per user

# How it works

🛠 **One-time setup**

1. Install the single-binary Cecil server process somewhere. 
1. Configure Cecil to monitor Cloudwatch Event streams of one or more AWS accounts by deploying a Cloudformation template.
1. Create an account administrator user with an email address via the Cecil REST API.
1. Configure Cecil via the REST API to tell it how to associate newly launched EC2 instances with owners by adding mappings between AWS Key Pairs and Email addresses.

🚀 **Each time an EC2 instance is launched**

1. Users can add a `CecilOwner` tag that contains their email address when they launch instances, which means they will be assigned the lease.  (not needed if the mappings between AWS Key Pairs and Email addresses are provided in previous step)
1. When a new instance is detected on the CloudWatch Event stream, a lease will be created and assigned to the person who launched it, or the admin user if the owner can't be identified.
1. When the lease is about to expire, the owner is notified by email twice and given a chance to extend the lease by clicking a link.
1. If the lease isn't extended and eventually expires, then the instance associated with the lease will get terminated.

![](docs/images/interaction-diagram.png)

# Documentation + Resources

1. 📓 [Cecil Manual](http://tleyden-misc.s3.amazonaws.com/cecil/index.html) -- primary documentation, start here.  ([up-to-date-version here](docs/index.asciidoc), but missing some formatting)
1. 📺 [Screencast: up and running (20 mins)](http://tleyden-misc.s3.amazonaws.com/cecil/CecilScreencastHD.mp4)
1. ⚙ [REST API reference](http://cecil.viewtheapi.io)
1. 📰 [Gitter Community](https://gitter.im/tleyden/cecil) - coming soon






