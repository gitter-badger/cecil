package main

import (
	"fmt"

	"github.com/tleyden/zerocloud/app"
	"github.com/tleyden/zerocloud/models"
)

// Create runs the create action.
func (c *CloudeventController) CreateImpl(ctx *app.CreateCloudeventContext) error {

	awsAccountId := *ctx.Payload.AwsAccountID
	logger.Info("Create CloudEvent", "aws_account_id", awsAccountId)

	// try to find the CloudAccount that has an upstream_account_id that matches param
	cloudAccount := models.CloudAccount{}
	cdb.Db.Where(&models.CloudAccount{UpstreamAccountID: awsAccountId}).First(&cloudAccount)
	logger.Info("Found CloudAccount", "CloudAccount", fmt.Sprintf("%+v", cloudAccount))

	// Make sure we found a valid CloudAccount, otherwise abort
	if cloudAccount.ID == 0 {
		ctx.ResponseData.Service.Send(
			ctx.Context,
			400,
			fmt.Sprintf("Could not find CloudAccount with upstream provider account id: %v", awsAccountId),
		)
		return nil
	}

	/*
				Field("sqs_payload", gorma.Text)
				Field("sqs_timestamp", gorma.Timestamp)
				Field("cw_event_detail_type", gorma.String)
				Field("cw_event_source", gorma.String)
				Field("cw_event_timestamp", gorma.Timestamp)
				Field("cw_event_region", gorma.String)
				Field("cw_event_detail_instance_id", gorma.String)
				Field("cw_event_detail_state", gorma.String)

		Example JSON for payload design:

		{
		    "Message":{
		        "account":"788612350743",
		        "detail":{
		            "instance-id":"i-0a74797fd283b53de",
		            "state":"running"
		        },
		        "detail-type":"EC2 Instance State-change Notification",
		        "id":"2ecfc931-d9f2-4b25-9c00-87e6431d09f7",
		        "region":"us-west-1",
		        "resources":[
		            "arn:aws:ec2:us-west-1:788612350743:instance/i-0a74797fd283b53de"
		        ],
		        "source":"aws.ec2",
		        "time":"2016-08-06T20:53:38Z",
		        "version":"0"
		    },
		    "MessageId":"fb7dad1a-ccee-5ac8-ac38-fd3a9c7dfe35",
		    "SQSPayloadBase64":"ewogICAgIkF0dHJpYnV0Z........5TlpRPT0iCn0=",
		    "Signature":"boAXkBZ7IP9AnuDrdTlGdF876jAgOGjyDbMzCSOroybhCD3hbNzhO6d8c8FFygttepj4Kc+sC28PEn3SogCYSBaRX93MuE1zg3kkY9335oVNUuDl4GRqCobmpVQEy+Va79IPL7zkMkDGEXfx/3dFdjCz6FRSK0ATtFq/tOEwzKG431gU0qB+/yboRBWztJYBUqICBi0DoEIyNutlXWCYuhQyuMKu17hSo7uIT19oHC8io+xiy6hj1muQwpZH7aKrqQqPuVKinKwDQCQx7svrA8shyHS9QOjTxMSLC7Q3sAd0OuNzY7D4drVi+8bEJsclohAZn+auDLzBCoJClivgeQ==",
		    "SignatureVersion":"1",
		    "SigningCertURL":"https://sns.us-west-1.amazonaws.com/SimpleNotificationService-bb750dd426d95ee9390147a5624348ee.pem",
		    "Timestamp":"2016-08-06T20:53:39.209Z",
		    "TopicArn":"arn:aws:sns:us-west-1:788612350743:BigDBEC2Events",
		    "Type":"Notification",
		    "UnsubscribeURL":"https://sns.us-west-1.amazonaws.com/?Action=Unsubscribe\u0026SubscriptionArn=arn:aws:sns:us-west-1:788612350743:BigDBEC2Events:0e41602b-834e-4e9f-b710-e043cb758754"
		}

	*/

	// Save the raw CloudEvent to the database
	e := models.CloudEvent{}
	e.AwsAccountID = awsAccountId
	e.CloudAccountID = cloudAccount.ID
	e.AccountID = cloudAccount.AccountID
	err := edb.Add(ctx.Context, &e)
	if err != nil {
		return ErrDatabaseError(err)
	}

	// Create a Lease object that references this (immutable) CloudEvent and expires
	// based on the settings in the Account
	// TODO: or can this be an AfterCreate callback on the CloudEvent?
	// file:///Users/tleyden/DevLibraries/gorm/callbacks.html
	err = createLease(e)
	if err != nil {
		return ErrDatabaseError(err)
	}

	// TODO: should this return the path to the cloudevent .. should there even be one?
	/// ctx.ResponseData.Header().Set("Location", app.CloudeventHref(ctx.AccountID, a.ID))

	return ctx.Created()

}

func createLease(cloudEvent models.CloudEvent) error {
	// TODO
	return nil
}
