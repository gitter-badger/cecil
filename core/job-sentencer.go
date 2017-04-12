package core

import (
	"time"

	"github.com/tleyden/cecil/models"
	"github.com/tleyden/cecil/tasks"
)

// SentencerJob polls the DB for expired leases and pushes them to the TerminatorQueue
func (s *Service) SentencerJob() error {

	var expiredLeases []models.Lease
	var expiredLeasesCount int64

	s.DB.Table("leases").
		Where("expires_at < ? AND terminated_at IS NULL", time.Now().UTC()).
		Or("approved_at IS NULL AND launched_at < ? AND terminated_at IS NULL", time.Now().UTC().Add(-s.Config().Lease.ApprovalTimeoutDuration)).
		Find(&expiredLeases).
		Count(&expiredLeasesCount)

	Logger.Info("SentencerJob(): Expired leases", "count", expiredLeasesCount)

	for _, expiredLease := range expiredLeases {
		Logger.Info("expired lease",
			"lease_id", expiredLease.ID,
			"group_type", expiredLease.GroupType.String(),
			"group_uid", expiredLease.GroupUID,
		)
		s.Queues().TerminatorQueue().PushTask(tasks.TerminatorTask{Lease: expiredLease})
	}

	return nil
}
