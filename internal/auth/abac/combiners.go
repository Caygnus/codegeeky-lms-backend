package abac

import (
	"sort"
)

// AllMustAllowCombiner requires all policies to allow access
type AllMustAllowCombiner struct{}

// NewAllMustAllowCombiner creates a new AllMustAllowCombiner
func NewAllMustAllowCombiner() PolicyCombiner {
	return &AllMustAllowCombiner{}
}

// Combine combines decisions where all must allow for access to be granted
func (c *AllMustAllowCombiner) Combine(decisions []Decision) Decision {
	if len(decisions) == 0 {
		return Decision{
			Allow:  true,
			Reason: "No policies to evaluate",
		}
	}

	// Check if all policies allow access
	for _, decision := range decisions {
		if !decision.Allow {
			return Decision{
				Allow:    false,
				Reason:   decision.Reason,
				Priority: decision.Priority,
				Metadata: decision.Metadata,
			}
		}
	}

	// All policies allowed
	return Decision{
		Allow:  true,
		Reason: "All policies allowed access",
	}
}

// GetName returns the combiner name
func (c *AllMustAllowCombiner) GetName() string {
	return "AllMustAllow"
}

// AnyCanAllowCombiner allows access if any policy allows
type AnyCanAllowCombiner struct{}

// NewAnyCanAllowCombiner creates a new AnyCanAllowCombiner
func NewAnyCanAllowCombiner() PolicyCombiner {
	return &AnyCanAllowCombiner{}
}

// Combine combines decisions where any allow grants access
func (c *AnyCanAllowCombiner) Combine(decisions []Decision) Decision {
	if len(decisions) == 0 {
		return Decision{
			Allow:  true,
			Reason: "No policies to evaluate",
		}
	}

	// Check if any policy allows access
	for _, decision := range decisions {
		if decision.Allow {
			return Decision{
				Allow:    true,
				Reason:   decision.Reason,
				Priority: decision.Priority,
				Metadata: decision.Metadata,
			}
		}
	}

	// All policies denied
	return Decision{
		Allow:  false,
		Reason: "All policies denied access",
	}
}

// GetName returns the combiner name
func (c *AnyCanAllowCombiner) GetName() string {
	return "AnyCanAllow"
}

// PriorityBasedCombiner uses policy priority to determine final decision
type PriorityBasedCombiner struct{}

// NewPriorityBasedCombiner creates a new PriorityBasedCombiner
func NewPriorityBasedCombiner() PolicyCombiner {
	return &PriorityBasedCombiner{}
}

// Combine combines decisions based on priority (highest priority wins)
func (c *PriorityBasedCombiner) Combine(decisions []Decision) Decision {
	if len(decisions) == 0 {
		return Decision{
			Allow:  true,
			Reason: "No policies to evaluate",
		}
	}

	// Sort decisions by priority (highest first)
	sort.Slice(decisions, func(i, j int) bool {
		return decisions[i].Priority > decisions[j].Priority
	})

	// Return the highest priority decision
	highestPriority := decisions[0]
	return Decision{
		Allow:    highestPriority.Allow,
		Reason:   highestPriority.Reason,
		Priority: highestPriority.Priority,
		Metadata: highestPriority.Metadata,
	}
}

// GetName returns the combiner name
func (c *PriorityBasedCombiner) GetName() string {
	return "PriorityBased"
}

// MajorityWinsCombiner uses majority vote to determine access
type MajorityWinsCombiner struct{}

// NewMajorityWinsCombiner creates a new MajorityWinsCombiner
func NewMajorityWinsCombiner() PolicyCombiner {
	return &MajorityWinsCombiner{}
}

// Combine combines decisions based on majority vote
func (c *MajorityWinsCombiner) Combine(decisions []Decision) Decision {
	if len(decisions) == 0 {
		return Decision{
			Allow:  true,
			Reason: "No policies to evaluate",
		}
	}

	allowCount := 0
	denyCount := 0

	for _, decision := range decisions {
		if decision.Allow {
			allowCount++
		} else {
			denyCount++
		}
	}

	if allowCount > denyCount {
		return Decision{
			Allow:  true,
			Reason: "Majority of policies allowed access",
		}
	} else if denyCount > allowCount {
		return Decision{
			Allow:  false,
			Reason: "Majority of policies denied access",
		}
	} else {
		// Tie - default to deny for security
		return Decision{
			Allow:  false,
			Reason: "Equal number of allow/deny decisions - defaulting to deny",
		}
	}
}

// GetName returns the combiner name
func (c *MajorityWinsCombiner) GetName() string {
	return "MajorityWins"
}
