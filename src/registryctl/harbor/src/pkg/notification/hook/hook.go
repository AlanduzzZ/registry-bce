package hook

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cJob "github.com/goharbor/harbor/src/common/job"
	"github.com/goharbor/harbor/src/common/job/models"
	cModels "github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/core/utils"
	"github.com/goharbor/harbor/src/lib/config"
	"github.com/goharbor/harbor/src/lib/log"
	"github.com/goharbor/harbor/src/pkg/notification/job"
	job_model "github.com/goharbor/harbor/src/pkg/notification/job/model"
	"github.com/goharbor/harbor/src/pkg/notifier/model"
)

// Manager send hook
type Manager interface {
	StartHook(context.Context, *model.HookEvent, *models.JobData) error
}

// DefaultManager ...
type DefaultManager struct {
	jobMgr job.Manager
	client cJob.Client
}

// NewHookManager ...
func NewHookManager() *DefaultManager {
	return &DefaultManager{
		jobMgr: job.NewManager(),
		client: utils.GetJobServiceClient(),
	}
}

// StartHook create a notification job record in database, and submit it to jobservice
func (hm *DefaultManager) StartHook(ctx context.Context, event *model.HookEvent, data *models.JobData) error {
	payload, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}

	t := time.Now()
	id, err := hm.jobMgr.Create(ctx, &job_model.Job{
		PolicyID:     event.PolicyID,
		EventType:    event.EventType,
		NotifyType:   event.Target.Type,
		Status:       cModels.JobPending,
		CreationTime: t,
		UpdateTime:   t,
		JobDetail:    string(payload),
	})
	if err != nil {
		return fmt.Errorf("failed to create the job record for notification based on policy %d: %v", event.PolicyID, err)
	}
	statusHookURL := fmt.Sprintf("%s/service/notifications/jobs/webhook/%d", config.InternalCoreURL(), id)
	data.StatusHook = statusHookURL

	log.Debugf("created a notification job %d for the policy %d", id, event.PolicyID)

	// submit hook job to jobservice
	jobUUID, err := hm.client.SubmitJob(data)
	if err != nil {
		log.Errorf("failed to submit job with notification event: %v", err)
		e := hm.jobMgr.Update(ctx, &job_model.Job{
			ID:     id,
			Status: cModels.JobError,
		}, "Status")
		if e != nil {
			log.Errorf("failed to update the notification job status %d: %v", id, e)
		}
		return err
	}

	if err = hm.jobMgr.Update(ctx, &job_model.Job{
		ID:   id,
		UUID: jobUUID,
	}, "UUID"); err != nil {
		log.Errorf("failed to update the notification job %d: %v", id, err)
		return err
	}
	return nil
}
