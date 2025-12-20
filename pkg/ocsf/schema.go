package ocsf

import "time"

// OCSF Base Event fields
type BaseEvent struct {
	ClassUID int       `json:"class_uid"`
	Category int       `json:"category_uid"`
	Activity int       `json:"activity_id"`
	Time     time.Time `json:"time"`
	Severity string    `json:"severity"`
	Message  string    `json:"message"`
}

// Device Inventory Object (part of many classes)
type Device struct {
	Hostname   string `json:"hostname"`
	OS         string `json:"os"`
	OSVersion  string `json:"os_version"`
	IP         string `json:"ip"`
	MAC        string `json:"mac"`
	InstanceID string `json:"instance_id,omitempty"`
}

// InventoryInfo Class (1001 - mocked ID for illustration)
// Represents a snapshot of system state
type InventoryInfo struct {
	BaseEvent
	Device Device `json:"device"`
}

// VulnerabilityFinding Class (2001)
type VulnerabilityFinding struct {
	BaseEvent
	Device        Device        `json:"device"`
	Vulnerability Vulnerability `json:"vulnerability"`
}

type Vulnerability struct {
	CVE         string   `json:"cve"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
	References  []string `json:"references,omitempty"`
}

// ProcessActivity Class (4001)
type ProcessActivity struct {
	BaseEvent
	Device  Device  `json:"device"`
	Process Process `json:"process"`
	Actor   Actor   `json:"actor"`
}

type Process struct {
	PID         int    `json:"pid"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	CommandLine string `json:"cmd_line"`
}

type Actor struct {
	User string `json:"user"`
}

// Function to create a new standard Inventory Event
func NewInventoryEvent(device Device) InventoryInfo {
	return InventoryInfo{
		BaseEvent: BaseEvent{
			ClassUID: 5001, // Example ID for Inventory
			Category: 5,    // Discovery
			Time:     time.Now(),
			Severity: "Informational",
		},
		Device: device,
	}
}
