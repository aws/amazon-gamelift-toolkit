package runner

import (
	"fmt"
	"log/slog"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/gamelift"
	"github.com/pterm/pterm"
)

// InstanceProgressWriter is used to track and display the progress of updating a single GameLift instance to the user
type InstanceProgressWriter struct {
	verbose             bool
	instanceId          string
	instanceIp          string
	instanceUpdateState InstanceUpdateState
	progressBar         *pterm.ProgressbarPrinter
}

// NewInstanceProgressWriter builds a new progress writer for the provided instance
func NewInstanceProgressWriter(instance *gamelift.Instance, verbose bool) (*InstanceProgressWriter, error) {
	if verbose {
		return &InstanceProgressWriter{verbose: verbose}, nil
	}

	instanceUpdateState := UpdateStateNotStarted

	// Set up the progress bar we'll be showing to the user
	progressBar, err := pterm.DefaultProgressbar.
		WithTotal(int(UpdateStateCount)).
		WithShowElapsedTime(false).
		WithTitle(stateString(instance.InstanceId, instance.IpAddress, instanceUpdateState)).Start()
	if err != nil {
		return nil, err
	}

	return &InstanceProgressWriter{
		instanceId:          instance.InstanceId,
		instanceIp:          instance.IpAddress,
		verbose:             verbose,
		instanceUpdateState: instanceUpdateState,
		progressBar:         progressBar,
	}, nil
}

// UpdateState update the instance progress with newState, and display any relevant information to the user
func (i *InstanceProgressWriter) UpdateState(newState InstanceUpdateState) {
	if i.verbose {
		return
	}

	// Calculate the amount to update the progress bar
	diff := newState - i.instanceUpdateState

	// Update the state and title with the new state
	i.instanceUpdateState = newState
	i.progressBar.UpdateTitle(stateString(i.instanceId, i.instanceIp, i.instanceUpdateState))

	// Actually update the progress bar (do this last, otherwise it causes display issues)
	i.progressBar.Add(int(diff))

	// If we're at the end, show a nice notice the to the user
	if newState == UpdateStateCount {
		pterm.Success.Println(i.instanceId)
	}
}

// UpdateFailed is used to alert the user that updating this instance failed
func (i *InstanceProgressWriter) UpdateFailed(err error) {
	if i.verbose {
		return
	}

	_, stopErr := i.progressBar.Stop()
	if stopErr != nil {
		slog.Debug("error stopping progress bar", "error", stopErr)
	}
}

func stateString(instanceId, instanceIp string, state InstanceUpdateState) string {
	return fmt.Sprintf("%s (%s) %s", instanceId, instanceIp, state.String())
}
