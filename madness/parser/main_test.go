package main

import (
	"testing"
	"strings"
)

func TestNewParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantRisk float64
	}{
		{
			name:     "empty input",
			input:    "",
			wantRisk: 0,
		},
		{
			name:     "starts with sun",
			input:    "🌞 hello",
			wantRisk: 0,
		},
		{
			name:     "no sun start",
			input:    "hello world",
			wantRisk: 0.05,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotRisk := NewParser(tt.input)
			if gotRisk != tt.wantRisk {
				t.Errorf("NewParser() risk = %v, want %v", gotRisk, tt.wantRisk)
			}
		})
	}
}

func TestParser_parseNodeStart(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantNode string
	}{
		{
			name:     "simple node",
			input:    "kit {",
			wantNode: "kit",
		},
		{
			name:     "node with spaces",
			input:    "  spacey node  {",
			wantNode: "spacey node",
		},
		{
			name:     "emoji node",
			input:    "🔵 {",
			wantNode: "🔵",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{system: kit.New()}
			p.parseNodeStart(tt.input)
			if p.currentNode != tt.wantNode {
				t.Errorf("parseNodeStart() currentNode = %v, want %v", p.currentNode, tt.wantNode)
			}
		})
	}
}

func TestParser_parseProperty(t *testing.T) {
	tests := []struct {
		name       string
		currentNode string
		input      string
		wantRelation bool
	}{
		{
			name:        "is relation",
			currentNode: "kit",
			input:      "kit is good",
			wantRelation: true,
		},
		{
			name:        "has relation",
			currentNode: "kit",
			input:      "kit has power",
			wantRelation: true,
		},
		{
			name:        "from relation",
			currentNode: "kit",
			input:      "power from source",
			wantRelation: true,
		},
		{
			name:        "simple property",
			currentNode: "kit",
			input:      "be Kind",
			wantRelation: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				system:      kit.New(),
				currentNode: tt.currentNode,
			}
			p.parseProperty(tt.input)
			// Note: This is a basic test that just verifies the function runs
			// You may want to add more specific assertions based on your kit.System implementation
		})
	}
}

func TestParser_parseConversation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantSpeaker string
		wantMessage string
	}{
		{
			name:        "simple conversation",
			input:      "<Duane>: Hello world",
			wantSpeaker: "Duane",
			wantMessage: "Hello world",
		},
		{
			name:        "conversation with spaces",
			input:      "<kit >:  message with spaces  ",
			wantSpeaker: "kit",
			wantMessage: "message with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{system: kit.New()}
			p.parseConversation(tt.input)
			// Note: Add assertions here once you have a way to verify the conversations in your system
		})
	}
}

func TestParser_parseObservation(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple observation",
			input: "kit.observe(test)",
			want:  "test",
		},
		{
			name:  "observation with spaces",
			input: "kit.observe( complex test )",
			want:  "complex test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{system: kit.New()}
			p.parseObservation(tt.input)
			// Note: Add assertions here once you have a way to verify the observations in your system
		})
	}
}

func TestParser_ComplexInput(t *testing.T) {
	input := `
🔵 {
    rgb
    is matter
    has matter
    kit is inside
}

kit {
    be Kind
    Duane is good
    kind is good
}

<Duane>: Hello world
kit.observe(test observation)
`
	system, risk := NewParser(input)
	if risk != 0.05 {
		t.Errorf("Complex input risk = %v, want 0.05", risk)
	}
}