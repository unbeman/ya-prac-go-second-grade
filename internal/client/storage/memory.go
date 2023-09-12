package storage

import (
	"strings"
	"sync"

	uuid "github.com/satori/go.uuid"

	"github.com/unbeman/ya-prac-go-second-grade/internal/client/model"
)

type MemStore struct { //todo: sync.map
	sync.RWMutex
	credentials map[uuid.UUID]*model.Credential
}

func NewMemStore() (*MemStore, error) {
	return &MemStore{credentials: map[uuid.UUID]*model.Credential{}}, nil
}

func (s *MemStore) FillFromLocal(creds []*model.Credential) {
	s.Lock()
	defer s.Unlock()
	for _, cred := range creds {
		s.credentials[cred.ID] = cred
	}
}

func (s *MemStore) GetAll() map[uuid.UUID]model.Credential {
	s.RLock()
	defer s.RUnlock()

	copied := map[uuid.UUID]model.Credential{}

	for id, cred := range s.credentials {
		copied[id] = *cred
	}
	return copied
}

func (s *MemStore) SearchCredentials(query string) (map[uuid.UUID]model.Credential, error) {
	s.RLock()
	defer s.RUnlock()

	lowQuery := strings.ToLower(query)

	result := map[uuid.UUID]model.Credential{}

	for id, cred := range s.credentials {
		if strings.Contains(strings.ToLower(cred.MetaData), lowQuery) {
			result[id] = *cred
		}
	}

	return result, nil
}

func (s *MemStore) GetCredentialByID(credID uuid.UUID) (model.Credential, error) {
	s.RLock()
	defer s.RUnlock()
	var cred model.Credential

	if _, ok := s.credentials[credID]; !ok {
		return cred, ErrNotFound
	}
	cred = *s.credentials[credID]
	return cred, nil
}

func (s *MemStore) UpsertCredential(cred model.Credential) {
	s.Lock()
	defer s.Unlock()
	s.credentials[cred.ID] = &cred
}

func (s *MemStore) DeleteCredential(cred model.Credential) {
	s.Lock()
	defer s.Unlock()
	delete(s.credentials, cred.ID)
}
