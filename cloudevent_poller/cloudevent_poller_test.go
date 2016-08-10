package cloudevent_poller

import (
	"encoding/json"
	"log"
	"testing"
)

func TestTransformSQS2RestAPICloudEvent(t *testing.T) {

	test_json := `{
    "Attributes":null,
    "Body":"{\n  \"Type\" : \"Notification\",\n  \"MessageId\" : \"fb7dad1a-ccee-5ac8-ac38-fd3a9c7dfe35\",\n  \"TopicArn\" : \"arn:aws:sns:us-west-1:788612350743:BigDBEC2Events\",\n  \"Message\" : \"{\\\"version\\\":\\\"0\\\",\\\"id\\\":\\\"2ecfc931-d9f2-4b25-9c00-87e6431d09f7\\\",\\\"detail-type\\\":\\\"EC2 Instance State-change Notification\\\",\\\"source\\\":\\\"aws.ec2\\\",\\\"account\\\":\\\"788612350743\\\",\\\"time\\\":\\\"2016-08-06T20:53:38Z\\\",\\\"region\\\":\\\"us-west-1\\\",\\\"resources\\\":[\\\"arn:aws:ec2:us-west-1:788612350743:instance/i-0a74797fd283b53de\\\"],\\\"detail\\\":{\\\"instance-id\\\":\\\"i-0a74797fd283b53de\\\",\\\"state\\\":\\\"running\\\"}}\",\n  \"Timestamp\" : \"2016-08-06T20:53:39.209Z\",\n  \"SignatureVersion\" : \"1\",\n  \"Signature\" : \"boAXkBZ7IP9AnuDrdTlGdF876jAgOGjyDbMzCSOroybhCD3hbNzhO6d8c8FFygttepj4Kc+sC28PEn3SogCYSBaRX93MuE1zg3kkY9335oVNUuDl4GRqCobmpVQEy+Va79IPL7zkMkDGEXfx/3dFdjCz6FRSK0ATtFq/tOEwzKG431gU0qB+/yboRBWztJYBUqICBi0DoEIyNutlXWCYuhQyuMKu17hSo7uIT19oHC8io+xiy6hj1muQwpZH7aKrqQqPuVKinKwDQCQx7svrA8shyHS9QOjTxMSLC7Q3sAd0OuNzY7D4drVi+8bEJsclohAZn+auDLzBCoJClivgeQ==\",\n  \"SigningCertURL\" : \"https://sns.us-west-1.amazonaws.com/SimpleNotificationService-bb750dd426d95ee9390147a5624348ee.pem\",\n  \"UnsubscribeURL\" : \"https://sns.us-west-1.amazonaws.com/?Action=Unsubscribe\u0026SubscriptionArn=arn:aws:sns:us-west-1:788612350743:BigDBEC2Events:0e41602b-834e-4e9f-b710-e043cb758754\"\n}",
    "MD5OfBody":"05ae7bf2999e64a721b9ecb0aee0af21",
    "MD5OfMessageAttributes":null,
    "MessageAttributes":null,
    "MessageId":"e9b98b32-8b6a-4b62-90f6-61927c8de4a6",
    "ReceiptHandle":"AQEBO3rPq22vcT8B3mi1VxGo+/eYN4bdRi5p3H7NPMIFgMHrYJR0AMyJGdr+xKMvXbrNUlU4w70dwld3GGJfALS10QxKDqGvY6nQNMgw9CdwRr0/e+PBjmOAseqU07e2RNYihsm3WsH3VbmyqPMKDDVIk+oQZVqHAEdF/+l4j8pdLr6sp6oBgb5M89nyamsBeXX5tJeybmgiIIMYzRCB28MZYRoratNjU1ak0mqvaYTv0t2JNnwPDnSkehYX6/o/Vfp1W8KNHXlpaPf3PdHE/LEPeC3O2w0EupLdhQAcVApvDzDyBtEQaqtx2i9AP0mwy1Ldj2MuphOKR3PZ0YLdZpVoiIMEz0JDFDpEo0fZ2vAJkyAnWkpMk1H5hjMUVvzYOVu96p3MR9Fcm184XSG2Hg9NZQ=="
}`
	outputJsonStr, err := transformSQS2RestAPICloudEvent(test_json)
	if err != nil {
		t.Errorf("Got error transforming JSON: %v", err)
	}
	if len(outputJsonStr) == 0 {
		t.Errorf("Expected non-zero JSON string result")
	}

	// parse JSON string into a struct
	var outputJson map[string]interface{}
	err = json.Unmarshal([]byte(outputJsonStr), &outputJson)
	if err != nil {
		t.Errorf("Got error unmarshalling JSON into a struct: %v", err)
	}

	// the resulting JSON struct should have a top level "Type" field, which is present in the "Body" field of the original message
	typeField, ok := outputJson["Type"]
	if !ok {
		t.Errorf("Expected top-level Type field in JSON")
	}
	log.Printf("typeField: %v", typeField)
	log.Printf("typeField: %T", typeField)
	if typeField != "Notification" {
		t.Errorf("Expected Type field to have a value of Notification")
	}

	// TODO: should have a SQSPayload field with base64'd JSON of original SQS payload

}
