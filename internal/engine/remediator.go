package engine

import (
	"log"
	"strings"

	"nexus/internal/config"
	"nexus/internal/modules/patcher"
	"nexus/pkg/ocsf"
)

// Remediator decides if a finding should be patched.
type Remediator struct {
	Config  *config.RemediationConfig
	Patcher *patcher.Manager
}

func NewRemediator(cfg *config.RemediationConfig, patcher *patcher.Manager) *Remediator {
	return &Remediator{
		Config:  cfg,
		Patcher: patcher,
	}
}

// Evaluate checks if the finding meets the criteria for auto-remediation.
func (r *Remediator) Evaluate(finding ocsf.VulnerabilityFinding) bool {
	if !r.Config.Enabled {
		return false
	}

	// Check Severity
	severityMatch := false
	for _, s := range r.Config.AutoPatchSeverity {
		if strings.EqualFold(s, finding.Vulnerability.Severity) {
			severityMatch = true
			break
		}
	}

	if !severityMatch {
		return false
	}

	// TODO: Check if fix is available (needs OCSF enrichment)

	return true
}

// Remediate attempts to fix the vulnerability.
func (r *Remediator) Remediate(finding ocsf.VulnerabilityFinding) {
	log.Printf("[REMEDIATION] Attempting to fix %s (%s)", finding.Vulnerability.Title, finding.Vulnerability.CVE)

	// Infer package ID from title/description (Loose matching for prototype)
	// In production, Trivy gives PkgName which needs mapping to Winget ID
	pkgID := finding.Vulnerability.Title

	// Attempt patch
	if err := r.Patcher.InstallUpdate(pkgID); err != nil {
		log.Printf("[REMEDIATION FAILED] %v", err)
	} else {
		log.Printf("[REMEDIATION SUCCESS] Patch initiated for %s", pkgID)
	}
}
