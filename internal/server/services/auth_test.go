package services

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	pb "github.com/unbeman/ya-prac-go-second-grade/api/v1"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/database"
	mock_database "github.com/unbeman/ya-prac-go-second-grade/internal/server/database/mock"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/model"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/utils"
	mock_utils "github.com/unbeman/ya-prac-go-second-grade/internal/server/utils/mock"
	"github.com/unbeman/ya-prac-go-second-grade/internal/test_helpers"
)

func TestAuth_Register(t *testing.T) {
	user := model.User{Login: "test", MasterKey2Hash: "hashed key-hash"}
	user.ID = uuid.NewV4()

	tests := []struct {
		name       string
		buildStubs func(db *mock_database.MockDatabase)
		jwtGen     func(j *mock_utils.MockIJWT)
		input      *pb.RegisterRequest
		out        *pb.RegisterResponse
		wantErr    bool
	}{
		{
			name:    "OK",
			wantErr: false,
			input: &pb.RegisterRequest{
				Login:   "test",
				KeyHash: base64.StdEncoding.EncodeToString([]byte("key-hash")),
			},
			out: &pb.RegisterResponse{AccessToken: "valid-access-token"},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(user, nil).
					Times(1)
			},
			jwtGen: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("valid-access-token", nil).
					Times(1)
			},
		},
		{
			name:    "user already exists",
			wantErr: true,
			input: &pb.RegisterRequest{
				Login:   "test",
				KeyHash: base64.StdEncoding.EncodeToString([]byte("key-hash")),
			},
			out: &pb.RegisterResponse{AccessToken: ""},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(model.User{}, database.ErrUserAlreadyExists).
					Times(1)
			},
			jwtGen: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("", nil).
					Times(0)
			},
		},
		{
			name:    "invalid key hash encoding",
			wantErr: true,
			input: &pb.RegisterRequest{
				Login:   "test",
				KeyHash: "invalid key-hash",
			},
			out: &pb.RegisterResponse{AccessToken: ""},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(model.User{}, nil).
					Times(0)
			},
			jwtGen: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("", nil).
					Times(0)
			},
		},
		{
			name:    "jwt error",
			wantErr: true,
			input: &pb.RegisterRequest{
				Login:   "test",
				KeyHash: base64.StdEncoding.EncodeToString([]byte("key-hash")),
			},
			out: &pb.RegisterResponse{AccessToken: ""},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(user, nil).
					Times(1)
			},
			jwtGen: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("", errors.New("could not create signed token")).
					Times(1)
			},
		},
		{
			name:    "invalid input key hash",
			wantErr: true,
			input: &pb.RegisterRequest{
				Login:   "test",
				KeyHash: base64.StdEncoding.EncodeToString([]byte("$2a$10$5WhGIyumJh3atF4ED.LbSe/8xmauk0jDyF/co0rkYDD/69TaLhhykrLL0Y.GFw0Cfa===")),
			},
			out: &pb.RegisterResponse{AccessToken: ""},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(user, nil).
					Times(0)
			},
			jwtGen: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("", nil).
					Times(0)
			},
		},
		{
			name:    "db error",
			wantErr: true,
			input: &pb.RegisterRequest{
				Login:   "test",
				KeyHash: base64.StdEncoding.EncodeToString([]byte("key-hash")),
			},
			out: &pb.RegisterResponse{AccessToken: ""},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(model.User{}, database.ErrDatabase).
					Times(1)
			},
			jwtGen: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("", nil).
					Times(0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			jwt := test_helpers.SetupMockJWTManager(t, ctrl, tt.jwtGen)
			db := test_helpers.SetupMockDB(t, ctrl, tt.buildStubs)

			auth := NewAuthService(db, jwt)

			got, err := auth.Register(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Equal(t, tt.out.AccessToken, got.AccessToken)
			}
		})
	}
}

