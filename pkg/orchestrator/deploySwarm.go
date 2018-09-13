package orchestrator

// This will take a standardised service spec and deploy to a swarm cluster

import (
	"context"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// InvokeSwarm - w
func (s *Service) InvokeSwarm(autoReap bool) error {
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
	log.Infof("Service Created with ID [%s]", svc.ID)
	// If we're not auto reaping then return once the deployment has been set
	if autoReap == false {
		return nil
	}

	log.Infoln("Will reap the deployment (after 10 seconds) once all replicas are running")

	// Wait a minimum of ten seconds before beginning the auto reap
	time.Sleep(10 * time.Second)
	orcMarkTasks, err := c.TaskList(context.Background(), types.TaskListOptions{})
	if err != nil {
		return err
	}

	for s.tasksDeployed(svc.ID, orcMarkTasks) != true {
		time.Sleep(1 * time.Second)
		orcMarkTasks, err = c.TaskList(context.Background(), types.TaskListOptions{})
		if err != nil {
			return err
		}
	}
	err = c.ServiceRemove(context.Background(), svc.ID)
	if err != nil {
		return err
	}
	log.Infof("Succesfully removed service [%s]", svc.ID)
	return nil
}

// tasksDeployed will look at all of the tasks in the swarm and identify if the correct amount of replicas are present
func (s *Service) tasksDeployed(serviceID string, tasks []swarm.Task) bool {
	var replicaCount uint64
	for taskID := range tasks {
		if tasks[taskID].ServiceID == serviceID {
			if tasks[taskID].Status.State == swarm.TaskStateRunning {
				replicaCount++
			}
		}
	}
	log.Debugf("Currently [%d] running replicas", replicaCount)

	if replicaCount == s.Replicas {
		return true
	}
	return false
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
