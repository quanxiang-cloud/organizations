package department

/*
Copyright 2022 QuanxiangCloud Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
import (
	"strings"

	"github.com/quanxiang-cloud/organizations/internal/models/org"
)

//Node tree node
type Node struct {
	// /A/B/C
	Pattern string
	// ps: A
	DepName string
	// child node ps: B、C
	Children []*Node
	DepID    string
}

func (n *Node) matchChild(part string) *Node {
	for _, child := range n.Children {
		if child.DepName == part {
			return child
		}
	}
	return nil
}

func (n *Node) matchChildren(part string) []*Node {
	nodes := make([]*Node, 0)
	if n.DepName == part {
		return n.Children
	}
	for _, child := range n.Children {
		if child.DepName == part {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *Node) search(parts []string, height int) *Node {
	if len(parts) == height || len(parts) == 1 {
		if n.Pattern == "" {
			return nil
		}
		if n.DepName == parts[0] {
			return n
		}
		return nil
	}

	part := parts[height]
	if n.DepName == part {
		children := n.matchChildren(part)
		if children == nil {
			return n
		}
		for _, child := range children {
			result := child.search(parts, height+1)
			if result != nil {
				return result
			}
		}
	}
	return nil
}

//DepRouter Department tree router
type DepRouter struct {
	Roots *Node
}

//NewDepartmentRouter new
func NewDepartmentRouter() *DepRouter {
	return &DepRouter{
		Roots: &Node{},
	}
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
		}
	}
	return parts
}

//AddRoute add
func (r *DepRouter) AddRoute(deps []org.Department) {
	node := r.makeTrees(deps)
	r.Roots = node
}

//GetRoute get
func (r *DepRouter) GetRoute(path string) *Node {
	searchParts := parsePattern(path)
	root := r.Roots
	if root == nil {
		return nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		return n
	}
	return nil
}

/*
makeRoot
*/
func (r *DepRouter) makeTrees(deps []org.Department) *Node {
	var outs = &Node{}
	var mps = make(map[string][]Node)
	for k, v := range deps {
		if v.PID == "" {
			outs.Pattern = deps[k].Name
			outs.DepName = deps[k].Name
			outs.DepID = deps[k].ID
		} else {
			child := Node{}
			child.Pattern = deps[k].Name
			child.DepName = deps[k].Name
			child.DepID = deps[k].ID
			mps[v.PID] = append(mps[v.PID], child)
		}

	}
	if outs != nil {

		r.makeTree(outs, mps)

	}

	return outs
}

/*
makeTree
*/
func (r *DepRouter) makeTree(dep *Node, mps map[string][]Node) {
	for k := range mps {
		if k == dep.DepID {
			for k1 := range mps[k] {
				mps[k][k1].Pattern = dep.Pattern + "/" + mps[k][k1].DepName
				dep.Children = append(dep.Children, &mps[k][k1])
			}

		}
	}
	if len(dep.Children) > 0 {
		for i := 0; i < len(dep.Children); i++ {
			r.makeTree(dep.Children[i], mps)
		}
	}
}
