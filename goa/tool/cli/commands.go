package cli

import (
	"encoding/json"
	"fmt"
	"github.com/goadesign/goa"
	goaclient "github.com/goadesign/goa/client"
	uuid "github.com/goadesign/goa/uuid"
	"github.com/spf13/cobra"
	"github.com/tleyden/cecil/goa/client"
	"golang.org/x/net/context"
	"log"
	"net/url"
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

	// MailerConfigAccountCommand is the command line data structure for the mailerConfig action of account
	MailerConfigAccountCommand struct {
		Payload     string
		ContentType string
		// Account ID
		AccountID   int
		PrettyPrint bool
	}

	// RemoveSlackAccountCommand is the command line data structure for the removeSlack action of account
	RemoveSlackAccountCommand struct {
		// Account ID
		AccountID   int
		PrettyPrint bool
	}

	// ShowAccountCommand is the command line data structure for the show action of account
	ShowAccountCommand struct {
		// Account ID
		AccountID   int
		PrettyPrint bool
	}

	// SlackConfigAccountCommand is the command line data structure for the slackConfig action of account
	SlackConfigAccountCommand struct {
		Payload     string
		ContentType string
		// Account ID
		AccountID   int
		PrettyPrint bool
	}

	// VerifyAccountCommand is the command line data structure for the verify action of account
	VerifyAccountCommand struct {
		Payload     string
		ContentType string
		// Account ID
		AccountID   int
		PrettyPrint bool
	}

	// AddCloudaccountCommand is the command line data structure for the add action of cloudaccount
	AddCloudaccountCommand struct {
		Payload     string
		ContentType string
		// Account ID
		AccountID   int
		PrettyPrint bool
	}

	// AddEmailToWhitelistCloudaccountCommand is the command line data structure for the addEmailToWhitelist action of cloudaccount
	AddEmailToWhitelistCloudaccountCommand struct {
		Payload     string
		ContentType string
		// Account ID
		AccountID int
		// CloudAccount ID
		CloudaccountID int
		PrettyPrint    bool
	}

	// DownloadInitialSetupTemplateCloudaccountCommand is the command line data structure for the downloadInitialSetupTemplate action of cloudaccount
	DownloadInitialSetupTemplateCloudaccountCommand struct {
		// Account ID
		AccountID int
		// CloudAccount ID
		CloudaccountID int
		PrettyPrint    bool
	}

	// DownloadRegionSetupTemplateCloudaccountCommand is the command line data structure for the downloadRegionSetupTemplate action of cloudaccount
	DownloadRegionSetupTemplateCloudaccountCommand struct {
		// Account ID
		AccountID int
		// CloudAccount ID
		CloudaccountID int
		PrettyPrint    bool
	}

	// ListRegionsCloudaccountCommand is the command line data structure for the listRegions action of cloudaccount
	ListRegionsCloudaccountCommand struct {
		// Account ID
		AccountID int
		// CloudAccount ID
		CloudaccountID int
		PrettyPrint    bool
	}

	// SubscribeSNSToSQSCloudaccountCommand is the command line data structure for the subscribeSNSToSQS action of cloudaccount
	SubscribeSNSToSQSCloudaccountCommand struct {
		Payload     string
		ContentType string
		// Account ID
		AccountID int
		// CloudAccount ID
		CloudaccountID int
		PrettyPrint    bool
	}

	// UpdateCloudaccountCommand is the command line data structure for the update action of cloudaccount
	UpdateCloudaccountCommand struct {
		Payload     string
		ContentType string
		// Account ID
		AccountID int
		// CloudAccount ID
		CloudaccountID int
		PrettyPrint    bool
	}

	// ActionsEmailActionCommand is the command line data structure for the actions action of email_action
	ActionsEmailActionCommand struct {
		// Action to be peformed on the lease
		Action string
		// ID of the lease
		InstanceID string
		// UUID of the lease
		LeaseUUID string
		// The signature of this link
		Sig string
		// The token_once of this link
		Tok         string
		PrettyPrint bool
	}

	// DeleteFromDBLeasesCommand is the command line data structure for the deleteFromDB action of leases
	DeleteFromDBLeasesCommand struct {
		// Account ID
		AccountID int
		// CloudAccount ID
		CloudaccountID int
		// Lease ID
		LeaseID     int
		PrettyPrint bool
	}

	// ListLeasesForAccountLeasesCommand is the command line data structure for the listLeasesForAccount action of leases
	ListLeasesForAccountLeasesCommand struct {
		// Account ID
		AccountID   int
		Terminated  string
		PrettyPrint bool
	}

	// ListLeasesForCloudaccountLeasesCommand is the command line data structure for the listLeasesForCloudaccount action of leases
	ListLeasesForCloudaccountLeasesCommand struct {
		// Account ID
		AccountID int
		// CloudAccount ID
		CloudaccountID int
		Terminated     string
		PrettyPrint    bool
	}

	// SetExpiryLeasesCommand is the command line data structure for the setExpiry action of leases
	SetExpiryLeasesCommand struct {
		// Account ID
		AccountID int
		// CloudAccount ID
		CloudaccountID int
		// Lease ID
		LeaseID int
		// Target expiry datetime
		ExpiresAt   string
		PrettyPrint bool
	}

	// ShowLeasesCommand is the command line data structure for the show action of leases
	ShowLeasesCommand struct {
		// Account ID
		AccountID int
		// CloudAccount ID
		CloudaccountID int
		// Lease ID
		LeaseID     int
		PrettyPrint bool
	}

	// TerminateLeasesCommand is the command line data structure for the terminate action of leases
	TerminateLeasesCommand struct {
		// Account ID
		AccountID int
		// CloudAccount ID
		CloudaccountID int
		// Lease ID
		LeaseID     int
		PrettyPrint bool
	}

	// ShowRootCommand is the command line data structure for the show action of root
	ShowRootCommand struct {
		PrettyPrint bool
	}

	// DownloadCommand is the command line data structure for the download command.
	DownloadCommand struct {
		// OutFile is the path to the download output file.
		OutFile string
	}
)

