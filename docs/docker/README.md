# Docker

## Docker build

On your workstation, install the [docker toolbox](https://www.docker.com/products/docker-toolbox)

In $GOPATH/src/github.com/tleyden/cecil, run:

```bash
docker build -t cecil .
```

## Docker run

Make sure the following environment variables are set to the correct values:

* `AWS_ACCESS_KEY_ID`
* `AWS_SECRET_ACCESS_KEY`
* `AWS_ACCOUNT_ID`
* `AWS_REGION`


```
$ docker run -itd \
-e "AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID" \
-e "AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY" \
-e "AWS_ACCOUNT_ID=$AWS_ACCOUNT_ID" \
-e "AWS_REGION=$AWS_REGION" \
-v config.yml:/go/config.yml \
cecil
```

You can also optionally pass the following environment variables:

```
-e "MAILERDOMAIN=..." \
-e "MAILERAPIKEY=..." \
-e "MAILERPUBLICAPIKEY=..." \
```

## Docker PaaS deployment instructions

* [DockerCloud](DeployDockerCloud.md)
* [Amazon ECS](DeployAmazonECS.md)

