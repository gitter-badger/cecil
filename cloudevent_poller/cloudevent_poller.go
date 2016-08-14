package cloudevent_poller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
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
	SQSService      *sqs.SQS
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
	p.SQSService = sqsService

	for {
		err := p.pullItemsFromSQSPushToZeroCloud()
		if err != nil {
			logger.Error("Error pulling items from SQS and pushing to ZeroCloud", "error", err)
		}

	}

}

func (p *CloudEventPoller) pullItemsFromSQSPushToZeroCloud() error {

	// pull an item off queue
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(p.SQSQueueURL),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(1),
	}
	resp, err := p.SQSService.ReceiveMessage(params)

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
	outboundJson, err := transformSQS2RestAPICloudEvent(sqsMsgJson)
	if err != nil {
		return err
	}

	// enhance with things like instance tags (call out to AWS)
	instanceID, err := extractInstanceID(outboundJson)
	if err != nil {
		return err
	}
	logger.Info("Lookup instance-id", "instance-id", instanceID)

	// lookup the EC2 instance tag
	tags, err := p.lookupEC2InstanceTags(instanceID)
	if err != nil {
		return err
	}

	// TODO: might be good to grab the instance state (running, terminated, etc)
	// and attach to the JSON

	log.Printf("tags: %v", tags)
	outboundJson["Tags"] = tags

	// add tags to JSON

	// serialize json
	outboundJsonStr, err := json.Marshal(outboundJson)
	if err != nil {
		return err
	}
	logger.Info("outboundJsonStr", "outboundJsonStr", fmt.Sprintf("%v", string(outboundJsonStr)))

	// push to zerocloud rest API

	// upon succcessful push to zerocloud rest API, delete from SQS queue

	return nil

}

// Given an instance-id, hit the AWS EC2 api to find all the tags associated with
// the instance.
func (p CloudEventPoller) lookupEC2InstanceTags(instanceID string) ([]*ec2.TagDescription, error) {

	// TODO: can probably re-use existing session (p.AWSSession) here.  Need to test.
	session, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	stsService := sts.New(session, &aws.Config{Region: aws.String(p.AWSRegion)})

	// equivalent of CLI:  aws sts assume-role --role-arn arn:aws:iam::788612350743:role/ZeroCloud --role-session-name zerocloud2bigdb --external-id bigdb
	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String("arn:aws:iam::788612350743:role/ZeroCloud"), // TODO: lookup from DB
		RoleSessionName: aws.String("zerocloud2bigdb"),                          // TODO: generate something unique here
		ExternalId:      aws.String("bigdb"),                                    // TODO: lookup from DB
	}
	resp, err := stsService.AssumeRole(params)
	if err != nil {
		return nil, err
	}

	// TODO: rework this according to
	// https://github.com/aws/aws-sdk-go/issues/801#issuecomment-239519183
	provider := NewAssumeRoleCredentialsProvider(resp.Credentials)

	ec2Service := ec2.New(session,
		&aws.Config{
			Region:      aws.String(p.AWSRegion),
			Credentials: credentials.NewCredentials(provider),
		},
	)

	// Add a filter which will filter by that particular ec2 instance
	// It's the CLI equivalent of --filters "Name=resource-id,Values=instance-id"
	paramsDescribeTags := &ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("resource-id"),
				Values: []*string{
					aws.String(instanceID),
				},
			},
		},
	}

	output, err := ec2Service.DescribeTags(paramsDescribeTags)
	if err != nil {
		return nil, err
	}

	return output.Tags, nil
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
func extractInstanceID(cloudEventSQSJson map[string]interface{}) (string, error) {

	// We're passed a map with the parsed JSON, but in order to get this
	// into an app.CloudEventPayload{} instance, marshal to a string
	// and then parse directly into an app.CloudEventPayload{} instance.
	// TODO: there must be an easier way!
	cloudEventSQSJsonStr, err := json.Marshal(cloudEventSQSJson)
	if err != nil {
		return "", err
	}

	payload := app.CloudEventPayload{}
	err = json.Unmarshal([]byte(cloudEventSQSJsonStr), &payload)
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

// This changes the shape of the JSON to what is expected by the /cloudevent REST API
//
// See unit test for more details
func transformSQS2RestAPICloudEvent(inputJsonStr string) (outputJson map[string]interface{}, err error) {

	// parse the inputJSON into a struct
	var inputJson map[string]interface{}
	err = json.Unmarshal([]byte(inputJsonStr), &inputJson)
	if err != nil {
		return nil, err
	}

	// get the Body field, which is an embedded JSON doc
	bodyJsonInterface, ok := inputJson["Body"]
	if !ok {
		return nil, fmt.Errorf("Did not find Body field in %v", inputJsonStr)
	}
	bodyJsonStr := bodyJsonInterface.(string)

	// parse the Body field into JSON
	var bodyJson map[string]interface{}
	err = json.Unmarshal([]byte(bodyJsonStr), &bodyJson)
	if err != nil {
		return nil, err
	}

	// get the Body/Message field, which is an embedded JSON doc
	messageJsonInterface, ok := bodyJson["Message"]
	if !ok {
		return nil, fmt.Errorf("Did not find Body/Message field in %v", inputJsonStr)
	}
	messageJsonStr := messageJsonInterface.(string)

	// parse the Body/Message field into JSON
	var messageJson map[string]interface{}
	err = json.Unmarshal([]byte(messageJsonStr), &messageJson)
	if err != nil {
		return nil, err
	}

	// add it to the resulting payload (overwriting current value with embedded JSON doc)
	bodyJson["Message"] = messageJson

	// base64 encode the entire inputJSON and add as a field
	inputJsonBase64 := base64.StdEncoding.EncodeToString([]byte(inputJsonStr))
	bodyJson["SQSPayloadBase64"] = inputJsonBase64

	return bodyJson, nil
}
