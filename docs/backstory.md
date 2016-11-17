
## Why Cecil?

At Couchbase we have the following requirements:

* Allow developers to spin up large, emphemeral distributed database clusters on AWS for adhoc testing *without* needing to request resources from IT -- developers are given direct AWS access with high VM limits
* Minimize cost waste associated with AWS

The problem is that it's *way too easy* to forget about EC2 instances for days, weeks, or even months.  Whoops!!  `¯\_(ツ)_/¯`

While it's possible for an IT department to manually audit the instances and chase down the developers who created them, it's an extremely error prone process.  As the number of both developers and AWS accounts at Couchbase started to increase, managing this by hand became impractical.

## Avoiding an IT bottleneck

The knee jerk reaction to minimize the cost of AWS garbage might be to force developers to go through the IT department to create EC2 instances, and leave it up to the IT department to shut them down, but now you have other problems:

* This will slow developers down
* This doesn't play well with automation
* Since there is so much friction to getting EC2 instances, developers will *hoard* the ones they've been given, driving up costs

So in order to try to meet the requirements *without* the above drawbacks, Cecil was created.

## Cecil vs Netflix Janitor Monkey

Why build Cecil when Netflix Janitor Monkey already exists?

Several reasons:

1. Netflix Janitor Monkey (NJM) is not well documented
1. NJM allows for perennial leases which makes it easy for things to slip through the cracks and accumulate cost
1. NJM appears to maybe be deprecated, it does not seem to have survived the rewrite to Go
1. NJM uses perdiodic polling instead of subscribing to CloudWatch Events, which does not exist at the time it was being built
1. The NJM codebase is a bit unwieldy, as echoed on a recent [GoTime.fm](https://changelog.com/gotime/9) podcast interview w/ Scott Manfield from Netflix

## Cecil vs Capital One Cloud Custodian

Capital One Cloud Custodian is a new tool to pop on the scene .. haven't checked it thoroughly yet.





