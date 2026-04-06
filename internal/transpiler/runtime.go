package transpiler

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/radeqq007/sunbird/internal/runtime"
)

func WriteRuntime(outputDir string) error {
	dest := filepath.Join(outputDir, "sunbird-rt.ts")
	if _, err := os.Stat(dest); err == nil {
		return nil // already written this run
	}

	if err := os.WriteFile(dest, []byte(runtime.Runtime), 0o644); err != nil {
		return fmt.Errorf("writing sunbird-rt.ts: %w", err)
	}

	return nil
}
