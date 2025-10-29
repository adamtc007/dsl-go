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
		w("    (version %d)", req.Meta.Version)
		if !req.Meta.CreatedAt.IsZero() {
			w("\n    (created-at %q)", req.Meta.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"))
		}
		if !req.Meta.UpdatedAt.IsZero() {
			w("\n    (updated-at %q)", req.Meta.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"))
		}
		w(")\n")
	}
	// orchestrator
	if req.Orchestrator != nil {
		w("  (:orchestrator\n")
		if req.Orchestrator.Lifecycle != nil {
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
			w("      (transitions))\n")
		}

		// entities
		if len(req.Orchestrator.Entities) > 0 {
			w("    (:entities\n")
			for _, e := range req.Orchestrator.Entities {
				w("      (entity :id %q :type %s\n", e.ID, e.Typ)
				w("        (attrs\n")
				for _, attr := range e.Attrs {
					w("          (%s %s)\n", attr.Key, printValue(attr.Value))
				}
				w("        ))\n")
			}
			w("    )\n")
		}

		// resources
		if len(req.Orchestrator.Resources) > 0 {
			w("    (:resources\n")
			for _, r := range req.Orchestrator.Resources {
				w("      (resource :id %q :type %s)\n", r.ID, r.Typ)
			}
			w("    )\n")
		}

		// flows
		if len(req.Orchestrator.Flows) > 0 {
			w("    (:flows\n")
			for _, f := range req.Orchestrator.Flows {
				w("      (flow :id %q\n", f.ID)
				w("        (steps\n")
				for _, s := range f.Steps {
					if s.Task != nil {
						w("          (task :id %q :on %q :op %s)\n", s.Task.ID, s.Task.On, s.Task.Op)
					} else if s.Gate != nil {
						w("          (gate :id %q (when %q))\n", s.Gate.ID, s.Gate.Condition)
					}
				}
				w("        ))\n")
			}
			w("    ))\n")
		}
		w("  )\n")
	}

	w(")\n")
	return b.String()
}

func printValue(v *ast.Value) string {
	if v == nil {
		return ""
	}
	if v.String != nil {
		return fmt.Sprintf("%q", *v.String)
	} else if v.Int != nil {
		return fmt.Sprintf("%d", *v.Int)
	} else if v.Float != nil {
		return fmt.Sprintf("%g", *v.Float)
	} else if v.Bool != nil {
		return fmt.Sprintf("%t", *v.Bool)
	} else if v.Symbol != nil {
		return *v.Symbol
	}
	return ""
}
