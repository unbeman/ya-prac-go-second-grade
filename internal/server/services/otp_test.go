package services

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pquerna/otp/totp"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	pb "github.com/unbeman/ya-prac-go-second-grade/api/v1"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/config"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/database"
	mock_database "github.com/unbeman/ya-prac-go-second-grade/internal/server/database/mock"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/model"
	mock_utils "github.com/unbeman/ya-prac-go-second-grade/internal/server/utils/mock"
	"github.com/unbeman/ya-prac-go-second-grade/internal/test_helpers"
)

func setUserToContext(ctx context.Context, user model.User) context.Context {
	return context.WithValue(ctx, "auth-user", user)
}

func TestOTP_OTPGenerate(t *testing.T) {
	user := model.User{Base: model.Base{ID: uuid.NewV4()}, Login: "test"}

	tests := []struct {
		name       string
		user       model.User
		wantErr    bool
		buildStubs func(db *mock_database.MockDatabase)
	}{
		{
			name:    "OK",
			user:    model.User{Base: model.Base{ID: uuid.NewV4()}, Login: "test"},
			wantErr: false,
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(user, nil).
					Times(1)
			},
		},
		{
			name:    "no login",
			user:    model.User{Base: model.Base{ID: uuid.NewV4()}, Login: ""},
			wantErr: true,
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(user, nil).
					Times(0)
			},
		},
		{
			name:    "db err",
			user:    model.User{Base: model.Base{ID: uuid.NewV4()}, Login: "test"},
			wantErr: true,
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(user, database.ErrDatabase).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			jwt := mock_utils.NewMockIJWT(ctrl)
			db := test_helpers.SetupMockDB(t, ctrl, tt.buildStubs)

			s := NewOTPService(config.BuildTestOTPConfig(), db, jwt)

			ctx = setUserToContext(ctx, tt.user)

			got, err := s.OTPGenerate(ctx, &pb.OTPGenRequest{})
			if (err != nil) != tt.wantErr {
				t.Errorf("OTPGenerate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotEqual(t, "", got.SecretKey)
				assert.NotEqual(t, "", got.AuthUrl)
			}
		})
	}
}

func TestOTP_OTPVerify(t *testing.T) {
	userID := uuid.NewV4()
	user := model.User{Login: "test"}
	user.ID = userID
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "test",
		AccountName: user.Login,
	})
	require.NoError(t, err)
	user.OtpAuthUrl = key.URL()
	user.OtpSecret = key.Secret()
	otpEnable := true
	user.OtpEnabled = &otpEnable
	user.OtpVerified = &otpEnable

	token, err := totp.GenerateCode(key.Secret(), time.Now())
	require.NoError(t, err)

	tests := []struct {
		name       string
		input      *pb.OTPVerifyRequest
		out        *pb.OTPVerifyResponse
		user       model.User
		wantErr    bool
		buildStubs func(db *mock_database.MockDatabase)
		jwtVer     func(j *mock_utils.MockIJWT)
	}{
		{
			name:    "OK",
			input:   &pb.OTPVerifyRequest{Token: token},
			out:     &pb.OTPVerifyResponse{AccessToken: "access-token"},
			user:    user,
			wantErr: false,
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(user, nil).
					Times(1)
			},
			jwtVer: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("access-token", nil).
					Times(1)
			},
		},
		{
			name:    "empty user login",
			input:   &pb.OTPVerifyRequest{Token: token},
			out:     &pb.OTPVerifyResponse{},
			user:    model.User{},
			wantErr: true,
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(model.User{}, nil).
					Times(0)
			},
			jwtVer: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("access-token", nil).
					Times(0)
			},
		},
		{
			name:    "invalid otp code",
			input:   &pb.OTPVerifyRequest{Token: "invalid token"},
			out:     &pb.OTPVerifyResponse{},
			user:    user,
			wantErr: true,
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(user, nil).
					Times(0)
			},
			jwtVer: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("access-token", nil).
					Times(0)
			},
		},
		{
			name:    "update db error",
			input:   &pb.OTPVerifyRequest{Token: token},
			out:     &pb.OTPVerifyResponse{AccessToken: "access-token"},
			user:    user,
			wantErr: true,
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(user, database.ErrDatabase).
					Times(1)
			},
			jwtVer: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("access-token", nil).
					Times(0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			jwt := test_helpers.SetupMockJWTManager(t, ctrl, tt.jwtVer)
			db := test_helpers.SetupMockDB(t, ctrl, tt.buildStubs)

			s := NewOTPService(config.BuildTestOTPConfig(), db, jwt)

			ctx = setUserToContext(ctx, tt.user)

			got, err := s.OTPVerify(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("OTPVerify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Equal(t, tt.out.AccessToken, got.AccessToken)
			}

		})
	}
}

