package utils

import (
	"testing"
)

func TestHashToStore(t *testing.T) {
	tests := []struct {
		name    string
		key     []byte
		wantErr bool
	}{
		{
			name: "good",
			key:  []byte("scrypt hashed key"),
		},
		{
			name:    "key too long",
			key:     []byte("invalid too long key from request of unknown client that create who know who"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := HashToStore(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashToStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestValidateKey(t *testing.T) {

	tests := []struct {
		name       string
		storedHash []byte
		key        []byte
		wantErr    bool
	}{
		{
			name:       "good",
			key:        []byte("scrypt password"),
			storedHash: []byte("$2a$10$5WhGIyumJh3atF4ED.LbSe/8xmauk0jDyF/69TaLhLL0Y.GFw0Cfa==="),
			wantErr:    false,
		},
		{
			name:       "invalid key",
			key:        []byte("invalid password"),
			storedHash: []byte("$2a$10$5WhGIyumJh3atF4ED.LbSe/8xmauk0jDyF/69TaLhLL0Y.GFw0Cfa==="),
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateKey(tt.storedHash, tt.key); (err != nil) != tt.wantErr {
				t.Errorf("ValidateKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
