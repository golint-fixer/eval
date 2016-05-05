/*
 * Copyright 2015 Fabrício Godoy
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package eval

import (
	"fmt"
	"time"

	"gopkg.in/raiqub/docker.v0"
)

const (
	// ImageRedisName defines the default image name of Redis Docker image.
	ImageRedisName = "redis"

	// StartupTimeout defines the maximum duration of time waiting for container
	// startup
	StartupTimeout = 30 * time.Second
)

// A Environment represents a Docker testing environment.
type Environment struct {
	dockerBin *docker.Docker
	image     *docker.Image
	container *docker.Container
	run       func() (*docker.Container, error)
}

// NewMongoDBEnvironment creates a instance that allows to use MongoDB for
// testing
func NewMongoDBEnvironment() *Environment {
	d := docker.NewDocker()
	mongo := docker.NewImageMongoDB(d)

	return &Environment{
		dockerBin: d,
		image:     &mongo.Image,
		run: func() (*docker.Container, error) {
			cfg := docker.NewRunConfig()
			cfg.Detach()
			return mongo.RunLight(cfg)
		},
	}
}

// NewRedisEnvironment creates a instance that allows to use Redis for testing.
func NewRedisEnvironment() *Environment {
	d := docker.NewDocker()
	redis := docker.NewImage(d, ImageRedisName)

	return &Environment{
		dockerBin: d,
		image:     redis,
		run: func() (*docker.Container, error) {
			cfg := docker.NewRunConfig()
			cfg.Detach()
			return redis.Run(cfg)
		},
	}
}

// Applicability tests whether current testing environment can be run on current
// host.
func (s *Environment) Applicability() (bool, *ErrUser) {
	if !s.dockerBin.HasBin() {
		return false, &ErrUser{
			Warn,
			"Docker binary was not found",
		}
	}

	_, err := s.dockerBin.Run("ps")
	if err != nil {
		return false, &ErrUser{
			Warn,
			"Docker is installed but is not running or current user " +
				"is lacking permissions",
		}
	}

	return true, nil
}

// Network returns network information from current running container.
func (s *Environment) Network() ([]docker.NetworkNode, error) {
	if s.container == nil {
		return nil, NotRunningError(s.image.Name())
	}

	nodes, err := s.container.NetworkNodes()
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// Run starts a new Docker instance for testing environment.
func (s *Environment) Run() (bool, *ErrUser) {
	if err := s.image.Setup(); err != nil {
		return false, &ErrUser{
			Fatal,
			fmt.Sprintf("Error setting up Docker: %v", err),
		}
	}

	var err error
	s.container, err = s.run()
	if err != nil {
		s.Stop()
		return false, &ErrUser{
			Fatal,
			fmt.Sprintf("Error running a new Docker container: %v", err),
		}
	}

	if s.container.HasExposedPorts() {
		if err := s.container.WaitStartup(StartupTimeout); err != nil {
			s.Stop()
			return false, &ErrUser{
				Fatal,
				fmt.Sprintf("Timeout waiting Docker instance to respond: %v", err),
			}
		}
	} else {
		timeout := time.After(StartupTimeout)
		for {
			select {
			case <-timeout:
				inspect, err := s.container.Inspect()
				if err != nil || !inspect[0].State.Running {
					return false, &ErrUser{
						Fatal,
						fmt.Sprint("Timeout waiting container startup"),
					}
				}
			}
		}
	}

	return true, nil
}

// Stop removes current running testing environment.
func (s *Environment) Stop() {
	if s.container == nil {
		return
	}

	s.container.Kill()
	s.container.Remove()
}
