package database

import (
	"context"
	"github.com/glebarez/sqlite"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"testing"

	"github.com/unbeman/ya-prac-go-second-grade/internal/server/model"
)

func prepareDB(t *testing.T) *pg {
	conn, err := gorm.Open(sqlite.Open(
		"file::memory:?cache=shared"),
		&gorm.Config{TranslateError: true},
	)
	require.NoError(t, err)

	db := pg{conn: conn}
	err = db.conn.AutoMigrate(&model.User{}, &model.Credential{})
	require.NoError(t, err)
	return &db
}

func closeDB(t *testing.T, db *pg) {
	instance, err := db.conn.DB()
	require.NoError(t, err)
	err = instance.Close()
	require.NoError(t, err)

}

func Test_pg_CreateUser(t *testing.T) {
	tests := []struct {
		name       string
		user       model.User
		beforeExec func(db *pg)
		wantErr    bool
	}{
		{
			name:       "good",
			user:       model.User{Login: "test", MasterKey2Hash: "key hash"},
			beforeExec: func(db *pg) {},
			wantErr:    false,
		},
		{
			name: "login already exists",
			user: model.User{Login: "exist", MasterKey2Hash: "key hash"},
			beforeExec: func(db *pg) {
				db.conn.Save(&model.User{Login: "exist"})
			},
			wantErr: true,
		},
		{
			name: "db closed",
			user: model.User{Login: "test", MasterKey2Hash: "key hash"},
			beforeExec: func(db *pg) {
				instance, _ := db.conn.DB()
				instance.Close()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := prepareDB(t)
			defer closeDB(t, db)

			tt.beforeExec(db)

			ctx := context.Background()
			got, err := db.CreateUser(ctx, tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				require.NotEqual(t, uuid.Nil, got.ID)
				require.Equal(t, tt.user.Login, got.Login)
				require.Equal(t, tt.user.MasterKey2Hash, got.MasterKey2Hash)
			}
		})
	}
}

func Test_pg_GetUserByLogin(t *testing.T) {
	tests := []struct {
		name       string
		user       model.User
		beforeExec func(db *pg)
		wantErr    bool
	}{
		{
			name: "good",
			user: model.User{Login: "test"},
			beforeExec: func(db *pg) {
				db.conn.Save(&model.User{Login: "test"})
			},
			wantErr: false,
		},
		{
			name: "not found",
			user: model.User{Login: "not found"},
			beforeExec: func(db *pg) {
			},
			wantErr: true,
		},
		{
			name: "db closed",
			user: model.User{Login: "test"},
			beforeExec: func(db *pg) {
				instance, _ := db.conn.DB()
				instance.Close()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := prepareDB(t)
			defer closeDB(t, db)

			tt.beforeExec(db)

			ctx := context.Background()

			got, err := db.GetUserByLogin(ctx, tt.user.Login)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByLogin() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				require.NotEqual(t, uuid.Nil, got.ID)
				require.Equal(t, tt.user.Login, got.Login)
			}
		})
	}
}

func Test_pg_GetUserByID(t *testing.T) {
	userId := uuid.NewV4()

	tests := []struct {
		name       string
		user       model.User
		beforeExec func(db *pg)
		wantErr    bool
	}{
		{
			name: "good",
			user: model.User{Login: "test", Base: model.Base{ID: userId}},
			beforeExec: func(db *pg) {
				db.conn.Save(&model.User{Login: "test", Base: model.Base{ID: userId}})
			},
			wantErr: false,
		},
		{
			name: "not found",
			user: model.User{Login: "not found", Base: model.Base{ID: uuid.NewV4()}},
			beforeExec: func(db *pg) {
			},
			wantErr: true,
		},
		{
			name: "db closed",
			user: model.User{Login: "test"},
			beforeExec: func(db *pg) {
				instance, _ := db.conn.DB()
				instance.Close()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := prepareDB(t)
			defer closeDB(t, db)

			tt.beforeExec(db)

			ctx := context.Background()

			got, err := db.GetUserByID(ctx, tt.user.ID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				require.Equal(t, tt.user.ID, got.ID)
			}
		})
	}
}

