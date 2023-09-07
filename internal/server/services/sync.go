package services

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/unbeman/ya-prac-go-second-grade/api/v1"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/database"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/model"
)

// Sync service implements proto pb.SyncServiceServer interface.
type Sync struct {
	pb.UnimplementedSyncServiceServer
	db database.Database
}

// NewSyncService setups new Sync instance.
func NewSyncService(db database.Database) *Sync {
	return &Sync{db: db}
}

func (s *Sync) Save(ctx context.Context, input *pb.SaveRequest) (*pb.SaveResponse, error) {
	user, err := getUserFromContext(ctx)
	if err != nil {
		return nil, GenStatusError(err)
	}

	log.Infof("got (%d) creds from user (%s) for save", len(input.Credentials), user.Login)

	for _, ic := range input.Credentials {
		uid, err := uuid.FromString(ic.GetLocalId())
		if err != nil {
			return nil, GenStatusError(fmt.Errorf("%w local_id format, got (%s)", ErrInvalid, ic.GetLocalId()))
		}
		cred := model.Credential{
			UserID:   user.ID,
			Type:     model.CredentialType(ic.GetType()), //todo: check
			MetaData: ic.GetMetaData(),
			Secret:   ic.GetSecret(),
		}

		cred.ID = uid
		cred.CreatedAt = ic.GetCreatedAt().AsTime()
		cred.UpdatedAt = ic.GetUpdatedAt().AsTime()
		if ic.GetDeletedAt() != nil {
			t := ic.GetDeletedAt().AsTime()
			cred.DeletedAt = &t
		}
		user.Credentials = append(user.Credentials, cred)
	}

	err = s.db.SaveUserSecrets(ctx, user)
	if err != nil {
		return nil, GenStatusError(err)
	}

	out := pb.SaveResponse{}
	return &out, nil
}
func (s *Sync) Load(ctx context.Context, input *pb.LoadRequest) (*pb.LoadResponse, error) {
	user, err := getUserFromContext(ctx)
	if err != nil {
		return nil, GenStatusError(err)
	}

	creds, err := s.db.GetUserSecrets(ctx, user)
	if err != nil {
		return nil, GenStatusError(err)
	}

	out := pb.LoadResponse{}
	for _, cred := range creds {
		oc := pb.Credential{
			LocalId:   cred.LocalID.String(),
			Type:      string(cred.Type),
			MetaData:  cred.MetaData,
			Secret:    cred.Secret,
			CreatedAt: timestamppb.New(cred.CreatedAt),
			UpdatedAt: timestamppb.New(cred.UpdatedAt),
		}

		if cred.DeletedAt != nil {
			oc.DeletedAt = timestamppb.New(*cred.DeletedAt)
		}
		out.Credentials = append(out.Credentials, &oc)
	}

	return &out, nil
}
