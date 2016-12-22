Here is the system overview of what a Cecil deploy looks like:

![](architecture-flowcharts/system-overview-diagram.png)

In this document you will setup the *right hand side*.  If your company name is "Acme", you would likely call this the "Acme Cecil Service"

The instructions below are geared towards running Cecil directly on a machine or virtual machine.  Alternatively, you can [deploy Cecil in the cloud](DeployToCloud.md)


# Get code

```
go get -t github.com/tleyden/cecil/...
```

If this completes without errors, you will have a new binary in `$GOPATH/bin/cecil`

# Choose or create an AWS account for Cecil

If you have an existing account that is **separate from the AWS account(s) you want to monitor**, use that.  Otherwise, create a brand new AWS account.

You will need to create an Access Key for the root user of the AWS account if you don't already have one.

At this point you should have the following:

| Description | AWS Account ID        | AWS_KEY           | AWS_SECRET_KEY |  Root/IAM | Attached Policies 
| ------------- |:-------------:|:-----:|:-----:|:-----:|:-----:|
| Cecil AWS account Root User | 193822812427      | AKIAEXAMPLEWAGRHKOMEWQ | ********** | Root | N/A 


# Cecil AWS Setup (AWS CLI)

This step will create the following resources on your AWS account:

* An IAM User that the Cecil process will use (CecilRootUser)
* Assign policies to the CecilRootUser
    * STSAssumeRole
    * Access to the CecilQueue SQS queue
    * The ability to subscribe to SNS
* An SQS queue to receive CloudWatch Events (CecilQueue)

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

Create credentials (keys) for `CecilRootUser` using the same `AWS_ACCESS_KEY_ID` and secret as the previous step:

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
        "AccessKeyId": "AKIAIEXAMPLERQ4U4N67LE7A"
    }
}
```

At this point you should have the following:

| Description | AWS Account ID        | AWS_KEY           | AWS_SECRET_KEY |  Root/IAM | Attached Policies 
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

There are other optional [static configuration options](StaticConfig.md)

Run cecil:

```
$ cecil
```

Now you are ready to go to the next step: [Configure an AWS account](ConfigureAWSAccount.md)


