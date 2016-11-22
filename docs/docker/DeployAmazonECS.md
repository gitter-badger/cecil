
## AWS ECR setup

* Signup to ECR
* Create a docker repository called "cecil"
* Using your root AWS user for the account where you plan to deploy the container, run

```bash
$ aws ecr get-login --region us-east-1
```
* Run the `docker login` command returned from running `aws ecr get-login`

## Docker deploy

```
$ docker tag cecil:latest 193822812427.dkr.ecr.us-east-1.amazonaws.com/cecil:latest
$ docker push 193822812427.dkr.ecr.us-east-1.amazonaws.com/cecil:latest
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
-e "MAILERDOMAIN=mg.cecil.co" \
-e "MAILERAPIKEY=..." \
-e "MAILERPUBLICAPIKEY=..." \
-itd -v /tmp/config.yml:/go/config.yml 193822812427.dkr.ecr.us-east-1.amazonaws.com/cecil:latest cecil
```
