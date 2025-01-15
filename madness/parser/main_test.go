package main

import (
    "testing"
)

func TestParse(t *testing.T) {
    testcases := []struct {
    	name string
        in string
        wantRiskBelow float64
    }{
        {"Any day", "🌞, wake, bake, walk dog, breakfast, work, sleep, 🌝", 0},
        {"Most days", "🌝, wake, 🌞,  bake, walk dog, breakfast, work, 🌝, sleep", 0.1},
    }
    for _, tc := range testcases {
    	risk := NewParser(tc.in)
        if risk > tc.wantRiskBelow {
                t.Errorf("NewParser: %f, want %f", risk, tc.wantRiskBelow)
        }
    }
}