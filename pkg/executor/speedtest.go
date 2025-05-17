package executor

import (
	"context"
	"fmt"
	"time"

	"github.com/agnath18K/lumo/pkg/nlp"
	"github.com/agnath18K/lumo/pkg/speedtest"
	"github.com/agnath18K/lumo/pkg/utils"
)

// executeSpeedTest performs an internet speed test
func (e *Executor) executeSpeedTest(cmd *nlp.Command) (*Result, error) {
	// Check if there's an internet connection
	if !utils.CheckInternetConnectivity() {
		return &Result{
			Output:     "Error: No internet connection detected. Please check your network connection and try again.",
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// Create a speed tester
	tester := speedtest.NewSpeedTester()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(e.config.SpeedTestTimeout)*time.Second)
	defer cancel()

	// Determine which test to run based on the intent
	var result *speedtest.SpeedTestResult
	var err error

	intent := cmd.Intent
	if intent == "" || intent == "full" {
		// Run a full speed test
		result, err = tester.RunTest(ctx)
	} else if intent == "download" {
		// Run only a download test
		result, err = tester.RunDownloadTest(ctx)
	} else if intent == "upload" {
		// Run only an upload test
		result, err = tester.RunUploadTest(ctx)
	} else {
		// Default to full test for any other input
		result, err = tester.RunTest(ctx)
	}

	if err != nil {
		return &Result{
			Output:     fmt.Sprintf("Error performing speed test: %v", err),
			IsError:    true,
			CommandRun: cmd.RawInput,
		}, nil
	}

	// Format the result
	formattedResult := tester.FormatResult(result)

	return &Result{
		Output:     formattedResult,
		IsError:    false,
		CommandRun: cmd.RawInput,
	}, nil
}