// RegisterCommands registers the resource action CLI commands.
func RegisterCommands(app *cobra.Command, c *client.Client) {
	var command, sub *cobra.Command
	command = &cobra.Command{
		Use:   "actions",
		Short: `Perform an action on a lease`,
	}
	tmp1 := new(ActionsEmailActionCommand)
	sub = &cobra.Command{
		Use:   `email_action ["/email_action/leases/LEASE_UUID/INSTANCE_ID/ACTION"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp1.Run(c, args) },
	}
	tmp1.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp1.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "add",
		Short: `Add new cloudaccount`,
	}
	tmp2 := new(AddCloudaccountCommand)
	sub = &cobra.Command{
		Use:   `cloudaccount ["/accounts/ACCOUNT_ID/cloudaccounts"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp2.Run(c, args) },
	}
	tmp2.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp2.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "addEmailToWhitelist",
		Short: `Add new email to owner tag whitelist`,
	}
	tmp3 := new(AddEmailToWhitelistCloudaccountCommand)
	sub = &cobra.Command{
		Use:   `cloudaccount ["/accounts/ACCOUNT_ID/cloudaccounts/CLOUDACCOUNT_ID/owners"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp3.Run(c, args) },
	}
	tmp3.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp3.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "create",
		Short: `Create new account`,
	}
	tmp4 := new(CreateAccountCommand)
	sub = &cobra.Command{
		Use:   `account ["/accounts"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp4.Run(c, args) },
	}
	tmp4.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp4.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "deleteFromDB",
		Short: `Delete a lease from DB`,
	}
	tmp5 := new(DeleteFromDBLeasesCommand)
	sub = &cobra.Command{
		Use:   `leases [("/accounts/ACCOUNT_ID/cloudaccounts/CLOUDACCOUNT_ID/leases/LEASE_ID/delete"|"/accounts/ACCOUNT_ID/leases/LEASE_ID/delete")]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp5.Run(c, args) },
	}
	tmp5.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp5.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "downloadInitialSetupTemplate",
		Short: `Download AWS initial setup cloudformation template`,
	}
	tmp6 := new(DownloadInitialSetupTemplateCloudaccountCommand)
	sub = &cobra.Command{
		Use:   `cloudaccount ["/accounts/ACCOUNT_ID/cloudaccounts/CLOUDACCOUNT_ID/tenant-aws-initial-setup.template"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp6.Run(c, args) },
	}
	tmp6.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp6.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "downloadRegionSetupTemplate",
		Short: `Download AWS region setup cloudformation template`,
	}
	tmp7 := new(DownloadRegionSetupTemplateCloudaccountCommand)
	sub = &cobra.Command{
		Use:   `cloudaccount ["/accounts/ACCOUNT_ID/cloudaccounts/CLOUDACCOUNT_ID/tenant-aws-region-setup.template"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp7.Run(c, args) },
	}
	tmp7.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp7.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "listLeasesForAccount",
		Short: `List all leases for account`,
	}
	tmp8 := new(ListLeasesForAccountLeasesCommand)
	sub = &cobra.Command{
		Use:   `leases ["/accounts/ACCOUNT_ID/leases"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp8.Run(c, args) },
	}
	tmp8.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp8.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "listLeasesForCloudaccount",
		Short: `List all leases for cloudAccount`,
	}
	tmp9 := new(ListLeasesForCloudaccountLeasesCommand)
	sub = &cobra.Command{
		Use:   `leases ["/accounts/ACCOUNT_ID/cloudaccounts/CLOUDACCOUNT_ID/leases"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp9.Run(c, args) },
	}
	tmp9.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp9.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "listRegions",
		Short: `List all regions and their status`,
	}
	tmp10 := new(ListRegionsCloudaccountCommand)
	sub = &cobra.Command{
		Use:   `cloudaccount ["/accounts/ACCOUNT_ID/cloudaccounts/CLOUDACCOUNT_ID/regions"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp10.Run(c, args) },
	}
	tmp10.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp10.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "mailerConfig",
		Short: `Configure mailer`,
	}
	tmp11 := new(MailerConfigAccountCommand)
	sub = &cobra.Command{
		Use:   `account ["/accounts/ACCOUNT_ID/mailer_config"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp11.Run(c, args) },
	}
	tmp11.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp11.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "removeSlack",
		Short: `Remove slack`,
	}
	tmp12 := new(RemoveSlackAccountCommand)
	sub = &cobra.Command{
		Use:   `account ["/accounts/ACCOUNT_ID/slack_config"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp12.Run(c, args) },
	}
	tmp12.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp12.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "setExpiry",
		Short: `Set expiry of a lease`,
	}
	tmp13 := new(SetExpiryLeasesCommand)
	sub = &cobra.Command{
		Use:   `leases [("/accounts/ACCOUNT_ID/cloudaccounts/CLOUDACCOUNT_ID/leases/LEASE_ID/expiry"|"/accounts/ACCOUNT_ID/leases/LEASE_ID/expiry")]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp13.Run(c, args) },
	}
	tmp13.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp13.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "show",
		Short: `show action`,
	}
	tmp14 := new(ShowAccountCommand)
	sub = &cobra.Command{
		Use:   `account ["/accounts/ACCOUNT_ID"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp14.Run(c, args) },
	}
	tmp14.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp14.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	tmp15 := new(ShowLeasesCommand)
	sub = &cobra.Command{
		Use:   `leases [("/accounts/ACCOUNT_ID/cloudaccounts/CLOUDACCOUNT_ID/leases/LEASE_ID"|"/accounts/ACCOUNT_ID/leases/LEASE_ID")]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp15.Run(c, args) },
	}
	tmp15.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp15.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	tmp16 := new(ShowRootCommand)
	sub = &cobra.Command{
		Use:   `root ["/"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp16.Run(c, args) },
	}
	tmp16.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp16.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "slackConfig",
		Short: `Configure slack`,
	}
	tmp17 := new(SlackConfigAccountCommand)
	sub = &cobra.Command{
		Use:   `account ["/accounts/ACCOUNT_ID/slack_config"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp17.Run(c, args) },
	}
	tmp17.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp17.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "subscribeSNSToSQS",
		Short: `Subscribe SNS to SQS`,
	}
	tmp18 := new(SubscribeSNSToSQSCloudaccountCommand)
	sub = &cobra.Command{
		Use:   `cloudaccount ["/accounts/ACCOUNT_ID/cloudaccounts/CLOUDACCOUNT_ID/subscribe-sns-to-sqs"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp18.Run(c, args) },
	}
	tmp18.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp18.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "terminate",
		Short: `Terminate a lease`,
	}
	tmp19 := new(TerminateLeasesCommand)
	sub = &cobra.Command{
		Use:   `leases [("/accounts/ACCOUNT_ID/cloudaccounts/CLOUDACCOUNT_ID/leases/LEASE_ID/terminate"|"/accounts/ACCOUNT_ID/leases/LEASE_ID/terminate")]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp19.Run(c, args) },
	}
	tmp19.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp19.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "update",
		Short: `Update a cloudaccount`,
	}
	tmp20 := new(UpdateCloudaccountCommand)
	sub = &cobra.Command{
		Use:   `cloudaccount ["/accounts/ACCOUNT_ID/cloudaccounts/CLOUDACCOUNT_ID"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp20.Run(c, args) },
	}
	tmp20.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp20.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "verify",
		Short: `Verify account and get API token`,
	}
	tmp21 := new(VerifyAccountCommand)
	sub = &cobra.Command{
		Use:   `account ["/accounts/ACCOUNT_ID/api_token"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp21.Run(c, args) },
	}
	tmp21.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp21.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)

	dl := new(DownloadCommand)
	dlc := &cobra.Command{
		Use:   "download [PATH]",
		Short: "Download file with given path",
		RunE: func(cmd *cobra.Command, args []string) error {
			return dl.Run(c, args)
		},
	}
	dlc.Flags().StringVar(&dl.OutFile, "out", "", "Output file")
	app.AddCommand(dlc)
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
	t, err := time.Parse(time.RFC3339, val)
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

// Run downloads files with given paths.
func (cmd *DownloadCommand) Run(c *client.Client, args []string) error {
	var (
		fnf func(context.Context, string) (int64, error)
		fnd func(context.Context, string, string) (int64, error)

		rpath   = args[0]
		outfile = cmd.OutFile
		logger  = goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
		ctx     = goa.WithLogger(context.Background(), logger)
		err     error
	)

	if rpath[0] != '/' {
		rpath = "/" + rpath
	}
	if rpath == "/swagger.json" {
		fnf = c.DownloadSwaggerJSON
		if outfile == "" {
			outfile = "swagger.json"
		}
		goto found
	}
	return fmt.Errorf("don't know how to download %s", rpath)
found:
	ctx = goa.WithLogContext(ctx, "file", outfile)
	if fnf != nil {
		_, err = fnf(ctx, outfile)
	} else {
		_, err = fnd(ctx, rpath, outfile)
	}
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	return nil
}

// Run makes the HTTP request corresponding to the CreateAccountCommand command.
func (cmd *CreateAccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = "/accounts"
	}
	var payload client.CreateAccountPayload
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &payload)
		if err != nil {
			return fmt.Errorf("failed to deserialize payload: %s", err)
		}
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.CreateAccount(ctx, path, &payload)
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

// Run makes the HTTP request corresponding to the MailerConfigAccountCommand command.
func (cmd *MailerConfigAccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/mailer_config", cmd.AccountID)
	}
	var payload client.MailerConfigAccountPayload
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &payload)
		if err != nil {
			return fmt.Errorf("failed to deserialize payload: %s", err)
		}
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.MailerConfigAccount(ctx, path, &payload)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *MailerConfigAccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	cc.Flags().StringVar(&cmd.Payload, "payload", "", "Request body encoded in JSON")
	cc.Flags().StringVar(&cmd.ContentType, "content", "", "Request content type override, e.g. 'application/x-www-form-urlencoded'")
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
}

// Run makes the HTTP request corresponding to the RemoveSlackAccountCommand command.
func (cmd *RemoveSlackAccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/slack_config", cmd.AccountID)
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.RemoveSlackAccount(ctx, path)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *RemoveSlackAccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
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
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
}

// Run makes the HTTP request corresponding to the SlackConfigAccountCommand command.
func (cmd *SlackConfigAccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/slack_config", cmd.AccountID)
	}
	var payload client.SlackConfigAccountPayload
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &payload)
		if err != nil {
			return fmt.Errorf("failed to deserialize payload: %s", err)
		}
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.SlackConfigAccount(ctx, path, &payload)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *SlackConfigAccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	cc.Flags().StringVar(&cmd.Payload, "payload", "", "Request body encoded in JSON")
	cc.Flags().StringVar(&cmd.ContentType, "content", "", "Request content type override, e.g. 'application/x-www-form-urlencoded'")
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
}

// Run makes the HTTP request corresponding to the VerifyAccountCommand command.
func (cmd *VerifyAccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/api_token", cmd.AccountID)
	}
	var payload client.VerifyAccountPayload
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &payload)
		if err != nil {
			return fmt.Errorf("failed to deserialize payload: %s", err)
		}
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.VerifyAccount(ctx, path, &payload)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *VerifyAccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	cc.Flags().StringVar(&cmd.Payload, "payload", "", "Request body encoded in JSON")
	cc.Flags().StringVar(&cmd.ContentType, "content", "", "Request content type override, e.g. 'application/x-www-form-urlencoded'")
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
}

// Run makes the HTTP request corresponding to the AddCloudaccountCommand command.
func (cmd *AddCloudaccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts", cmd.AccountID)
	}
	var payload client.AddCloudaccountPayload
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &payload)
		if err != nil {
			return fmt.Errorf("failed to deserialize payload: %s", err)
		}
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.AddCloudaccount(ctx, path, &payload)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *AddCloudaccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	cc.Flags().StringVar(&cmd.Payload, "payload", "", "Request body encoded in JSON")
	cc.Flags().StringVar(&cmd.ContentType, "content", "", "Request content type override, e.g. 'application/x-www-form-urlencoded'")
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
}

// Run makes the HTTP request corresponding to the AddEmailToWhitelistCloudaccountCommand command.
func (cmd *AddEmailToWhitelistCloudaccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts/%v/owners", cmd.AccountID, cmd.CloudaccountID)
	}
	var payload client.AddEmailToWhitelistCloudaccountPayload
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &payload)
		if err != nil {
			return fmt.Errorf("failed to deserialize payload: %s", err)
		}
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.AddEmailToWhitelistCloudaccount(ctx, path, &payload)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *AddEmailToWhitelistCloudaccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	cc.Flags().StringVar(&cmd.Payload, "payload", "", "Request body encoded in JSON")
	cc.Flags().StringVar(&cmd.ContentType, "content", "", "Request content type override, e.g. 'application/x-www-form-urlencoded'")
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
	var cloudaccountID int
	cc.Flags().IntVar(&cmd.CloudaccountID, "cloudaccount_id", cloudaccountID, `CloudAccount ID`)
}

// Run makes the HTTP request corresponding to the DownloadInitialSetupTemplateCloudaccountCommand command.
func (cmd *DownloadInitialSetupTemplateCloudaccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts/%v/tenant-aws-initial-setup.template", cmd.AccountID, cmd.CloudaccountID)
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.DownloadInitialSetupTemplateCloudaccount(ctx, path)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *DownloadInitialSetupTemplateCloudaccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
	var cloudaccountID int
	cc.Flags().IntVar(&cmd.CloudaccountID, "cloudaccount_id", cloudaccountID, `CloudAccount ID`)
}

// Run makes the HTTP request corresponding to the DownloadRegionSetupTemplateCloudaccountCommand command.
func (cmd *DownloadRegionSetupTemplateCloudaccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts/%v/tenant-aws-region-setup.template", cmd.AccountID, cmd.CloudaccountID)
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.DownloadRegionSetupTemplateCloudaccount(ctx, path)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *DownloadRegionSetupTemplateCloudaccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
	var cloudaccountID int
	cc.Flags().IntVar(&cmd.CloudaccountID, "cloudaccount_id", cloudaccountID, `CloudAccount ID`)
}

// Run makes the HTTP request corresponding to the ListRegionsCloudaccountCommand command.
func (cmd *ListRegionsCloudaccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts/%v/regions", cmd.AccountID, cmd.CloudaccountID)
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.ListRegionsCloudaccount(ctx, path)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *ListRegionsCloudaccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
	var cloudaccountID int
	cc.Flags().IntVar(&cmd.CloudaccountID, "cloudaccount_id", cloudaccountID, `CloudAccount ID`)
}

// Run makes the HTTP request corresponding to the SubscribeSNSToSQSCloudaccountCommand command.
func (cmd *SubscribeSNSToSQSCloudaccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts/%v/subscribe-sns-to-sqs", cmd.AccountID, cmd.CloudaccountID)
	}
	var payload client.SubscribeSNSToSQSCloudaccountPayload
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &payload)
		if err != nil {
			return fmt.Errorf("failed to deserialize payload: %s", err)
		}
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.SubscribeSNSToSQSCloudaccount(ctx, path, &payload)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *SubscribeSNSToSQSCloudaccountCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	cc.Flags().StringVar(&cmd.Payload, "payload", "", "Request body encoded in JSON")
	cc.Flags().StringVar(&cmd.ContentType, "content", "", "Request content type override, e.g. 'application/x-www-form-urlencoded'")
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
	var cloudaccountID int
	cc.Flags().IntVar(&cmd.CloudaccountID, "cloudaccount_id", cloudaccountID, `CloudAccount ID`)
}

// Run makes the HTTP request corresponding to the UpdateCloudaccountCommand command.
func (cmd *UpdateCloudaccountCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts/%v", cmd.AccountID, cmd.CloudaccountID)
	}
	var payload client.UpdateCloudaccountPayload
	if cmd.Payload != "" {
		err := json.Unmarshal([]byte(cmd.Payload), &payload)
		if err != nil {
			return fmt.Errorf("failed to deserialize payload: %s", err)
		}
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.UpdateCloudaccount(ctx, path, &payload)
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
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
	var cloudaccountID int
	cc.Flags().IntVar(&cmd.CloudaccountID, "cloudaccount_id", cloudaccountID, `CloudAccount ID`)
}

// Run makes the HTTP request corresponding to the ActionsEmailActionCommand command.
func (cmd *ActionsEmailActionCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/email_action/leases/%v/%v/%v", cmd.LeaseUUID, url.QueryEscape(cmd.InstanceID), url.QueryEscape(cmd.Action))
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.ActionsEmailAction(ctx, path, cmd.Sig, cmd.Tok)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *ActionsEmailActionCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var action string
	cc.Flags().StringVar(&cmd.Action, "action", action, `Action to be peformed on the lease`)
	var instanceID string
	cc.Flags().StringVar(&cmd.InstanceID, "instance_id", instanceID, `ID of the lease`)
	var leaseUUID string
	cc.Flags().StringVar(&cmd.LeaseUUID, "lease_uuid", leaseUUID, `UUID of the lease`)
	var sig string
	cc.Flags().StringVar(&cmd.Sig, "sig", sig, `The signature of this link`)
	var tok string
	cc.Flags().StringVar(&cmd.Tok, "tok", tok, `The token_once of this link`)
}

// Run makes the HTTP request corresponding to the DeleteFromDBLeasesCommand command.
func (cmd *DeleteFromDBLeasesCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts/%v/leases/%v/delete", cmd.AccountID, cmd.CloudaccountID, cmd.LeaseID)
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.DeleteFromDBLeases(ctx, path)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *DeleteFromDBLeasesCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
	var cloudaccountID int
	cc.Flags().IntVar(&cmd.CloudaccountID, "cloudaccount_id", cloudaccountID, `CloudAccount ID`)
	var leaseID int
	cc.Flags().IntVar(&cmd.LeaseID, "lease_id", leaseID, `Lease ID`)
}

// Run makes the HTTP request corresponding to the ListLeasesForAccountLeasesCommand command.
func (cmd *ListLeasesForAccountLeasesCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/leases", cmd.AccountID)
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	var tmp22 *bool
	if cmd.Terminated != "" {
		var err error
		tmp22, err = boolVal(cmd.Terminated)
		if err != nil {
			goa.LogError(ctx, "failed to parse flag into *bool value", "flag", "--terminated", "err", err)
			return err
		}
	}
	resp, err := c.ListLeasesForAccountLeases(ctx, path, tmp22)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *ListLeasesForAccountLeasesCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
	var terminated string
	cc.Flags().StringVar(&cmd.Terminated, "terminated", terminated, ``)
}

// Run makes the HTTP request corresponding to the ListLeasesForCloudaccountLeasesCommand command.
func (cmd *ListLeasesForCloudaccountLeasesCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts/%v/leases", cmd.AccountID, cmd.CloudaccountID)
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	var tmp23 *bool
	if cmd.Terminated != "" {
		var err error
		tmp23, err = boolVal(cmd.Terminated)
		if err != nil {
			goa.LogError(ctx, "failed to parse flag into *bool value", "flag", "--terminated", "err", err)
			return err
		}
	}
	resp, err := c.ListLeasesForCloudaccountLeases(ctx, path, tmp23)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *ListLeasesForCloudaccountLeasesCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
	var cloudaccountID int
	cc.Flags().IntVar(&cmd.CloudaccountID, "cloudaccount_id", cloudaccountID, `CloudAccount ID`)
	var terminated string
	cc.Flags().StringVar(&cmd.Terminated, "terminated", terminated, ``)
}

// Run makes the HTTP request corresponding to the SetExpiryLeasesCommand command.
func (cmd *SetExpiryLeasesCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts/%v/leases/%v/expiry", cmd.AccountID, cmd.CloudaccountID, cmd.LeaseID)
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	var tmp24 *time.Time
	if cmd.ExpiresAt != "" {
		var err error
		tmp24, err = timeVal(cmd.ExpiresAt)
		if err != nil {
			goa.LogError(ctx, "failed to parse flag into *time.Time value", "flag", "--expires_at", "err", err)
			return err
		}
	}
	resp, err := c.SetExpiryLeases(ctx, path, *tmp24)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *SetExpiryLeasesCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
	var cloudaccountID int
	cc.Flags().IntVar(&cmd.CloudaccountID, "cloudaccount_id", cloudaccountID, `CloudAccount ID`)
	var leaseID int
	cc.Flags().IntVar(&cmd.LeaseID, "lease_id", leaseID, `Lease ID`)
	var expiresAt string
	cc.Flags().StringVar(&cmd.ExpiresAt, "expires_at", expiresAt, `Target expiry datetime`)
}

// Run makes the HTTP request corresponding to the ShowLeasesCommand command.
func (cmd *ShowLeasesCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts/%v/leases/%v", cmd.AccountID, cmd.CloudaccountID, cmd.LeaseID)
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.ShowLeases(ctx, path)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *ShowLeasesCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
	var cloudaccountID int
	cc.Flags().IntVar(&cmd.CloudaccountID, "cloudaccount_id", cloudaccountID, `CloudAccount ID`)
	var leaseID int
	cc.Flags().IntVar(&cmd.LeaseID, "lease_id", leaseID, `Lease ID`)
}

// Run makes the HTTP request corresponding to the TerminateLeasesCommand command.
func (cmd *TerminateLeasesCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/accounts/%v/cloudaccounts/%v/leases/%v/terminate", cmd.AccountID, cmd.CloudaccountID, cmd.LeaseID)
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.TerminateLeases(ctx, path)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *TerminateLeasesCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
	var accountID int
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account ID`)
	var cloudaccountID int
	cc.Flags().IntVar(&cmd.CloudaccountID, "cloudaccount_id", cloudaccountID, `CloudAccount ID`)
	var leaseID int
	cc.Flags().IntVar(&cmd.LeaseID, "lease_id", leaseID, `Lease ID`)
}

// Run makes the HTTP request corresponding to the ShowRootCommand command.
func (cmd *ShowRootCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = "/"
	}
	logger := goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	ctx := goa.WithLogger(context.Background(), logger)
	resp, err := c.ShowRoot(ctx, path)
	if err != nil {
		goa.LogError(ctx, "failed", "err", err)
		return err
	}

	goaclient.HandleResponse(c.Client, resp, cmd.PrettyPrint)
	return nil
}

// RegisterFlags registers the command flags with the command line.
func (cmd *ShowRootCommand) RegisterFlags(cc *cobra.Command, c *client.Client) {
}
