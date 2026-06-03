package planner

import (
	"reflect"
	"testing"

	"github.com/Gabrielvilabracho/selops-ai/internal/model"
)

// TestResolverAddsMissingDependenciesInOrder verifies that Resolve() adds
// missing transitive dependencies. OPS fork (Phase 0e): ComponentSDD removed;
// uses ComponentSDDOps (dep: Engram) to test the same resolver behavior.
func TestResolverAddsMissingDependenciesInOrder(t *testing.T) {
	resolver := NewResolver(MVPGraph())

	selection := model.Selection{
		Agents:     []model.AgentID{model.AgentClaudeCode, model.AgentOpenCode},
		Components: []model.ComponentID{model.ComponentSDDOps},
	}

	plan, err := resolver.Resolve(selection)
	if err != nil {
		t.Fatalf("Resolve() returned error: %v", err)
	}

	if !reflect.DeepEqual(plan.Agents, []model.AgentID{model.AgentClaudeCode, model.AgentOpenCode}) {
		t.Fatalf("Resolve() agents = %v", plan.Agents)
	}

	// ComponentSDDOps has dep on ComponentEngram — Engram must be auto-added.
	if !reflect.DeepEqual(plan.OrderedComponents, []model.ComponentID{model.ComponentEngram, model.ComponentSDDOps}) {
		t.Fatalf("Resolve() ordered components = %v", plan.OrderedComponents)
	}

	if !reflect.DeepEqual(plan.AddedDependencies, []model.ComponentID{model.ComponentEngram}) {
		t.Fatalf("Resolve() added dependencies = %v", plan.AddedDependencies)
	}
}

// TestResolverPersonaOperatorOrderedBeforeSDDOpsWhenSelected verifies soft ordering.
// OPS fork (Phase 0e): ComponentSDD removed; uses PersonaOperator+SDDOps pair.
func TestResolverPersonaOperatorOrderedBeforeSDDOpsWhenSelected(t *testing.T) {
	resolver := NewResolver(MVPGraph())

	selection := model.Selection{
		Components: []model.ComponentID{model.ComponentPersonaOperator, model.ComponentSDDOps},
	}

	plan, err := resolver.Resolve(selection)
	if err != nil {
		t.Fatalf("Resolve() returned error: %v", err)
	}

	// PersonaOperator must appear before SDDOps due to soft ordering.
	// Engram is auto-added as SDDOps dependency.
	personaIdx, sddOpsIdx := -1, -1
	for i, c := range plan.OrderedComponents {
		if c == model.ComponentPersonaOperator {
			personaIdx = i
		}
		if c == model.ComponentSDDOps {
			sddOpsIdx = i
		}
	}
	if personaIdx < 0 || sddOpsIdx < 0 {
		t.Fatalf("Resolve() ordered components = %v, missing PersonaOperator or SDDOps", plan.OrderedComponents)
	}
	if personaIdx > sddOpsIdx {
		t.Fatalf("PersonaOperator (%d) must be before SDDOps (%d), got %v", personaIdx, sddOpsIdx, plan.OrderedComponents)
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

// TestResolverSDDOpsOnlyDoesNotForcePersonaOperator verifies that selecting
// only ComponentSDDOps does not force ComponentPersonaOperator.
// OPS fork (Phase 0e): ComponentSDD removed; equivalent test with SDDOps.
func TestResolverSDDOpsOnlyDoesNotForcePersonaOperator(t *testing.T) {
	resolver := NewResolver(MVPGraph())

	selection := model.Selection{
		Components: []model.ComponentID{model.ComponentSDDOps},
	}

	plan, err := resolver.Resolve(selection)
	if err != nil {
		t.Fatalf("Resolve() returned error: %v", err)
	}

	for _, dep := range plan.AddedDependencies {
		if dep == model.ComponentPersonaOperator {
			t.Fatalf("SDDOps-only selection should NOT force PersonaOperator, got AddedDependencies=%v", plan.AddedDependencies)
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

// TestKnowledgeBaseHasNilDepsInMVPGraph verifies that ComponentKnowledgeBase
// has nil (no) dependencies in MVPGraph — it must NOT transitively pull in
// ComponentSDD, ComponentSkills, or any other DEV component.
func TestKnowledgeBaseHasNilDepsInMVPGraph(t *testing.T) {
	g := MVPGraph()

	// ComponentKnowledgeBase must be registered in the graph.
	if !g.Has(model.ComponentKnowledgeBase) {
		t.Fatalf("MVPGraph must contain ComponentKnowledgeBase")
	}

	// It must have nil (zero) dependencies — no transitive DEV pull.
	deps := g.DependenciesOf(model.ComponentKnowledgeBase)
	if len(deps) != 0 {
		t.Fatalf("ComponentKnowledgeBase must have no deps in MVPGraph (got %v)", deps)
	}
}

// TestKnowledgeBaseDoesNotTransitivelyPullDEVComponents verifies that selecting
// ComponentKnowledgeBase alone does not auto-add any DEV components.
func TestKnowledgeBaseDoesNotTransitivelyPullDEVComponents(t *testing.T) {
	resolver := NewResolver(MVPGraph())

	selection := model.Selection{
		Components: []model.ComponentID{model.ComponentKnowledgeBase},
	}

	plan, err := resolver.Resolve(selection)
	if err != nil {
		t.Fatalf("Resolve() returned error: %v", err)
	}

	// AddedDependencies must be empty — no transitive pull.
	if len(plan.AddedDependencies) != 0 {
		t.Fatalf("ComponentKnowledgeBase must not transitively pull any deps; got %v", plan.AddedDependencies)
	}

	// OPS fork (Phase 0e): ComponentSDD removed from devComponents check.
	devComponents := []model.ComponentID{
		model.ComponentSkills,
		model.ComponentPermission,
		model.ComponentGGA,
	}

	for _, dev := range devComponents {
		for _, comp := range plan.OrderedComponents {
			if comp == dev {
				t.Errorf("ComponentKnowledgeBase must not pull DEV component %q; ordered = %v", dev, plan.OrderedComponents)
			}
		}
	}
}
