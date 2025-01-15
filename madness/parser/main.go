package main

import (
	"github.com/dominikbraun/graph"
	"github.com/iopred/kit/kit"

	"fmt"
	"strings"
)

func NewParser(input string) (system kit.System, riskCost float64) {
	fmt.Println(input)
	if strings.Index(input, "🌞") != 0 {
		riskCost = 0.05	
	}

	system = kit.New()
	
	return system, riskCost
}


func main() {
	
	risk, system := NewParser("🌝, wake, 🌞,  bake, walk dog, breakfast, work, 🌝, sleep")
	fmt.Println(risk)


}
