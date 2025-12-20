package telemetry

import (
	"math/rand"
	"time"

	"nexus/pkg/ocsf"
)

// SimulatedCollector generates fake process activity for testing.
type SimulatedCollector struct {
	events chan Event
	stop   chan struct{}
}

func NewSimulatedCollector() *SimulatedCollector {
	return &SimulatedCollector{
		events: make(chan Event, 100),
		stop:   make(chan struct{}),
	}
}

func (s *SimulatedCollector) Start() error {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-s.stop:
				return
			case <-ticker.C:
				// Generate a fake process event
				evt := s.generateEvent()
				s.events <- Event{Type: "ProcessActivity", Data: evt}
			}
		}
	}()
	return nil
}

func (s *SimulatedCollector) Stop() error {
	close(s.stop)
	close(s.events)
	return nil
}

func (s *SimulatedCollector) Events() <-chan Event {
	return s.events
}

func (s *SimulatedCollector) generateEvent() ocsf.ProcessActivity {
	// Randomly pick a process name
	procs := []string{"chrome.exe", "svchost.exe", "powershell.exe", "cmd.exe"}
	name := procs[rand.Intn(len(procs))]

	return ocsf.ProcessActivity{
		BaseEvent: ocsf.BaseEvent{
			ClassUID: 4001,
			Activity: 1, // Start
			Time:     time.Now(),
			Severity: "Informational",
		},
		Process: ocsf.Process{
			PID:  rand.Intn(10000),
			Name: name,
			Path: "C:\\Windows\\System32\\" + name,
		},
		Actor: ocsf.Actor{
			User: "SYSTEM",
		},
	}
}
