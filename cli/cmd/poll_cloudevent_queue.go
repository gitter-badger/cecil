package cmd

import (
	"log"

	"github.com/inconshreveable/log15"
	"github.com/spf13/cobra"
	"github.com/tleyden/zerocloud/cloudevent_poller"
)

var (
	SQSQueueURL     string
	ZeroCloudAPIURL string
	AWSRegion       string
	logger          log15.Logger
)

// poll_cloudevent_queueCmd respresents the poll_cloudevent_queue command
var poll_cloudevent_queueCmd = &cobra.Command{
	Use:   "poll_cloudevent_queue",
	Short: "Looks for new events on CloudEvent SQS queue and pushes to ZeroCloud REST API",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		logger = log15.New()

		if len(SQSQueueURL) == 0 {
			log.Fatalf("SQSQueueURL argument required")
		}
		if len(ZeroCloudAPIURL) == 0 {
			log.Fatalf("ZeroCloudAPIURL argument required")
		}
		if len(AWSRegion) == 0 {
			log.Fatalf("AWSRegion argument required")
		}

		cloudEventPoller := cloudevent_poller.CloudEventPoller{
			SQSQueueURL:     SQSQueueURL,
			ZeroCloudAPIURL: ZeroCloudAPIURL,
			AWSRegion:       AWSRegion,
		}

		err := cloudEventPoller.Run()
		log.Fatalf("Error running cloudEventPoller: %v", err)

	},
}

func init() {

	RootCmd.AddCommand(poll_cloudevent_queueCmd)

	poll_cloudevent_queueCmd.PersistentFlags().StringVar(
		&SQSQueueURL,
		"SQSQueueURL",
		"https://sqs.us-west-1.amazonaws.com/193822812427/ZeroCloudBigDBEvents",
		"The URL of the SQS queue to pull from",
	)
	poll_cloudevent_queueCmd.PersistentFlags().StringVar(
		&ZeroCloudAPIURL,
		"ZeroCloudAPIURL",
		"http://localhost:8080",
		"The URL of the ZeroCloud REST API",
	)
	poll_cloudevent_queueCmd.PersistentFlags().StringVar(
		&AWSRegion,
		"AWSRegion",
		"us-west-1",
		"The AWS region of the SQS queue",
	)

}
