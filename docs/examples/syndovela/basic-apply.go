// +build ignore

// Example: Basic Syndovela Apply Operation
// Demonstrates the /apply endpoint integration with Syndovela facade.
//
// This example shows how to:
// 1. Create a Syndovela RuntimeFacade
// 2. Apply a bundle binding via HTTP
// 3. Handle the response and state transition
//
// Prerequisites:
// - PRAXOVELA v1.3.0+
// - Syndovela facade running on http://localhost:8765/syndovela
// - Valid BundleBinding with resolution lock reference

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/axisrobo/praxovela/adapters/syndovela"
	"github.com/axisrobo/praxovela/pkg/syndovela/contract"
)

// BundleBinding 定义用于 /apply 的绑定结构
type BundleBinding struct {
	BundleID          string `json:"bundleId"`
	Version           string `json:"version"`
	Digest            string `json:"digest"`
	ResolutionLockRef string `json:"resolutionLockRef"`
	IsolationRequired string `json:"isolationRequired,omitempty"`
	ABI               string `json:"abi,omitempty"`
	PolicyRefs        []string `json:"policyRefs,omitempty"`
}

// 主函数 - 运行基本 apply 操作演示
func main() {
	// 1. 创建 Syndovela facade
	descriptor := syndovela.RuntimeDescriptor{
		RuntimeID:             "praxovela",
		Implementation:        "axon-core",
		ImplementationVersion: "1.3.0",
		ProtocolVersions:      []string{"sbrp/v1"},
		Isolation:             []string{"wasm"},
		ABIs:                  []string{"wasi/preview2"},
		Platform:              "windows/amd64",
		Features:              []string{},
		Limits:                map[string]string{"memory": "16Gi"},
	}

	// 创建锁 (模拟已 materialized 的 resolution lock)
	locks := []syndovela.ResolutionLock{
		{
			LockID:          "lock-1",
			ResolverVersion: "r1",
			Digest:          "sha256:manifest",
			Selected: []syndovela.SelectedBundle{
				{BundleID: "my-skill", Version: "2.0.0", Digest: "sha256:abc123"},
			},
		},
	}

	facade := syndovela.NewRuntimeFacade(descriptor, locks)

	// 2. 构造 BundleBinding 请求
	binding := BundleBinding{
		BundleID:          "my-skill",
		Version:           "2.0.0",
		Digest:            "sha256:abc123",
		ResolutionLockRef: "lock-1",
		IsolationRequired: "wasm",
	}

	// 3. 发送 POST /apply 请求
	body, _ := json.Marshal(binding)
	resp, err := http.Post(
		"http://localhost:8765/syndovela/apply",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// 4. 处理响应
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Apply failed: status %d\n", resp.StatusCode)
		os.Exit(1)
	}

	// 5. 解析并显示响应
	var instance contract.BundleInstance
	if err := json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding response: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Apply successful!\n")
	fmt.Printf("  Instance ID: %s\n", instance.InstanceID)
	fmt.Printf("  Bundle ID: %s\n", instance.BundleID)
	fmt.Printf("  Version: %s\n", instance.Version)
	fmt.Printf("  State: %s\n", instance.State)
	fmt.Printf("  Health: %s\n", instance.Health)
	fmt.Printf("  Active Invocations: %d\n", instance.ActiveInvocations)

	if instance.State == syndovela.InstanceActive {
		fmt.Println("✅ Instance is ACTIVE - ready for skill invocations")
	} else {
		fmt.Fprintf(os.Stderr, "⚠️  Instance state: %s (expected ACTIVE)\n", instance.State)
	}
}

// 辅助函数 - 导入 bytes
import "bytes"