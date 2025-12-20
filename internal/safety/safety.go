package safety

import (
	"fmt"
	"time"
)

// Check is a single safety condition.
type Check interface {
	Name() string
	Check() error
}

// Chain is a collection of checks.
type Chain struct {
	checks []Check
}

func NewChain(checks ...Check) *Chain {
	return &Chain{checks: checks}
}

// Execute runs all checks. Returns error on the first failure.
func (c *Chain) Execute() error {
	for _, check := range c.checks {
		if err := check.Check(); err != nil {
			return fmt.Errorf("safety check failed [%s]: %v", check.Name(), err)
		}
	}
	return nil
}

// --- Specific Checks ---

type TimeWindowCheck struct {
	AllowedStartHour int
	AllowedEndHour   int
}

func (t *TimeWindowCheck) Name() string { return "MaintenanceTimestamp" }
func (t *TimeWindowCheck) Check() error {
	now := time.Now().Hour()
	if now < t.AllowedStartHour || now > t.AllowedEndHour {
		return fmt.Errorf("current time %d is outside allowed window [%d-%d]", now, t.AllowedStartHour, t.AllowedEndHour)
	}
	return nil
}

type CPUCheck struct {
	MaxUsage float64
}

func (c *CPUCheck) Name() string { return "CPUUtilization" }
func (c *CPUCheck) Check() error {
	// TODO: Implement real WMI/Performance Counter check.
	// For now, simulate a low load.
	currentLoad := 10.0 // Mocked
	if currentLoad > c.MaxUsage {
		return fmt.Errorf("cpu usage %.2f%% exceeds limit %.2f%%", currentLoad, c.MaxUsage)
	}
	return nil
}
