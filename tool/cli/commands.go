package cli

import (
	"encoding/json"
	"fmt"
	"github.com/goadesign/goa"
	goaclient "github.com/goadesign/goa/client"
	uuid "github.com/goadesign/goa/uuid"
	"github.com/spf13/cobra"
	"github.com/tleyden/zerocloud/client"
	"golang.org/x/net/context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	// CreateAccountCommand is the command line data structure for the create action of account
	CreateAccountCommand struct {
		Payload     string
		ContentType string
		PrettyPrint bool
	}

	// DeleteAccountCommand is the command line data structure for the delete action of account
	DeleteAccountCommand struct {
		// Account ID
		AccountID   int
		PrettyPrint bool
	}

	// ListAccountCommand is the command line data structure for the list action of account
	ListAccountCommand struct {
		PrettyPrint bool
	}

	// ShowAccountCommand is the command line data structure for the show action of account
	ShowAccountCommand struct {
		// Account ID
		AccountID   int
		PrettyPrint bool
	}

	// UpdateAccountCommand is the command line data structure for the update action of account
	UpdateAccountCommand struct {
		Payload     string
		ContentType string
		// Account ID
		AccountID   int
		PrettyPrint bool
	}

	// ShowAwsCommand is the command line data structure for the show action of aws
	ShowAwsCommand struct {
		AwsAccountID string
		PrettyPrint  bool
	}

	// CreateCloudaccountCommand is the command line data structure for the create action of cloudaccount
	CreateCloudaccountCommand struct {
		Payload     string
		ContentType string
		// Account ID
		AccountID   int
		PrettyPrint bool
	}

	// DeleteCloudaccountCommand is the command line data structure for the delete action of cloudaccount
	DeleteCloudaccountCommand struct {
		// Account ID
		AccountID      int
		CloudAccountID int
		PrettyPrint    bool
	}

	// ListCloudaccountCommand is the command line data structure for the list action of cloudaccount
	ListCloudaccountCommand struct {
		// Account ID
		AccountID   int
		PrettyPrint bool
	}

	// ShowCloudaccountCommand is the command line data structure for the show action of cloudaccount
	ShowCloudaccountCommand struct {
		// Account ID
		AccountID      int
		CloudAccountID int
		PrettyPrint    bool
	}

	// UpdateCloudaccountCommand is the command line data structure for the update action of cloudaccount
	UpdateCloudaccountCommand struct {
		Payload     string
		ContentType string
		// Account ID
		AccountID      int
		CloudAccountID int
		PrettyPrint    bool
	}

	// CreateCloudeventCommand is the command line data structure for the create action of cloudevent
	CreateCloudeventCommand struct {
		Payload     string
		ContentType string
		PrettyPrint bool
	}

	// ListLeaseCommand is the command line data structure for the list action of lease
	ListLeaseCommand struct {
		PrettyPrint bool
	}
)

