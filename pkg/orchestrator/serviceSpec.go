package orchestrator

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

// Service is a "general" design for micro-service that can be applied to different orchestrators
type Service struct {

	// Which orchestrator to deploy to
	Orchestrator string `json:"orchestrator" yaml:"orchestrator"`

	// Container image to run
	Image string `json:"image" yaml:"image"`
	// Number of containers to start
	Replicas int `json:"replicas" yaml:"replicas"`
	// Command to run inside the container
	CMD string `json:"cmd" yaml:"cmd"`
}

// ParseBenchmarkSpec will parse the spec and begin the benchmark
func ParseBenchmarkSpec(spec []byte) (*Service, error) {

	var benchmarkSpec Service

	// Attempt to unmarshal JSON
	err := json.Unmarshal(spec, &benchmarkSpec)

	if err != nil {
		// If there is an error Attempt to unmarshall YAML
		err = yaml.Unmarshal(spec, &benchmarkSpec)
		if err != nil {
			// Unable to parse the YAML
			return nil, fmt.Errorf("Unable to parse the file")
		}
	}

	return &benchmarkSpec, nil
}

// ExampleOutput - Create example output is a format (either JSON or YAML)
func ExampleOutput(j, y bool) []byte {
	var exampleSpec Service
	exampleSpec.Orchestrator = "swarm"
	exampleSpec.Image = "nginx:latest"
	exampleSpec.Replicas = 10
	exampleSpec.CMD = "sleep 100"

	// Create JSON output
	if j == true {
		b, _ := json.Marshal(exampleSpec)
		return b
	}

	// Create YAML output
	if y == true {
		b, _ := yaml.Marshal(exampleSpec)
		return b
	}

	// In theory this should never be reached
	return nil
}
