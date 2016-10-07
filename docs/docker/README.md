# Docker

## Docker build

On your workstation, install the [docker toolbox](https://www.docker.com/products/docker-toolbox)

In $GOPATH/src/github.com/tleyden/zerocloud, run:

```
docker build -t zerocloud .
```

## AWS ECR setup

* Signup to ECR
* Create a docker repository called "zerocloud"
* Using your root AWS user for the account where you plan to deploy the container, run

```
$ aws ecr get-login --region us-east-1
```
* Run the `docker login` command returned from running `aws ecr get-login`

## Docker deploy

```
$ docker tag zerocloud:latest 193822812427.dkr.ecr.us-east-1.amazonaws.com/zerocloud:latest
$ docker push 193822812427.dkr.ecr.us-east-1.amazonaws.com/zerocloud:latest
```

## Docker run

On the machine you want to run the docker image:

* Run the `docker login` command returned from running `aws ecr get-login`
* Open `/tmp/config.yml` and add your config
* Docker run

```
$ docker run \
-e "AWS_ACCESS_KEY_ID=..." \
-e "AWS_SECRET_ACCESS_KEY=..." \
-e "AWS_ACCOUNT_ID=..." \
-e "AWS_REGION=us-east-1" \
-e "MAILERDOMAIN=mg.zerocloud.co" \
-e "MAILERAPIKEY=..." \
-e "MAILERPUBLICAPIKEY=..." \
-itd -v /tmp/config.yml:/go/config.yml 193822812427.dkr.ecr.us-east-1.amazonaws.com/zerocloud:latest zerocloud
```
