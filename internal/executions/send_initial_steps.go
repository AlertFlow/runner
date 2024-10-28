package executions

import (
	"alertflow-runner/config"
	"alertflow-runner/pkg/models"
	"time"
)

var initialSteps = []models.ExecutionSteps{
	{
		ActionID:       "runner_pick_up",
		ActionName:     "Runner Pick Up",
		ActionMessages: []string{"Runner Picked Up Execution"},
		Icon:           "solar:rocket-2-bold-duotone",
		RunnerID:       config.Config.RunnerID,
		StartedAt:      time.Now(),
		Pending:        false,
		Finished:       true,
		FinishedAt:     time.Now(),
	},
	{
		ActionID:       "execution_start",
		ActionName:     "Execution Start",
		ActionMessages: []string{"Execution Started"},
		Icon:           "solar:rocket-2-bold-duotone",
		RunnerID:       config.Config.RunnerID,
		StartedAt:      time.Now(),
		Pending:        false,
		Finished:       true,
		FinishedAt:     time.Now(),
	},
	{
		ActionID:   "collect_data",
		ActionName: "Collect Data",
		Icon:       "solar:inbox-archive-linear",
		Pending:    true,
	},
	{
		ActionID:   "pattern_check",
		ActionName: "Pattern Check",
		Icon:       "solar:list-check-minimalistic-bold",
		Pending:    true,
	},
	{
		ActionID:   "flow_actions_check",
		ActionName: "Flow Actions Check",
		Icon:       "solar:minimalistic-magnifer-linear",
		Pending:    true,
	},
}

// SendInitialSteps sends initial steps to alertflow
func SendInitialSteps(execution models.Execution) (stepsWithIDs []models.ExecutionSteps, err error) {
	for i, step := range initialSteps {
		step.ExecutionID = execution.ID.String()
		stepID, err := SendStep(execution, step)
		if err != nil {
			return nil, err
		}
		step.ID = stepID.ID
		initialSteps[i] = step
	}
	return initialSteps, nil
}
