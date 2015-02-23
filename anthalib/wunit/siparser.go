package wunit

type NodeType uint8

const (
	TopNode = iota
	UnitPlusPrefixNode
	LeafNode
)

type PNode struct {
	Name     string
	Type     NodeType
	Up       *PNode
	Children []*PNode
	Value    interface{}
}

type SIPrefixedUnit struct {
	Stack   []*PNode
	Text    string
	TreeTop *PNode
	CurNode *PNode
}

func (p *SIPrefixedUnit) Init(text []byte) {
	p.Text = string(text)
	node := NewNode("Top", UnitPlusPrefixNode, 2)
	p.TreeTop = node
	p.CurNode = node
}

// couple of utility functions

func (p *SIPrefixedUnit) AddNodeToStack(node *PNode) {
	p.Stack = append(p.Stack, node)

	/*
		for _, n := range p.Stack {
			fmt.Print(n.Name, ":")
			fmt.Println()
		}
	*/
}

func (p *SIPrefixedUnit) AddStackToNode(node *PNode) {
	node.Children = p.Stack

	for _, n := range node.Children {
		n.Up = node
	}
}

func (p *SIPrefixedUnit) PopStack() (node *PNode) {
	lastItem := len(p.Stack) - 1

	if lastItem == -1 {
		return nil
	}
	node = p.Stack[lastItem]

	if lastItem == 0 {
		p.NewStack()
	} else {
		p.Stack = p.Stack[:lastItem]
	}

	return node
}

func (p *SIPrefixedUnit) PopStackAndAddTo(node *PNode) {
	n := p.PopStack()
	if n != nil {
		node.AddChild(n)
	}
}

func (pn *PNode) AddChild(child *PNode) {
	if child != nil {
		child.Up = pn
		pn.Children = append(pn.Children, child)
	}
}

func (p *SIPrefixedUnit) NewStack() {
	p.Stack = make([]*PNode, 0, 4)
}

func (p *SIPrefixedUnit) AddNodeToCurrent(node *PNode) {
	p.CurNode.Children = append(p.CurNode.Children, node)
}

func NewNode(name string, typ NodeType, cap uint8) *PNode {
	var children []*PNode
	children = nil
	if cap != 0 {
		children = make([]*PNode, 0, cap)
	}
	node := PNode{name, typ, nil, children, nil}
	return &node
}

// Functions for building the tree

func (p *SIPrefixedUnit) AddUnit(s string) {
	//	fmt.Println("Adding Unit", s)
	node := NewNode("Unit", LeafNode, 0)
	node.Value = s
	p.AddNodeToStack(node)
}

func (p *SIPrefixedUnit) AddUnitPrefix(s string) {
	//	fmt.Println("Adding unit prefix", s)
	node := NewNode("UnitPrefix", LeafNode, 0)
	node.Value = s
	p.AddNodeToStack(node)
}

func (p *SIPrefixedUnit) AddUnitPlusPrefixNode() {
	node := NewNode("UnitPlusPrefix", UnitPlusPrefixNode, 2)
	pref := p.PopStack()
	p.PopStackAndAddTo(node)
	node.AddChild(pref)
	p.TreeTop = node
	p.CurNode = node
}
