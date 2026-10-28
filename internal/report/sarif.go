package report

import (
	"encoding/json"
	"io"

	"github.com/arjunuxd/uxdguard/internal/scan"
)

// SARIFReport represents a minimal SARIF 2.1.0 report
type SARIFReport struct {
	Schema  string `json:"$schema"`
	Version string `json:"version"`
	Runs    []Run  `json:"runs"`
}

type Run struct {
	Tool    Tool     `json:"tool"`
	Results []Result `json:"results"`
}

type Tool struct {
	Driver Driver `json:"driver"`
}

type Driver struct {
	Name            string `json:"name"`
	InformationURI  string `json:"informationUri"`
	SemanticVersion string `json:"semanticVersion"`
	Rules           []Rule `json:"rules"`
}

type Rule struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	ShortDescription Message     `json:"shortDescription"`
	FullDescription  Message     `json:"fullDescription"`
	DefaultConfig    DefaultConf `json:"defaultConfiguration"`
	Properties       Props       `json:"properties"`
}

type Message struct {
	Text string `json:"text"`
}

type DefaultConf struct {
	Level string `json:"level"`
}

type Props struct {
	Tags []string `json:"tags"`
}

type Result struct {
	RuleID    string  `json:"ruleId"`
	Level     string  `json:"level"`
	Message   Message `json:"message"`
	Locations []Loc   `json:"locations"`
	Fixes     []Fix   `json:"fixes,omitempty"`
}

type Loc struct {
	PhysicalLocation PhysLoc `json:"physicalLocation"`
}

type PhysLoc struct {
	ArtifactLocation ArtLoc `json:"artifactLocation"`
	Region           Region `json:"region"`
}

type ArtLoc struct {
	URI string `json:"uri"`
}

type Region struct {
	StartLine int `json:"startLine"`
}

type Fix struct {
	Description Message `json:"description"`
}

// SARIFReporter writes SARIF output
type SARIFReporter struct {
	Out io.Writer
}

// Report writes SARIF report
func (s *SARIFReporter) Report(res scan.RunResult, root string, version string) error {
	rep := SARIFReport{
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Version: "2.1.0",
		Runs: []Run{{
			Tool: Tool{Driver: Driver{
				Name:            "uxdguard",
				InformationURI:  "https://github.com/arjunuxd/uxdguard",
				SemanticVersion: version,
				Rules:           []Rule{},
			}},
			Results: []Result{},
		}},
	}

	levelMap := map[scan.Severity]string{
		scan.SeverityCritical: "error",
		scan.SeverityHigh:     "error",
		scan.SeverityMedium:   "warning",
		scan.SeverityLow:      "note",
		scan.SeverityInfo:     "note",
	}

	for _, f := range res.Findings {
		level := "warning"
		if l, ok := levelMap[f.Severity]; ok {
			level = l
		}
		msg := f.Message
		if msg == "" {
			msg = f.Title
		}
		resItem := Result{
			RuleID:  f.RuleID,
			Level:   level,
			Message: Message{Text: msg},
			Locations: []Loc{{
				PhysicalLocation: PhysLoc{
					ArtifactLocation: ArtLoc{URI: f.File},
					Region:           Region{StartLine: f.Line},
				},
			}},
		}
		if f.Remediation != "" {
			resItem.Fixes = []Fix{{Description: Message{Text: f.Remediation}}}
		}
		rep.Runs[0].Results = append(rep.Runs[0].Results, resItem)
	}

	enc := json.NewEncoder(s.Out)
	enc.SetIndent("", "  ")
	return enc.Encode(rep)
}
