package utils

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/unbeman/ya-prac-go-second-grade/internal/server/config"
)

var cert = []byte(`-----BEGIN CERTIFICATE-----
MIIB0zCCAX2gAwIBAgIJAI/M7BYjwB+uMA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTIwOTEyMjE1MjAyWhcNMTUwOTEyMjE1MjAyWjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBANLJ
hPHhITqQbPklG3ibCVxwGMRfp/v4XqhfdQHdcVfHap6NQ5Wok/4xIA+ui35/MmNa
rtNuC+BdZ1tMuVCPFZcCAwEAAaNQME4wHQYDVR0OBBYEFJvKs8RfJaXTH08W+SGv
zQyKn0H8MB8GA1UdIwQYMBaAFJvKs8RfJaXTH08W+SGvzQyKn0H8MAwGA1UdEwQF
MAMBAf8wDQYJKoZIhvcNAQEFBQADQQBJlffJHybjDGxRMqaRmDhX0+6v02TUKZsW
r5QuVbpQhH6u+0UgcW0jp9QwpxoPTLTWGXEWBBBurxFwiCBhkQ+V
-----END CERTIFICATE-----
`)

var key = []byte(`-----BEGIN PRIVATE KEY-----
MIIBOwIBAAJBANLJhPHhITqQbPklG3ibCVxwGMRfp/v4XqhfdQHdcVfHap6NQ5Wo
k/4xIA+ui35/MmNartNuC+BdZ1tMuVCPFZcCAwEAAQJAEJ2N+zsR0Xn8/Q6twa4G
6OB1M1WO+k+ztnX/1SvNeWu8D6GImtupLTYgjZcHufykj09jiHmjHx8u8ZZB/o1N
MQIhAPW+eyZo7ay3lMz1V01WVjNKK9QSn1MJlb06h/LuYv9FAiEA25WPedKgVyCW
SmUwbPw8fnTcpqDWE3yTO3vKcebqMSsCIBF3UmVue8YU3jybC3NxuXq3wNm34R8T
xVLHwDXh/6NJAiEAl2oHGGLz64BuAfjKrqwz7qMYr9HCLIe/YsoWq/olzScCIQDi
D2lWusoe2/nEqfDVVWGWlyJ7yOmqaVm/iNUN9B2N2g==
-----END PRIVATE KEY-----
`)

func TestLoadTLSCredentials(t *testing.T) {
	dir := t.TempDir()

	certPath := dir + "/cert.crt"

	certFile, err := os.Create(certPath)
	require.NoError(t, err)
	_, err = certFile.Write(cert)
	require.NoError(t, err)
	err = certFile.Close()
	require.NoError(t, err)

	keyPath := dir + "/key.key"

	keyFile, err := os.Create(keyPath)
	require.NoError(t, err)
	_, err = keyFile.Write(key)
	require.NoError(t, err)
	err = keyFile.Close()
	require.NoError(t, err)

	tests := []struct {
		name    string
		cfg     config.TLS
		wantErr bool
	}{
		{
			name:    "good",
			cfg:     config.BuildTestTLSConfig(certPath, keyPath),
			wantErr: false,
		},
		{
			name:    "no files",
			cfg:     config.BuildTestTLSConfig(dir+"no.crt", dir+"no.key"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadTLSCredentials(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadTLSCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
