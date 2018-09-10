package orchestrator

// This will take a standardised service spec and deploy to a swarm cluster

import (
	"context"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// InvokeDockerClient - w
func (s *Service) InvokeDockerClient() error {
	c, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		return err
	}
	log.Debugf("Docker Version: %s", c.ClientVersion())
	// Turn the general service into a swarm spec
	spec := s.setSwarmSpec()
	createOptions := types.ServiceCreateOptions{}
	svc, err := c.ServiceCreate(context.Background(), spec, createOptions)
	if err != nil {
		return err
	}
	log.Printf("Service Created with ID [%s]", svc.ID)
	return nil
}

func (s *Service) setSwarmSpec() swarm.ServiceSpec {
	var spec swarm.ServiceSpec
	var service swarm.ReplicatedService
	var container swarm.ContainerSpec
	service.Replicas = &s.Replicas
	spec.Mode.Replicated = &service

	if len(s.CMD) >= 0 {
		if s.CMD[0] != "" {
			container.Command = s.CMD
		}
	}

	container.Image = s.Image
	spec.TaskTemplate.ContainerSpec = &container
	return spec
}
