package planner

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Gabrielvilabracho/selops-ai/internal/model"
)

// TestTopologicalSortOrdersDependenciesFirst verifies that topological sort
// correctly orders components with their dependencies first.
// OPS fork (Phase 0e): ComponentSDD removed; uses ComponentSDDOps with its dep on Engram.
func TestTopologicalSortOrdersDependenciesFirst(t *testing.T) {
	deps := map[model.ComponentID][]model.ComponentID{
		model.ComponentSkills:          nil,
		model.ComponentSDDOps:          {model.ComponentEngram},
		model.ComponentEngram:          nil,
		model.ComponentPersona:         nil,
		model.ComponentContext7:         nil,
	}

	ordered, err := TopologicalSort(deps)
	if err != nil {
		t.Fatalf("TopologicalSort() returned error: %v", err)
	}

	// Engram must appear before SDDOps (it's a dependency).
	engramIdx, sddOpsIdx := -1, -1
	for i, c := range ordered {
		if c == model.ComponentEngram {
			engramIdx = i
		}
		if c == model.ComponentSDDOps {
			sddOpsIdx = i
		}
	}
	if engramIdx < 0 || sddOpsIdx < 0 {
		t.Fatalf("missing components in result: %v", ordered)
	}
	if engramIdx > sddOpsIdx {
		t.Fatalf("Engram (%d) must be before SDDOps (%d) in sorted order, got %v", engramIdx, sddOpsIdx, ordered)
	}
}

func TestApplySoftOrderingReordersWithoutAddingDependencies(t *testing.T) {
	ordered := []model.ComponentID{
		model.ComponentContext7,
		model.ComponentEngram,
		model.ComponentPersona,
		model.ComponentSDDOps,
	}

	result := applySoftOrdering(ordered, [][2]model.ComponentID{{model.ComponentPersona, model.ComponentEngram}})

	if !reflect.DeepEqual(result, []model.ComponentID{
		model.ComponentContext7,
		model.ComponentPersona,
		model.ComponentEngram,
		model.ComponentSDDOps,
	}) {
		t.Fatalf("applySoftOrdering() = %v", result)
	}

	// If the first component is absent, nothing should be added.
	result = applySoftOrdering([]model.ComponentID{model.ComponentEngram}, [][2]model.ComponentID{{model.ComponentPersona, model.ComponentEngram}})
	if !reflect.DeepEqual(result, []model.ComponentID{model.ComponentEngram}) {
		t.Fatalf("applySoftOrdering() should not add missing components (first absent), got %v", result)
	}
}

func TestApplySoftOrderingEdgeCases(t *testing.T) {
	pair := [][2]model.ComponentID{{model.ComponentPersona, model.ComponentEngram}}

	// Second absent — no-op, no panic
	result := applySoftOrdering([]model.ComponentID{model.ComponentPersona}, pair)
	if !reflect.DeepEqual(result, []model.ComponentID{model.ComponentPersona}) {
		t.Fatalf("second absent: expected [persona], got %v", result)
	}

	// Both absent — no-op (use ComponentSkills as neutral component)
	result = applySoftOrdering([]model.ComponentID{model.ComponentSkills}, pair)
	if !reflect.DeepEqual(result, []model.ComponentID{model.ComponentSkills}) {
		t.Fatalf("both absent: expected [skills], got %v", result)
	}

	// Already correct order — no-op (must not mutate)
	already := []model.ComponentID{model.ComponentPersona, model.ComponentEngram}
	result = applySoftOrdering(already, pair)
	if !reflect.DeepEqual(result, []model.ComponentID{model.ComponentPersona, model.ComponentEngram}) {
		t.Fatalf("already correct: expected [persona, engram], got %v", result)
	}

	// Input slice must NOT be mutated
	input := []model.ComponentID{model.ComponentEngram, model.ComponentPersona}
	_ = applySoftOrdering(input, pair)
	if !reflect.DeepEqual(input, []model.ComponentID{model.ComponentEngram, model.ComponentPersona}) {
		t.Fatalf("input slice was mutated")
	}
}

// TestApplySoftOrderingBothMVPPairsWithFullOPSSelection verifies soft ordering
// for OPS preset: PersonaOperator must be before SDDOps.
// OPS fork (Phase 0e): ComponentSDD removed; uses OPS components.
func TestApplySoftOrderingBothMVPPairsWithFullOPSSelection(t *testing.T) {
	// Simulates the OPS preset scenario.
	ordered := []model.ComponentID{
		model.ComponentContext7,
		model.ComponentEngram,
		model.ComponentSDDOps,
		model.ComponentPersonaOperator,
	}

	result := applySoftOrdering(ordered, SoftOrderingConstraints())

	// PersonaOperator must appear before SDDOps.
	personaIdx, sddOpsIdx := -1, -1
	for i, c := range result {
		switch c {
		case model.ComponentPersonaOperator:
			personaIdx = i
		case model.ComponentSDDOps:
			sddOpsIdx = i
		}
	}

	if personaIdx < 0 || sddOpsIdx < 0 {
		t.Fatalf("missing components in result: %v", result)
	}
	if personaIdx > sddOpsIdx {
		t.Fatalf("PersonaOperator (%d) must be before SDDOps (%d), got %v", personaIdx, sddOpsIdx, result)
	}
}

func TestSoftOrderingContainsPersonaOperatorBeforeSDDOps(t *testing.T) {
	pairs := SoftOrderingConstraints()

	found := false
	for _, pair := range pairs {
		if pair[0] == model.ComponentPersonaOperator && pair[1] == model.ComponentSDDOps {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("SoftOrderingConstraints() does not contain {ComponentPersonaOperator, ComponentSDDOps} pair, got: %v", pairs)
	}
}

func TestTopologicalSortDetectsCycles(t *testing.T) {
	deps := map[model.ComponentID][]model.ComponentID{
		model.ComponentEngram: {model.ComponentSDDOps},
		model.ComponentSDDOps: {model.ComponentEngram},
	}

	_, err := TopologicalSort(deps)
	if err == nil {
		t.Fatalf("TopologicalSort() expected cycle error")
	}

	if !errors.Is(err, ErrDependencyCycle) {
		t.Fatalf("TopologicalSort() error = %v, want ErrDependencyCycle", err)
	}
}
