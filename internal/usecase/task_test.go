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
			defer ctrl.Finish()

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

func TestTaskService_CreateTask(t *testing.T) {

	tests := []struct {
		name     string
		title    string
		mockFunc func(m *mocks.MockTaskRepository)
		wantErr  bool
	}{
		{
			name:  "successfully create task",
			title: "test title",
			mockFunc: func(m *mocks.MockTaskRepository) {
				expectedTask := &domain.Task{
					Title: "test title",
				}
				m.EXPECT().CreateTask(gomock.Any(), expectedTask).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "unsuccessfully create task with length < 1",
			title:    "",
			mockFunc: func(m *mocks.MockTaskRepository) {},
			wantErr:  true,
		},
		{
			name:     "unsuccessfully create task with length > 20",
			title:    "fivefivefivefivefive1",
			mockFunc: func(m *mocks.MockTaskRepository) {},
			wantErr:  true,
		},
		{
			name:  "successfully create task with length = 20",
			title: "test title with max!",
			mockFunc: func(m *mocks.MockTaskRepository) {
				expectedTask := &domain.Task{
					Title: "test title with max!",
				}
				m.EXPECT().CreateTask(gomock.Any(), expectedTask).
					Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockTaskRepository(ctrl)
			tt.mockFunc(mockRepo)

			svc := usecase.NewTaskService(mockRepo)
			got, err := svc.CreateTask(context.TODO(), tt.title)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.title, got.Task.Title)
			}

		})
	}
}