// RegisterCommands registers the resource action CLI commands.
func RegisterCommands(app *cobra.Command, c *client.Client) {
	var command, sub *cobra.Command
	command = &cobra.Command{
		Use:   "create",
		Short: `create action`,
	}
	tmp1 := new(CreateAccountCommand)
	sub = &cobra.Command{
		Use:   `account ["/accounts"]`,
		Short: ``,
		Long: `

Payload example:

{
   "lease_expires_in": 3,
   "lease_expires_in_units": "days",
   "name": "BigDB"
}`,
		RunE: func(cmd *cobra.Command, args []string) error { return tmp1.Run(c, args) },
	}
	tmp1.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp1.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	tmp2 := new(CreateCloudaccountCommand)
	sub = &cobra.Command{
		Use:   `cloudaccount ["/accounts/ACCOUNTID/cloudaccounts"]`,
		Short: ``,
		Long: `

Payload example:

{
   "assume_role_arn": "arn:aws:iam::788612350743:role/ZeroCloud",
   "assume_role_external_id": "bigdb",
   "cloudprovider": "AWS",
   "name": "BigDB.co's perf testing AWS account",
   "upstream_account_id": "98798079879"
}`,
		RunE: func(cmd *cobra.Command, args []string) error { return tmp2.Run(c, args) },
	}
	tmp2.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp2.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	tmp3 := new(CreateCloudeventCommand)
	sub = &cobra.Command{
		Use:   `cloudevent ["/cloudevent"]`,
		Short: ``,
		Long: `

Payload example:

{
   "Message": {
      "account": "868768768",
      "detail": {
         "instance-id": "i-0a74797fd283b53de",
         "state": "running"
      },
      "detail-type": "EC2 Instance State-change Notification",
      "id": "2ecfc931-d9f2-4b25-9c00-87e6431d09f7",
      "region": "us-west-1",
      "source": "aws.ec2",
      "time": "2016-08-06T20:53:38Z",
      "version": "0"
   },
   "MessageId": "fb7dad1a-ccee-5ac8-ac38-fd3a9c7dfe35",
   "SQSPayloadBase64": "ewogICAgIkF0dHJpYnV0Z........5TlpRPT0iCn0=",
   "Timestamp": "2016-08-06T20:53:39.209Z",
   "TopicArn": "arn:aws:sns:us-west-1:788612350743:BigDBEC2Events",
   "Type": "Notification"
}`,
		RunE: func(cmd *cobra.Command, args []string) error { return tmp3.Run(c, args) },
	}
	tmp3.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp3.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "delete",
		Short: `delete action`,
	}
	tmp4 := new(DeleteAccountCommand)
	sub = &cobra.Command{
		Use:   `account ["/accounts/ACCOUNTID"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp4.Run(c, args) },
	}
	tmp4.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp4.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	tmp5 := new(DeleteCloudaccountCommand)
	sub = &cobra.Command{
		Use:   `cloudaccount ["/accounts/ACCOUNTID/cloudaccounts/CLOUDACCOUNTID"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp5.Run(c, args) },
	}
	tmp5.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp5.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "list",
		Short: `list action`,
	}
	tmp6 := new(ListAccountCommand)
	sub = &cobra.Command{
		Use:   `account ["/accounts"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp6.Run(c, args) },
	}
	tmp6.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp6.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	tmp7 := new(ListCloudaccountCommand)
	sub = &cobra.Command{
		Use:   `cloudaccount ["/accounts/ACCOUNTID/cloudaccounts"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp7.Run(c, args) },
	}
	tmp7.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp7.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	tmp8 := new(ListLeaseCommand)
	sub = &cobra.Command{
		Use:   `lease ["/leases"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp8.Run(c, args) },
	}
	tmp8.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp8.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "show",
		Short: `show action`,
	}
	tmp9 := new(ShowAccountCommand)
	sub = &cobra.Command{
		Use:   `account ["/accounts/ACCOUNTID"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp9.Run(c, args) },
	}
	tmp9.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp9.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	tmp10 := new(ShowAwsCommand)
	sub = &cobra.Command{
		Use:   `aws ["/aws/AWSACCOUNTID"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp10.Run(c, args) },
	}
	tmp10.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp10.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	tmp11 := new(ShowCloudaccountCommand)
	sub = &cobra.Command{
		Use:   `cloudaccount ["/accounts/ACCOUNTID/cloudaccounts/CLOUDACCOUNTID"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp11.Run(c, args) },
	}
	tmp11.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp11.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "update",
		Short: `update action`,
	}
	tmp12 := new(UpdateAccountCommand)
	sub = &cobra.Command{
		Use:   `account ["/accounts/ACCOUNTID"]`,
		Short: ``,
		Long: `

Payload example:

{
   "name": "test"
}`,
		RunE: func(cmd *cobra.Command, args []string) error { return tmp12.Run(c, args) },
	}
	tmp12.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp12.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	tmp13 := new(UpdateCloudaccountCommand)
	sub = &cobra.Command{
		Use:   `cloudaccount ["/accounts/ACCOUNTID/cloudaccounts/CLOUDACCOUNTID"]`,
		Short: ``,
		Long: `

Payload example:

{
   "assume_role_arn": "arn:aws:iam::788612350743:role/ZeroCloud",
   "assume_role_external_id": "bigdb",
   "cloudprovider": "AWS",
   "name": "BigDB.co's perf testing AWS account",
   "upstream_account_id": "98798079879"
}`,
		RunE: func(cmd *cobra.Command, args []string) error { return tmp13.Run(c, args) },
	}
	tmp13.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp13.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
}

