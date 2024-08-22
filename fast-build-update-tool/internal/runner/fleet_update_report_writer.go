package runner

import (
	"strings"

	"github.com/pterm/pterm"
)

// FleetUpdateReportWriter is used to print out a report of results when a fleet update is in progress,
// and when it is done updating (either successfully or in an error state).
type FleetUpdateReportWriter struct {
	fleetId string
	verbose bool
}

// NewFleetUpdateReportWriter will generate a new FleetUpdateReportWriter
func NewFleetUpdateReportWriter(fleetId string, verbose bool) *FleetUpdateReportWriter {
	return &FleetUpdateReportWriter{fleetId: fleetId, verbose: verbose}
}

// Preparing will print any relevant messaging around the tool entering the initial "Preparing" state
func (f *FleetUpdateReportWriter) Preparing() {
	if f.verbose {
		return
	}

	pterm.Info.Printf("Preparing fleet for update: %s\n", f.fleetId)
}

// StartUpdatingInstances will print any relevant messaging before the tool starts to update instances in the fleet
func (f *FleetUpdateReportWriter) StartUpdatingInstances(instanceCount int) {
	if f.verbose {
		return
	}

	pterm.Info.Printf("Starting update process for %d instance(s)\n", instanceCount)
}

// ReportResults will print data round the results of a fleet update (either successful or failed)
func (f *FleetUpdateReportWriter) ReportResults(results *FleetUpdateResults) {
	if f.verbose {
		return
	}

	if len(results.InstancesFailedUpdate) == 0 {
		pterm.Success.Printf("Fleet Update Succeeded! Updated %d instance(s)\n", results.InstancesUpdated)
	} else {
		pterm.Error.Printf("Fleet Update Failed. Failed to update %d instance(s)\n", len(results.InstancesFailedUpdate))
		pterm.Error.Printf("Instance(s) failed: %s\n", strings.Join(results.InstancesFailedUpdate, ", "))
		pterm.Printf("Instance(s) Successfully Updated: %d\n", results.InstancesUpdated)
		pterm.Printf("Total Instance(s) Found: %d\n", results.InstancesFound)
	}
}
