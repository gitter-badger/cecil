[![CircleCI](https://circleci.com/gh/tleyden/cecil.svg?style=svg&circle-token=95a33d3c7729a0423eb4acdf306a8ebf398647d3)](https://circleci.com/gh/tleyden/cecil) [![Docs](https://img.shields.io/badge/Docs-latest-brightgreen.svg)](http://cecil-assets.s3-website-us-east-1.amazonaws.com/asciidoc/) [![ViewTheAPI](https://img.shields.io/badge/REST%20API-latest-brightgreen.svg)](http://cecil-assets.s3-website-us-east-1.amazonaws.com/swagger/)  [![Golang](https://img.shields.io/badge/Go-1.8-blue.svg)](https://golang.org/) [![Apache 2](https://img.shields.io/badge/license-Apache%202-blue.svg )](https://www.apache.org/licenses/LICENSE-2.0) [![Screencast](https://img.shields.io/badge/screencast-20mins-yellow.svg )](http://cecil-assets.s3.amazonaws.com/screencast/CecilScreencastHD.mp4) 

[![Launch Cecil](https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png)](https://console.aws.amazon.com/cloudformation/home?region=us-east-1#/stacks/new?stackName=CecilRootStack&templateURL=http://cecil-assets.s3.amazonaws.com/cloudformation/cecil-root.template) [![Deploy to Docker Cloud](https://files.cloud.docker.com/images/deploy-to-dockercloud.svg)](https://cloud.docker.com/stack/deploy/?repo=https://github.com/tleyden/cecil) 


# 🤖 Cecil - an AWS EC2 instance garbage collector

Cecil is a monitoring and clean-up tool designed to make it as hard as possible to let EC2 instance garbage accumulate and rack up pointless AWS charges.  It's geared towards doing **development and testing** in AWS.

It uses automation and a self-serve approach to minimizing EC2 waste:

1. Whenever you start a new EC2 instance, Cecil assigns a short-term lease to you (3 days by default) for that instance and notifies you via email.
1. When the lease is about to expire, Cecil will notify you by email and give you a chance to renew the lease if you're still actually using it.
1. Unless you renew the lease, Cecil will automatically terminate the EC2 instance.

Cecil was developed at [Couchbase](http://www.couchbase.com) [![Couchbase](docs/images/couchbase.png)](http://www.couchbase.com) to help control AWS costs related to large-scale performance testing of it's [open source distributed NoSQL database](https://developer.couchbase.com/documentation/server/current/architecture/architecture-intro.html).


# Features

* ✅ Configurable lease expiration times, number of renewals allowed, maximum number of leases per user
* ✅ Assign leases based on SSH keypair or EC2 tag
* ✅ Explicitly group instances into a single lease via a custom tag
* ✅ Monitor multiple AWS accounts from a single Cecil instance via STS role assumption
* ✅ Low latency Stream based approach via Cloudwatch Events
* ✅ Treats Cloudformations and AutoScalingGroups as individual units


# Roadmap

* 💡 Offhours support -- stop and restart EBS-backed instances based on a schedule
* 💡 Slackbot / Hipchat bot support ([work in progress](https://github.com/tleyden/cecil/blob/master/docs/index.asciidoc#slack-integration))
* 💡 Usage Reports ([work in progress](https://github.com/tleyden/cecil/issues/122)) 
* 💡 [Add a feature request!](https://github.com/tleyden/cecil/issues/new)

# Getting started

To learn more, start reading about the [Design Philosophy](http://cecil-assets.s3-website-us-east-1.amazonaws.com/asciidoc/index.html#_cecil_design) or [How it Works](http://cecil-assets.s3-website-us-east-1.amazonaws.com/asciidoc/index.html#_cecil_for_administrators#_how_it_works).

Or if you just want to get up and running, jump right to [Cecil for Administrators](http://cecil-assets.s3-website-us-east-1.amazonaws.com/asciidoc/index.html#_cecil_for_administrators) to install it and have it monitor your AWS account(s)

# Documentation

1. 📓 [Cecil Manual](http://cecil-assets.s3-website-us-east-1.amazonaws.com/asciidoc/) -- primary documentation, start here.  ([up-to-date-version here](docs/index.asciidoc), but missing some formatting)
1. ⚙ [REST API reference](http://cecil-assets.s3-website-us-east-1.amazonaws.com/swagger/)
1. 📺 [Screencast: up and running (20 mins)](http://cecil-assets.s3.amazonaws.com/screencast/CecilScreencastHD.mp4)

# Community

1. 📰 [Gitter Community](https://gitter.im/tleyden/cecil) - coming soon
1. 📮 [Google Group]() - coming soon
1. 🐦 [Twitter feed]() - coming soon


