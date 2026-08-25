package types

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestHeuristicsYAMLUnmarshaling(t *testing.T) {
	// 1. Mock YAML data that matches your struct tags
	mockYAML := []byte(`
named_patterns:
  obj_c_pattern:
    - "^\\s*@interface"
    - "^\\s*#import"
disambiguations:
  - extensions:
      - "m"
      - "h"
    rules:
      - language:
          - "Objective-C"
        pattern:
          - "^\\s*@interface"
        negative_pattern:
          - "^:- module"
        named_pattern: "obj_c_pattern"
        and:
          - language:
              - "C"
            pattern:
              - "#include <stdio.h>"
`)

	// 2. Unmarshal into your struct
	var h Heuristics
	err := yaml.Unmarshal(mockYAML, &h)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// 3. Verify NamedPatterns mapped correctly
	expectedPatterns, ok := h.NamedPatterns["obj_c_pattern"]
	if !ok {
		t.Fatalf("Expected 'obj_c_pattern' to exist in NamedPatterns map")
	}
	if len(expectedPatterns) != 2 || expectedPatterns[0] != "^\\s*@interface" {
		t.Errorf("NamedPatterns parsed incorrectly. Got: %v", expectedPatterns)
	}

	// 4. Verify Disambiguations mapped correctly
	if len(h.Disambiguations) != 1 {
		t.Fatalf("Expected 1 Disambiguation, got %d", len(h.Disambiguations))
	}

	disambig := h.Disambiguations[0]
	expectedExts := []string{"m", "h"}
	if !reflect.DeepEqual(disambig.Extensions, expectedExts) {
		t.Errorf("Extensions = %v, want %v", disambig.Extensions, expectedExts)
	}

	// 5. Verify Rules mapped correctly
	if len(disambig.Rules) != 1 {
		t.Fatalf("Expected 1 Rule, got %d", len(disambig.Rules))
	}

	rule := disambig.Rules[0]
	if rule.Language[0] != "Objective-C" {
		t.Errorf("Rule Language = %v, want [Objective-C]", rule.Language)
	}
	if rule.Pattern[0] != "^\\s*@interface" {
		t.Errorf("Rule Pattern = %v, want [^\\s*@interface]", rule.Pattern)
	}
	if rule.NegativePattern[0] != "^:- module" {
		t.Errorf("Rule NegativePattern = %v, want [^:- module]", rule.NegativePattern)
	}
	if rule.NamedPattern != "obj_c_pattern" {
		t.Errorf("Rule NamedPattern = %s, want obj_c_pattern", rule.NamedPattern)
	}

	// 6. Verify nested 'And' rules mapped correctly
	if len(rule.And) != 1 {
		t.Fatalf("Expected 1 nested 'And' rule, got %d", len(rule.And))
	}
	if rule.And[0].Language[0] != "C" {
		t.Errorf("Nested And Rule Language = %v, want [C]", rule.And[0].Language)
	}
}

func TestCompiledRuleAndHeuristicRuleAssignment(t *testing.T) {
	// A simple sanity check that the internal structs can be assigned properly
	// Since they lack YAML tags, they are likely populated internally by your engine.
	cRule := CompiledRule{
		Language: "Go",
		Patterns: []string{"^package main"},
		Negative: []string{"^package notmain"},
	}

	if cRule.Language != "Go" {
		t.Errorf("CompiledRule Assignment failed")
	}

	hRule := HeuristicRule{
		Language: "Rust",
		Patterns: nil, // Would normally be []*regexp.Regexp
		Negative: nil,
	}

	if hRule.Language != "Rust" {
		t.Errorf("HeuristicRule Assignment failed")
	}
}
