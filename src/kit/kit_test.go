package kit

import "testing"

func systemMatches(t *testing.T, expected, got System) {
	t.Helper()
	kit := expected.(*Kit)
	kat := got.(*Kit)
	for i, node := range kit.nodes {
		if !node.Equal(kat.nodes[i]) {
			t.Errorf("Non matching nodes %v and %v", node.String(), kat.nodes[i].String())
		}
	}
}

func valueMatches(t *testing.T, expected, got float64) {
	t.Helper()
	if got != expected {
		t.Errorf("Non matching value %v. expected: %v", got, expected)
	}
}

func TestNew(t *testing.T) {
	system := New()
	system.AddNode(Node{
		X:  0,
		Dx: 2,
		Dt: 1,
	})
	system.AddNode(Node{
		X:  1,
		Dx: -2,
		Dt: -1,
	})
	got := system.Nodes()
	if got != 2 {
		t.Errorf("system.Nodes() = %d; want 2", got)
	}
}

func TestValue(t *testing.T) {
	system := New()
	system.AddNode(Node{
		Dt:     1,
		Dvalue: 1,
		Name:   "test",
	})
	type testcase struct {
		Scalar float64
		Value  float64
	}
	cases := []testcase{
		{
			Scalar: 0,
			Value:  0,
		},
		{
			Scalar: 0.5,
			Value:  0.5,
		},
		{
			Scalar: 1,
			Value:  1,
		},
	}

	for _, c := range cases {
		n := system.Next(c.Scalar).Node("test")
		valueMatches(t, c.Value, n.Value)
	}
}

func TestNext(t *testing.T) {
	system := New()
	system.AddNode(Node{
		Dt:    1,
		Value: 1,
	})
	system.AddNode(Node{
		Dt:    1,
		Value: 1,
	})
	systemMatches(t, system, system.Next(0))
}

func TestKit(t *testing.T) {
	kit := New()
	kit.AddNode(Node{
		X:     0,
		Dx:    1,
		T:     1,
		Dt:    0.5,
		Value: 1,
	})
	kit.AddNode(Node{
		X:     0,
		Dx:    1,
		T:     1,
		Dt:    0.5,
		Value: 1,
	})

}
