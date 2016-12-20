This is an ad-hoc description of (most) of the Cecil REST API.  It might be out of date, so please refer to the Swagger REST API specification for the most up-to-date version of the API docs.


## Swagger REST API docs

This is the definitive description of the Cecil REST API and will always be in sync with the code.

You can find the swagger file at `/swagger.json`

E.g.

```bash
curl http://0.0.0.0:8080/swagger.json
```

The best way to view this is via the [Swagger Editor](http://editor.swagger.io/).  You can use the **File/Open URL** menu to display the Swagger Spec in the HTML viewer. 

## Update cloudaccount configuration

```bash
curl -X PATCH \
-H "Authorization: Bearer eyJhbGci" \
-H "Content-Type: application/json" \
-H "Cache-Control: no-cache" \
-d '{
	"default_lease_duration":"45h10s"
}' \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1"
```

`default_lease_duration` allows you to specify the duration for all leases that will be created under that cloudaccount.

To go back to the global setting, just set `default_lease_duration` to `0`.

## Add email to owner tag whitelist

```bash
curl -X POST \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
-H "Content-Type: application/json" \
-d '{"email":"someone.legit@example.com"}' \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/owners"
```

Response:

```json
{
  "message": "owner added successfully to whitelist"
}
```

## List regions and their status

```bash
curl -X GET \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/regions"
```

Response:

```json
{
  "ap-northeast-1": {
    "topic": "not_exists",
    "subscription": "not_active"
  },
  "ap-northeast-2": {
    "topic": "not_exists",
    "subscription": "not_active"
  },
  "ap-south-1": {
    "topic": "not_exists",
    "subscription": "not_active"
  },
  "ap-southeast-1": {
    "topic": "not_exists",
    "subscription": "not_active"
  },
  "ap-southeast-2": {
    "topic": "not_exists",
    "subscription": "not_active"
  },
  "eu-central-1": {
    "topic": "not_exists",
    "subscription": "not_active"
  },
  "eu-west-1": {
    "topic": "not_exists",
    "subscription": "not_active"
  },
  "sa-east-1": {
    "topic": "not_exists",
    "subscription": "not_active"
  },
  "us-east-1": {
    "topic": "exists",
    "subscription": "active"
  },
  "us-east-2": {
    "topic": "not_exists",
    "subscription": "not_active"
  },
  "us-west-1": {
    "topic": "not_exists",
    "subscription": "not_active"
  },
  "us-west-2": {
    "topic": "not_exists",
    "subscription": "not_active"
  }
}
```

As you can see, on the `us-east-1` region the SNS topic exists, and Cecil is subscribed to it.


## Force subscription to topics (to do after you successfully set up the tenant stack on a region)

You can provide a list of regions you want to force subscription (i.e. try subscribing), or just use `["all"]` to force subscription on all regions.

```bash
curl -X POST \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
-H "Content-Type: application/json" \
-d '{
   "regions": ["us-east-1","us-east-2","some-invalid-region-name"]
}' \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/subscribe-sns-to-sqs"
```

Response:

```json
{
  "us-east-1": {
    "topic": "exists",
    "subscription": "active"
  },
  "us-east-2": {
    "topic": "not_exists",
    "subscription": "not_active"
  }
}
```

The response contains the results:

- `us-east-1`: the SNS topic exists and the subscription was successfully setup. This region is monitored.
- `us-east-2`: the SNS topic does not exists yet. It was not possible to setup the subscription.
- `some-invalid-region-name`: this region is just ignore.

## List leases of account

List all leases (new, expired, deleted, all).

Request:

```bash
curl -X GET \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/leases"
```

Response:

```json
[
  {
    "account_id": 1,
    "cloud_account_id": 1,
    "owner_id": 1,
    "aws_account_id": "8932879238795",
    "instance_id": "i-0aa66cc3f2086345543",
    "region": "us-east-1",
    "availability_zone": "us-east-1b",
    "instance_type": "t2.micro",
    "terminated": true,
    "launched_at": "2016-12-04T15:58:29Z",
    "expires_at": "2016-12-05T04:01:20.563096142Z",
    "terminated_at": "2016-12-04T16:03:21Z"
  },
  {
    "account_id": 1,
    "cloud_account_id": 1,
    "owner_id": 1,
    "aws_account_id": "8932879238795",
    "instance_id": "i-0fefcb2718f833247",
    "region": "us-east-1",
    "availability_zone": "us-east-1b",
    "instance_type": "t2.micro",
    "terminated": true,
    "launched_at": "2016-12-04T16:05:01Z",
    "expires_at": "2016-12-05T04:05:08.234707913Z",
    "terminated_at": "2016-12-04T16:06:47Z"
  },
  {
    "account_id": 1,
    "cloud_account_id": 2,
    "owner_id": 1,
    "aws_account_id": "24326257445",
    "instance_id": "i-0fefcb2q32433247",
    "region": "us-east-1",
    "availability_zone": "us-east-1b",
    "instance_type": "t2.micro",
    "terminated": true,
    "launched_at": "2016-12-04T16:05:01Z",
    "expires_at": "2016-12-05T04:05:08.234707913Z",
    "terminated_at": "2016-12-04T16:06:47Z"
  }
]
```

Optionally, you can provide the `terminated` url query parameter to select only terminated or non-terminated leases.

E.g. To get all active leases owned by an account (this might include leases that are being terminated):

```bash
curl -X GET \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/leases?terminated=false"
```

E.g. To get all terminated leases owned by an account:

```bash
curl -X GET \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/leases?terminated=true"
```

## List leases of cloudaccount

List all leases (new, expired, deleted, all).

Request:

```bash
curl -X GET \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/leases"
```

Response:

```json
[
  {
    "account_id": 1,
    "cloud_account_id": 1,
    "owner_id": 1,
    "aws_account_id": "8932879238795",
    "instance_id": "i-0aa66cc3f208632f9",
    "region": "us-east-1",
    "availability_zone": "us-east-1b",
    "instance_type": "t2.micro",
    "terminated": true,
    "launched_at": "2016-12-04T15:58:29Z",
    "expires_at": "2016-12-05T04:01:20.563096142Z",
    "terminated_at": "2016-12-04T16:03:21Z"
  },
  {
    "account_id": 1,
    "cloud_account_id": 1,
    "owner_id": 1,
    "aws_account_id": "8932879238795",
    "instance_id": "i-0fefcb2718f833247",
    "region": "us-east-1",
    "availability_zone": "us-east-1b",
    "instance_type": "t2.micro",
    "terminated": true,
    "launched_at": "2016-12-04T16:05:01Z",
    "expires_at": "2016-12-05T04:05:08.234707913Z",
    "terminated_at": "2016-12-04T16:06:47Z"
  }
]
```

Optionally, you can provide the `terminated` url query parameter to select only terminated or non-terminated leases.

E.g. To get all **active leases** owned by a cloudaccount (this might include leases that are being terminated):

```bash
curl -X GET \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/leases?terminated=false"
```

E.g. To get all **terminated leases** owned by a cloudaccount:

```bash
curl -X GET \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/leases?terminated=true"
```

## Configure Slack

To add Slack as a mean of comunication between you and Cecil, use this endpoint.

```bash
curl -X POST \
-H "Authorization: Bearer eyJhbGc" \
-H "Content-Type: application/json" \
-H "Cache-Control: no-cache" \
-d '{
	"token":"xoxb-000000000-aaaaaaaaaaaaa",
	"channel_id":"#general"
}' \
"http://0.0.0.0:8080/accounts/1/slack_config"
```

Cecil will send messages to the specified channel, and you will be able to issue commands to Cecil.

E.g. To list all available commands, post this in the channel specified in the config, or to the Cecil bot user directly:

```
@cecil help
```

## Remove Slack

```bash
curl -X DELETE \
-H "Authorization: Bearer eyJhbGc" \
-H "Content-Type: application/json" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/slack_config"
```


## Configure Mailer (Mailgun only for the moment)

To add a custom mailer (instead of using the default one of Cecil), use this endpoint:

```bash
curl -X POST \
-H "Authorization: Bearer eyJhbGci" \
-H "Content-Type: application/json" \
-H "Cache-Control: no-cache" \
-d '{
	"domain":"example.com",
	"api_key":"key-000000001a1a1a1a1a1a1a1a1a11a1",
	"public_api_key":"pubkey-2b2b2b2b2b2b22b2b2b2b2b2b2",
	"from_name":"Cecil Account"
}' \
"http://0.0.0.0:8080/accounts/1/mailer_config"
```


## Show a specific lease

Specific lease under an account:

```bash
curl -X GET \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/leases/1"
```


Specific lease under a cloudaccount:

```bash
curl -X GET \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/leases/1"
```

## Terminate a specific lease

Specific lease under an account:

```bash
curl -X POST \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/leases/1/terminate"
```


Specific lease under a cloudaccount:

```bash
curl -X POST \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/leases/1/terminate"
```

## Delete a specific lease from DB

Specific lease under an account:

```bash
curl -X POST \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/leases/1/delete"
```


Specific lease under a cloudaccount:

```bash
curl -X POST \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/leases/1/delete"
```

## Set specific lease's expiration

Specific lease under an account:

```bash
curl -X POST \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/leases/1/expiry?expires_at=2016-12-17T10:46:30Z"
```


Specific lease under a cloudaccount:

```bash
curl -X POST \
-H "Authorization: Bearer eyJhbGc" \
-H "Cache-Control: no-cache" \
"http://0.0.0.0:8080/accounts/1/cloudaccounts/1/leases/1/expiry?expires_at=2016-12-17T10:46:30Z"
```
