// main is the main entrypoint of the application
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/config"
	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-script/internal/runner"
)

func main() {
	appContext := context.Background()

	/*
	 * Parse command line arguments from the user
	 */
	args, err := config.ParseAndValidateCLIArgs(os.Args)
	if err != nil {
		if err == flag.ErrHelp {
			return
		}
		fmt.Println("error passing arguments:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	/*
	 * Set up the application logger
	 */
	appLogger, err := config.InitializeLogger(args.Verbose)
	if err != nil {
		fmt.Println("error initializing the logger: ", err)
		os.Exit(1)
	}
	defer appLogger.Close()

	/*
	 * Initialize the fleet updater
	 */
	updater, err := runner.NewFleetUpdater(appContext, appLogger, args)
	if err != nil {
		slog.Error("error building a fleet updater", "error", strings.Replace(err.Error(), "\n", ", ", -1))
		os.Exit(1)
	}
	defer updater.Cleanup()

	/*
	 * Update the instances in the fleet
	 */
	_, err = updater.UpdateInstances(appContext)
	if err != nil {
		if err != runner.UpdateFailedError {
			slog.Error("error updating instances", "error", err)
		}

		os.Exit(1)
	}
}
