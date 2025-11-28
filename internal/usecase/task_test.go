package usecase_test

import (
	"context"
	"errors"
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
		{
			name:  "successfully create task",
			title: "test title",
			mockFunc: func(m *mocks.MockTaskRepository) {
				expectedTask := &domain.Task{
					Title: "test title",
				}
				m.EXPECT().CreateTask(gomock.Any(), expectedTask).
					Return(errors.New("internal error"))
			},
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

func TestTaskService_ListTasks(t *testing.T) {
	in := usecase.ListTasksIn{}
	tests := []struct {
		name         string
		input        usecase.ListTasksIn
		setupMock    func(mockRepo *mocks.MockTaskRepository)
		wantTasksLen int
		wantErr      bool
	}{
		{
			name: "successfully return tasks list",
			setupMock: func(mockRepo *mocks.MockTaskRepository) {

				expectedTasks := &domain.PaginatedTasks{
					Tasks: []domain.Task{
						{
							ID:    "1",
							Title: "First task",
						},
						{
							ID:    "2",
							Title: "Second task",
						},
					},
				}

				mockRepo.EXPECT().
					ListTasks(gomock.Any(), in.Limit, in.Offset).
					Return(expectedTasks, nil)
			},
			wantTasksLen: 2,
			wantErr:      false,
		},
		{
			name: "return empty list",
			setupMock: func(mockRepo *mocks.MockTaskRepository) {

				var emptyTasks domain.PaginatedTasks

				mockRepo.EXPECT().
					ListTasks(gomock.Any(), in.Limit, in.Offset).
					Return(&emptyTasks, nil)
			},
			wantTasksLen: 0,
			wantErr:      false,
		},
		{
			name: "return error from repository",
			setupMock: func(mockRepo *mocks.MockTaskRepository) {
				mockRepo.EXPECT().
					ListTasks(gomock.Any(), in.Limit, in.Offset).
					Return(nil, errors.New("database connection error"))
			},
			wantTasksLen: 0,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockTaskRepository(ctrl)
			tt.setupMock(mockRepo)

			svc := usecase.NewTaskService(mockRepo)
			tasks, err := svc.ListTasks(context.TODO(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, tasks.Tasks, tt.wantTasksLen)
		})
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(mockRepo *mocks.MockTaskRepository)
		wantErr   bool
	}{
		{
			name: "successfully delete task with valid UUID",
			id:   "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(mockRepo *mocks.MockTaskRepository) {
				mockRepo.EXPECT().
					DeleteTask(gomock.Any(), domain.ID("550e8400-e29b-41d4-a716-446655440000")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "return error for invalid UUID format",
			id:        "invalid-uuid",
			setupMock: func(mockRepo *mocks.MockTaskRepository) {},
			wantErr:   true,
		},
		{
			name: "return error when repository fails",
			id:   "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(mockRepo *mocks.MockTaskRepository) {
				mockRepo.EXPECT().
					DeleteTask(gomock.Any(), domain.ID("550e8400-e29b-41d4-a716-446655440000")).
					Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockTaskRepository(ctrl)
			tt.setupMock(mockRepo)

			svc := usecase.NewTaskService(mockRepo)
			err := svc.DeleteTask(context.TODO(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTaskService_CompleteTask(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(mockRepo *mocks.MockTaskRepository)
		wantErr   bool
	}{
		{
			name: "successfully complete task with valid UUID",
			id:   "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(mockRepo *mocks.MockTaskRepository) {
				mockRepo.EXPECT().
					CompleteTask(gomock.Any(), domain.ID("550e8400-e29b-41d4-a716-446655440000"), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "return error for invalid UUID format",
			id:        "invalid-uuid",
			setupMock: func(mockRepo *mocks.MockTaskRepository) {},
			wantErr:   true,
		},
		{
			name: "return error when repository fails",
			id:   "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(mockRepo *mocks.MockTaskRepository) {
				mockRepo.EXPECT().
					CompleteTask(gomock.Any(), domain.ID("550e8400-e29b-41d4-a716-446655440000"), gomock.Any()).
					Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockTaskRepository(ctrl)
			tt.setupMock(mockRepo)

			svc := usecase.NewTaskService(mockRepo)
			err := svc.CompleteTask(context.TODO(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})

	}
}

func TestTaskService_UnCompleteTask(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(mockRepo *mocks.MockTaskRepository)
		wantErr   bool
	}{
		{
			name: "successfully complete task with valid UUID",
			id:   "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(mockRepo *mocks.MockTaskRepository) {
				mockRepo.EXPECT().
					UnCompleteTask(gomock.Any(), domain.ID("550e8400-e29b-41d4-a716-446655440000")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "return error for invalid UUID format",
			id:        "invalid-uuid",
			setupMock: func(mockRepo *mocks.MockTaskRepository) {},
			wantErr:   true,
		},
		{
			name: "return error when repository fails",
			id:   "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(mockRepo *mocks.MockTaskRepository) {
				mockRepo.EXPECT().
					UnCompleteTask(gomock.Any(), domain.ID("550e8400-e29b-41d4-a716-446655440000")).
					Return(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name: "return error when repository fails",
			id:   "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(mockRepo *mocks.MockTaskRepository) {
				mockRepo.EXPECT().
					UnCompleteTask(gomock.Any(), domain.ID("550e8400-e29b-41d4-a716-446655440000")).
					Return(domain.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockTaskRepository(ctrl)
			tt.setupMock(mockRepo)

			svc := usecase.NewTaskService(mockRepo)
			err := svc.UnCompleteTask(context.TODO(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})

	}
}
