package tools

import (
	"io"

	"github.com/aymanbagabas/go-pty"
)

//go:generate moq -skip-ensure -out ./moq_pty_test.go . PTY

// PTY is an interface used to start and interact with a pseudo terminal
type PTY interface {
	// Start a new PTY session with the provided command, command args, and environment
	Start(cmdName string, args []string, env []string) error

	// Run a command on the PTY session after it has been started
	RunCommand(cmd string) error

	// Wait for the session to finish
	Wait() error

	// Clean-up any resources
	Cleanup()

	// Get the reader for the PTY session
	Reader() io.Reader
}

// ptyRunner is used for running commands in a pseudo terminal environment
type ptyRunner struct {
	pty pty.Pty
	cmd *pty.Cmd
}

func newPtyCommandRunner() (*ptyRunner, error) {
	p, err := pty.New()
	if err != nil {
		return nil, err
	}
	return &ptyRunner{pty: p}, nil
}

func (p *ptyRunner) Start(cmdName string, args []string, env []string) error {
	p.cmd = p.pty.Command(cmdName, args...)
	p.cmd.Env = env
	err := p.cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (p *ptyRunner) RunCommand(cmd string) error {
	_, err := p.pty.Write([]byte(cmd))
	return err
}

func (p *ptyRunner) Wait() error {
	return p.cmd.Wait()
}

func (p *ptyRunner) Cleanup() {
	p.pty.Close()
}

func (p *ptyRunner) Reader() io.Reader {
	return p.pty
}
