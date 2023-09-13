package test_helpers

import (
	"github.com/golang/mock/gomock"
	"testing"

	"github.com/unbeman/ya-prac-go-second-grade/internal/server/database/mock"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/utils"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/utils/mock"
)

func SetupMockDB(t *testing.T, ctrl *gomock.Controller, buildStubs func(db *mock_database.MockDatabase)) *mock_database.MockDatabase {
	database := mock_database.NewMockDatabase(ctrl)
	buildStubs(database)
	return database
}

func SetupMockJWTManager(t *testing.T, ctrl *gomock.Controller, stub func(j *mock_utils.MockIJWT)) utils.IJWT {
	jwt := mock_utils.NewMockIJWT(ctrl)
	stub(jwt)
	return jwt
}
