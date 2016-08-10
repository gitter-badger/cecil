package cloudevent_poller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/inconshreveable/log15"
)

var logger log15.Logger

func init() {
	logger = log15.New()
}

type CloudEventPoller struct {
	SQSQueueURL     string
	ZeroCloudAPIURL string
	AWSRegion       string
}

func (p *CloudEventPoller) Run() error {
	logger.Info("Run() called.", "poller", fmt.Sprintf("%+v", p))

	// connect to SQS queue
	session, err := session.NewSession()
	if err != nil {
		return err
	}
	logger.Info("Session", "session", fmt.Sprintf("%+v", session))
	logger.Info("Session", "config", fmt.Sprintf("%+v", session.Config))
	logger.Info("Session", "config region", fmt.Sprintf("%+v", *session.Config.Region))
	logger.Info("Session", "config credentials", fmt.Sprintf("%+v", session.Config.Credentials))

	sqsService := sqs.New(session, &aws.Config{Region: aws.String(p.AWSRegion)})
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

	// enhance with things like instance tags (call out to AWS)

	// push to zerocloud rest API

	// upon succcessful push to zerocloud rest API, delete from SQS queue

	// }

	return nil

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
