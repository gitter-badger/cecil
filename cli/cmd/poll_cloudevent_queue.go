package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tleyden/zerocloud/cloudevent_poller"
)

var (
	SQSQueueTopicARN string
	ZeroCloudAPIURL  string
)

// poll_cloudevent_queueCmd respresents the poll_cloudevent_queue command
var poll_cloudevent_queueCmd = &cobra.Command{
	Use:   "poll_cloudevent_queue",
	Short: "Looks for new events on CloudEvent SQS queue and pushes to ZeroCloud REST API",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("poll_cloudevent_queue called")
		log.Printf("ZeroCloudAPIURL: %v", ZeroCloudAPIURL)
		log.Printf("SQSQueueTopicARN: %v", SQSQueueTopicARN)
		log.Printf("args: %+v", args)
		if len(SQSQueueTopicARN) == 0 {
			log.Fatalf("SQSQueueTopicARN argument required")
		}
		if len(ZeroCloudAPIURL) == 0 {
			log.Fatalf("ZeroCloudAPIURL argument required")
		}

		cloudEventPoller := cloudevent_poller.CloudEventPoller{
			SQSQueueTopicARN: SQSQueueTopicARN,
			ZeroCloudAPIURL:  ZeroCloudAPIURL,
		}
		cloudEventPoller.Run()

	},
}

func init() {
	RootCmd.AddCommand(poll_cloudevent_queueCmd)

	poll_cloudevent_queueCmd.PersistentFlags().StringVar(
		&SQSQueueTopicARN,
		"SQSQueueTopicARN",
		"arn:aws:sns:us-west-1:788612350743:BigDBEC2Events",
		"The ARN of the SQS queue to pull from",
	)
	poll_cloudevent_queueCmd.PersistentFlags().StringVar(
		&ZeroCloudAPIURL,
		"ZeroCloudAPIURL",
		"http://localhost:8080",
		"The URL of the ZeroCloud REST API",
	)

}
