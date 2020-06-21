package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/matrix-org/complement/internal"
	"github.com/matrix-org/complement/internal/config"
	"github.com/sirupsen/logrus"
)

/*
This is the main entry point for Complement. TestMain governs:
 - Loading blueprints.
 - Creating homeserver base containers.
 - Running blueprints on containers.
 - Committing the containers as new images with well-defined names $blueprintName:$hsName
Tests will then ask for a deployment of a blueprint by name which will deploy potentially
multiple servers (if testing Federation). Those servers can then be poked until the deployment
is destroyed.

setup (before tests are run)                      +---------------------+
                                                  |              Docker |
 +------------+          +---------+    runs      |  +--------+         |
 | Blueprints | -------> | Builder | -----------> |  | Images |         |
 +------------+          +---------+   commits    |  +--------+         |
                                                  |                     |
                                                  |                     |
---------------------------------------------------------------------------------
tests                                             |                     |
                                                  |                     |
 +-------+                +----------+            |  +------------+     |
 | Tests | -------------> | Deployer | ---------> |  | Containers |     |
 +-------+                +----------+   runs     |  +------------+     |
                                                  +---------------------+

*/

// TestMain is the main entry point for Complement. It will process COMPLEMENT_ env vars and build blueprints
// to images before executing the tests.
func TestMain(m *testing.M) {
	cfg := config.NewConfigFromEnvVars()
	log.Printf("config: %+v", cfg)
	builder, err := internal.Start(cfg)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	// we use GMSL which uses logrus by default. We don't want those logs in our test output unless they are Serious.
	logrus.SetLevel(logrus.ErrorLevel)

	exitCode := m.Run()
	builder.Cleanup()
	os.Exit(exitCode)
}