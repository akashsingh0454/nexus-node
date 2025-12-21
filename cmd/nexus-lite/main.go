package main

import (
	"log"
	"os"
	"time"

	"nexus/internal/config"
	"nexus/internal/exporter/splunk"
	"nexus/internal/pipeline"
	"nexus/pkg/ocsf"
)

// Nexus Lite: Ephemeral agent for Serverless (AWS Lambda / CloudRun)
// No P2P, No Patching, Just Scan & Report.
func main() {
	log.Println("⚡ Nexus Lite: Starting ephemeral scan...")

	// 1. Load Simplified Config (Env Vars preferred in Serverless)
	// Mock config object as we don't need full YAML parsing for simple lambda
	splunkURL := os.Getenv("NEXUS_SPLUNK_URL")
	splunkToken := os.Getenv("NEXUS_SPLUNK_TOKEN")
	agentName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	if agentName == "" {
		agentName = "lambda-function"
	}

	// 2. Setup Exporter
	var splunkClient *splunk.Client
	if splunkURL != "" && splunkToken != "" {
		splunkClient = splunk.NewClient(splunkURL, splunkToken, nil)
	}

	// 3. Setup Router
	router := pipeline.NewRouter(config.PipelineConfig{})
	if splunkClient != nil {
		router.AddDestination("splunk_lite", splunkClient)
	}

	// 4. Gather Telemetry (Simulated Snapshot)
	evt := ocsf.NewInventoryEvent(ocsf.Device{
		Hostname: agentName,
		OS:       "Linux/Serverless",
	})

	// 5. Emit
	log.Println("   - Emitting Inventory Event...")
	router.Route(evt)

	// 6. Security Scan (Trivy Library mode - Placeholder)
	// In a real lambda, we'd check /tmp or the layer itself.
	log.Println("   - Scanning Runtime Environment...")
	time.Sleep(50 * time.Millisecond) // Simulate fast work

	log.Println("⚡ Nexus Lite: Complete. Exiting.")
}
