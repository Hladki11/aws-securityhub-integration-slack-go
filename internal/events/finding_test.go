package events

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSecurityHubV2FindingParsing(t *testing.T) {
	path := filepath.Join("..", "..", "fixtures", "samples.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read samples: %v", err)
	}

	var findings []json.RawMessage
	if err := json.Unmarshal(raw, &findings); err != nil {
		t.Fatalf("failed to unmarshal samples: %v", err)
	}

	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}

	// test first finding (GuardDuty)
	f1, err := NewSecurityHubFinding(findings[0])
	if err != nil {
		t.Fatalf("failed to parse finding 1: %v", err)
	}

	if f1.Metadata.Product.Name != "GuardDuty" {
		t.Errorf("expected GuardDuty, got %s", f1.Metadata.Product.Name)
	}
	if f1.Severity != "Medium" {
		t.Errorf("expected Medium severity, got %s", f1.Severity)
	}
	if f1.FindingInfo.Title != "A container has executed a newly created binary file." {
		t.Errorf("unexpected title: %s", f1.FindingInfo.Title)
	}
	if !f1.IsAlertable() {
		t.Error("GuardDuty finding should be alertable")
	}

	// test second finding (Security Hub CSPM)
	f2, err := NewSecurityHubFinding(findings[1])
	if err != nil {
		t.Fatalf("failed to parse finding 2: %v", err)
	}

	if f2.Metadata.Product.Name != "Security Hub" {
		t.Errorf("expected Security Hub, got %s", f2.Metadata.Product.Name)
	}
	if f2.Severity != "Critical" {
		t.Errorf("expected Critical severity, got %s", f2.Severity)
	}
	if f2.Compliance == nil {
		t.Error("expected compliance data")
	}
	if !f2.IsAlertable() {
		t.Error("Failed compliance finding should be alertable")
	}
}
