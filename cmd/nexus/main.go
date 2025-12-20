package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nexus/internal/agent"
	"nexus/internal/config"
	"nexus/internal/logger"
)

func main() {
	configFile := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	logger.Setup()

	// Load config from *configFile
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Loaded config for agent: %s", cfg.Agent.Name)

	log.Println("Starting Nexus Node...")

	// Create a context that acts as a signal to cancel all workers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the Agent Chassis
	ag := agent.NewAgent(cfg)

	// Run the agent in a goroutine
	go func() {
		if err := ag.Run(ctx); err != nil {
			log.Fatalf("Agent runtime error: %v", err)
		}
	}()

	// Wait for OS signals (SIGINT, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received signal: %v. Shutting down...", sig)

	// Trigger graceful shutdown
	cancel()

	// Give components time to clean up
	time.Sleep(2 * time.Second)
	log.Println("Nexus Node shutdown complete.")
}
