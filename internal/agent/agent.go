package agent

import (
	"context"
	"log"
	"time"

	"nexus/internal/config"
	"nexus/internal/engine"
	"nexus/internal/exporter/splunk"
	"nexus/internal/module"
	"nexus/internal/modules/patcher"
	"nexus/internal/modules/trivy"
	"nexus/internal/telemetry"
	"nexus/pkg/ocsf"
)

// Agent allows for the orchestration of various security modules.
type Agent struct {
	Config        *config.Config
	ModuleManager *module.Manager
	SplunkClient  *splunk.Client
	Telemetry     telemetry.Collector
	TrivyManager  *trivy.Manager
	Patcher       *patcher.Manager
	Remediator    *engine.Remediator
}

// NewAgent creates a new instance of the Nexus Node Agent.
func NewAgent(cfg *config.Config) *Agent {
	mgr := module.NewManager()

	// Register modules from config
	if cfg.Modules.Osquery.Enabled {
		mod := module.NewModule("osquery", cfg.Modules.Osquery.BinaryPath, []string{"--pidfile", "osquery.pid"})
		mgr.Register(mod)
	}

	// Initialize Splunk Client
	var splunkClient *splunk.Client
	if cfg.Exporter.Splunk.Enabled {
		splunkClient = splunk.NewClient(cfg.Exporter.Splunk.URL, cfg.Exporter.Splunk.Token)
	}

	// Initialize Telemetry
	tel := telemetry.NewSimulatedCollector()

	// Initialize Patcher (Dry Run)
	patcherMgr := patcher.NewManager(true)

	// Initialize Remediator
	remediator := engine.NewRemediator(&cfg.Remediation, patcherMgr)

	return &Agent{
		Config:        cfg,
		ModuleManager: mgr,
		SplunkClient:  splunkClient,
		Telemetry:     tel,
		// Assuming trivy is in ./bin/trivy.exe for now
		TrivyManager: trivy.NewManager("./bin/trivy.exe"),
		Patcher:      patcherMgr,
		Remediator:   remediator,
	}
}

// Run starts the main event loop of the agent.
// It blocks until the context is cancelled or a fatal error occurs.
func Run(a *Agent, ctx context.Context) error {
	log.Println("Agent Chassis initialized and running.")

	// Start modules
	a.ModuleManager.StartAll(ctx)
	defer a.ModuleManager.StopAll()

	// Start Telemetry
	if a.Telemetry != nil {
		if err := a.Telemetry.Start(); err != nil {
			log.Printf("Failed to start telemetry: %v", err)
		} else {
			defer a.Telemetry.Stop()
		}
	}

	// Test Trigger: Check for updates on startup (async)
	go func() {
		time.Sleep(5 * time.Second) // Wait for boot
		if _, err := a.Patcher.ListUpdates(); err == nil {
			log.Println("Patcher: ListUpdates completed successfully (Dry Run).")
		}

		// TEST: Simulate Auto-Remediation
		time.Sleep(2 * time.Second)
		log.Println("[TEST] Simulating Critical Vulnerability Discovery...")
		fakeVuln := ocsf.VulnerabilityFinding{
			Vulnerability: ocsf.Vulnerability{
				Title:    "Mozilla.Firefox", // ID compatible with winget for test
				Severity: "CRITICAL",
				CVE:      "CVE-2025-MOCK",
			},
		}
		if a.Remediator.Evaluate(fakeVuln) {
			a.Remediator.Remediate(fakeVuln)
		}
	}()

	// Heartbeat ticker
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping agent main loop...")
			return nil

		// Telemetry Events
		case evt := <-a.Telemetry.Events():
			if a.SplunkClient != nil {
				// Fire and forget for speed
				// In production: use a worker pool or buffer
				go a.SplunkClient.SendEvent(evt.Data)
			}

		case <-ticker.C:
			// Send heartbeat
			log.Println("Heartbeat: Agent is alive.")

			if a.SplunkClient != nil {
				hb := ocsf.NewInventoryEvent(ocsf.Device{
					Hostname: a.Config.Agent.Name,
					OS:       "Windows", // Placeholder functionality
				})
				if err := a.SplunkClient.SendEvent(hb); err != nil {
					log.Printf("Failed to send heartbeat to Splunk: %v", err)
				}
			}
		}
	}
}

// Run is an instance method wrapper for the logic
func (a *Agent) Run(ctx context.Context) error {
	return Run(a, ctx)
}
