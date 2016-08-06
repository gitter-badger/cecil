## Mission

Allow your devs and testers unfettered access to create AWS instances yet make it impossible for them to forget about unused resources and let your AWS bill spin out of control.

You define your policies, we enforce them.

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

## Open Questions

- Polling AWS vs setting up CloudTrail Events
     - Idea: polling by default, but later add option for CloudTrail which might require more work / permissions on customer’s part

