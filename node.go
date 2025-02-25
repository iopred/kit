package main

import (
	"fmt"
	"strings"
)


type Node struct {
	X       bool
	Y       bool
	Z       bool
	Gravity bool
	Nodes   rune
}

func (n Node) x() string {
	if n.X {
		return fmt.Sprintf("Node{X:%v, Y:%v, Z:%v, Gravity:%v, Nodes:%d} intersects with spacetime vector (x, 0, 0, 0)", n.X, n.Y, n.Z, n.Gravity, n.Nodes)
	}
	return "0"
}

func (n Node) y() string {
	if n.Y {
		return fmt.Sprintf("Node{X:%v, Y:%v, Z:%v, Gravity:%v, Nodes:%d} intersects with spacetime vector (0, y, 0, 0)", n.X, n.Y, n.Z, n.Gravity, n.Nodes)
	}
	return "0"
}

func (n Node) z() string {
	if n.Z {
		return fmt.Sprintf("Node{X:%v, Y:%v, Z:%v, Gravity:%v, Nodes:%d} intersects with spacetime vector (0, 0, z, 0)", n.X, n.Y, n.Z, n.Gravity, n.Nodes)
	}
	return "0"
}

func (n Node) g() string {
	if n.Gravity {
		return fmt.Sprintf("Node{X:%v, Y:%v, Z:%v, Gravity:%v, Nodes:%d} intersects with spacetime vector (0, 0, 0, g)", n.X, n.Y, n.Z, n.Gravity, n.Nodes)
	}
	return "0"
}

func (n Node) kat() Node {
	if n.Gravity {
		return Node{n.X, true, n.Z, n.Gravity, n.Nodes}
	}
	return Node{n.X, n.Y, n.Z, n.Gravity, n.Nodes}
}

func (n Node) now() string {
	return fmt.Sprintf("%03d9", n.Nodes)
}

func (n Node) next() string {
	switch n.now() {
	case "0000":
		return "0001"
	case "0001":
		return "0002"
	case "0002":
		return "0003"
	case "0003":
		return "0004"
	case "0004":
		return "0005"
	case "0005":
		return "0006"
	case "0006":
		return "0007"
	case "0007":
		return "0008"
	case "0008":
		return "0009"
	case "0009":
		return "🌞"
	}

	return string(n.Nodes)
}

func (n Node) renderFrame() string {
	// Create a recursive tree structure using emojis
	// Root emoji is the current time/state
	root := n.now()

	// Build branches using different emojis based on node properties
	var branches []string

	// Add branches based on node properties
	if n.X {
		branches = append(branches, "→🔴") // Red for X
	}
	if n.Y {
		branches = append(branches, "↑🟡") // Yellow for Y
	}
	if n.Z {
		branches = append(branches, "↗️🔵") // Blue for Z
	}
	if n.Gravity {
		branches = append(branches, "↓⚫") // Black for Gravity
	}

	// Add sub-nodes based on Nodes byte value

	// todo: append sub nodes to string following the path.
	// for i := byte(0); i < n.Nodes; i++ {
	//     branches = append(branches, "⤵️"+string(rune('0'+i)))
	// }

	// Combine root with branches
	result := root
	if len(branches) > 0 {
		result += " " + strings.Join(branches, " ")
	}

	// Add next state
	// result += " → " + n.kat().next()

	return result
}

// Frame: 0009 →🔴 ↗️🔵 ↓⚫ ⤵️0 ⤵️1 → 🌞

/*
func (n Node) renderFrame() THREE.Group {
	// Create a group to hold the nodes
	group := THREE.NewGroup()

	// Create a box for each node
	box := THREE.NewBoxGeometry(1, 1, 1)

	// Create a material for the nodes
	material := THREE.NewMeshBasicMaterial(THREE.Color(0x00ff00))

	// Create a mesh for each node
	mesh := THREE.NewMesh(box, material)

	mesh.Position().Set(0, 0, 0)
	group.Add(mesh)

	// Define a recursive function to build the node tree.
	// This function adds a child node (as a red box) at a given position and scale,
	// then recurses to add child nodes in all six cardinal directions.
	var addChildNodes func(parent THREE.Group, depth int, x, y, z, scale float64)
	addChildNodes = func(parent THREE.Group, depth int, x, y, z, scale float64) {
		if depth <= 0 {
			return
		}
		// Create a child node mesh with a scaled box and red material.
		childBox := THREE.NewBoxGeometry(scale, scale, scale)
		childMaterial := THREE.NewMeshBasicMaterial(THREE.Color(0xff0000))
		childMesh := THREE.NewMesh(childBox, childMaterial)
		childMesh.Position().Set(x, y, z)
		parent.Add(childMesh)

		// Set offset and next scale for deeper recursion.
		offset := scale * 2.0
		nextScale := scale * 0.5

		// Recursively add child nodes in six directions: +X, -X, +Y, -Y, +Z, -Z.
		addChildNodes(parent, depth-1, x+offset, y, z, nextScale)
		addChildNodes(parent, depth-1, x-offset, y, z, nextScale)
		addChildNodes(parent, depth-1, x, y+offset, z, nextScale)
		addChildNodes(parent, depth-1, x, y-offset, z, nextScale)
		addChildNodes(parent, depth-1, x, y, z+offset, nextScale)
		addChildNodes(parent, depth-1, x, y, z-offset, nextScale)
	}

	// Start the recursive generation with a chosen depth and initial parameters.
	// Here, we use depth 3 and begin from position (2.0, 2.0, 2.0) with an initial scale of 0.5.
	addChildNodes(group, 3, 2.0, 2.0, 2.0, 0.5)

	// Return the complete group representing the recursive threejs scene.
	return group


}
*/

func loadKit() (Node, error) {
	return Node{false, false, false, true, 0001}, nil
}
