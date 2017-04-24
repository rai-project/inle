package kernel

import (
	"bytes"
	"io"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/rai-project/docker"
)

type sessionManager struct {
	cache   *cache.Cache
	created time.Time
}

type session struct {
	name            string
	buf             *bytes.Buffer
	dockerClient    *docker.Client
	dockerContainer *docker.Container
	created         time.Time
}

var (
	DefaultSessionManager *sessionManager
)

func NewSessionManager() *sessionManager {
	return &sessionManager{
		cache:   cache.New(cache.NoExpiration, cache.NoExpiration),
		created: time.Now(),
	}
}

func (s *sessionManager) Add(name string, val *session) {
	s.cache.Set(name, val, cache.NoExpiration)
}

func (s *sessionManager) Remove(name string) {
	s.cache.Delete(name)
}

func (s *sessionManager) Has(name string) bool {
	_, ok := s.cache.Get(name)
	return ok
}

func (s *sessionManager) Get(name string) (*session, error) {
	e, ok := s.cache.Get(name)
	if !ok {
		return nil, errors.Errorf("the session %s was not found", name)
	}
	q, ok := e.(*session)
	if !ok {
		return nil, errors.New("invalid session type")
	}
	return q, nil
}

func (s *sessionManager) Flush() {
	for _, it := range s.cache.Items() {
		if rc, ok := it.Object.(io.Closer); ok {
			rc.Close()
		}
	}
}

func NewSession(name string) (*session, error) {
	imageName := "ubuntu:17.04"
	srcDir := "/src"
	buildDir := "/build"

	buf := bytes.NewBufferString("")

	clnt, err := docker.NewClient(
		docker.Stdout(buf),
		docker.Stderr(buf),
		docker.Stdin(nil),
	)
	if err != nil {
		return nil, err
	}

	containerOpts := []docker.ContainerOption{
		docker.Image(imageName),
		docker.AddEnv("IMAGE_NAME", imageName),
		docker.AddVolume(srcDir),
		docker.AddVolume(buildDir),
		docker.WorkingDirectory(buildDir),
		docker.Shell([]string{"/bin/bash"}),
		docker.Entrypoint([]string{}),
	}

	cont, err := docker.NewContainer(clnt, containerOpts...)
	if err != nil {
		return nil, err
	}

	if err := cont.Start(); err != nil {
		return nil, err
	}

	sess := &session{
		name:            name,
		buf:             buf,
		dockerClient:    clnt,
		dockerContainer: cont,
		created:         time.Now(),
	}
	DefaultSessionManager.Add(name, sess)
	return sess, nil
}

func (s *session) Exec(cmd string) error {
	exec, err := docker.NewExecutionFromString(s.dockerContainer, cmd)
	if err != nil {
		return err
	}
	return exec.Run()
}

func (s *session) Close() error {
	s.dockerContainer.Close()
	s.dockerClient.Close()
	DefaultSessionManager.Remove(s.name)
	return nil
}

func init() {
	DefaultSessionManager = NewSessionManager()
}
