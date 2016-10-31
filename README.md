[![CircleCI](https://circleci.com/gh/tleyden/zerocloud.svg?style=svg&circle-token=0b966949f6517187f0a2cece8aac8be59e0182a3)](https://circleci.com/gh/tleyden/zerocloud)

Cecil is a [C]ustodian for your [CL]oud.

Cecil is similar in spirit to [Netflix Janitor Monkey](http://techblog.netflix.com/2013/01/janitor-monkey-keeping-cloud-tidy-and.html).  It mops up AWS EC2 instances using a leasing mechanism, to make it as hard as possible for developers to spin up EC2 instances and forget about them and cause your AWS bill to unnecessarily bloat.  Because of it's cloud-mopping abilities, Cecil is also affectionately known as "Mopster".

It encourages the use of "owner tagging" to track which developers in your organization are responsible for which EC2 instances.  Any instances without owner tags will default to have a lease being assigned to the account admin.

# Flow

1. One-time setup process
1. A developer spins up a 10-node cluster to run some performance tests overnight
1. Cecil detects the new instances, sees the `Owner` EC2 instance tag that the developer (hopefully) added when launching the instances, and assigns the leases to the developer.  (or defaults to admin if no owner tag present)
1. After the configurable lease period has expired, assuming they are still up, Cecil sends and email to the lease owner asking if they want to renew the lease.
   - If the owner renews, then the lease will be extended and the instances will be left alone
   - If the owner fails to respond in time, then the instances will be terminated


# Packaging

To run Cecil you need to:

1. Setup up a Cecil "service", which requires it's own AWS account credentials.
1. Register as many AWS accounts + regions as you want, grouping them as separate "Tenants" of the Cecil service.
1. Connect each AWS account to the Cecil AWS account and service so that Cecil can monitor events and take actions.

# Get code

```
go get -t github.com/tleyden/zerocloud/...
```

# Cecil AWS Setup (AWS CLI)

This assumes you already have keys to access your root AWS account to create this stack.

```
aws cloudformation create-stack --stack-name "CecilRootStack" \
--template-body "file://./docs/cloudformation-templates/cecil-root.template" \
--capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
--region us-east-1
```

Create credentials (keys) for `CecilRootUser`:

```bash
$ aws iam create-access-key --user-name CecilRootUser
```

This will return something like

```json
{
    "AccessKey": {
        "SecretAccessKey": "je7MtGbClwBF/2Zp9Utk/h3yCo8nvbEXAMPLEKEY",
        "Status": "Active",
        "CreateDate": "2013-01-02T22:44:12.897Z",
        "UserName": "CecilRootUser",
        "AccessKeyId": "AKIAI44QH8DHBEXAMPLE"
    }
}
```

Alternatively, you can setup the stacks using the AWS web GUI instead of the CLI.


# Service setup

- Open a terminal tab/window
- cd to `github.com/tleyden/zerocloud/`
- Enter each of the following commands in the terminal (with a leading space), completing with the proper values

```
 export AWS_REGION=<fill in here>

 export AWS_ACCESS_KEY_ID=<fill in here>
 export AWS_SECRET_ACCESS_KEY=<fill in here>
 export AWS_ACCOUNT_ID=<fill in here>

 export MAILERDOMAIN=mg.zerocloud.co
 export MAILERAPIKEY=key-<fill in here>
 export MAILERPUBLICAPIKEY=pubkey-<fill in here>
```

You can find the mailer (Mailgun) API keys at [mailgun.com/app/account/security](https://mailgun.com/app/account/security)  For `MAILERAPIKEY` use the value in `Active API Key` and for `MAILERPUBLICAPIKEY` use `Email Validation Key`


# Run

Run `go run main.go` or use the [docker container](docs/docker/README.md)


## Endpoint usage examples

## Create account

```bash
curl -X POST \
-H "Cache-Control: no-cache" \
-d '{
	"email":"example@example.com",
	"name":"John",
	"surname":"Smith"
}' \
"http://127.0.0.1:8080/accounts"
```

Response:

```json
{
  "email": "example@example.com",
  "account_id": 1,
  "response": "An email has been sent to the specified address with a verification token and instructions.",
  "verified": false
}
```
You will receive an email with a vefication code.


## Verify account and get api token

```bash
curl -X POST \
-H "Cache-Control: no-cache" \
-d '{"verification_token":"0d78a4e0-9922-4b55-93d7-5adfd0f589be7b9a0fa6-c5bc-4991-9f8e-b8bdbc429343322e200c-ab6c-4189-9e81-453ab0b34d56"}' \
"http://0.0.0.0:8080/accounts/1/api_token"
```

Response:

```json
{
  "account_id": 1,
  "api_token": "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoxLCJpYXQiOjE0Nzc0MDg1MzJ9.tr5Ark32AIQyYfM4AnQuC4I6ROQsP7PUSuz6hMR5EOMjDEHQ74A6JKxxR08OkdIgA8NCLw7a8oUyKqDc4XalrQKIq--FCZzf47dswMsJNjtwZPPFTX1hLjhsvuuQiVvtm39jjJL_t4l-ICa0oKX8nrJNGmB5epVR3KMPySlXXShUx-vc77P6My4WOpLIZV8lyeVlobRvLxfCKyXtqxKSRiu0-oJ1rXxCDkcGVvGFMk8vVjYeXDHM4dITuoweb_1TVHxRelePKtpuw5BEyakYXJmLI7m3eQYk8Pv9sBpviS2KhGjq9qPG6kweopGNCuYsrF0L1x5YZ3jWcBL0-KpK2g",
  "email": "example@example.com",
  "verified": true
}
```

Use the api token to manage your account.

## Add CloudAccount

```bash
curl -X POST \
-H "Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoxLCJpYXQiOjE0Nzc0MDg1MzJ9.tr5Ark32AIQyYfM4AnQuC4I6ROQsP7PUSuz6hMR5EOMjDEHQ74A6JKxxR08OkdIgA8NCLw7a8oUyKqDc4XalrQKIq--FCZzf47dswMsJNjtwZPPFTX1hLjhsvuuQiVvtm39jjJL_t4l-ICa0oKX8nrJNGmB5epVR3KMPySlXXShUx-vc77P6My4WOpLIZV8lyeVlobRvLxfCKyXtqxKSRiu0-oJ1rXxCDkcGVvGFMk8vVjYeXDHM4dITuoweb_1TVHxRelePKtpuw5BEyakYXJmLI7m3eQYk8Pv9sBpviS2KhGjq9qPG6kweopGNCuYsrF0L1x5YZ3jWcBL0-KpK2g" \
-H "Cache-Control: no-cache" \
-d '{
	"aws_id":"0123456789"
}' \
"http://0.0.0.0:8080/accounts/1/cloudaccounts"
```

Response:

```json
{
  "aws_id": "0123456789",
  "cloudaccount_id": 1,
  "initial_setup_cloudformation_url": "/accounts/1/cloudaccounts/1/cecil-aws-initial-setup.template",
  "region_setup_cloudformation_url": "/accounts/1/cloudaccounts/1/cecil-aws-region-setup.template"
}
```

Before this cloudaccount is active, you need to setup the Cecil stacks on your AWS account:

1. The first stack is the **initial stack**. It's a one-time only setup, and will be valid for the whole AWS account.
2.  The second stack is the **region stack**. This stack is to be created on each region you want to monitor with Cecil.

To setup the stacks, download them from the urls provided in this response. And then use AWS cli or AWS web gui to set them up.


## Cloudformation template for initial setup

First download it:

```bash
curl -X GET \
-H "Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoxLCJpYXQiOjE0Nzc0MDg1MzJ9.tr5Ark32AIQyYfM4AnQuC4I6ROQsP7PUSuz6hMR5EOMjDEHQ74A6JKxxR08OkdIgA8NCLw7a8oUyKqDc4XalrQKIq--FCZzf47dswMsJNjtwZPPFTX1hLjhsvuuQiVvtm39jjJL_t4l-ICa0oKX8nrJNGmB5epVR3KMPySlXXShUx-vc77P6My4WOpLIZV8lyeVlobRvLxfCKyXtqxKSRiu0-oJ1rXxCDkcGVvGFMk8vVjYeXDHM4dITuoweb_1TVHxRelePKtpuw5BEyakYXJmLI7m3eQYk8Pv9sBpviS2KhGjq9qPG6kweopGNCuYsrF0L1x5YZ3jWcBL0-KpK2g" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/cecil-aws-initial-setup.template" > cecil-aws-initial-setup.template
```

Then install it:

```bash
export AWS_ACCESS_KEY_ID=<access key id of cecil user aws account>
export AWS_SECRET_ACCESS_KEY=<access key id of cecil user aws account>
aws cloudformation create-stack --stack-name "CecilStack" --template-body "file://cecil-aws-initial-setup.template" --region us-east-1 --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM
```

Or alternatively you can upload this in the Cloudformation section of the AWS web UI.

## Cloudformation template for REGION setup

```bash
curl -X GET \
-H "Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoxLCJpYXQiOjE0Nzc0MDg1MzJ9.tr5Ark32AIQyYfM4AnQuC4I6ROQsP7PUSuz6hMR5EOMjDEHQ74A6JKxxR08OkdIgA8NCLw7a8oUyKqDc4XalrQKIq--FCZzf47dswMsJNjtwZPPFTX1hLjhsvuuQiVvtm39jjJL_t4l-ICa0oKX8nrJNGmB5epVR3KMPySlXXShUx-vc77P6My4WOpLIZV8lyeVlobRvLxfCKyXtqxKSRiu0-oJ1rXxCDkcGVvGFMk8vVjYeXDHM4dITuoweb_1TVHxRelePKtpuw5BEyakYXJmLI7m3eQYk8Pv9sBpviS2KhGjq9qPG6kweopGNCuYsrF0L1x5YZ3jWcBL0-KpK2g" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/cecil-aws-region-setup.template" > cecil-aws-region-setup.template
```

Then install it:

```bash
aws cloudformation create-stack --stack-name "CecilUSEastStack" --template-body "file://cecil-aws-region-setup.template" --region us-east-1
```

After this has been successfully setup by AWS, you will receive an email from Cecil.

## Add email to owner tag whitelist

```bash
curl -X POST \
-H "Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoxLCJpYXQiOjE0Nzc0MDg1MzJ9.tr5Ark32AIQyYfM4AnQuC4I6ROQsP7PUSuz6hMR5EOMjDEHQ74A6JKxxR08OkdIgA8NCLw7a8oUyKqDc4XalrQKIq--FCZzf47dswMsJNjtwZPPFTX1hLjhsvuuQiVvtm39jjJL_t4l-ICa0oKX8nrJNGmB5epVR3KMPySlXXShUx-vc77P6My4WOpLIZV8lyeVlobRvLxfCKyXtqxKSRiu0-oJ1rXxCDkcGVvGFMk8vVjYeXDHM4dITuoweb_1TVHxRelePKtpuw5BEyakYXJmLI7m3eQYk8Pv9sBpviS2KhGjq9qPG6kweopGNCuYsrF0L1x5YZ3jWcBL0-KpK2g" \
-H "Cache-Control: no-cache" \
-d '{"email":"someone.legit@example.com"}' \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/owners"
```

Response:

```json
{
  "message": "owner added successfully to whitelist"
}
```

## Additional docs

* [Listing of code directories/files and their purposes](docs/CodeInventory.md)