func intFlagVal(name string, parsed int) *int {
	if hasFlag(name) {
		return &parsed
	}
	return nil
}

func float64FlagVal(name string, parsed float64) *float64 {
	if hasFlag(name) {
		return &parsed
	}
	return nil
}

func boolFlagVal(name string, parsed bool) *bool {
	if hasFlag(name) {
		return &parsed
	}
	return nil
}

func stringFlagVal(name string, parsed string) *string {
	if hasFlag(name) {
		return &parsed
	}
	return nil
}

func hasFlag(name string) bool {
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "--"+name) {
			return true
		}
	}
	return false
}

func jsonVal(val string) (*interface{}, error) {
	var t interface{}
	err := json.Unmarshal([]byte(val), &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func jsonArray(ins []string) ([]interface{}, error) {
	if ins == nil {
		return nil, nil
	}
	var vals []interface{}
	for _, id := range ins {
		val, err := jsonVal(id)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return vals, nil
}

func timeVal(val string) (*time.Time, error) {
	t, err := time.Parse("RFC3339", val)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func timeArray(ins []string) ([]time.Time, error) {
	if ins == nil {
		return nil, nil
	}
	var vals []time.Time
	for _, id := range ins {
		val, err := timeVal(id)
		if err != nil {
			return nil, err
		}
		vals = append(vals, *val)
	}
	return vals, nil
}

func uuidVal(val string) (*uuid.UUID, error) {
	t, err := uuid.FromString(val)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func uuidArray(ins []string) ([]uuid.UUID, error) {
	if ins == nil {
		return nil, nil
	}
	var vals []uuid.UUID
	for _, id := range ins {
		val, err := uuidVal(id)
		if err != nil {
			return nil, err
		}
		vals = append(vals, *val)
	}
	return vals, nil
}

func float64Val(val string) (*float64, error) {
	t, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func float64Array(ins []string) ([]float64, error) {
	if ins == nil {
		return nil, nil
	}
	var vals []float64
	for _, id := range ins {
		val, err := float64Val(id)
		if err != nil {
			return nil, err
		}
		vals = append(vals, *val)
	}
	return vals, nil
}

func boolVal(val string) (*bool, error) {
	t, err := strconv.ParseBool(val)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func boolArray(ins []string) ([]bool, error) {
	if ins == nil {
		return nil, nil
	}
	var vals []bool
	for _, id := range ins {
		val, err := boolVal(id)
		if err != nil {
			return nil, err
		}
		vals = append(vals, *val)
	}
	return vals, nil
}

// Run makes the HTTP request corresponding to the CreateAccountCommand command.
func (cmd *CreateAccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = "/accounts"
	}
	var payload client.AccountPayload
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &payload)
		if err != nil {
			return fmt.Errorf("failed to deserialize payload: %s", err)
		}
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.CreateAccount(ctx, path, &payload, cmd.ContentType)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *CreateAccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	cc.Flags().StringVar(&cmd.Payload, "payload", "", "Request body encoded in JSON")
	cc.Flags().StringVar(&cmd.ContentType, "content", "", "Request content type override, e.g. 'application/x-www-form-urlencoded'")
}

// Run makes the HTTP request corresponding to the DeleteAccountCommand command.
func (cmd *DeleteAccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v", cmd.AccountID)
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.DeleteAccount(ctx, path)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *DeleteAccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "accountID", accountID, `Account ID`)
}

// Run makes the HTTP request corresponding to the ListAccountCommand command.
func (cmd *ListAccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = "/accounts"
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.ListAccount(ctx, path)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *ListAccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
}

// Run makes the HTTP request corresponding to the ShowAccountCommand command.
func (cmd *ShowAccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v", cmd.AccountID)
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.ShowAccount(ctx, path)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *ShowAccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "accountID", accountID, `Account ID`)
}

// Run makes the HTTP request corresponding to the UpdateAccountCommand command.
func (cmd *UpdateAccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v", cmd.AccountID)
	}
	var payload client.UpdateAccountPayload
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &payload)
		if err != nil {
			return fmt.Errorf("failed to deserialize payload: %s", err)
		}
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.UpdateAccount(ctx, path, &payload, cmd.ContentType)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *UpdateAccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	cc.Flags().StringVar(&cmd.Payload, "payload", "", "Request body encoded in JSON")
	cc.Flags().StringVar(&cmd.ContentType, "content", "", "Request content type override, e.g. 'application/x-www-form-urlencoded'")
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "accountID", accountID, `Account ID`)
}

