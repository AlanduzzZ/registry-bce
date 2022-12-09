package notification

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/goharbor/harbor/src/common/job/models"
	"github.com/goharbor/harbor/src/pkg/notification"
	policy_model "github.com/goharbor/harbor/src/pkg/notification/policy/model"
	"github.com/goharbor/harbor/src/pkg/notifier/event"
	"github.com/goharbor/harbor/src/pkg/notifier/model"
)

type fakedHookManager struct {
}

func (f *fakedHookManager) StartHook(ctx context.Context, event *model.HookEvent, job *models.JobData) error {
	return nil
}

func TestHTTPHandler_Handle(t *testing.T) {
	hookMgr := notification.HookManager
	defer func() {
		notification.HookManager = hookMgr
	}()
	notification.HookManager = &fakedHookManager{}

	handler := &HTTPHandler{}

	type args struct {
		event *event.Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "HTTPHandler_Handle Want Error 1",
			args: args{
				event: &event.Event{
					Topic: "http",
					Data:  nil,
				},
			},
			wantErr: true,
		},
		{
			name: "HTTPHandler_Handle Want Error 2",
			args: args{
				event: &event.Event{
					Topic: "http",
					Data:  &model.EventData{},
				},
			},
			wantErr: true,
		},
		{
			name: "HTTPHandler_Handle 1",
			args: args{
				event: &event.Event{
					Topic: "http",
					Data: &model.HookEvent{
						PolicyID:  1,
						EventType: "pushImage",
						Target: &policy_model.EventTarget{
							Type:    "http",
							Address: "http://127.0.0.1:8080",
						},
						Payload: &model.Payload{
							OccurAt: time.Now().Unix(),
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.Handle(context.TODO(), tt.args.event.Data)
			if tt.wantErr {
				require.NotNil(t, err, "Error: %s", err)
				return
			}
		})
	}
}

func TestHTTPHandler_IsStateful(t *testing.T) {
	handler := &HTTPHandler{}
	assert.False(t, handler.IsStateful())
}

func TestHTTPHandler_Name(t *testing.T) {
	handler := &HTTPHandler{}
	assert.Equal(t, "HTTP", handler.Name())
}
