//go:build linux

package telemetry

import "log"

// EBPFCollector is a placeholder for real Linux eBPF logic.
type EBPFCollector struct {
	// eBPF maps and links would go here
}

func NewEBPFCollector() *EBPFCollector {
	return &EBPFCollector{}
}

func (e *EBPFCollector) Start() error {
	log.Println("Starting eBPF probes (Stub)...")
	return nil
}

func (e *EBPFCollector) Stop() error {
	return nil
}

func (e *EBPFCollector) Events() <-chan Event {
	return nil
}
