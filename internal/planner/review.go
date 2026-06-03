package planner

import "github.com/Gabrielvilabracho/selops-ai/internal/model"

func BuildReviewPayload(selection model.Selection, resolved ResolvedPlan) ReviewPayload {
	autoAdded := make(map[model.ComponentID]struct{}, len(resolved.AddedDependencies))
	for _, component := range resolved.AddedDependencies {
		autoAdded[component] = struct{}{}
	}

	components := make([]ComponentAction, 0, len(resolved.OrderedComponents))
	for _, component := range resolved.OrderedComponents {
		action := "selected"
		if _, ok := autoAdded[component]; ok {
			action = "auto-dependency"
		}
		components = append(components, ComponentAction{ID: component, Action: action})
	}

	return ReviewPayload{
		Agents:            resolved.Agents,
		UnsupportedAgents: resolved.UnsupportedAgents,
		Persona:           selection.Persona,
		Preset:            selection.Preset,
		Components:        components,
		AddedDependencies: resolved.AddedDependencies,
		PlatformDecision:  resolved.PlatformDecision,
		// Issue #145: pass skills from selection.
		Skills: selection.Skills,
		// Issue #149: pass StrictTDD. HasSDD always false in OPS fork (Phase 0e).
		StrictTDD: selection.StrictTDD,
		HasSDD:    false,
	}
}
