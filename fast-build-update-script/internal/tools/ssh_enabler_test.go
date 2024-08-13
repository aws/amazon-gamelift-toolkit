package tools

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/config"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/gamelift"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func testGenerateKey(t *testing.T) ssh.PublicKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(t, err)

	key, err := ssh.NewPublicKey(&privateKey.PublicKey)
	assert.Nil(t, err)

	return key
}

func testGenerateED25519Key(t *testing.T) ssh.PublicKey {
	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	assert.Nil(t, err)

	key, err := ssh.NewPublicKey(publicKey)
	assert.Nil(t, err)

	return key
}

func NewTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
}

func TestNewSSHEnablerWindows(t *testing.T) {
	instance := &gamelift.Instance{OperatingSystem: config.OperatingSystemWindows}

	enabler, err := NewSSHEnabler(NewTestLogger(), instance, &GameLiftInstanceAccessGetterMock{}, testGenerateKey(t), 22)
	assert.Nil(t, err)
	assert.True(t, len(enabler.commandsToRun) > 0)
}

func TestNewSSHEnablerLinux(t *testing.T) {
	instance := &gamelift.Instance{OperatingSystem: config.OperatingSystemLinux}

	enabler, err := NewSSHEnabler(NewTestLogger(), instance, &GameLiftInstanceAccessGetterMock{}, testGenerateKey(t), 22)
	assert.Nil(t, err)
	assert.True(t, len(enabler.commandsToRun) > 0)
}

func TestNewSSHEnablerUnknownOS(t *testing.T) {
	instance := &gamelift.Instance{}

	_, err := NewSSHEnabler(NewTestLogger(), instance, &GameLiftInstanceAccessGetterMock{}, testGenerateKey(t), 22)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "unknown operating system")
}

func TestIsNewCommandOutputLinux(t *testing.T) {
	assert.True(t, IsNewCommandOutputLinux("sh-5.2$ "))
	assert.True(t, IsNewCommandOutputLinux("abunchofrandominput\r\nsh-5.3$ "))
	assert.True(t, IsNewCommandOutputLinux("abunch\r\nofrand\r\nominput\nsh-4.6$ "))
	assert.True(t, IsNewCommandOutputLinux("abunchofrandominput\nsh-5.2$ "))
	assert.True(t, IsNewCommandOutputLinux(`abunchofrandominput
sh-5.2$ `))

	assert.False(t, IsNewCommandOutputLinux("abunchofrandominput"))
	assert.False(t, IsNewCommandOutputLinux("abunchofrandominput\r\n$"))
	assert.False(t, IsNewCommandOutputLinux("abunchofrandominput\n$"))
}

func TestIsNewCommandOutputWindows(t *testing.T) {
	assert.True(t, IsNewCommandOutputWindows("PS C:\\Windows\\system32>"))
	assert.True(t, IsNewCommandOutputWindows("abunchofrandominput\r\nPS C:\\Windows\\system32>"))
	assert.True(t, IsNewCommandOutputWindows("abunch\r\nofrand\r\nominput\nPS C:\\Windows\\system32>"))
	assert.True(t, IsNewCommandOutputWindows(`abunchofrandominput
PS C:\Windows\system32> `))

	assert.False(t, IsNewCommandOutputWindows("abunchofrandominput"))
	assert.False(t, IsNewCommandOutputWindows("abunchofrandominput\r\nPS"))
	assert.False(t, IsNewCommandOutputWindows("abunchofrandominput\nPS"))
}

func TestFindED25519PublicKey(t *testing.T) {
	expectedKeyStr := "ssh-ed25519 ABCDE0FghI1jKLM1NOP5RSTUVW04XYzA0BCDefghIJKlMNOPQRSTuVWXYzabCdEf11GhI"
	assert.Equal(t, expectedKeyStr, FindED25519PublicKey(expectedKeyStr))
	// Make sure we can filter out additional garbage
	assert.Equal(t, expectedKeyStr, FindED25519PublicKey(expectedKeyStr+`\x1b[?25l`))
}

// mockSSMReader mocks how a remote SSM session would read output from a remote terminal session
type mockSSMReader struct {
	publicKey string
}

func (t *mockSSMReader) Read(p []byte) (n int, err error) {
	commandSeparator := "sh-5.2$"
	// write the public key and the terminal output so we can exit the session properly.
	bytes := copy(p, []byte(fmt.Sprintf("%s\n%s", t.publicKey, commandSeparator)))
	return bytes, nil
}

func TestEnable(t *testing.T) {
	expectedAccessKey := "AccessKeyId"
	expectedSecretAccessKey := "SecretAccessKey"
	expectedSessionToken := "SessionToken"

	instanceAccessGetter := &GameLiftInstanceAccessGetterMock{}
	instanceAccessGetter.GetInstanceAccessFunc = func(ctx context.Context, fleetId string, instanceId string) (*gamelift.InstanceAccessCredentials, error) {
		return &gamelift.InstanceAccessCredentials{
			AccessKeyId:     expectedAccessKey,
			SecretAccessKey: expectedSecretAccessKey,
			SessionToken:    expectedSessionToken,
		}, nil
	}

	publicKey := testGenerateED25519Key(t)
	reader := &mockSSMReader{
		publicKey: convertPublicKeyToString(publicKey),
	}

	mockedSSMCommandRunner := &PTYMock{
		CleanupFunc: func() {},
		ReaderFunc: func() io.Reader {
			return reader
		},
		RunCommandFunc: func(cmd string) error {
			return nil
		},
		StartFunc: func(cmdName string, args []string, env []string) error {
			return nil
		},
		WaitFunc: func() error {
			return nil
		},
	}

	instanceId := "i-1234"
	fleetId := "f-1234"

	enabler := &SSHEnabler{
		logger:               NewTestLogger(),
		instance:             &gamelift.Instance{FleetId: fleetId, InstanceId: instanceId},
		instanceAccessGetter: instanceAccessGetter,
		clientPublicKey:      "",
		isNewCommandOutput:   IsNewCommandOutputLinux,
		pty:                  mockedSSMCommandRunner,
		commandsToRun:        []string{"ls -lah"},
	}

	_, err := enabler.Enable(context.Background())
	assert.Nil(t, err)

	assert.Len(t, instanceAccessGetter.GetInstanceAccessCalls(), 1)
	assert.Equal(t, fleetId, instanceAccessGetter.GetInstanceAccessCalls()[0].FleetId)
	assert.Equal(t, instanceId, instanceAccessGetter.GetInstanceAccessCalls()[0].InstanceId)

	assert.Len(t, mockedSSMCommandRunner.StartCalls(), 1)
	assert.Equal(t, mockedSSMCommandRunner.StartCalls()[0].CmdName, "aws")
	assert.Equal(t, mockedSSMCommandRunner.StartCalls()[0].Args, []string{"ssm", "start-session", "--target", instanceId})
	assert.Contains(t, mockedSSMCommandRunner.StartCalls()[0].Env, "AWS_ACCESS_KEY_ID="+expectedAccessKey)
	assert.Contains(t, mockedSSMCommandRunner.StartCalls()[0].Env, "AWS_SECRET_ACCESS_KEY="+expectedSecretAccessKey)
	assert.Contains(t, mockedSSMCommandRunner.StartCalls()[0].Env, "AWS_SESSION_TOKEN="+expectedSessionToken)
}
