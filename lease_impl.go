package main

import (
	"fmt"

	"github.com/tleyden/zerocloud/app"
)

// List runs the list action.
func (c *LeaseController) ListImpl(ctx *app.ListLeaseContext) error {

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

		// lookup associated account manually
		// TODO: get goa/gorma to do this automatically
		account, err := adb.OneAccount(ctx.Context, leaseModel.AccountID)
		if err != nil {
			return err
		}
		lease.Account = account

		// lookup associated cloudaccount
		cloudAccount, err := cdb.OneCloudaccount(ctx.Context, leaseModel.CloudAccountID, leaseModel.AccountID)
		if err != nil {
			return err
		}
		lease.CloudAccount = cloudAccount

		// lookup associated cloudevent
		cloudEvent, err := edb.OneCloudevent(ctx.Context, leaseModel.CloudEventID, leaseModel.CloudAccountID, leaseModel.AccountID)
		if err != nil {
			return err
		}
		lease.CloudEvent = cloudEvent

		// lookup

		logger.Info("lease", "lease", fmt.Sprintf("%+v", lease))
		leaseCollection = append(leaseCollection, lease)
	}

	return ctx.OK(leaseCollection)

}
