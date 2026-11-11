package ota

import (
	"sort"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Phase is where one vehicle is in an update campaign.
type Phase string

const (
	PhasePending     Phase = "pending"
	PhaseDownloading Phase = "downloading"
	PhaseInstalling  Phase = "installing"
	PhaseComplete    Phase = "complete"
	PhaseFailed      Phase = "failed"
)

var phaseOrder = []Phase{PhasePending, PhaseDownloading, PhaseInstalling, PhaseComplete, PhaseFailed}

// Campaign is one over-the-air software release and the vehicles it targets.
type Campaign struct {
	ID       string
	Release  string
	Region   string
	Vehicles map[string]Phase // VIN to phase
}

// PhaseCount is the number of vehicles in one phase.
type PhaseCount struct {
	Phase string `json:"phase"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

// Summary is the dashboard view of a campaign.
type Summary struct {
	ID              string       `json:"id"`
	Release         string       `json:"release"`
	Region          string       `json:"region"`
	Total           int          `json:"total"`
	CompletePercent int          `json:"complete_percent"`
	Halted          bool         `json:"halted"`
	Phases          []PhaseCount `json:"phases"`
	FailedVINs      []string     `json:"failed_vins"`
}

// haltThresholdPercent is the share of failed vehicles at which a campaign is
// halted for review.
const haltThresholdPercent = 5

var titleCaser = cases.Title(language.English)

// Summarize builds the dashboard view of a campaign.
func Summarize(c Campaign) Summary {
	counts := make(map[Phase]int, len(phaseOrder))
	failed := []string{}
	for vin, phase := range c.Vehicles {
		counts[phase]++
		if phase == PhaseFailed {
			failed = append(failed, vin)
		}
	}
	sort.Strings(failed)

	total := len(c.Vehicles)
	s := Summary{ID: c.ID, Release: c.Release, Region: c.Region, Total: total, FailedVINs: failed}
	for _, p := range phaseOrder {
		s.Phases = append(s.Phases, PhaseCount{Phase: string(p), Label: titleCaser.String(string(p)), Count: counts[p]})
	}
	if total > 0 {
		s.CompletePercent = counts[PhaseComplete] * 100 / total
		s.Halted = counts[PhaseFailed]*100 >= haltThresholdPercent*total
	}
	return s
}

// SeedCampaigns returns invented campaigns. The VINs are not real.
func SeedCampaigns() []Campaign {
	return []Campaign{
		{
			ID: "cmp-2026-11-a", Release: "2.6.2", Region: "us-west",
			Vehicles: map[string]Phase{
				"5LDEMAA0000000001": PhaseComplete,
				"5LDEMAA0000000002": PhaseComplete,
				"5LDEMAA0000000003": PhaseInstalling,
				"5LDEMAA0000000004": PhaseDownloading,
				"5LDEMAA0000000005": PhasePending,
				"5LDEMAA0000000006": PhaseComplete,
			},
		},
		{
			ID: "cmp-2026-11-b", Release: "2.6.2", Region: "ca-central",
			Vehicles: map[string]Phase{
				"5LDEMAA0000000101": PhaseComplete,
				"5LDEMAA0000000102": PhaseFailed,
				"5LDEMAA0000000103": PhaseInstalling,
				"5LDEMAA0000000104": PhasePending,
			},
		},
	}
}