func TestOTP_OTPValidate(t *testing.T) {
	userID := uuid.NewV4()
	user := model.User{Login: "test"}
	user.ID = userID
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "test",
		AccountName: user.Login,
	})
	require.NoError(t, err)
	user.OtpAuthUrl = key.URL()
	user.OtpSecret = key.Secret()

	token, err := totp.GenerateCode(key.Secret(), time.Now())
	require.NoError(t, err)

	tests := []struct {
		name    string
		input   *pb.OTPValidateRequest
		out     *pb.OTPValidateResponse
		user    model.User
		wantErr bool
		jwtVer  func(j *mock_utils.MockIJWT)
	}{
		{
			name:    "OK",
			input:   &pb.OTPValidateRequest{Token: token},
			out:     &pb.OTPValidateResponse{AccessToken: "access-token"},
			user:    user,
			wantErr: false,
			jwtVer: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("access-token", nil).
					Times(1)
			},
		},
		{
			name:    "empty user login",
			input:   &pb.OTPValidateRequest{Token: token},
			out:     &pb.OTPValidateResponse{},
			user:    model.User{},
			wantErr: true,
			jwtVer: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("access-token", nil).
					Times(0)
			},
		},
		{
			name:    "invalid otp code",
			input:   &pb.OTPValidateRequest{Token: "invalid token"},
			out:     &pb.OTPValidateResponse{},
			user:    user,
			wantErr: true,
			jwtVer: func(j *mock_utils.MockIJWT) {
				j.EXPECT().
					Generate(gomock.Any(), gomock.Any()).
					Return("access-token", nil).
					Times(0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			jwt := test_helpers.SetupMockJWTManager(t, ctrl, tt.jwtVer)
			db := mock_database.NewMockDatabase(ctrl)

			s := NewOTPService(config.BuildTestOTPConfig(), db, jwt)

			ctx = setUserToContext(ctx, tt.user)

			got, err := s.OTPValidate(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("OTPValidate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Equal(t, tt.out.AccessToken, got.AccessToken)
			}

		})
	}
}

func TestOTP_OTPDisable(t *testing.T) {
	userID := uuid.NewV4()
	user := model.User{Login: "test"}
	user.ID = userID
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "test",
		AccountName: user.Login,
	})
	require.NoError(t, err)
	user.OtpAuthUrl = key.URL()
	user.OtpSecret = key.Secret()
	otpEnable := true
	user.OtpEnabled = &otpEnable
	user.OtpVerified = &otpEnable

	require.NoError(t, err)

	tests := []struct {
		name       string
		user       model.User
		wantErr    bool
		buildStubs func(db *mock_database.MockDatabase)
	}{
		{
			name:    "OK",
			user:    user,
			wantErr: false,
			buildStubs: func(db *mock_database.MockDatabase) {
				disable := false
				updUser := user
				updUser.OtpEnabled = &disable

				db.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(updUser, nil).
					Times(1)
			},
		},
		{
			name:    "update db error",
			user:    user,
			wantErr: true,
			buildStubs: func(db *mock_database.MockDatabase) {
				db.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(user, database.ErrDatabase).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			jwt := mock_utils.NewMockIJWT(ctrl)
			db := test_helpers.SetupMockDB(t, ctrl, tt.buildStubs)

			s := NewOTPService(config.BuildTestOTPConfig(), db, jwt)

			ctx = setUserToContext(ctx, tt.user)

			_, err = s.OTPDisable(ctx, &pb.OTPDisableRequest{})
			if (err != nil) != tt.wantErr {
				t.Errorf("OTPVerify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func Test_getUserFromContext(t *testing.T) {
	ctx := context.Background()
	user := model.User{}

	tests := []struct {
		name    string
		ctx     context.Context
		want    model.User
		wantErr bool
	}{
		{
			name:    "good",
			ctx:     setUserToContext(ctx, user),
			wantErr: false,
			want:    user,
		},
		{
			name: "not user",
			ctx: func() context.Context {
				notUser := model.Credential{}
				return context.WithValue(ctx, "auth-user", notUser)
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getUserFromContext(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUserFromContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equalf(t, tt.want, got, "getUserFromContext(%v)", tt.ctx)
			}
		})
	}
}
