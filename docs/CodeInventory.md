
## core package contents

- `add-owner-handler.go` -- Contains the handler function for adding a new owner to owner's whitelist for a cloudaccount.
- `aws.go` -- Contains SQS structs and DefaultEc2ServiceFactory.
- `common.go` -- Contains common utility functions.
- `core.go` -- Contains the all the initialization code for the core package.
- `core_test.go` -- core package test.
- `db-models.go` -- Contains the database models.
- `email-action-handler.go` -- Contains the handler function for lease approval|extension|termination link endpoints.
- `email-templates.go` -- Will contain the templates of the emails sent out for specific scenarios (new lease, lease expired, instance terminated, etc.).
- `mock_ec2.go` -- Contains a mock of the EC2 API.
- `mock_mailgun.go` -- Contains a mock of the Mailgun API.
- `mock_sqs.go` -- Contains a mock of the SQS API.
- `new-lease-queue-consumer.go` -- Contains the consumer function for the NewLeaseQueue.
- `periodic-jobs.go` -- Contains the periodic job functions
- `service.go` -- Contains the Service struct and the initialization methods (to setup queues, db, external services, etc.)
- `task-consumers.go` -- Contains some of the functions that consume tasks from queues; some got their own file because are big.
- `task-structs.go` -- Contains the structs of the tasks passed in-out of queues.
- `transmission.go` -- Contains the `Transmission` and its methods; `Transmission` is what an SQS message is parsed to.
