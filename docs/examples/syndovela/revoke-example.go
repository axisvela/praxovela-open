// +build ignore

// Example: S6 Revoke Operation
// Demonstrates the /revoke endpoint with state machine transitions.
//
// This example shows:
// 1. Applying a bundle to ACTIVE state
// 2. Revoking the instance via /revoke
// 3. State transitions: ACTIVE → DRAINING → STOPPED
//
// Prerequisites:
// - PRAXOVELA v1.2.0+ (S6)
// - Syndovela facade running
// - Valid instance created via /apply

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/axisrobo/praxovela/adapters/syndovela"
	"github.com/axisrobo/praxovela/pkg/syndovela/contract"
)

// BundleBinding 结构
type BundleBinding struct {
	BundleID          string `json:"bundleId"`
	Version           string `json:"version"`
	Digest            string `json:"digest"`
	ResolutionLockRef string `json:"resolutionLockRef"`
	IsolationRequired string `json:"isolationRequired,omitempty"`
}

// SkillInvocation 结构 (仅作演示)
type SkillInvocation struct {
	InvocationID string `json:"invocationId"`
	Caller       string `json:"caller"`
	InstanceID   string `json:"instanceId"`
	SkillContract string `json:"skillContract"`
}

func main() {
	// 1. 创建 Syndovela facade with S6 state machine
	descriptor := syndovela.RuntimeDescriptor{
		RuntimeID:             "praxovela",
		Implementation:        "axon-core",
		ImplementationVersion: "1.2.0",
		ProtocolVersions:      []string{"sbrp/v1"},
		Isolation:             []string{"wasm"},
		ABIs:                []string{"wasi/preview2"},
		Platform:            "windows/amd64",
	}

	locks := []syndovela.ResolutionLock{
		{
			LockID:          "lock-1",
			ResolverVersion: "r1",
			Digest:          "sha256:manifest",
			Selected: []syndovela.SelectedBundle{
				{BundleID: "my-skill", Version: "1.0.0", Digest: "sha256:def456"},
			},
		},
	}

	facade := syndovela.NewRuntimeFacade(descriptor, locks)

	// 2. Apply bundle to create ACTIVE instance
	binding := BundleBinding{
		BundleID:          "my-skill",
		Version:           "1.0.0",
		Digest:            "sha256:def456",
		ResolutionLockRef: "lock-1",
	}

	body, _ := json.Marshal(binding)
	resp, err := http.Post(
		"http://localhost:8765/syndovela/apply",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Apply error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Apply failed: status %d\n", resp.StatusCode)
		os.Exit(1)
	}

	var activeInstance contract.BundleInstance
	if err := json.NewDecoder(resp.Body).Decode(&activeInstance); err != nil {
		fmt.Fprintf(os.Stderr, "Decode error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Bundle applied - Instance is ACTIVE\n")
	fmt.Printf("   Instance ID: %s\n", activeInstance.InstanceID)
	fmt.Printf("   State: %s\n", activeInstance.State)

	// 2. Revoke the instance
	revokeBody, _ := json.Marshal(map[string]string{
		"instanceId": activeInstance.InstanceID,
	})

	revokeResp, err := http.Post(
		"http://localhost:8765/syndovela/revoke",
		"application/json",
		bytes.NewReader(revokeBody),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Revoke error: %v\n", err)
		os.Exit(1)
	}
	defer revokeResp.Body.Close()

	if revokeResp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Revoke failed: status %d\n", revokeResp.StatusCode)
		os.Exit(1)
	}

	var revokedInstance contract.BundleInstance
	if err := json.NewDecoder(revokeResp.Body).Decode(&revokedInstance); err != nil {
		fmt.Fprintf(os.Stderr, "Decode error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Revoke successful!\n")
	fmt.Printf("   Instance ID: %s\n", revokedInstance.InstanceID)
	fmt.Printf("   New State: %s\n", revokedInstance.State)

	// State machine verification
	switch {
	case revokedInstance.State == syndovela.InstanceStopped:
		fmt.Println("✅ Instance fully stopped (DRAINING completed, no in-flight invocations)")
	case revokedInstance.State == syndovela.InstanceDraining:
		fmt.Println("⏳ Instance in DRAINING state (waiting for in-flight invocations to complete)")
	default:
		fmt.Printf("⚠️  Unexpected state: %s\n", revokedInstance.State)
	}
}

import "bytes"