package cmd

import (
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/orcmark/pkg/orchestrator"
)

// Used by the example function to define the output type
var jsonflag, yamlflag bool

// This string is the path to a benchmark spec
var filepath string

// Automatically reap the deployments once succesfully deployed
var autoReap bool

func init() {
	benchmarkCMD.Flags().StringVar(&filepath, "path", "", "The path to a benchmark spec")
	benchmarkCMD.Flags().BoolVar(&autoReap, "reap", false, "Automatically reap the service once succesfully deployed")

	benchmarkExampleCMD.Flags().BoolVar(&jsonflag, "json", false, "Create JSON output")
	benchmarkExampleCMD.Flags().BoolVar(&yamlflag, "yaml", false, "Create YAML output")

	benchmarkCMD.AddCommand(benchmarkExampleCMD)

	// Add subcommands to the main application
	orcmarkCmd.AddCommand(benchmarkCMD)
}

var benchmarkCMD = &cobra.Command{
	Use:   "benchmark",
	Short: "Run a benchmark",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(loglevel))

		if filepath == "" {
			cmd.Help()
			log.Fatalln("No benchmark spec file was specified")
		}
		//path := fmt.Sprintf("%s/.ucptoken", os.Getenv("HOME"))
		log.Debugf("Reading benchmark spec from [%s]", filepath)
		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Fatalf("%v", err)
		}
		spec, err := orchestrator.ParseBenchmarkSpec(data)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Infof("Creating [%d] replicas of image [%s] through orchestrator [%s]", spec.Replicas, spec.Image, spec.Orchestrator)
		switch spec.Orchestrator {
		case "swarm":
			err = spec.InvokeSwarm(autoReap)
			if err != nil {
				log.Fatalf("%v", err)
			}
		case "kubernetes":
			err = spec.InvokeKubernetes(autoReap)
			if err != nil {
				log.Fatalf("%v", err)
			}
		default:
			log.Fatalf("Unknown orchestrator [%s]", spec.Orchestrator)
		}
	},
}

var benchmarkExampleCMD = &cobra.Command{
	Use:   "example",
	Short: "Create an example benchmark spec in either yaml or json",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(loglevel))

		if yamlflag == false && jsonflag == false {
			cmd.Help()
			log.Fatalln("At lease one output type is required")
		}
		if jsonflag == true {
			b := orchestrator.ExampleOutput(true, false)
			err := ioutil.WriteFile("example.json", b, 0644)
			if err != nil {
				log.Fatalf("%v", err)
			}
		}
		if yamlflag == true {
			b := orchestrator.ExampleOutput(false, true)
			err := ioutil.WriteFile("example.yaml", b, 0644)
			if err != nil {
				log.Fatalf("%v", err)
			}
		}
	},
}
