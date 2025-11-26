package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/darkjian/simpletodolist/internal/adapters/database/mocks"
	"github.com/darkjian/simpletodolist/internal/domain"
	"github.com/darkjian/simpletodolist/internal/usecase"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"
)

func TestTaskService_GetTask(t *testing.T) {

	var (
		taskID = uuid.New().String()
		now    = time.Now().UTC()
		want   = &domain.Task{
			ID:        domain.ID(taskID),
			Title:     "Test Task",
			CreatedAt: now,
			UpdatedAt: now,
		}
	)

	tests := []struct {
		name     string
		id       string
		mockFunc func(m *mocks.MockTaskRepository)
		want     domain.Task
		wantErr  bool
	}{
		{
			name: "successfully finds task by id",
			id:   taskID,
			mockFunc: func(m *mocks.MockTaskRepository) {
				m.EXPECT().GetTask(gomock.Any(), domain.ID(taskID)).
					Return(want, nil)
			},
			want:    *want,
			wantErr: false,
		},
		{
			name: "doesn't find task by id",
			id:   taskID,
			mockFunc: func(m *mocks.MockTaskRepository) {
				m.EXPECT().GetTask(gomock.Any(), domain.ID(taskID)).
					Return(nil, domain.ErrNotFound)
			},
			want:    domain.Task{},
			wantErr: true,
		},
		{
			name:     "passed invalid id",
			id:       "not-a-uuid",
			mockFunc: func(m *mocks.MockTaskRepository) {},
			want:     domain.Task{},
			wantErr:  true,
		},
		{
			name: "correct data internal error",
			id:   taskID,
			mockFunc: func(m *mocks.MockTaskRepository) {
				m.EXPECT().GetTask(gomock.Any(), domain.ID(taskID)).
					Return(want, domain.ErrInternal)
			},
			want:    domain.Task{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mocks.NewMockTaskRepository(ctrl)
			tt.mockFunc(mockRepo)

			svc := usecase.NewTaskService(mockRepo)
			got, err := svc.GetTask(context.TODO(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got.Task)
		})
	}
}
