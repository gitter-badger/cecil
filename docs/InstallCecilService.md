Here is the system overview of what a Cecil deploy looks like:

![](architecture-flowcharts/system-overview-diagram.png)

In this document you will setup the *right hand side*.  If your company name is "Acme", you would likely call this the "Acme Cecil Service"

# Get code

```
go get -t github.com/tleyden/cecil/...
```

If this completes without errors, you will have a new binary in `$GOPATH/bin/cecil`

# Choose or create an AWS account for Cecil

If you have an existing account that is **separate from the AWS account(s) you want to monitor**, use that.  Otherwise, create a brand new AWS account.

You will need to have an Access Key defined for the root user

It might also be possible to create a new IAM user in the account that has the built-in AdministratorAccess policy attached to it.  This example assumes you have created an Access Key for the AWS Account Root User.

At this point you should have the following:

Description | AWS Account ID        | AWS_KEY           | AWS_SECRET_KEY |  Root/IAM | Attached Policies 
| ------------- |:-------------:|:-----:|:-----:|:-----:|:-----:|
| Cecil AWS account Root User | 193822812427      | AKIAEXAMPLEWAGRHKOMEWQ | ********** | Root | N/A 


# Cecil AWS Setup (AWS CLI)

This assumes you already have keys to access your root AWS account to create this stack.

```
$ export AWS_ACCESS_KEY_ID=AKIAEXAMPLEWAGRHKOMEWQ 
$ export AWS_SECRET_ACCESS_KEY=***** 
$ aws cloudformation create-stack --stack-name "CecilRootStack" \
--template-body "file://./docs/cloudformation-templates/cecil-root.template" \
--capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
--region us-east-1
```

You should see output similar to:

```
{
    "StackId": "arn:aws:cloudformation:us-east-1:193822812427:stack/CecilRootStack/fff31310-ab37-11e6-94ba-50d5cafe7636"
}
```

Create credentials (keys) for `CecilRootUser`:

```bash
$ export AWS_ACCESS_KEY_ID=AKIAEXAMPLEWAGRHKOMEWQ 
$ export AWS_SECRET_ACCESS_KEY=***** 
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
        "AccessKeyId": "AKIAIEXAMPLERQ4U4N67LE7A"
    }
}
```

At this point you should have the following:

Description | AWS Account ID        | AWS_KEY           | AWS_SECRET_KEY |  Root/IAM | Attached Policies 
| ------------- |:-------------:|:-----:|:-----:|:-----:|:-----:|
| Cecil AWS account Root User | 193822812427      | AKIAEXAMPLEWAGRHKOMEWQ | ********** | Root | N/A
| Cecil AWS account Cecil Root User | 193822812427      | AKIAIEXAMPLERQ4U4N67LE7A | ********** | IAM: CecilRootUser | allowassumerole,giveaccesstoqueueonly |  



Alternatively, you can setup the stacks using the AWS web GUI instead of the CLI.

1. Go to [CloudFormation Console](https://console.aws.amazon.com/cloudformation/home) to create the stack using the `docs/cloudformation-templates/cecil-root.template` CloudFormation Template.

1. And to [IAM Console](https://console.aws.amazon.com/iam/home?#/users/CecilRootUser) to create and download credentials for Cecil.

# Run

Set env variables for the `CecilRootUser` AWS Access Key:

```
$ export AWS_ACCESS_KEY_ID=AKIAIEXAMPLERQ4U4N67LE7A 
$ export AWS_SECRET_ACCESS_KEY=***** 
$ export AWS_REGION=us-east-1 
$ export AWS_ACCOUNT_ID=193822812427 
```

Run cecil:

```
$ cecil
```

Now you are ready to [Monitor AWS accounts](ConfigureAWSAccount.md)

Alternatively, you can run Cecil in a [docker container](docs/docker/README.md)



# Customize MailGun Settings (optional)

```
$ export MAILERDOMAIN=mg.yourdomain.co
$ export MAILERAPIKEY=key-<fill in here>
$ export MAILERPUBLICAPIKEY=pubkey-<fill in here>
```

You can find the mailer (Mailgun) API keys at [mailgun.com/app/account/security](https://mailgun.com/app/account/security)  For `MAILERAPIKEY` use the value in `Active API Key` and for `MAILERPUBLICAPIKEY` use `Email Validation Key`

# Customize keypair for signing JWT tokens (optional)

Cecil uses JWT tokens in a few places to verify the authenticity of links sent to users via email.  In order for this to work, it needs an RSA keypair.

If not provided, it will generate a keypair on it's own and use it, and emit it in the logs.  However, if you want to restart the `cecil` process and re-use the generated keypair, check the logs from the first run and capture the emitted private key into an environment variable named `CECIL_RSA_PRIVATE`:

$ export CECIL_RSA_PRIVATE='-----BEGIN RSA PRIVATE KEY----- MIIEowIBAAKCAQEAt ... -----END RSA PRIVATE KEY-----`

