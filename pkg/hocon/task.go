package hocon

import (
	"github.com/falcosecurity/kilt/pkg/kilt"
	"github.com/go-akka/configuration"
)

func extractTask(config *configuration.Config) (*kilt.Task, error) {
	var task = new(kilt.Task)

	if config.HasPath("task.pid_mode") {
		task.PidMode = config.GetString("task.pid_mode")
	}

	return task, nil
}
