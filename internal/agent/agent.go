package agent

import (
	"context"
	"log"
	"net/http"
	"time"

	"nexus/internal/config"
	"nexus/internal/engine"
	"nexus/internal/exporter/splunk"
	"nexus/internal/module"
	"nexus/internal/modules/hive"
	"nexus/internal/modules/patcher"
	"nexus/internal/modules/trivy"
	"nexus/internal/pipeline"
	"nexus/internal/telemetry"
	"nexus/internal/transport"
	"nexus/pkg/ocsf"
)

// Agent allows for the orchestration of various security modules.
type Agent struct {
	Config        *config.Config
	ModuleManager *module.Manager
	Router        *pipeline.Router
	Telemetry     telemetry.Collector
	TrivyManager  *trivy.Manager
	Patcher       *patcher.Manager
	HiveManager   *hive.Manager
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

	// 1. Initialize Secure Transport (mTLS / CA)
	secureClient, err := transport.NewClient(cfg.Security)
	if err != nil {
		log.Printf("Failed to create secure transport: %v. Falling back to default.", err)
		secureClient = &http.Client{Timeout: 30 * time.Second}
	}

	// 2. Initialize Exporters (Splunk)
	var splunkClient *splunk.Client
	if cfg.Pipeline.Exporters.Splunk.Enabled {
		splunkClient = splunk.NewClient(
			cfg.Pipeline.Exporters.Splunk.URL,
			cfg.Pipeline.Exporters.Splunk.Token,
			secureClient,
		)
	}

	// 3. Initialize Pipeline Router
	router := pipeline.NewRouter(cfg.Pipeline, splunkClient)

	// Initialize Telemetry
	tel := telemetry.NewSimulatedCollector()

	// Initialize Patcher (Dry Run)
	patcherMgr := patcher.NewManager(true)

	// Initialize Remediator
	remediator := engine.NewRemediator(&cfg.Remediation, patcherMgr)

	// Initialize Hive (P2P)
	hiveMgr := hive.NewManager(cfg.Modules.Hive, cfg.Agent.Name)

	return &Agent{
		Config:        cfg,
		ModuleManager: mgr,
		Router:        router,
		Telemetry:     tel,
		// Assuming trivy is in ./bin/trivy.exe for now
		TrivyManager: trivy.NewManager("./bin/trivy.exe"),
		Patcher:      patcherMgr,
		HiveManager:  hiveMgr,
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

	// Start Hive
	if a.HiveManager != nil {
		a.HiveManager.Start(ctx)
		defer a.HiveManager.Stop()
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

		// Telemetry Events calling Router
		case evt := <-a.Telemetry.Events():
			// Fire and forget routing
			go a.Router.Route(evt)

		case <-ticker.C:
			// Send heartbeat via Router
			log.Println("Heartbeat: Agent is alive.")

			hb := ocsf.NewInventoryEvent(ocsf.Device{
				Hostname: a.Config.Agent.Name,
				OS:       "Windows", // Placeholder functionality
			})

			// We wrap OCSF event if needed, or Router handles raw OCSF
			// For now, assuming ocsf.Event matches what Router expects (which is just interface{} or ocsf.Event base)
			// Actually Router expects ocsf.Event, InventoryEvent embeds BaseEvent so it satisfies interfaces?
			// Let's assume yes or cast.
			go a.Router.Route(hb)
		}
	}
}

// Run is an instance method wrapper for the logic
func (a *Agent) Run(ctx context.Context) error {
	return Run(a, ctx)
}
