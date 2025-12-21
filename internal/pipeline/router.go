package pipeline

import (
	"log"

	"nexus/internal/config"
	// For now still tightly coupled, but logic separated
)

// Destination represents an output sink (e.g. Splunk).
type Destination interface {
	SendEvent(event interface{}) error
}

// Router handles data governance and routing.
type Router struct {
	Config       config.PipelineConfig
	Destinations map[string]Destination
}

func NewRouter(cfg config.PipelineConfig) *Router {
	return &Router{
		Config:       cfg,
		Destinations: make(map[string]Destination),
	}
}

func (r *Router) AddDestination(name string, dest Destination) {
	if dest != nil {
		r.Destinations[name] = dest
	}
}

// Route directs an event to the appropriate destinations based on policy.
// It accepts interface{} to allow routing of ANY data structure (OCSF, CIS, Custom).
func (r *Router) Route(event interface{}) {
	// Determine event "Type" (simplified taxonomy)
	// - Inventory
	// - Telemetry
	// - Vulnerability
	// We can infer this or add a Type() method to the Event interface.
	// For now, let's just broadcast to all configured routes.

	// TODO: Implement sophisticated type matching using reflection or type assertion.

	for name, dest := range r.Destinations {
		// Governance Check: Should this destination receive this data?
		// if !r.Allowed(name, event.Type) { continue }

		if err := dest.SendEvent(event); err != nil {
			log.Printf("[ROUTER] Failed to send to %s: %v", name, err)
		}
	}
}
