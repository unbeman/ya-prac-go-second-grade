package transport

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/metadata"
	"testing"

	"github.com/unbeman/ya-prac-go-second-grade/internal/server/database"
	mock_database "github.com/unbeman/ya-prac-go-second-grade/internal/server/database/mock"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/model"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/services"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/utils"
	mock_utils "github.com/unbeman/ya-prac-go-second-grade/internal/server/utils/mock"
	"github.com/unbeman/ya-prac-go-second-grade/internal/test_helpers"
)

func TestAuthInterceptor_authorize(t *testing.T) {
	userID := uuid.NewV4()

	tests := []struct {
		name        string
		wantErr     bool
		ctx         context.Context
		inputMethod string
		accessToken string
		dbStubs     func(db *mock_database.MockDatabase)
		jwtStubs    func(j *mock_utils.MockIJWT)
	}{
		{
			name:    "good",
			wantErr: false,
			ctx: metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
				"authorization": "valid access token",
			})),
			dbStubs: func(db *mock_database.MockDatabase) {
				user := model.User{}
				user.ID = userID
				db.EXPECT().
					GetUserByID(gomock.Any(), gomock.Any()).
					Return(user, nil)
			},
			jwtStubs: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Verify(gomock.Any()).
					Return(&utils.UserClaims{UserID: userID}, nil).
					Times(1)
			},
		},
		{
			name:    "no metadata",
			wantErr: true,
			ctx:     context.Background(),
			dbStubs: func(db *mock_database.MockDatabase) {
				user := model.User{}
				user.ID = userID
				db.EXPECT().
					GetUserByID(gomock.Any(), gomock.Any()).
					Return(user, nil).
					Times(0)
			},
			jwtStubs: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Verify(gomock.Any()).
					Return(&utils.UserClaims{UserID: userID}, nil).
					Times(0)
			},
		},
		{
			name:    "no authorization metadata",
			wantErr: true,
			ctx: metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
				"invalid header": "valid access token",
			})),
			dbStubs: func(db *mock_database.MockDatabase) {
				user := model.User{}
				user.ID = userID
				db.EXPECT().
					GetUserByID(gomock.Any(), gomock.Any()).
					Return(user, nil).
					Times(0)
			},
			jwtStubs: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Verify(gomock.Any()).
					Return(&utils.UserClaims{UserID: userID}, nil).
					Times(0)
			},
		},
		{
			name:    "invalid access token",
			wantErr: true,
			ctx: metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
				"authorization": "invalid access token",
			})),
			dbStubs: func(db *mock_database.MockDatabase) {
				user := model.User{}
				user.ID = userID
				db.EXPECT().
					GetUserByID(gomock.Any(), gomock.Any()).
					Return(user, nil).Times(0)
			},
			jwtStubs: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Verify(gomock.Any()).
					Return(nil, errors.New("verify error")).
					Times(1)
			},
		},
		{
			name:    "2fa not validated before request",
			wantErr: true,
			ctx: metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
				"authorization": "invalid access token",
			})),
			dbStubs: func(db *mock_database.MockDatabase) {
				user := model.User{}
				user.ID = userID
				db.EXPECT().
					GetUserByID(gomock.Any(), gomock.Any()).
					Return(user, nil).Times(0)
			},
			jwtStubs: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Verify(gomock.Any()).
					Return(&utils.UserClaims{UserID: userID, OtpEnforce: true}, nil).
					Times(1)
			},
		},
		{
			name:    "user not found",
			wantErr: true,
			ctx: metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
				"authorization": "valid access token",
			})),
			dbStubs: func(db *mock_database.MockDatabase) {
				user := model.User{}
				user.ID = userID
				db.EXPECT().
					GetUserByID(gomock.Any(), gomock.Any()).
					Return(user, database.ErrUserNotFound).
					Times(1)
			},
			jwtStubs: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Verify(gomock.Any()).
					Return(&utils.UserClaims{UserID: userID}, nil).
					Times(1)
			},
		},
		{
			name:    "db error",
			wantErr: true,
			ctx: metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
				"authorization": "valid access token",
			})),
			dbStubs: func(db *mock_database.MockDatabase) {
				user := model.User{}
				user.ID = userID
				db.EXPECT().
					GetUserByID(gomock.Any(), gomock.Any()).
					Return(user, database.ErrDatabase).
					Times(1)
			},
			jwtStubs: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Verify(gomock.Any()).
					Return(&utils.UserClaims{UserID: userID}, nil).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db := test_helpers.SetupMockDB(t, ctrl, tt.dbStubs)
			jwt := test_helpers.SetupMockJWTManager(t, ctrl, tt.jwtStubs)

			auth := services.NewAuthService(db, jwt)

			i := NewAuthInterceptor(auth, map[string]bool{}, map[string]bool{})

			_, err := i.authorize(tt.ctx, tt.inputMethod)
			if (err != nil) != tt.wantErr {
				t.Errorf("authorize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//todo: check got user
		})
	}
}
