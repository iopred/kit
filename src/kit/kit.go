package kit

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Node struct {
	X, Y, Z, T, Dx, Dy, Dz, Dt float64 // t = time, d = delta
	Value                      float64 // value is essentially radius, anything that points at this
	Dvalue                     float64
	Name                       string // Name is the name of the system, this should be unique at the top level, system local variables don't really mean something in kit, variables have a location, so system variables are just inside the system itself, like a big circle diagram
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
	return n.X == n2.X && n.Y == n2.Y && n.Z == n2.Z && n.Value == n2.Value && n.T == n2.T
}

func (n Node) Length() float64 {
	return math.Sqrt(math.Pow(n.X-n.Dx, 2) + math.Pow(n.Y-n.Dy, 2) + math.Pow(n.Z-n.Dz, 2) + math.Pow(n.T-n.Tz, 2))
}

func (n Node) Sub(n2 Node) Node {
	n.X -= n2.X
	n.Dx -= n2.Dx
	n.Y -= n2.Y
	n.Dy -= n2.Dy
	n.Z -= n2.Z
	n.Dz -= n2.Dz
	n.T -= n2.T
	n.Dt -= n2.Dt
	return n
}

func (n Node) Distance(n2 Node) float64 {
	return n.Sub(n2).Length()
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

func (k *Kit) Node(name string) (Node, error) {
	for _, node := range k.nodes {
		if node.Name == name {
			return node, nil
		}
	}
	return Node{}, errors.New("Could not find node")
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
	AddNode(node Node) Node
	Node(name string) (Node, error)
	Next(scalar float64) System
	String() string
}

func New() System {
	return &Kit{}
}
