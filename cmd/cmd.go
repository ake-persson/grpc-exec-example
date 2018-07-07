package cmd

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Command struct {
	Cmd      string
	Args     []string
	Env      []string
	User     string
	Group    string
	Dir      string
	Redirect bool
	Timeout  *time.Duration

	proc   *os.Process
	status *Status
	stdout *pipe
	stderr *pipe
}

type CommandOption func(*Command)

type Status struct {
	Running    bool                `json:"running"`
	Terminated bool                `json:"terminated"`
	Started    *time.Time          `json:"started,omitempty"`
	Finished   *time.Time          `json:"finished,omitempty"`
	Duration   *time.Duration      `json:"duration,omitempty"`
	Error      *string             `json:"error,omitempty"`
	ExitCode   *syscall.WaitStatus `json:"exitCode,omitempty"`
	Signal     *syscall.Signal     `json:"signal,omitempty"`
}

type pipe struct {
	reader *os.File
	writer *os.File
}

func New(cmd string, args []string, options ...CommandOption) *Command {
	c := &Command{
		Cmd:      cmd,
		Args:     args,
		Redirect: false,

		status: &Status{},
		stdout: &pipe{},
		stderr: &pipe{},
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func WithEnv(env []string) CommandOption {
	return func(c *Command) {
		c.Env = env
	}
}

func WithRedirect(c *Command) {
	c.Redirect = true
}

func WithUser(user string) CommandOption {
	return func(c *Command) {
		c.User = user
	}
}

func WithGroup(group string) CommandOption {
	return func(c *Command) {
		c.Group = group
	}
}

func WithDir(dir string) CommandOption {
	return func(c *Command) {
		c.Dir = dir
	}
}

func WithTimeout(timeout time.Duration) CommandOption {
	return func(c *Command) {
		c.Timeout = &timeout
	}
}

func getCredentials(userName string, groupName string) (*syscall.Credential, error) {
	var uid, gid string

	if userName != "" {
		u, err := user.Lookup(userName)
		if err != nil {
			return nil, fmt.Errorf("get user: %v", err)
		}
		uid = u.Uid
		gid = u.Gid
	}

	if groupName != "" {
		g, err := user.LookupGroup(groupName)
		if err != nil {
			return nil, fmt.Errorf("get group: %v", err)
		}
		gid = g.Gid
	}

	if userName == "" && groupName == "" {
		return nil, nil
	}

	c := syscall.Credential{}

	var err error
	var u64 uint64
	if u64, err = strconv.ParseUint(uid, 10, 32); err != nil {
		return nil, fmt.Errorf("conv string uid: %v", err)
	}
	c.Uid = uint32(u64)

	if u64, err = strconv.ParseUint(gid, 10, 32); err != nil {
		return nil, fmt.Errorf("conv string gid: %v", err)
	}
	c.Gid = uint32(u64)

	return &c, nil
}

func (c *Command) String() string {
	var opts []string
	if c.Env != nil {
		opts = append(opts, "env: "+strings.Join(c.Env, ","))
	}
	if c.User != "" {
		opts = append(opts, "user: "+c.User)
	}
	if c.Group != "" {
		opts = append(opts, "group: "+c.Group)
	}
	if c.Dir != "" {
		opts = append(opts, "dir: "+c.Dir)
	}
	if c.Timeout != nil {
		opts = append(opts, "tmout: "+fmt.Sprint(c.Timeout))
	}
	if len(opts) > 1 {
		return fmt.Sprintf("%s command: %s %s", strings.Join(opts, " "), c.Cmd, strings.Join(c.Args, " "))
	}
	return fmt.Sprintf("command: %s %s", c.Cmd, strings.Join(c.Args, " "))
}

func (c *Command) Start() (*Status, error) {
	var err error
	c.stdout.reader, c.stdout.writer, err = os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("create stdout pipe: %v", err)
	}

	if c.Redirect {
		c.stderr.writer = c.stdout.writer
	} else {
		c.stderr.reader, c.stderr.writer, err = os.Pipe()
		if err != nil {
			return nil, fmt.Errorf("create stderr pipe: %v", err)
		}
	}

	cred, err := getCredentials(c.User, c.Group)
	if err != nil {
		return nil, err
	}

	// First entry in args is ignored
	args := append([]string{""}, c.Args...)

	c.proc, err = os.StartProcess(c.Cmd, args, &os.ProcAttr{
		Env:   c.Env,
		Dir:   c.Dir,
		Files: []*os.File{nil, c.stdout.writer, c.stderr.writer},
		Sys: &syscall.SysProcAttr{
			Credential: cred,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("start process: %v", err)
	}

	timer := time.NewTimer(*c.Timeout)
	go func() {
		<-timer.C
		c.Kill()
	}()

	started := time.Now()
	c.status.Started = &started
	c.status.Running = true

	c.stdout.writer.Close()
	if !c.Redirect {
		c.stderr.writer.Close()
	}

	return c.status, nil
}

func (c *Command) Wait() (*Status, error) {
	if c.proc == nil {
		return nil, fmt.Errorf("command not started")
	}

	state, err := c.proc.Wait()

	exitCode := state.Sys().(syscall.WaitStatus)
	c.status.ExitCode = &exitCode

	finished := time.Now()
	c.status.Finished = &finished
	duration := finished.Sub(*c.status.Started)
	c.status.Duration = &duration
	c.status.Running = false

	if err != nil {
		e := err.Error()
		c.status.Error = &e

		return c.status, fmt.Errorf("wait process: %v", err)
	}

	c.stdout.reader.Close()
	if !c.Redirect {
		c.stderr.reader.Close()
	}

	return c.status, nil
}

func (c *Command) kill(signal syscall.Signal) (*Status, error) {
	if c.proc == nil {
		return nil, fmt.Errorf("command not started")
	}

	err := syscall.Kill(c.proc.Pid, signal)
	if err != nil {
		return nil, err
	}

	c.status.Terminated = true
	c.status.Signal = &signal

	return c.status, nil
}

func (c *Command) Stop() (*Status, error) {
	return c.kill(syscall.SIGTERM)
}

func (c *Command) Kill() (*Status, error) {
	return c.kill(syscall.SIGKILL)
}

func (c *Command) Status() *Status {
	return c.status
}

func (c *Command) Running() bool {
	return c.status.Running
}

func (c *Command) Terminated() bool {
	return c.status.Terminated
}

func (c *Command) Stdout() *os.File {
	return c.stdout.reader
}

func (c *Command) Stderr() *os.File {
	return c.stderr.reader
}
