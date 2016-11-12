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

