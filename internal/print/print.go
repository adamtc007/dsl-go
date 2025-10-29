package print

import (
	"fmt"
	"strings"

	"github.com/example/dsl-go/internal/ast"
)

func ToSexpr(req *ast.Request) string {
	var b strings.Builder
	w := func(s string, args ...interface{}) { fmt.Fprintf(&b, s, args...) }
	w("(onboarding-request\n")
	// meta
	if req.Meta != nil {
		w("  (:meta\n")
		w("    (request-id %q)\n", req.Meta.RequestID)
		w("    (version %d)\n", req.Meta.Version)
		if !req.Meta.CreatedAt.IsZero() {
			w("    (created-at %q)\n", req.Meta.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"))
		}
		if !req.Meta.UpdatedAt.IsZero() {
			w("    (updated-at %q)\n", req.Meta.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"))
		}
		w("  )\n")
	}
	// orchestrator
	w("  (:orchestrator\n")
	w("    (:lifecycle\n")
	w("      (states")
	if len(req.Orchestrator.Lifecycle.States) == 0 {
		w(" draft validated compiled executing completed failed")
	} else {
		for _, st := range req.Orchestrator.Lifecycle.States {
			w(" %s", st)
		}
	}
	w(")\n")
	if req.Orchestrator.Lifecycle.Initial == "" {
		w("      (initial draft)\n")
	} else {
		w("      (initial %s)\n", req.Orchestrator.Lifecycle.Initial)
	}
	w("      (transitions)\n") // TODO
	w("    )\n")

	// entities
	w("    (:entities")
	if len(req.Orchestrator.Entities) == 0 {
		w(")\n")
	} else {
		w("\n")
		for _, e := range req.Orchestrator.Entities {
			w("      (entity :id %q :type %s\n", e.ID, e.Typ)
			w("        (attrs\n")
			for k, v := range e.Attrs {
				w("          (%s %v)\n", k, printValue(v.Value)) // Basic print
			}
			w("        ))\n")
		}
		w("    )\n")
	}

	// resources
	w("    (:resources")
	if len(req.Orchestrator.Resources) == 0 {
		w(")\n")
	} else {
		w("\n")
		for _, r := range req.Orchestrator.Resources {
			w("      (resource :id %q :type %s)\n", r.ID, r.Typ) // Basic print
		}
		w("    )\n")
	}

	// flows
	w("    (:flows")
	if len(req.Orchestrator.Flows) == 0 {
		w(")\n")
	} else {
		w("\n")
		for _, f := range req.Orchestrator.Flows {
			w("      (flow :id %q\n", f.ID)
			w("        (steps\n")
			for _, s := range f.Steps {
				switch s.Kind {
				case "task":
					w("          (task :id %q :on %q :op %s)\n", s.Task.ID, s.Task.On, s.Task.Op) // Basic
				case "gate":
					w("          (gate :id %q (when %q))\n", s.Gate.ID, s.Gate.Condition)
				}
			}
			w("        ))\n")
		}
		w("    )\n")
	}

	w("  )\n")
	w(")\n")
	return b.String()
}

func printValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case uint64, int, float64:
		return fmt.Sprintf("%d", val)
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
