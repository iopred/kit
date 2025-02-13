package main

import (
	"github.com/dominikbraun/graph"
	"github.com/iopred/kit/kit"

	"fmt"
	"strings"
	"unicode"
)

type Parser struct {
	input     string
	pos       int
	system    kit.System
	riskCost  float64
	currentNode string
	bracketDepth int
}

func NewParser(input string) (kit.System, float64) {
	p := &Parser{
		input:    input,
		system:   kit.New(),
		riskCost: 0,
	}

	// Initial risk assessment
	if strings.Index(input, "🌞") != 0 {
		p.riskCost = 0.05
	}

	p.parse()
	return p.system, p.riskCost
}

func (p *Parser) parse() {
	lines := strings.Split(p.input, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// Handle conversation markers
		if strings.HasPrefix(line, "<") && strings.Contains(line, ">:") {
			p.parseConversation(line)
			continue
		}

		// Handle kit.observe() calls
		if strings.HasPrefix(line, "kit.observe") {
			p.parseObservation(line)
			continue
		}

		// Handle node definitions
		if strings.Contains(line, "{") {
			p.parseNodeStart(line)
			continue
		}

		if strings.Contains(line, "}") {
			p.bracketDepth--
			p.currentNode = ""
			continue
		}

		// Handle properties within nodes
		if p.currentNode != "" {
			p.parseProperty(line)
		}
	}
}

func (p *Parser) parseNodeStart(line string) {
	nodeName := strings.TrimSpace(strings.Split(line, "{")[0])
	p.currentNode = nodeName
	p.bracketDepth++

	// Add node to the system
	p.system.AddNode(nodeName)
}

func (p *Parser) parseProperty(line string) {
	// Handle different property types
	if strings.Contains(line, "is") {
		parts := strings.Split(line, "is")
		if len(parts) == 2 {
			subject := strings.TrimSpace(parts[0])
			object := strings.TrimSpace(parts[1])
			p.system.AddRelation(p.currentNode, "is", object)
		}
	} else if strings.Contains(line, "has") {
		parts := strings.Split(line, "has")
		if len(parts) == 2 {
			subject := strings.TrimSpace(parts[0])
			object := strings.TrimSpace(parts[1])
			p.system.AddRelation(p.currentNode, "has", object)
		}
	} else if strings.Contains(line, "from") {
		// Handle 'from' relationships
		parts := strings.Split(line, "from")
		if len(parts) == 2 {
			source := strings.TrimSpace(parts[1])
			p.system.AddRelation(p.currentNode, "from", source)
		}
	} else {
		// Handle simple properties or time coordinates
		p.system.AddProperty(p.currentNode, line)
	}
}

func (p *Parser) parseConversation(line string) {
	// Extract speaker and message
	parts := strings.SplitN(line, ":", 2)
	if len(parts) == 2 {
		speaker := strings.Trim(parts[0], "<>")
		message := strings.TrimSpace(parts[1])
		p.system.AddConversation(speaker, message)
	}
}

func (p *Parser) parseObservation(line string) {
	// Extract observation parameters
	start := strings.Index(line, "(")
	end := strings.Index(line, ")")
	if start != -1 && end != -1 {
		params := strings.TrimSpace(line[start+1:end])
		p.system.AddObservation(params)
	}
}

func main() {
	// Example usage
	input := `🌝, wake, 🌞, bake, walk dog, breakfast, work, 🌝, sleep`
	system, risk := NewParser(input)
	fmt.Printf("Risk Cost: %f\n", risk)
	
	// Add more test cases here
	complexInput := `
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
}`
	system, risk = NewParser(complexInput)
	fmt.Printf("Complex Parse Risk Cost: %f\n", risk)
}
