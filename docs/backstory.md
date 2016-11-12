
## Why Cecil?

I'm a developer at Couchbase and I was kind of shocked when our IT System Administrator handed me the keys to our AWS kingdom at one point, and basically said "Have a blast!".

I created automation around spinning up **very large** Couchbase database clusters on AWS.  I would usually run some functional and performance tests, and then shut them down.  Most of the time at least.  Sometimes I would forget to shut them down, and then I'd feel guilty that I just needlessly burned through a bunch of cash for the company.  Whoops!!  `¯\_(ツ)_/¯`

## Taming Dev Clouds By Hand

Every once in a while, I'd get a tap on the shoulder via email -- "Hey, you still need those 20 EC2 instances?  Not seeing much CPU usage on them".  And I'd go shut them down.

This process was repeated across the company.  Certain folks would keep an eye on our AWS usage and chase down forgetful developers and testers.

It worked, sorta.  The problem is that we have lots of AWS accounts at Couchbase, and it's just too much tedious work for people to keep tabs on all of them.  So no, it doesn't really work, it's a manual and highly error-prone process.  That started to get me thinking.

## Growing Anxiety About Spreadsheets

Through conversations with our IT department I learned that we may be facing a possible AWS lock-down one day.  Meaning that if you needed an EC2 instance, you need to file a Jira ticket with IT and they will record an entry in a spreadsheet to keep track of it and other details of the instance: how long you are planning to use it, what you are using it for, etc.

This freaked me out!  This would really hinder my ability to do awesome stuff in the Cloud.  This .. cannot .. happen ..

As I started talking to folks at other companies, I got even more freaked out, because this is a thing!  My sample set is far too small to draw any sweeping conclusions, but let's just say it was **not uncommon** for developers to need to file IT tickets when they wanted cloud resources.  

## I need a Netflix Janitor Monkey -- Or not

I started Googling for a solution and stumbled across a nifty looking little janitor monkey written in Java.  I started reading the docs and ... well .. they sucked.  (and no offense to Netflix, I'm sure they know and don't really care, especially since they rewrote the Chaos Monkey in Go)

And then I looked over the netflix janitor monkey feature set, and it only had about 1/3rd of what I wanted.  It had a concept of leases, but it made it *way* too easy for someone to grant themselves a permanent lease and then forget about.

## Cecil is born

I wanted something far more proactive.  I wanted to make the programmatic equivalent of a finance department intern who was **on your case** about every single EC2 instance that you fired up.

And I wanted to write it in Go.  And I wanted it to be perfect, and reliable, and re-usable, and well documented, and well tested ...