func TestAuth_Login(t *testing.T) {
	key := "key-hash"
	storedKey, err := utils.HashToStore([]byte(key))
	require.NoError(t, err)

	otpEnabled := true
	user := model.User{Login: "test", MasterKey2Hash: base64.StdEncoding.EncodeToString(storedKey)}
	user.ID = uuid.NewV4()
	user.OtpEnabled = &otpEnabled

	tests := []struct {
		name       string
		buildStubs func(db *mock_database.MockDatabase)
		jwtGen     func(j *mock_utils.MockIJWT)
		input      *pb.LoginRequest
		out        *pb.LoginResponse
		wantErr    bool
	}{
		{
			name:    "OK",
			wantErr: false,
			input: &pb.LoginRequest{
				Login:   "test",
				KeyHash: base64.StdEncoding.EncodeToString([]byte(key)),
			},
			out: &pb.LoginResponse{AccessToken: "valid-access-token", Enforce_2FA: true},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					GetUserByLogin(gomock.Any(), gomock.Any()).
					Return(user, nil).
					Times(1)
			},
			jwtGen: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("valid-access-token", nil).
					Times(1)
			},
		},
		{
			name:    "user not found",
			wantErr: true,
			input: &pb.LoginRequest{
				Login:   "test",
				KeyHash: base64.StdEncoding.EncodeToString([]byte(key)),
			},
			out: &pb.LoginResponse{AccessToken: "", Enforce_2FA: false},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					GetUserByLogin(gomock.Any(), gomock.Any()).
					Return(model.User{}, database.ErrUserNotFound).
					Times(1)
			},
			jwtGen: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("", nil).
					Times(0)
			},
		},
		{
			name:    "invalid key hash",
			wantErr: true,
			input: &pb.LoginRequest{
				Login:   "test",
				KeyHash: base64.StdEncoding.EncodeToString([]byte("invalid-key-hash")),
			},
			out: &pb.LoginResponse{AccessToken: ""},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					GetUserByLogin(gomock.Any(), gomock.Any()).
					Return(model.User{}, nil).
					Times(1)
			},
			jwtGen: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("", nil).
					Times(0)
			},
		},
		{
			name:    "invalid key hash encoding",
			wantErr: true,
			input: &pb.LoginRequest{
				Login:   "test",
				KeyHash: "invalid key-hash",
			},
			out: &pb.LoginResponse{AccessToken: ""},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					GetUserByLogin(gomock.Any(), gomock.Any()).
					Return(user, nil).
					Times(1)
			},
			jwtGen: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("", nil).
					Times(0)
			},
		},
		{
			name:    "jwt error",
			wantErr: true,
			input: &pb.LoginRequest{
				Login:   "test",
				KeyHash: base64.StdEncoding.EncodeToString([]byte("key-hash")),
			},
			out: &pb.LoginResponse{},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					GetUserByLogin(gomock.Any(), gomock.Any()).
					Return(user, nil).
					Times(1)
			},
			jwtGen: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("", errors.New("could not create signed token")).
					Times(1)
			},
		},
		{
			name:    "invalid input key hash",
			wantErr: true,
			input: &pb.LoginRequest{
				Login:   "test",
				KeyHash: base64.StdEncoding.EncodeToString([]byte("$2a$10$5WhGIyumJh3atF4ED.LbSe/8xmauk0jDyF/co0rkYDD/69TaLhhykrLL0Y.GFw0Cfa===")),
			},
			out: &pb.LoginResponse{AccessToken: ""},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					GetUserByLogin(gomock.Any(), gomock.Any()).
					Return(user, nil).
					Times(1)
			},
			jwtGen: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("", nil).
					Times(0)
			},
		},
		{
			name:    "db error",
			wantErr: true,
			input: &pb.LoginRequest{
				Login:   "test",
				KeyHash: base64.StdEncoding.EncodeToString([]byte("key-hash")),
			},
			out: &pb.LoginResponse{},
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					GetUserByLogin(gomock.Any(), gomock.Any()).
					Return(model.User{}, database.ErrDatabase).
					Times(1)
			},
			jwtGen: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("", nil).
					Times(0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			jwt := test_helpers.SetupMockJWTManager(t, ctrl, tt.jwtGen)
			db := test_helpers.SetupMockDB(t, ctrl, tt.buildStubs)

			auth := NewAuthService(db, jwt)

			got, err := auth.Login(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Equal(t, tt.out.AccessToken, got.AccessToken)
				assert.Equal(t, tt.out.Enforce_2FA, got.Enforce_2FA)
			}
		})
	}
}

func TestAuth_GetUserByID(t *testing.T) {
	userId := uuid.NewV4()
	tests := []struct {
		name       string
		wantErr    bool
		userID     uuid.UUID
		buildStubs func(db *mock_database.MockDatabase)
	}{
		{
			name:    "good",
			wantErr: false,
			userID:  userId,
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					GetUserByID(gomock.Any(), gomock.Any()).
					Return(model.User{Base: model.Base{ID: userId}}, nil).
					Times(1)
			},
		},
		{
			name:    "not found",
			wantErr: true,
			userID:  userId,
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					GetUserByID(gomock.Any(), gomock.Any()).
					Return(model.User{}, database.ErrUserNotFound).
					Times(1)
			},
		},
		{
			name:    "db err",
			wantErr: true,
			userID:  userId,
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().
					GetUserByID(gomock.Any(), gomock.Any()).
					Return(model.User{}, database.ErrDatabase).
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
			jwt := mock_utils.NewMockIJWT(ctrl)

			auth := NewAuthService(db, jwt)

			_, err := auth.GetUserByID(ctx, tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
