package util

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goharbor/harbor/src/lib/config"
	"github.com/goharbor/harbor/src/lib/log"
	"github.com/goharbor/harbor/src/lib/orm"
	"github.com/goharbor/harbor/src/pkg/distribution"
	policy_model "github.com/goharbor/harbor/src/pkg/notification/policy/model"
	"github.com/goharbor/harbor/src/pkg/notifier/event"
	notifyModel "github.com/goharbor/harbor/src/pkg/notifier/model"
)

// SendHookWithPolicies send hook by publishing topic of specified target type(notify type)
func SendHookWithPolicies(policies []*policy_model.Policy, payload *notifyModel.Payload, eventType string) error {
	// if global notification configured disabled, return directly
	if !config.NotificationEnable(orm.Context()) {
		log.Debug("notification feature is not enabled")
		return nil
	}

	errRet := false
	for _, ply := range policies {
		targets := ply.Targets
		for _, target := range targets {
			evt := &event.Event{}
			hookMetadata := &event.HookMetaData{
				EventType: eventType,
				PolicyID:  ply.ID,
				Payload:   payload,
				Target:    &target,
			}
			// It should never affect evaluating other policies when one is failed, but error should return
			if err := evt.Build(hookMetadata); err == nil {
				if err := evt.Publish(); err != nil {
					errRet = true
					log.Errorf("failed to publish hook notify event: %v", err)
				}
			} else {
				errRet = true
				log.Errorf("failed to build hook notify event metadata: %v", err)
			}
			log.Debugf("published image event %s by topic %s", payload.Type, target.Type)
		}
	}
	if errRet {
		return errors.New("failed to trigger some of the events")
	}
	return nil
}

// GetNameFromImgRepoFullName gets image name from repo full name with format `repoName/imageName`
func GetNameFromImgRepoFullName(repo string) string {
	idx := strings.Index(repo, "/")
	return repo[idx+1:]
}

// BuildImageResourceURL ...
func BuildImageResourceURL(repoName, reference string) (string, error) {
	extURL, err := config.ExtURL()
	if err != nil {
		return "", fmt.Errorf("get external endpoint failed: %v", err)
	}

	if distribution.IsDigest(reference) {
		return fmt.Sprintf("%s/%s@%s", extURL, repoName, reference), nil
	}

	return fmt.Sprintf("%s/%s:%s", extURL, repoName, reference), nil
}