func Test_pg_UpdateUser(t *testing.T) {
	userId := uuid.NewV4()

	tests := []struct {
		name       string
		user       model.User
		beforeExec func(db *pg)
		wantErr    bool
	}{
		{
			name: "good",
			user: model.User{Login: "new Name", Base: model.Base{ID: userId}},
			beforeExec: func(db *pg) {
				db.conn.Save(&model.User{Login: "old Name", Base: model.Base{ID: userId}})
			},
			wantErr: false,
		},
		{
			name: "db closed",
			user: model.User{Login: "test"},
			beforeExec: func(db *pg) {
				instance, _ := db.conn.DB()
				instance.Close()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := prepareDB(t)
			defer closeDB(t, db)

			tt.beforeExec(db)

			ctx := context.Background()

			got, err := db.UpdateUser(ctx, tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				require.Equal(t, tt.user.Login, got.Login)
				require.NotEqual(t, tt.user.UpdatedAt, got.UpdatedAt)
			}
		})
	}
}

func Test_pg_GetUserSecrets(t *testing.T) {
	userId := uuid.NewV4()

	user := model.User{
		Login: "test",
		Base:  model.Base{ID: userId},
		Credentials: []model.Credential{
			{
				Type:     model.Login,
				MetaData: "google.com:test@gmail.com",
				Secret:   []byte("encrypted secret"),
			},
			{
				Type:     model.Note,
				MetaData: "First",
				Secret:   []byte("encrypted Very secret note"),
			},
			{
				Type:     model.Bank,
				MetaData: "Sberbank",
				Secret:   []byte("card credentials"),
			},
		},
	}

	tests := []struct {
		name       string
		user       model.User
		beforeExec func(db *pg)
		wantErr    bool
	}{
		{
			name: "good",
			user: user,
			beforeExec: func(db *pg) {
				db.conn.Save(&user)
			},
			wantErr: false,
		},
		{
			name: "db closed",
			user: model.User{Login: "test"},
			beforeExec: func(db *pg) {
				instance, _ := db.conn.DB()
				instance.Close()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := prepareDB(t)
			defer closeDB(t, db)

			tt.beforeExec(db)

			ctx := context.Background()

			got, err := db.GetUserSecrets(ctx, tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserSecrets() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				require.Equal(t, len(got), len(tt.user.Credentials))
			}
		})
	}
}

func Test_pg_SaveUserSecrets(t *testing.T) {
	userId := uuid.NewV4()

	user := model.User{
		Login: "test",
		Base:  model.Base{ID: userId},
		Credentials: []model.Credential{
			{
				Type:     model.Login,
				MetaData: "google.com:test@gmail.com",
				Secret:   []byte("encrypted secret"),
			},
			{
				Type:     model.Note,
				MetaData: "First",
				Secret:   []byte("encrypted Very secret note"),
			},
			{
				Type:     model.Bank,
				MetaData: "Sberbank",
				Secret:   []byte("card credentials"),
			},
		},
	}

	tests := []struct {
		name       string
		user       model.User
		beforeExec func(db *pg)
		wantErr    bool
	}{
		{
			name: "good",
			user: user,
			beforeExec: func(db *pg) {
				db.conn.Save(&user)
			},
			wantErr: false,
		},
		{
			name: "db closed",
			user: model.User{Login: "test"},
			beforeExec: func(db *pg) {
				instance, _ := db.conn.DB()
				instance.Close()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := prepareDB(t)
			defer closeDB(t, db)

			tt.beforeExec(db)

			ctx := context.Background()

			err := db.SaveUserSecrets(ctx, tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveUserSecrets() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func Test_pg_DeleteUserSecrets(t *testing.T) {
	userId := uuid.NewV4()

	user := model.User{
		Login: "test",
		Base:  model.Base{ID: userId},
		Credentials: []model.Credential{
			{
				Type:     model.Login,
				MetaData: "google.com:test@gmail.com",
				Secret:   []byte("encrypted secret"),
			},
			{
				Type:     model.Note,
				MetaData: "First",
				Secret:   []byte("encrypted Very secret note"),
			},
			{
				Type:     model.Bank,
				MetaData: "Sberbank",
				Secret:   []byte("card credentials"),
			},
		},
	}

	tests := []struct {
		name       string
		user       model.User
		beforeExec func(db *pg)
		wantErr    bool
	}{
		{
			name: "good",
			user: user,
			beforeExec: func(db *pg) {
				db.conn.Save(&user)
			},
			wantErr: false,
		},
		{
			name: "db closed",
			user: model.User{Login: "test"},
			beforeExec: func(db *pg) {
				instance, _ := db.conn.DB()
				instance.Close()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := prepareDB(t)
			defer closeDB(t, db)

			tt.beforeExec(db)

			ctx := context.Background()

			err := db.DeleteUserSecrets(ctx, tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUserSecrets() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				creds, err := db.GetUserSecrets(ctx, tt.user)
				require.NoError(t, err)
				require.Equal(t, 0, len(creds))
			}

		})
	}
}
