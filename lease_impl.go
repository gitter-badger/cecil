package main

import (
	"fmt"

	"github.com/tleyden/zerocloud/app"
)

// List runs the list action.
func (c *LeaseController) ListImpl(ctx *app.ListLeaseContext) error {

	// ListLease(ctx context.Context, accountID int, cloudAccountID int, cloudEventID int) []*app.Lease {

	leaseModels, err := ldb.List(ctx.Context)
	if err != nil {
		return err
	}
	leaseCollection := app.LeaseCollection{}
	for _, leaseModel := range leaseModels {
		logger.Info("leaseModel", "leaseModel", fmt.Sprintf("%+v", leaseModel))
		lease := leaseModel.LeaseToLease()
		lease.AccountID = leaseModel.AccountID
		lease.CloudEventID = leaseModel.CloudEventID
		lease.CloudAccountID = leaseModel.CloudAccountID

		// 	account, err := adb.OneAccount(ctx.Context, ctx.AccountID)
		account, err := adb.OneAccount(ctx.Context, leaseModel.AccountID)
		if err != nil {
			return err
		}
		lease.Account = account

		logger.Info("lease", "lease", fmt.Sprintf("%+v", lease))
		leaseCollection = append(leaseCollection, lease)
	}

	return ctx.OK(leaseCollection)

}
