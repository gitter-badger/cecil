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

	// ShowAccountCommand is the command line data structure for the show action of account
	ShowAccountCommand struct {
		// Account Id
		AccountID   int
		PrettyPrint bool
	}

	// VerifyAccountCommand is the command line data structure for the verify action of account
	VerifyAccountCommand struct {
		Payload     string
		ContentType string
		// Account Id
		AccountID   int
		PrettyPrint bool
	}

	// AddCloudaccountCommand is the command line data structure for the add action of cloudaccount
	AddCloudaccountCommand struct {
		Payload     string
		ContentType string
		// Account Id
		AccountID   int
		PrettyPrint bool
	}

	// AddEmailToWhitelistCloudaccountCommand is the command line data structure for the addEmailToWhitelist action of cloudaccount
	AddEmailToWhitelistCloudaccountCommand struct {
		Payload     string
		ContentType string
		// Account Id
		AccountID int
		// CloudAccount Id
		CloudaccountID int
		PrettyPrint    bool
	}

	// DownloadInitialSetupTemplateCloudaccountCommand is the command line data structure for the downloadInitialSetupTemplate action of cloudaccount
	DownloadInitialSetupTemplateCloudaccountCommand struct {
		// Account Id
		AccountID int
		// CloudAccount Id
		CloudaccountID int
		PrettyPrint    bool
	}

	// DownloadRegionSetupTemplateCloudaccountCommand is the command line data structure for the downloadRegionSetupTemplate action of cloudaccount
	DownloadRegionSetupTemplateCloudaccountCommand struct {
		// Account Id
		AccountID int
		// CloudAccount Id
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
		Use:   "downloadInitialSetupTemplate",
		Short: `Download AWS initial setup cloudformation template`,
	}
	tmp5 := new(DownloadInitialSetupTemplateCloudaccountCommand)
	sub = &cobra.Command{
		Use:   `cloudaccount ["/accounts/ACCOUNT_ID/cloudaccounts/CLOUDACCOUNT_ID/tenant-aws-initial-setup.template"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp5.Run(c, args) },
	}
	tmp5.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp5.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "downloadRegionSetupTemplate",
		Short: `Download AWS region setup cloudformation template`,
	}
	tmp6 := new(DownloadRegionSetupTemplateCloudaccountCommand)
	sub = &cobra.Command{
		Use:   `cloudaccount ["/accounts/ACCOUNT_ID/cloudaccounts/CLOUDACCOUNT_ID/tenant-aws-region-setup.template"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp6.Run(c, args) },
	}
	tmp6.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp6.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "show",
		Short: `Show account`,
	}
	tmp7 := new(ShowAccountCommand)
	sub = &cobra.Command{
		Use:   `account ["/accounts/ACCOUNT_ID"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp7.Run(c, args) },
	}
	tmp7.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp7.PrettyPrint, "pp", false, "Pretty print response body")
	command.AddCommand(sub)
	app.AddCommand(command)
	command = &cobra.Command{
		Use:   "verify",
		Short: `Verify account and get API token`,
	}
	tmp8 := new(VerifyAccountCommand)
	sub = &cobra.Command{
		Use:   `account ["/accounts/ACCOUNT_ID/api_token"]`,
		Short: ``,
		RunE:  func(cmd *cobra.Command, args []string) error { return tmp8.Run(c, args) },
	}
	tmp8.RegisterFlags(sub, c)
	sub.PersistentFlags().BoolVar(&tmp8.PrettyPrint, "pp", false, "Pretty print response body")
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
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account Id`)
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
	resp, err := c.VerifyAccount(ctx, path, &payload, cmd.ContentType)
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
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account Id`)
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
	resp, err := c.AddCloudaccount(ctx, path, &payload, cmd.ContentType)
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
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account Id`)
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
	resp, err := c.AddEmailToWhitelistCloudaccount(ctx, path, &payload, cmd.ContentType)
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
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account Id`)
	var cloudaccountID int
	cc.Flags().IntVar(&cmd.CloudaccountID, "cloudaccount_id", cloudaccountID, `CloudAccount Id`)
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
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account Id`)
	var cloudaccountID int
	cc.Flags().IntVar(&cmd.CloudaccountID, "cloudaccount_id", cloudaccountID, `CloudAccount Id`)
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
	cc.Flags().IntVar(&cmd.AccountID, "account_id", accountID, `Account Id`)
	var cloudaccountID int
	cc.Flags().IntVar(&cmd.CloudaccountID, "cloudaccount_id", cloudaccountID, `CloudAccount Id`)
}

// Run makes the HTTP request corresponding to the ActionsEmailActionCommand command.
func (cmd *ActionsEmailActionCommand) Run(c *client.Client, args []string) error {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = fmt.Sprintf("/email_action/leases/%v/%v/%v", cmd.Action, cmd.InstanceID, cmd.LeaseUUID)
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
