package qawaylinter

import (
	"go/types"
	"testing"
)

func TestGetMatchingTarget(t *testing.T) {
	settings := Settings{
		Targets: []Rules{
			{Packages: []string{"example.com/foo"}},
			{Packages: []string{"example.com/foo/bar"}},
			{Packages: []string{"example.com/foo/bar/baz"}},
		},
	}

	pkg := types.NewPackage("example.com/foo/bar/baz/qux", "qux")
	matchingTarget := settings.GetMatchingTarget(pkg)
	expected := "example.com/foo/bar/baz"

	if matchingTarget == nil || matchingTarget.Packages[0] != expected {
		t.Errorf("Expected most concrete target to be %s, but got %v", expected, matchingTarget)
	}
}

func TestFindMostConcreteTarget(t *testing.T) {
	targets := []Rules{
		{Packages: []string{"example.com/foo"}},
		{Packages: []string{"example.com/foo/bar"}},
		{Packages: []string{"example.com/foo/bar/baz"}},
	}

	matchingTargets := make(map[*Rules]string)
	for i := range targets {
		matchingTargets[&targets[i]] = targets[i].Packages[0]
	}

	mostConcreteTarget := findMostConcreteTarget(matchingTargets)
	expected := "example.com/foo/bar/baz"

	if mostConcreteTarget == nil || mostConcreteTarget.Packages[0] != expected {
		t.Errorf("Expected most concrete target to be %s, but got %v", expected, mostConcreteTarget)
	}
}
