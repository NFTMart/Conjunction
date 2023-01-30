package main

import (
	"testing"
)

func TestLoadNodes(t *testing.T) {
	newNodes := LoadNodes("nodes_test.son")
	if len(newNodes) != 1 {
		t.Errorf("len(newNodes) = %d; want 1", len(newNodes))
	}
	firstNode := newNodes[0]
	if firstNode.Name != "US Node 1" {
		t.Errorf("Node Name Wrong")
	}
	if firstNode.Address != "http://[2001:db8::1]:5000/" {
		t.Errorf("Node Address Wrong")
	}
	if firstNode.Type != FULL {
		t.Errorf("Node Type Wrong")
	}
	if firstNode.Active != true {
		t.Errorf("Node Active Wrong")
	}
	if firstNode.Features[LiveState] != true {
		t.Errorf("Node LiveState invalid")
	}
	if firstNode.Features[FullTransactionHistory] != true {
		t.Errorf("Node FullTransactionHistory invalid")
	}
	if firstNode.Features[AccountHistory] != false {
		t.Errorf("Node AccountHistory invalid")
	}
	if firstNode.Features[MarketHistory] != false {
		t.Errorf("Node MarketHistory invalid")
	}
	if firstNode.Features[NftHistory] != false {
		t.Errorf("Node NftHistory invalid")
	}
}

func TestHealthCheck(t *testing.T) {
	newNodes := LoadNodes("nodes.json")
	if len(newNodes) != 1 {
		t.Errorf("Error loading nodes")
	}
	firstNode := newNodes[0]
	healthCheckResult := firstNode.HealthCheck()
	if healthCheckResult != false {
		t.Errorf("Health check result invalid")
	}
	if firstNode.Active != false {
		t.Errorf("Node Active Wrong")
	}
}
