package main

import (
	"flag"
	"fmt"
	"gopkg.in/launchdarkly/go-server-sdk.v5/ldcomponents"
	"gopkg.in/launchdarkly/go-server-sdk.v5/ldfiledata"
	"gopkg.in/launchdarkly/go-server-sdk.v5/ldfilewatch"
	"os"
	"sync"
	"time"

	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
)

// Set sdkKey to your LaunchDarkly SDK key.
const sdkKey = "sdk-8c498d63-c1b4-43cb-8423-29180cd52efc"

// Set featureFlagKey to the feature flag key you want to evaluate.
const featureFlagKey = "test-flag"

func showMessage(s string) { fmt.Printf("*** %s\n\n", s) }

func main() {
	userName := flag.String("name", "default-user", "User Name")
	environment := flag.String("env", "local", "Environment")
	flag.Parse()

	if sdkKey == "" {
		showMessage("Please edit main.go to set sdkKey to your LaunchDarkly SDK key first")
		os.Exit(1)
	}

	var ldClient *ld.LDClient

	if *environment == "local" {
		var config ld.Config
		config.DataSource = ldfiledata.DataSource().
			FilePaths("feature-flags.json").
			Reloader(ldfilewatch.WatchFiles)
		config.Events = ldcomponents.NoEvents()

		ldClient, _ = ld.MakeCustomClient(sdkKey, config, 5*time.Second)
	} else {
		ldClient, _ = ld.MakeClient(sdkKey, 5*time.Second)
	}

	if ldClient.Initialized() {
		showMessage("SDK successfully initialized!")
	} else {
		showMessage("SDK failed to initialize")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go printFlagValueAndUserName(ldClient, *userName, &wg)
	wg.Wait()

	// Here we ensure that the SDK shuts down cleanly and has a chance to deliver analytics
	// events to LaunchDarkly before the program exits. If analytics events are not delivered,
	// the user properties and flag usage statistics will not appear on your dashboard. In a
	// normal long-running application, the SDK would continue running and events would be
	// delivered automatically in the background.
	ldClient.Close()
}

func printFlagValueAndUserName(ldClient *ld.LDClient, userName string, wg *sync.WaitGroup) {
	defer wg.Done()

	user := lduser.NewUserBuilder(userName).
		Name(userName).
		Build()

	for i := 0; i < 60; i++ {
		flagValue, err := ldClient.BoolVariation(featureFlagKey, user, false)
		if err != nil {
			showMessage("error: " + err.Error())
		}
		showMessage(fmt.Sprintf("Feature flag '%s' is %t for %s user", featureFlagKey, flagValue, userName))
		time.Sleep(time.Duration(1) * time.Second)
	}
}
