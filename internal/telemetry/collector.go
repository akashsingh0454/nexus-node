package telemetry

// Event represents a generic telemetry event mapped to OCSF.
type Event struct {
	Type string
	Data interface{} // Should be an OCSF struct
}

// Collector is the interface for any telemetry source.
type Collector interface {
	// Start begins the collection process (non-blocking).
	Start() error
	// Stop halts the collection.
	Stop() error
	// Events returns the channel to read normalized events from.
	Events() <-chan Event
}
