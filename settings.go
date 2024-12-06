package qawaylinter

import (
	"go/types"
	"strings"
)

// Settings is the root configuration object for the linter.
type Settings struct {
	Targets []Rules `json:"rules"`
}

// Rules defines rules and to which packages they apply
// Filters allow users to customize to which nodes a rule should apply to.
// For example, interfaces in the domain package may require comments, but interfaces in an internal dev package may not.
//
// The object is divided into individual rule attributes instead of one generic `rules` object.
// This makes JSON deserialization of the configuration easier.
// In addition, it works around limitations in Generics support in Go.
type Rules struct {
	Packages []string `json:"packages"`

	FunctionRule  *FunctionRule[FunctionRuleResults]   `json:"functions"`
	InterfaceRule *InterfaceRule[InterfaceRuleResults] `json:"interfaces"`
	StructRule    *StructRule[StructRuleResults]       `json:"structs"`
}

// MatchesPackage checks if the given package matches the target.
// Returns true if the full package path starts with any of the target packages.
// Also returns the package that matched.
// For example, if the target is `["example.com/foo"]`, the package `example.com/foo/bar` will match.
func (t Rules) MatchesPackage(pkg *types.Package) (bool, string) {
	for _, p := range t.Packages {
		if strings.HasPrefix(pkg.Path(), p) {
			return true, p
		}
	}
	return false, ""
}

// GetMatchingTarget finds the most specific target that matches the given package.
func (s Settings) GetMatchingTarget(pkg *types.Package) *Rules {
	var matchingTargets = make(map[*Rules]string)

	for _, t := range s.Targets {
		if matches, p := t.MatchesPackage(pkg); matches {
			matchingTargets[&t] = p
		}
	}

	return findMostConcreteTarget(matchingTargets)
}

// findMostConcreteTarget finds the most specific target from a list of matching targets.
// A target is more specific if the package that matches the given node are more specific than the
// matching package from the other node.
// A subpackage is more specific than a package.
func findMostConcreteTarget(matchingTargets map[*Rules]string) *Rules {
	var mostConcreteTarget *Rules

	for target, pkg := range matchingTargets {
		if mostConcreteTarget == nil {
			mostConcreteTarget = target
			continue
		}

		if len(pkg) > len(matchingTargets[mostConcreteTarget]) {
			mostConcreteTarget = target
		}
	}

	return mostConcreteTarget
}
