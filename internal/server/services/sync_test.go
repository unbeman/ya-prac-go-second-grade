package services

import (
	"context"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"

	pb "github.com/unbeman/ya-prac-go-second-grade/api/v1"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/database"
	mock_database "github.com/unbeman/ya-prac-go-second-grade/internal/server/database/mock"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/model"
	"github.com/unbeman/ya-prac-go-second-grade/internal/test_helpers"
)

func TestSync_Save(t *testing.T) {
	userID := uuid.NewV4()

	tests := []struct {
		name       string
		wantErr    bool
		user       model.User
		input      *pb.SaveRequest
		out        *pb.SaveResponse
		buildStubs func(db *mock_database.MockDatabase)
	}{
		{
			name:    "OK",
			wantErr: false,
			user:    model.User{Base: model.Base{ID: userID}},
			input: &pb.SaveRequest{
				Credentials: []*pb.Credential{
					{
						LocalId:   uuid.NewV4().String(),
						Type:      string(model.Login),
						MetaData:  "site:email",
						Secret:    []byte("encrypted secret"),
						CreatedAt: timestamppb.New(time.Now().Add(-time.Hour)),
						UpdatedAt: timestamppb.New(time.Now().Add(-time.Minute)),
						DeletedAt: nil,
					},
					{
						LocalId:   uuid.NewV4().String(),
						Type:      string(model.Note),
						MetaData:  "MyNote",
						Secret:    []byte("encrypted note"),
						CreatedAt: timestamppb.New(time.Now().Add(-time.Hour * 2)),
						UpdatedAt: timestamppb.New(time.Now().Add(-time.Minute * 3)),
						DeletedAt: timestamppb.New(time.Now().Add(-time.Minute)),
					},
				}},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					SaveUserSecrets(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:    "invalid uuid",
			wantErr: true,
			user:    model.User{Base: model.Base{ID: userID}},
			input: &pb.SaveRequest{
				Credentials: []*pb.Credential{
					{
						LocalId:   "invalid uuid",
						Type:      string(model.Login),
						MetaData:  "site:email",
						Secret:    []byte("encrypted secret"),
						CreatedAt: timestamppb.New(time.Now().Add(-time.Hour)),
						UpdatedAt: timestamppb.New(time.Now().Add(-time.Minute)),
						DeletedAt: nil,
					},
				},
			},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					SaveUserSecrets(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(0)
			},
		},
		{
			name:    "db err",
			wantErr: true,
			user:    model.User{Base: model.Base{ID: userID}},
			input: &pb.SaveRequest{
				Credentials: []*pb.Credential{
					{
						LocalId:   uuid.NewV4().String(),
						Type:      string(model.Login),
						MetaData:  "site:email",
						Secret:    []byte("encrypted secret"),
						CreatedAt: timestamppb.New(time.Now().Add(-time.Hour)),
						UpdatedAt: timestamppb.New(time.Now().Add(-time.Minute)),
						DeletedAt: nil,
					},
				},
			},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					SaveUserSecrets(gomock.Any(), gomock.Any()).
					Return(database.ErrDatabase).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db := test_helpers.SetupMockDB(t, ctrl, tt.buildStubs)

			sync := NewSyncService(db)

			ctx = setUserToContext(ctx, tt.user)

			_, err := sync.Save(ctx, tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestSync_Load(t *testing.T) {
	userID := uuid.NewV4()
	deleteNow := time.Now()

	tests := []struct {
		name       string
		wantErr    bool
		user       model.User
		buildStubs func(db *mock_database.MockDatabase)
	}{
		{
			name:    "OK",
			wantErr: false,
			user:    model.User{Base: model.Base{ID: userID}},
			buildStubs: func(db *mock_database.MockDatabase) {
				creds := []model.Credential{
					{Base: model.Base{
						ID:        uuid.NewV4(),
						CreatedAt: time.Now().Add(-time.Hour),
						UpdatedAt: time.Now().Add(-time.Hour),
					},
						Type:     model.Login,
						MetaData: "site:email",
						Secret:   []byte("encrypted secret"),
					},
					{Base: model.Base{
						ID:        uuid.NewV4(),
						CreatedAt: time.Now().Add(-time.Hour * 2),
						UpdatedAt: time.Now().Add(-time.Minute),
						DeletedAt: &deleteNow,
					},
						Type:     model.Note,
						MetaData: "Title",
						Secret:   []byte("encrypted note"),
					},
				}
				db.EXPECT().
					GetUserSecrets(gomock.Any(), gomock.Any()).
					Return(creds, nil).
					Times(1)
			},
		},
		{
			name:    "db err",
			wantErr: true,
			user:    model.User{Base: model.Base{ID: userID}},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					GetUserSecrets(gomock.Any(), gomock.Any()).
					Return(nil, database.ErrDatabase).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db := test_helpers.SetupMockDB(t, ctrl, tt.buildStubs)

			sync := NewSyncService(db)

			ctx = setUserToContext(ctx, tt.user)

			_, err := sync.Load(ctx, &pb.LoadRequest{})

			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
