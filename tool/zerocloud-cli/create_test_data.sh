
go run main.go create account --payload '{"name": "bigdb", "lease_expires_in_units": "days", "lease_expires_in": 3}' && \
    go run main.go create cloudaccount --accountID 1 --payload '{"cloudprovider": "AWS", "name": "BigDB perf testing AWS account", "upstream_account_id": "788612350743", "assume_role_arn": "arn:aws:iam::788612350743:role/ZeroCloud", "assume_role_external_id": "bigdb"}' && \
    go run main.go create cloudevent --payload '{"Message":{"account":"788612350743","detail":{"instance-id":"i-0a74797fd283b53de","state":"running"},"detail-type":"EC2 Instance State-change Notification","id":"2ecfc931-d9f2-4b25-9c00-87e6431d09f7","region":"us-west-1","source":"aws.ec2","time":"2016-08-06T20:53:38Z","version":"0"},"MessageId":"fb7dad1a-ccee-5ac8-ac38-fd3a9c7dfe35","SQSPayloadBase64":"ewogICAgIkF0dHJpYnV0Z........5TlpRPT0iCn0=","Timestamp":"2016-08-06T20:53:39.209Z","TopicArn":"arn:aws:sns:us-west-1:788612350743:BigDBEC2Events","Type":"Notification"}' && \
    curl http://localhost:8080/leases?state=expired | jq .

    

