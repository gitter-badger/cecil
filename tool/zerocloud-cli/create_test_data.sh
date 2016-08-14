
go run main.go create account --payload '{"name": "bigdb"}' && \
go run main.go create cloudaccount --accountID 1 --payload '{"cloudprovider": "AWS", "name": "BigDB perf testing AWS account", "upstream_account_id": "788612350743"}'

