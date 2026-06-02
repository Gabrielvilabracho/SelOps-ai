package planner

import (
	"reflect"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/model"
)

func TestResolverAddsMissingDependenciesInOrder(t *testing.T) {
	resolver := NewResolver(MVPGraph())

	selection := model.Selection{
		Agents:     []model.AgentID{model.AgentClaudeCode, model.AgentOpenCode},
		Components: []model.ComponentID{model.ComponentSkills},
	}

	plan, err := resolver.Resolve(selection)
	if err != nil {
		t.Fatalf("Resolve() returned error: %v", err)
	}

	if !reflect.DeepEqual(plan.Agents, []model.AgentID{model.AgentClaudeCode, model.AgentOpenCode}) {
		t.Fatalf("Resolve() agents = %v", plan.Agents)
	}

	if !reflect.DeepEqual(plan.OrderedComponents, []model.ComponentID{model.ComponentEngram, model.ComponentSDD, model.ComponentSkills}) {
		t.Fatalf("Resolve() ordered components = %v", plan.OrderedComponents)
	}

	if !reflect.DeepEqual(plan.AddedDependencies, []model.ComponentID{model.ComponentEngram, model.ComponentSDD}) {
		t.Fatalf("Resolve() added dependencies = %v", plan.AddedDependencies)
	}
}

func TestResolverPersonaOrderedBeforeEngramAndSDDWhenSelected(t *testing.T) {
	resolver := NewResolver(MVPGraph())

	selection := model.Selection{
		Components: []model.ComponentID{model.ComponentPersona, model.ComponentSDD},
	}

	plan, err := resolver.Resolve(selection)
	if err != nil {
		t.Fatalf("Resolve() returned error: %v", err)
	}

	if !reflect.DeepEqual(plan.OrderedComponents, []model.ComponentID{model.ComponentPersona, model.ComponentEngram, model.ComponentSDD}) {
		t.Fatalf("Resolve() ordered components = %v", plan.OrderedComponents)
	}

	if !reflect.DeepEqual(plan.AddedDependencies, []model.ComponentID{model.ComponentEngram}) {
		t.Fatalf("Resolve() added dependencies = %v", plan.AddedDependencies)
	}
}

func TestResolverEngramOnlyDoesNotForcePersona(t *testing.T) {
	resolver := NewResolver(MVPGraph())

	selection := model.Selection{
		Components: []model.ComponentID{model.ComponentEngram},
	}

	plan, err := resolver.Resolve(selection)
	if err != nil {
		t.Fatalf("Resolve() returned error: %v", err)
	}

	if !reflect.DeepEqual(plan.OrderedComponents, []model.ComponentID{model.ComponentEngram}) {
		t.Fatalf("Resolve() ordered components = %v", plan.OrderedComponents)
	}

	if len(plan.AddedDependencies) != 0 {
		t.Fatalf("Resolve() added dependencies = %v, want none", plan.AddedDependencies)
	}
}

func TestResolverSDDOnlyDoesNotForcePersona(t *testing.T) {
	resolver := NewResolver(MVPGraph())

	selection := model.Selection{
		Components: []model.ComponentID{model.ComponentSDD},
	}

	plan, err := resolver.Resolve(selection)
	if err != nil {
		t.Fatalf("Resolve() returned error: %v", err)
	}

	for _, dep := range plan.AddedDependencies {
		if dep == model.ComponentPersona {
			t.Fatalf("SDD-only selection should NOT force Persona, got AddedDependencies=%v", plan.AddedDependencies)
		}
	}
}

func TestResolverPersonaAndEngramWithoutSDD(t *testing.T) {
	resolver := NewResolver(MVPGraph())

	selection := model.Selection{
		Components: []model.ComponentID{model.ComponentPersona, model.ComponentEngram},
	}

	plan, err := resolver.Resolve(selection)
	if err != nil {
		t.Fatalf("Resolve() returned error: %v", err)
	}

	if !reflect.DeepEqual(plan.OrderedComponents, []model.ComponentID{model.ComponentPersona, model.ComponentEngram}) {
		t.Fatalf("Resolve() ordered components = %v, want [persona, engram]", plan.OrderedComponents)
	}

	if len(plan.AddedDependencies) != 0 {
		t.Fatalf("Resolve() added dependencies = %v, want none", plan.AddedDependencies)
	}
}

func TestResolverPresetSelOpsOperationalResolvesToFourComponents(t *testing.T) {
	resolver := NewResolver(MVPGraph())

	// PresetSelOpsOperational includes ComponentPersonaOperator, ComponentSDDOps,
	// ComponentOperationalMCP, and their transitive dependency ComponentEngram.
	selection := model.Selection{
		Components: []model.ComponentID{
			model.ComponentPersonaOperator,
			model.ComponentSDDOps,
			model.ComponentOperationalMCP,
		},
	}

	plan, err := resolver.Resolve(selection)
	if err != nil {
		t.Fatalf("Resolve() returned error: %v", err)
	}

	// ComponentEngram must be auto-added as a dependency of ComponentSDDOps.
	if len(plan.OrderedComponents) != 4 {
		t.Fatalf("expected 4 ordered components, got %d: %v", len(plan.OrderedComponents), plan.OrderedComponents)
	}

	// Verify all four are present.
	componentSet := make(map[model.ComponentID]struct{}, len(plan.OrderedComponents))
	for _, c := range plan.OrderedComponents {
		componentSet[c] = struct{}{}
	}
	for _, expected := range []model.ComponentID{
		model.ComponentPersonaOperator,
		model.ComponentSDDOps,
		model.ComponentOperationalMCP,
		model.ComponentEngram,
	} {
		if _, ok := componentSet[expected]; !ok {
			t.Fatalf("expected component %q in plan, got %v", expected, plan.OrderedComponents)
		}
	}

	// ComponentPersonaOperator must be ordered before ComponentSDDOps (soft constraint).
	personaIdx, sddOpsIdx := -1, -1
	for i, c := range plan.OrderedComponents {
		switch c {
		case model.ComponentPersonaOperator:
			personaIdx = i
		case model.ComponentSDDOps:
			sddOpsIdx = i
		}
	}
	if personaIdx < 0 || sddOpsIdx < 0 {
		t.Fatalf("ComponentPersonaOperator or ComponentSDDOps not found in plan: %v", plan.OrderedComponents)
	}
	if personaIdx > sddOpsIdx {
		t.Fatalf("ComponentPersonaOperator (%d) must be before ComponentSDDOps (%d), got %v", personaIdx, sddOpsIdx, plan.OrderedComponents)
	}
}

func TestResolverExcludesUnsupportedAgents(t *testing.T) {
	resolver := NewResolver(MVPGraph())

	selection := model.Selection{
		Agents: []model.AgentID{model.AgentClaudeCode, model.AgentCursor, model.AgentID("unknown-agent")},
	}

	plan, err := resolver.Resolve(selection)
	if err != nil {
		t.Fatalf("Resolve() returned error: %v", err)
	}

	if !reflect.DeepEqual(plan.Agents, []model.AgentID{model.AgentClaudeCode, model.AgentCursor}) {
		t.Fatalf("Resolve() agents = %v", plan.Agents)
	}

	if !reflect.DeepEqual(plan.UnsupportedAgents, []model.AgentID{model.AgentID("unknown-agent")}) {
		t.Fatalf("Resolve() unsupported agents = %v", plan.UnsupportedAgents)
	}
}
