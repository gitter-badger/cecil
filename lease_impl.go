package main

import "github.com/tleyden/zerocloud/app"

// List runs the list action.
func (c *LeaseController) ListImpl(ctx *app.ListLeaseContext) error {

	// ListLease(ctx context.Context, accountID int, cloudAccountID int, cloudEventID int) []*app.Lease {

	leaseModels, err := ldb.List(ctx.Context)
	if err != nil {
		return err
	}
	leaseCollection := app.LeaseCollection{}
	for _, leaseModel := range leaseModels {
		lease := leaseModel.LeaseToLease()
		leaseCollection = append(leaseCollection, lease)
	}

	return ctx.OK(leaseCollection)

}
