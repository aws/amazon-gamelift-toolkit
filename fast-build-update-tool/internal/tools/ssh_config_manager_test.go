package tools

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/config"
	"github.com/stretchr/testify/assert"
)

// TestDeterminePortLinuxDefault ensures that we don't allow custom ports set for Linux instances
func TestDeterminePortLinuxDefault(t *testing.T) {
	configMgr := NewSSHConfigManager(NewTestLogger(), "fake-path", 1045)

	port, err := configMgr.DeterminePort(config.OperatingSystemLinux)

	assert.Nil(t, err)
	assert.Equal(t, config.DefaultPortLinux, port)
}

// TestDeterminePortWindowsInvalid verifies that we don't allow the user to set an invalid port for a Windows instance
func TestDeterminePortWindowsInvalid(t *testing.T) {
	configMgr := NewSSHConfigManager(NewTestLogger(), "fake-path", 22)

	_, err := configMgr.DeterminePort(config.OperatingSystemWindows)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "ssh port must be greater than or equal to")
}

// TestDeterminePortWindowsDefault ensures that we return the expected default port for Windows instances
func TestDeterminePortWindowsDefault(t *testing.T) {
	configMgr := NewSSHConfigManager(NewTestLogger(), "fake-path", 00)

	port, err := configMgr.DeterminePort(config.OperatingSystemWindows)

	assert.Nil(t, err)
	assert.Equal(t, config.DefaultPortWindows, port)
}

// TestDeterminePortWindowCustom ensures that we allow a user to set a custom port for Windows instances when valid
func TestDeterminePortWindowCustom(t *testing.T) {
	configMgr := NewSSHConfigManager(NewTestLogger(), "fake-path", 1500)

	port, err := configMgr.DeterminePort(config.OperatingSystemWindows)

	assert.Nil(t, err)
	assert.Equal(t, int32(1500), port)
}

// TestLoadKey ensures that we can properly load an SSH key from the filesystem
func TestLoadKey(t *testing.T) {
	// Generate a temporary private key file
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(t, err)

	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	})

	keyFile, err := os.CreateTemp("", "key")
	assert.Nil(t, err)
	defer keyFile.Close()
	defer os.Remove(keyFile.Name())
	_, err = keyFile.Write(pemBytes)
	assert.Nil(t, err)

	configMgr := NewSSHConfigManager(NewTestLogger(), keyFile.Name(), 1500)

	signer, err := configMgr.LoadKey(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, "ssh-rsa", signer.PublicKey().Type())
}
