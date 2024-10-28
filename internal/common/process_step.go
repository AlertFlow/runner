package common

import (
	"alertflow-runner/internal/actions"
	"alertflow-runner/internal/executions"
	"alertflow-runner/pkg/models"
	"time"

	log "github.com/sirupsen/logrus"
)

func processStep(flow models.Flows, payload models.Payload, steps []models.ExecutionSteps, step models.ExecutionSteps, execution models.Execution) (data map[string]interface{}, finished bool, canceled bool, failed bool, err error) {
	// set step to running
	step.Pending = false
	step.Running = true
	step.StartedAt = time.Now()
	step.RunnerID = execution.RunnerID

	if err := executions.UpdateStep(execution.ID.String(), step); err != nil {
		log.Error(err)
		return nil, false, false, false, err
	}

	action, found := actions.SearchAction(step.ActionID)

	if !found {
		log.Warnf("Action %s not found", step.ActionID)

		step.ActionMessages = append(step.ActionMessages, "Action not found")
		step.Running = false
		step.Error = true
		step.Finished = true
		step.FinishedAt = time.Now()

		if err := executions.UpdateStep(execution.ID.String(), step); err != nil {
			log.Error(err)
			return nil, false, false, false, err
		}

		return nil, false, false, true, nil
	}

	if fn, ok := action.Function.(func(execution models.Execution, flow models.Flows, payload models.Payload, steps []models.ExecutionSteps, step models.ExecutionSteps, action models.Actions) (data map[string]interface{}, finished bool, canceled bool, failed bool)); ok {
		data, finished, canceled, failed := fn(execution, flow, payload, steps, step, models.Actions{})

		if failed {
			return nil, false, false, true, nil
		} else if canceled {
			return nil, false, true, false, nil
		} else if finished {
			return data, true, false, false, nil
		}
	}

	return data, true, false, false, nil
}
