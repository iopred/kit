package kit

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Noder interface {
	Node() Node
}

// Nodes exist in spacetime.
type Node struct {
	X, Y, Z, T, Dx, Dy, Dz, Dt float64 // t = time, d = delta
	Value                      float64 // value is essentially radius, this node will provide value of this magnitude
	Dvalue                     float64 // the delta of value, positive values mean this node is growing itself, it's useful to model this value as the value of mutually beneficial collaboration
	Name                       string // Name is the name of the system, this should be unique at the top level, system local variables don't really mean something in kit, variables have a location, so system variables are just inside the system itself, like a big circle diagram. In fact Name is a filter, which means 'give me all the nodes within me'
	Dimension                  Noder  // Dimension  is the place in spacetime to execute this system.
}

func (n Node) Node() Node {
	return n
}

func (n Node) Next(scalar float64) Node {
	n.X += n.Dx * n.Dt * scalar
	n.Y += n.Dy * n.Dt * scalar
	n.Z += n.Dz * n.Dt * scalar
	n.Value += n.Dvalue * n.Dt * scalar

	n.T += n.Dt * scalar // How long have I been alive?
	return n
}

func (n Node) Equal(n2 Node) bool {
	return n.X == n2.X && n.Y == n2.Y && n.Z == n2.Z && n.Value == n2.Value
}

func (n Node) Speed() float64 {
	return math.Sqrt(math.Pow(n.X-n.Dx, 2) + math.Pow(n.Y-n.Dy, 2) + math.Pow(n.Z-n.Dz, 2))
}

func (n Node) Length() float64 {
	return math.Sqrt(math.Pow(n.X, 2) + math.Pow(n.Y, 2) + math.Pow(n.Z, 2))
}

func (n Node) Add(n2 Node) Node {
	n.X += n2.X
	n.Dx += n2.Dx
	n.Y += n2.Y
	n.Dy += n2.Dy
	n.Z += n2.Z
	n.Dz += n2.Dz
	n.Value += n2.Value
	n.Dvalue += n2.Dvalue
	return n
}

func (n Node) Sub(n2 Node) Node {
	n.X -= n2.X
	n.Dx -= n2.Dx
	n.Y -= n2.Y
	n.Dy -= n2.Dy
	n.Z -= n2.Z
	n.Dz -= n2.Dz
	n.Value -= n2.Value
	n.Dvalue -= n2.Dvalue
	return n
}

func (n Node) Mul(n2 Node) Node {
	n.X *= n2.X
	n.Dx *= n2.Dx
	n.Y *= n2.Y
	n.Dy *= n2.Dy
	n.Z *= n2.Z
	n.Dz *= n2.Dz
	n.Value *= n2.Value
	n.Dvalue *= n2.Dvalue
	return n
}

func (n Node) Distance(n2 Node) float64 {
	return n.Sub(n2).Length()
}

func (n Node) ExistsWith(n2 Node) bool {
	return n.T+n.Dt >= n2.T && n.T <= n2.T+n2.Dt
}

func (n Node) String() string {
	return fmt.Sprintf("[%s <%f, %f, %f> time(%f) value(%f) length(%f)]", n.Name, n.X, n.Y, n.Z, n.T, n.Value, n.Length())
}

type Kit struct {
	nodes []Node
}

func (k *Kit) Nodes() int {
	return len(k.nodes)
}

func (k *Kit) AddNode(node Node) Node {
	k.nodes = append(k.nodes, node)
	return node
}

// Creates or gets an empty node with the specified name and spacetime, or returns the closest in spacetime.
func (k *Kit) Node(name string) Node {
	for _, node := range k.nodes {
		if node.Name == name {
			return node
		}
	}

	return Node{
		Name: name,
		T:    0,
	}
}

// This method is the special sauce.
// Given a node, return a system only containing the nodes that overlap `node` in space and time.
func (k *Kit) FilterToSpacetime(node Node) System {
	system := New()
	for _, n := range k.nodes {

		// Exists in the same time.
		if n.ExistsWith(node) {
			if n.Distance(node) <= n.Value+node.Value {
				system.AddNode(n)
			}
		}
	}
	return system
}

func (k *Kit) Next(scalar float64) System {
	system := New()
	for _, n := range k.nodes {
		next := n.Next(scalar)
		if next.T >= 0 {
			system.AddNode(next)
		}
	}
	return system
}

// Resolves the value for a node at a specific spacetime.
// Many nodes can have incongruent values so the algorithm tries to reduce error by inserting nodes into the simulation that minimize incongruity.
func (k *Kit) Resolve(name string, spacetime float64) float64 {
	n := k.Node(name)
	n.T = spacetime
	return k.FilterToSpacetime(n).Node(name).Value
}

func (k *Kit) String() string {
	s := "kit["
	for _, n := range k.nodes {
		s += fmt.Sprintf(n.String())
	}
	s += "]kat"
	return s
}

type System interface {
	Nodes() int
	AddNode(node Node) Node // Add a Node into the system, it will be placed at a position in time and space, 
	Node(name string) Node
	Resolve(name string, spacetime float64) float64
	Next(scalar float64) System
	FilterToSpacetime(node Node) System
	String() string
}

func New() System {
	return &Kit{}
}
