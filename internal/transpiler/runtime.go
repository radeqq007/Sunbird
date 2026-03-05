package transpiler

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed runtime/shared.ts
var sharedRuntimeSrc string

//go:embed runtime/node.ts
var nodeRuntimeSrc string

//go:embed runtime/deno.ts
var denoRuntimeSrc string

//go:embed runtime/bun.ts
var bunRuntimeSrc string

//go:embed runtime/web.ts
var webRuntimeSrc string

type Target string

const (
	TargetNode Target = "node"
	TargetDeno Target = "deno"
	TargetBun  Target = "bun"
	TargetWeb  Target = "web"
)

const DefaultTarget = TargetNode

func targetSrc(t Target) string {
	switch t {
	case TargetDeno:
		return denoRuntimeSrc
	case TargetBun:
		return bunRuntimeSrc
	case TargetWeb:
		return webRuntimeSrc
	default:
		return nodeRuntimeSrc
	}
}

func ParseTarget(s string) (Target, error) {
	switch Target(s) {
	case TargetNode, TargetDeno, TargetBun, TargetWeb:
		return Target(s), nil
	default:
		return "", fmt.Errorf(
			"unknown target %q — valid targets are: node, deno, bun, web", s,
		)
	}
}

func WriteRuntime(outputDir string, target Target) error {
	dest := filepath.Join(outputDir, "sunbird-rt.ts")
	if _, err := os.Stat(dest); err == nil {
		return nil // already written this run
	}

	combined := sharedRuntimeSrc +
		"\n// --- " + string(target) + " target ---\n\n" +
		targetSrc(target)

	if err := os.WriteFile(dest, []byte(combined), 0o644); err != nil {
		return fmt.Errorf("writing sunbird-rt.ts: %w", err)
	}

	return nil
}

