
# Create Cecil Root CloudFormation Stack

Before you can deploy to Docker Cloud, you will first need to follow the steps under the **Cecil AWS Setup (AWS CLI)** section in [InstallCecilService.md](InstallCecilService.md) to create the `CecilRootStack` CloudFormation stack.

# Docker Cloud Setup

## Setup a node cluster

* Single node (t2.small should work)
* It should work with any underlying IaaS provider

## Create service

Go to the **Services** section of Docker Cloud and create a new service.

* Choose `tleyden5iwx/cecil` (hosted on dockerhub) -- or you could create your own dockerhub repo and point to that.
* Use the `latest` tag
* Containers: 1
* Ports
    * Add a port with container port 8080 and node port 8080, and check the **Published** checkbox
* Set the following environment variables
    * `AWS_ACCESS_KEY_ID` (from `aws iam create-access-key` command in [InstallCecilService.md](InstallCecilService.md))
    * `AWS_ACCOUNT_ID` (the same AWS account id used in [InstallCecilService.md](InstallCecilService.md))
    * `AWS_REGION` (the same AWS account id used in [InstallCecilService.md](InstallCecilService.md))
    * `AWS_SECRET_ACCESS_KEY` (from `aws iam create-access-key` command in [InstallCecilService.md](InstallCecilService.md))
    * `MAILERAPIKEY` (optional, see [StaticConfig.md](StaticConfig.md))
    * `MAILERDOMAIN` (optional, see [StaticConfig.md](StaticConfig.md))
    * `MAILERPUBLICAPIKEY` (optional, see [StaticConfig.md](StaticConfig.md))
    

## Launch

Hit the **Create & Deploy** button to launch the service

## Verify REST API

In the Service definition, look for something like `http://cecil.f24cd253.svc.dockerapp.io:8080`.  Click that link, and you should see:

```
{
  "name": "Cecil REST API",
  "uptime": "12.880620179s",
  "time": "2016-12-23T00:40:39Z"
}
```

## Update SERVERHOSTNAME

The following step is needed for URLs in notification emails to point to valid URLs.

In the Service definition, add a `SERVERHOSTNAME` environment variable which corresponds to the hostname under the "Service Endpoints" section of the Service dashboard.  It should be something.svc.dockerapp.io.

After this, save your changes and redeploy.

## Configure tenant(s) + AWS account(s)

At this point, you are ready to [Configure one or more of your AWS accounts to be monitored by Cecil](ConfigureAWSAccount.md)

## Verify everything works

Follow the steps in the **Verify everything works**  section in [Configure one or more of your AWS accounts to be monitored by Cecil](ConfigureAWSAccount.md) 

## Preserving the database across restarts

Most of the time you will want to preserve the data across redeploys of the Cecil Docker Cloud service.  Here are the steps to do that:

### ssh into node

Follow the [SSH into a Docker Cloud-managed node](https://docs.docker.com/docker-cloud/infrastructure/ssh-into-a-node/) instructions to add your SSH key.

### Copy the database file to a file on the host

After you have ssh'd into the Docker Cloud host, run these steps to copy the database file

```
$ CONTAINER_ID=$(docker ps | grep -i cecil | awk '{print $1}')
$ docker cp $CONTAINER_ID:/go/src/github.com/tleyden/cecil/cecil.db .
```

### Update the service 

In the volumes section, hit the plus button to the right of the second line **Add volumes**, and use:

* Container path: `/go/src/github.com/tleyden/cecil/cecil.db`
* Host path: `/root/cecil.db`

You can now redeploy the service and your data will be preserved.

## Continuous Deployment

If you want it to automatically redeploy Cecil whenever a new build on dockerhub completes, simply check the **autoredeploy** toggle.  You'll most likely want to first preserve the database across restarts, otherwise you will lose all of your data each time.

## Parallel staging environment

* Create another Docker Cloud Service and call it `cecil-staging`
* Using _a different AWS account_ (eg, 788612350743), go through the steps under the **Cecil AWS Setup (AWS CLI)** section in [InstallCecilService.md](InstallCecilService.md) to create a parallel `CecilRootStack` CloudFormation stack.
* In the Docker Cloud service config `SERVERPORT` environment variable set to `:8081`
* Container Port mapping 8081 -> 8081
* Using the the same staging AWS account (eg, 788612350743) follow the steps in [Configure one or more of your AWS accounts to be monitored by Cecil](ConfigureAWSAccount.md) (note: this will use STS cross-account authentication within the _same_ AWS account) -- also note that you will _not_ receive a confirmation email after creating the `tenant-aws-region-setup.template` stack, not sure why.

For Preserving the database across restarts:

* Update the volume mapping to point to a _different database file on the host_ -- eg `/root/staging/cecil.db`