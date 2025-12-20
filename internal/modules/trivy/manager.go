package trivy

import (
	"encoding/json"
	"log"
	"os/exec"
)

// Manager handles the embedded Trivy process.
type Manager struct {
	BinaryPath string
}

func NewManager(binaryPath string) *Manager {
	return &Manager{
		BinaryPath: binaryPath,
	}
}

// TrivyReport is a simplified struct to capture Trivy's JSON output
type TrivyReport struct {
	Results []struct {
		Target          string `json:"Target"`
		Vulnerabilities []struct {
			VulnerabilityID  string   `json:"VulnerabilityID"`
			PkgName          string   `json:"PkgName"`
			InstalledVersion string   `json:"InstalledVersion"`
			FixedVersion     string   `json:"FixedVersion"`
			Title            string   `json:"Title"`
			Description      string   `json:"Description"`
			Severity         string   `json:"Severity"`
			References       []string `json:"References"`
		} `json:"Vulnerabilities"`
	} `json:"Results"`
}

// ScanFilesystem runs a vulnerability scan on the filesystem.
func (m *Manager) ScanFilesystem(path string) (*TrivyReport, error) {
	// trivy fs --format json --output - <path>
	// Note: In production you might want to stream this or write to a temp file
	cmd := exec.Command(m.BinaryPath, "fs", path, "--format", "json", "-q")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var report TrivyReport
	if err := json.Unmarshal(output, &report); err != nil {
		log.Printf("Failed to parse Trivy output: %v", err)
		return nil, err
	}

	return &report, nil
}
