package cloudevent_poller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/inconshreveable/log15"
	"github.com/tleyden/zerocloud/app"
)

var logger log15.Logger

func init() {
	logger = log15.New()
}

type CloudEventPoller struct {
	SQSQueueURL     string
	ZeroCloudAPIURL string
	AWSRegion       string
	AWSSession      *session.Session
}

func (p *CloudEventPoller) Run() error {
	logger.Info("Run() called.", "poller", fmt.Sprintf("%+v", p))

	// connect to SQS queue
	session, err := session.NewSession()
	if err != nil {
		return err
	}
	p.AWSSession = session

	sqsService := sqs.New(p.AWSSession, &aws.Config{Region: aws.String(p.AWSRegion)})
	logger.Info("sqs service", "sqs", fmt.Sprintf("%+v", sqsService))

	// pull an item off queue
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(p.SQSQueueURL),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(1),
	}
	resp, err := sqsService.ReceiveMessage(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return nil
	}
	logger.Info("resp", "resp", fmt.Sprintf("%+v", resp))

	// get the one and only message
	sqsMessage := extractOnlySQSMessage(resp)

	// serialize to json
	sqsMsgJson := serializeToJson(sqsMessage)

	log.Printf("sqsMsgJson: %v", sqsMsgJson)

	// transform input json to outbound JSON
	outboundJsonStr, err := transformSQS2RestAPICloudEvent(sqsMsgJson)
	if err != nil {
		return err
	}

	// enhance with things like instance tags (call out to AWS)
	instanceID, err := extractInstanceID(outboundJsonStr)
	if err != nil {
		return err
	}
	logger.Info("Lookup instance-id", "instance-id", instanceID)

	// lookup the EC2 instance tags
	tags, err := p.lookupEC2InstanceTags(instanceID)
	if err != nil {
		return err
	}

	log.Printf("tags: %v", tags)

	// push to zerocloud rest API

	// upon succcessful push to zerocloud rest API, delete from SQS queue

	// }

	return nil

}

// Given an instance-id, hit the AWS EC2 api to find all the tags associated with
// the instance.
// NOTE: this isn't working because it's using the WRONG AWS creds.  It's picking up
// the ZeroCloud creds from the environment, but it needs to use the BigDB (customer)
// creds via a call to AssumeRole
func (p CloudEventPoller) lookupEC2InstanceTags(instanceID string) (map[string]string, error) {

	//    DescribeTagsRequest(*ec2.DescribeTagsInput) (*request.Request, *ec2.DescribeTagsOutput)

	// TODO: AssumeRole stuff to use appropriate creds
	// aws sts assume-role --role-arn arn:aws:iam::788612350743:role/ZeroCloud --role-session-name zerocloud2bigdb --external-id bigdb
	// We're going to need to know the role-arn, which will need to be stored in
	// the CloudAccount record, along with the external-id

	session, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	stsService := sts.New(session, &aws.Config{Region: aws.String(p.AWSRegion)})

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String("arn:aws:iam::788612350743:role/ZeroCloud"), // TODO: lookup from DB
		RoleSessionName: aws.String("zerocloud2bigdb"),                          // TODO: generate something unique here
		ExternalId:      aws.String("bigdb"),                                    // TODO: lookup from DB
	}
	resp, err := stsService.AssumeRole(params)
	if err != nil {
		return nil, err
	}
	log.Printf("assumeRole resp: %+v", resp)
	/*
		2016/08/12 10:15:16 assumeRole resp: {
		  AssumedRoleUser: {
		    Arn: "arn:aws:sts::788612350743:assumed-role/ZeroCloud/zerocloud2bigdb",
		    AssumedRoleId: "AROAIUDLOTKTGLMKQZQTU:zerocloud2bigdb"
		  },
		  Credentials: {
		    AccessKeyId: "ASIAJ2DSXOPMDHVZ7BOA",
		    Expiration: 2016-08-12 18:15:15 +0000 UTC,
		    SecretAccessKey: "uSCiUK8LrUwxT4+N2t6WaeJoiVCTa3zbP3Kwo90t",
		    SessionToken: "FQoDYXdzEGIaDP2Pq9bzbEiUr4W+cyLSAdpegfLXVulXJUIWNkm74JI9GilGnNrLbIcVcGr4urM03aFEwf7nhgz/AwBvSvVIN4WKHp2v2xxdYszAYFJujSf3Ac7Jw2NvEvby7QhxjzXivOzUI1W60Fd6cg+bzov1cpV7t3InHmooSxpU0TErZ3gzAy/Fi0HeRllNhzZIad9d5bgTCO8HKXHOG40HnQpT0Pt8pjPxAB28oIA3bzvOA9dgMXGwbZsh217FS60SsrYBCuYIQ6V43/5JGp/PIqvYpmEJ8T07561PkayFm+7lyk2XWSijiLi9BQ=="
		  }
		}
	*/

	ec2Service := ec2.New(p.AWSSession, &aws.Config{Region: aws.String(p.AWSRegion)})

	params2 := &ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("resource-id"),
				Values: []*string{
					aws.String(instanceID),
				},
			},
		},
	}

	output, err := ec2Service.DescribeTags(params2)
	if err != nil {
		return nil, err
	}

	log.Printf("describe tags output: %+v", output)

	return nil, nil
}

