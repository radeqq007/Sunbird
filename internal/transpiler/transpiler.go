package transpiler

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"sunbird/internal/ast"
)

type Transpiler struct {
	buf bytes.Buffer
	indent int
	imports map[string]string // alias => import name
}

func New() *Transpiler {
	return &Transpiler{
		imports: make(map[string]string),
	}
}

func (t *Transpiler) Transpile(program *ast.Program) (string, error){
	var body bytes.Buffer
	for _, stmt := range program.Statements {
		s, err := t.transpileStatement(stmt)
		if err != nil {
			return "", err
		}
		body.WriteString(s)
		body.WriteByte('\n')
	}

	var out bytes.Buffer
	out.WriteString(t.buildImports())
	out.WriteString(body.String())

	return out.String(), nil
}

func (t *Transpiler) buildImports() string {
	if len(t.imports) == 0 {
  	return ""
  }

  seen := map[string]bool{}
  parts := []string{}
  for alias, mod := range t.imports {
  	if seen[mod] {
      continue
    }
    seen[mod] = true
    if alias == mod {
    	parts = append(parts, mod)
    } else {
      parts = append(parts, fmt.Sprintf("%s as %s", mod, alias))
		}
	}

	sort.Strings(parts)
	return fmt.Sprintf("import { %s } from \"./sunbird-rt.js\";\n\n", strings.Join(parts, ", "))
}

func (t *Transpiler) indentStr() string {
	return strings.Repeat("  ", t.indent)
}

func (t *Transpiler) pushIndent() { t.indent++ }
func (t *Transpiler) popIndent()  { t.indent-- }
