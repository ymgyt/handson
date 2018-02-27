package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/juju/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	dockerStatusExited  = "exited"
	dockerStatusRunning = "running"
)

type LocalDockerDB struct {
	Image     string
	Container string
	MountDir  string
	IPAddr    string
	Port      string
	User      string
	Password  string
	Logger    *zap.Logger
}

type LocalDockerDBOption func(*LocalDockerDB)

func NewLocalDockerDB(options ...LocalDockerDBOption) *LocalDockerDB {
	d := &LocalDockerDB{
		Image:     "postgres:10.1",
		Container: "webdevgo-postgres",
		MountDir:  "/var/tmp/webdevgo_db",
		Port:      "5432",
		User:      "me",
		Password:  "",
		Logger:    zap.NewNop(),
	}
	for _, option := range options {
		option(d)
	}
	return d
}

func (d *LocalDockerDB) Start(ctx context.Context) error {
	d.Logger.Info("local_docker_db", zap.String("status", "starting"))

	// docker daemon must be running
	if err := exec.Command("docker", "container", "ls").Run(); err != nil {
		return errors.Trace(err)
	}

	if err := os.MkdirAll(d.MountDir, 0755); err != nil {
		return errors.Annotatef(err, "failed to create mount dir %q", d.MountDir)
	}

	info, err := dockerContainerInfo(d.Container)

	if info != nil && info.status == dockerStatusRunning {
		d.Logger.Debug("local_docker_db",
			zap.String("container_info", "already running"),
			zap.String("container", d.Container),
		)
		d.IPAddr, err = decodeIPPort(info.mappings)
		return errors.Trace(err)
	}

	var cmd *exec.Cmd
	// start or resume container
	if info == nil && errors.IsNotFound(err) {
		d.Logger.Debug("local_docker_db",
			zap.String("container_info", "not_found"),
			zap.String("container", d.Container),
		)
		volumeMapping := "type=bind,source=" + d.MountDir + ",target=/var/lib/postgresql/data"
		portMapping := d.Port + ":5432"
		cmd = exec.Command("docker", "container", "run", "--detach",
			"--name", d.Container,
			"--publish", portMapping,
			"--mount", volumeMapping,
			d.Image,
		)
	} else if info != nil && info.status != "" {
		d.Logger.Debug("local_docker_db",
			zap.String("container_info", info.status),
			zap.String("container", d.Container),
		)
		cmd = exec.Command("docker", "container", "start", info.id)
	} else {
		return errors.New("failed to start container")
	}

	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	err = cmd.Run()
	d.Logger.Debug("local_docker_db",
		zap.String("exec docker cmd", strings.Join(cmd.Args, " ")),
		zap.String("stdout", stdOut.String()),
		zap.String("stderr", stdErr.String()),
	)
	if err != nil {
		return errors.Trace(err)
	}

loop:
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			break loop
		default:
			info, err := dockerContainerInfo(d.Container)
			if err != nil {
				err = errors.Trace(err)
				break loop
			}
			if info != nil && info.status == dockerStatusRunning {
				d.IPAddr, err = decodeIPPort(info.mappings)
				break loop
			}
			time.Sleep(time.Second)
		}
	}

	return errors.Trace(err)
}

type containerInfo struct {
	id       string
	name     string
	mappings string
	status   string
}

func decodeContainerStatus(status string) string {
	// convert "Exited(0) 2 days ago" into statusExited
	if strings.HasPrefix(status, "Exited") {
		return dockerStatusExited
	}

	// convert "Up <time>" into statusRunning
	if strings.HasPrefix(status, "Up") {
		return dockerStatusRunning
	}
	return strings.ToLower(status)
}

func dockerContainerInfo(containerName string) (*containerInfo, error) {
	cmd := exec.Command("docker", "container", "ls", "-a", "--format", "{{.ID}}|{{.Status}}|{{.Ports}}|{{.Names}}")
	stdOutErr, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Annotate(err, string(stdOutErr))
	}

	s := string(stdOutErr)

	s = strings.TrimSpace(s)
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) != 4 {
			return nil, errors.Errorf("unexpected output from docker container ls %s. expected 4 parts, got %d (%v)", line, len(parts), parts)
		}
		id, status, mappings, name := parts[0], parts[1], parts[2], parts[3]
		if containerName == name {
			return &containerInfo{
				id:       id,
				name:     name,
				mappings: mappings,
				status:   decodeContainerStatus(status),
			}, nil
		}
	}
	return nil, errors.NotFoundf(containerName)
}

// given:
// 0.0.0.0:3307->3306/tcp
func decodeIPPort(mappings string) (string, error) {
	parts := strings.Split(mappings, "->")
	if len(parts) != 2 {
		return "", errors.Errorf("invalid mappings string: %q", mappings)
	}
	parts = strings.Split(parts[0], ":")
	if len(parts) != 2 {
		return "", errors.Errorf("invalid mappings string: %q", mappings)
	}
	return parts[0], nil
}

func NewLogger(level int) (*zap.Logger, error) {
	cfg := &zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.Level(int8(level))),
		Development:      true,
		Encoding:         "console", // or json
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	zapOption := zap.AddStacktrace(zapcore.ErrorLevel)
	return cfg.Build(zapOption)
}

func RunLocalDockerDB() {
	logger, _ := NewLogger(-1)
	d := NewLocalDockerDB(func(d *LocalDockerDB) {
		d.Logger = logger
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err := d.Start(ctx); err != nil {
		log.Fatal(errors.ErrorStack(err))
	}
}

func main() {
	RunLocalDockerDB()
}