//  Extract instance-id from JSON with structure
//
//   {
//   "Message": {
//      "account": "868768768",
//      "detail": {
//         "instance-id": "i-0a74797fd283b53de",
//         "state": "running"
//      },
//
func extractInstanceID(cloudEventSQSJsonStr string) (string, error) {
	payload := app.CloudEventPayload{}
	err := json.Unmarshal([]byte(cloudEventSQSJsonStr), &payload)
	if err != nil {
		return "", err
	}
	return payload.Message.Detail.InstanceID, nil
}

func extractOnlySQSMessage(resp *sqs.ReceiveMessageOutput) *sqs.Message {
	if len(resp.Messages) != 1 {
		log.Panicf("Expected 1 message in SQS response, got %v messages", len(resp.Messages))
	}
	return resp.Messages[0]
}

func serializeToJson(msg *sqs.Message) string {
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Panicf("Error marshalling JSON to string.  Msg: %+v", msg)
	}
	return string(bytes)
}

func transformSQS2RestAPICloudEvent(inputJsonStr string) (outputJson string, err error) {

	// parse the inputJSON into a struct
	var inputJson map[string]interface{}
	err = json.Unmarshal([]byte(inputJsonStr), &inputJson)
	if err != nil {
		return "", err
	}

	// get the Body field, which is an embedded JSON doc
	bodyJsonInterface, ok := inputJson["Body"]
	if !ok {
		return "", fmt.Errorf("Did not find Body field in %v", inputJsonStr)
	}
	bodyJsonStr := bodyJsonInterface.(string)

	// parse the Body field into JSON
	var bodyJson map[string]interface{}
	err = json.Unmarshal([]byte(bodyJsonStr), &bodyJson)
	if err != nil {
		return "", err
	}

	// get the Body/Message field, which is an embedded JSON doc
	messageJsonInterface, ok := bodyJson["Message"]
	if !ok {
		return "", fmt.Errorf("Did not find Body/Message field in %v", inputJsonStr)
	}
	messageJsonStr := messageJsonInterface.(string)

	// parse the Body/Message field into JSON
	var messageJson map[string]interface{}
	err = json.Unmarshal([]byte(messageJsonStr), &messageJson)
	if err != nil {
		return "", err
	}

	// add it to the resulting payload (overwriting current value with embedded JSON doc)
	bodyJson["Message"] = messageJson

	// base64 encode the entire inputJSON and add as a field
	inputJsonBase64 := base64.StdEncoding.EncodeToString([]byte(inputJsonStr))
	bodyJson["SQSPayloadBase64"] = inputJsonBase64

	// serialize json and return
	resultJsonBytes, err := json.Marshal(bodyJson)
	if err != nil {
		return "", err
	}

	return string(resultJsonBytes), nil
}
