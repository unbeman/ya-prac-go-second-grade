package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/unbeman/ya-prac-go-second-grade/api/v1"
	"github.com/unbeman/ya-prac-go-second-grade/internal/client/model"
	"github.com/unbeman/ya-prac-go-second-grade/internal/client/storage"
	"github.com/unbeman/ya-prac-go-second-grade/internal/client/utils"
)

// todo: move to cfg
const SyncDurationDefault = time.Minute

type SyncService struct {
	client     pb.SyncServiceClient
	auth       *AuthService
	vault      storage.IStorage
	memStorage *storage.MemStore
	syncOnce   sync.Once
	syncDur    time.Duration
	stop       chan struct{}
	waitStop   chan struct{}
}

func NewSyncService(conn grpc.ClientConnInterface, vault storage.IStorage, mem *storage.MemStore, auth *AuthService) *SyncService {
	return &SyncService{
		client:     pb.NewSyncServiceClient(conn),
		vault:      vault,
		memStorage: mem,
		auth:       auth,
		syncDur:    SyncDurationDefault,
		syncOnce:   sync.Once{},
		stop:       make(chan struct{}),
		waitStop:   make(chan struct{}),
	}
}

//todo: запретить методы до авторизации

func (s *SyncService) StartSync() {
	go s.syncOnce.Do(func() {
		log.Info("starting sync")
		ticker := time.NewTicker(s.syncDur)
		s.LoadFromServer()
		for {
			select {
			case <-ticker.C:
				log.Info("sync now")
				go s.SaveOnServer()
				go s.LoadFromServer()
			case <-s.stop:
				s.SaveOnServer()
				log.Info("sync stopped")
				s.waitStop <- struct{}{}
				return
			}
		}
	})
}

func (s *SyncService) StopSync() {
	close(s.stop)
	<-s.waitStop
}

// SaveOnServer uploads all user's secrets from local storage to server vault.
func (s *SyncService) SaveOnServer() {
	ctx := context.TODO()

	locals, err := s.vault.GetAllCredentials(ctx)
	if err != nil {
		log.Error("SaveOnServer: ", err)
	}

	saveInput := pb.SaveRequest{}
	saveInput.Credentials = make([]*pb.Credential, 0, len(locals))
	for _, cred := range locals {
		gCred := pb.Credential{
			LocalId:   cred.ID.String(),
			Type:      string(cred.Type),
			MetaData:  cred.MetaData,
			Secret:    cred.Encrypted,
			CreatedAt: timestamppb.New(cred.CreatedAt),
			UpdatedAt: timestamppb.New(cred.UpdatedAt),
		}

		if cred.DeletedAt != nil {
			gCred.DeletedAt = timestamppb.New(*cred.DeletedAt)
		}
		saveInput.Credentials = append(saveInput.Credentials, &gCred)
	}

	_, err = s.client.Save(ctx, &saveInput)
	if err != nil {
		log.Error("SaveOnServer: ", err)
		return
	}
	log.Info("saved on server")
}

// LoadFromServer downloads all user's secrets from server vault to local storage.
func (s *SyncService) LoadFromServer() {
	ctx := context.TODO()
	loadInput := pb.LoadRequest{}

	out, err := s.client.Load(ctx, &loadInput)
	if err != nil {
		log.Error(err)
		return
	}
	for _, serverCred := range out.GetCredentials() {
		uuID, err := uuid.FromString(serverCred.GetLocalId())
		if err != nil {
			log.Error("LoadFromServer: ", err)
			return
		}

		cred := model.Credential{
			ID:        uuID,
			MetaData:  serverCred.GetMetaData(),
			Type:      model.CredentialType(serverCred.GetType()),
			Encrypted: serverCred.GetSecret(),
			CreatedAt: serverCred.GetCreatedAt().AsTime(),
			UpdatedAt: serverCred.GetUpdatedAt().AsTime(),
		}

		if serverCred.GetDeletedAt() != nil {
			t := serverCred.DeletedAt.AsTime()
			cred.DeletedAt = &t
		}

		_, err = s.vault.SaveCredential(ctx, cred)
		if err != nil {
			log.Error("LoadFromServer: ", err)
			return
		}

	}

	vCreds, err := s.vault.GetAllCredentials(ctx)
	if err != nil {
		log.Error("LoadFromServer: ", err)
		return
	}
	s.memStorage.FillFromLocal(vCreds)
	log.Info("load from server")
}

//todo: replace to controller layer

func (s *SyncService) CreateCred(cType model.CredentialType, metadata, sensitiveData string) error {
	key := s.auth.GetMasterKey()

	cred := &model.Credential{}
	cred.Type = cType
	cred.MetaData = metadata

	encryptedData, err := utils.Encrypt(key, []byte(sensitiveData))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInternal, err)
	}
	cred.Encrypted = encryptedData
	err = s.vault.AddCredential(context.TODO(), cred)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInternal, err)
	}

	s.memStorage.UpsertCredential(*cred)

	return nil
}

func (s *SyncService) EditCred(cred model.Credential, metadata, sensitiveData string) error {
	key := s.auth.GetMasterKey()

	cred.MetaData = metadata

	decrypted := []byte(sensitiveData)
	encryptedData, err := utils.Encrypt(key, decrypted)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInternal, err)
	}

	cred.Encrypted = encryptedData

	cred, err = s.vault.SaveCredential(context.TODO(), cred)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInternal, err)
	}

	cred.Decrypted = decrypted

	s.memStorage.UpsertCredential(cred)

	return nil
}

func (s *SyncService) DeleteCred(cred model.Credential) error {
	err := s.vault.DeleteCredential(context.TODO(), cred)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInternal, err)
	}

	s.memStorage.DeleteCredential(cred)
	return nil
}

func (s *SyncService) GetAll() (map[uuid.UUID]model.Credential, error) {
	return s.memStorage.GetAll(), nil
}

func (s *SyncService) Search(query string) (map[uuid.UUID]model.Credential, error) {
	return s.memStorage.SearchCredentials(query)
}

func (s *SyncService) GetCredByID(credID uuid.UUID) (model.Credential, error) {
	cred, err := s.memStorage.GetCredentialByID(credID)
	if err != nil {
		return cred, err
	}

	if cred.Decrypted == nil {
		data, err := utils.Decrypt(s.auth.GetMasterKey(), cred.Encrypted)
		if err != nil {
			return cred, fmt.Errorf("%w: %s", ErrInternal, err)
		}
		cred.Decrypted = data
	}

	s.memStorage.UpsertCredential(cred)

	return cred, nil
}

func (s *SyncService) GetTypes() map[model.CredentialType]bool {
	return map[model.CredentialType]bool{
		model.Login: true,
		model.Note:  true,
		model.Bank:  true,
	}
}
