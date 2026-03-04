package transpiler

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed runtime/runtime.ts
var runtimeSrc string

func WriteRuntime(outputDir string) error {
	dest := filepath.Join(outputDir, "sunbird-rt.ts")

	// Don't overwrite if already there from this run
	if _, err := os.Stat(dest); err == nil {
		return nil
	}

	return os.WriteFile(dest, []byte(runtimeSrc), 0o644)
}

