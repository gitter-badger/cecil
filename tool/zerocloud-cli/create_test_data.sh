
go run main.go create account --payload '{"name": "bigdb"}' && \
go run main.go create cloudaccount --accountID 1 --payload '{"cloudprovider": "AWS", "name": "BigDB perf testing AWS account", "upstream_account_id": "788612350743", "assume_role_arn": "arn:aws:iam::788612350743:role/ZeroCloud", "assume_role_external_id": "bigdb"}'