// Run makes the HTTP request corresponding to the ShowAwsCommand command.
func (cmd *ShowAwsCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/aws/%v", cmd.AwsAccountID)
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.ShowAws(ctx, path)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *ShowAwsCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var awsAccountID string
	cc.Flags().StringVar(&cmd.AwsAccountID, "awsAccountID", awsAccountID, ``)
}

// Run makes the HTTP request corresponding to the CreateCloudaccountCommand command.
func (cmd *CreateCloudaccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts", cmd.AccountID)
	}
	var payload client.CloudAccountPayload
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &payload)
		if err != nil {
			return fmt.Errorf("failed to deserialize payload: %s", err)
		}
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.CreateCloudaccount(ctx, path, &payload, cmd.ContentType)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *CreateCloudaccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	cc.Flags().StringVar(&cmd.Payload, "payload", "", "Request body encoded in JSON")
	cc.Flags().StringVar(&cmd.ContentType, "content", "", "Request content type override, e.g. 'application/x-www-form-urlencoded'")
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "accountID", accountID, `Account ID`)
}

// Run makes the HTTP request corresponding to the DeleteCloudaccountCommand command.
func (cmd *DeleteCloudaccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts/%v", cmd.AccountID, cmd.CloudAccountID)
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.DeleteCloudaccount(ctx, path)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *DeleteCloudaccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "accountID", accountID, `Account ID`)
	var cloudAccountID int
	cc.Flags().IntVar(&cmd.CloudAccountID, "cloudAccountID", cloudAccountID, ``)
}

// Run makes the HTTP request corresponding to the ListCloudaccountCommand command.
func (cmd *ListCloudaccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts", cmd.AccountID)
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.ListCloudaccount(ctx, path)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *ListCloudaccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "accountID", accountID, `Account ID`)
}

// Run makes the HTTP request corresponding to the ShowCloudaccountCommand command.
func (cmd *ShowCloudaccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts/%v", cmd.AccountID, cmd.CloudAccountID)
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.ShowCloudaccount(ctx, path)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *ShowCloudaccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "accountID", accountID, `Account ID`)
	var cloudAccountID int
	cc.Flags().IntVar(&cmd.CloudAccountID, "cloudAccountID", cloudAccountID, ``)
}

// Run makes the HTTP request corresponding to the UpdateCloudaccountCommand command.
func (cmd *UpdateCloudaccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts/%v", cmd.AccountID, cmd.CloudAccountID)
	}
	var payload client.CloudAccountPayload
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &payload)
		if err != nil {
			return fmt.Errorf("failed to deserialize payload: %s", err)
		}
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.UpdateCloudaccount(ctx, path, &payload, cmd.ContentType)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *UpdateCloudaccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	cc.Flags().StringVar(&cmd.Payload, "payload", "", "Request body encoded in JSON")
	cc.Flags().StringVar(&cmd.ContentType, "content", "", "Request content type override, e.g. 'application/x-www-form-urlencoded'")
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "accountID", accountID, `Account ID`)
	var cloudAccountID int
	cc.Flags().IntVar(&cmd.CloudAccountID, "cloudAccountID", cloudAccountID, ``)
}

// Run makes the HTTP request corresponding to the CreateCloudeventCommand command.
func (cmd *CreateCloudeventCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = "/cloudevent"
	}
	var payload client.CloudEventPayload
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &payload)
		if err != nil {
			return fmt.Errorf("failed to deserialize payload: %s", err)
		}
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.CreateCloudevent(ctx, path, &payload, cmd.ContentType)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *CreateCloudeventCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	cc.Flags().StringVar(&cmd.Payload, "payload", "", "Request body encoded in JSON")
	cc.Flags().StringVar(&cmd.ContentType, "content", "", "Request content type override, e.g. 'application/x-www-form-urlencoded'")
}

// Run makes the HTTP request corresponding to the ListLeaseCommand command.
func (cmd *ListLeaseCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = "/leases"
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.ListLease(ctx, path)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *ListLeaseCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
}
