package wunit

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const end_symbol rune = 4

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleunit_plus_prefix
	rulesi_prefix
	ruleunit
	ruleAction0
	rulePegText
	ruleAction1
	ruleAction2

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"unit_plus_prefix",
	"si_prefix",
	"unit",
	"Action0",
	"PegText",
	"Action1",
	"Action2",

	"Pre_",
	"_In_",
	"_Suf",
}

type tokenTree interface {
	Print()
	PrintSyntax()
	PrintSyntaxTree(buffer string)
	Add(rule pegRule, begin, end, next, depth int)
	Expand(index int) tokenTree
	Tokens() <-chan token32
	AST() *node32
	Error() []token32
	trim(length int)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(depth int, buffer string) {
	for node != nil {
		for c := 0; c < depth; c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(buffer[node.begin:node.end]))
		if node.up != nil {
			node.up.print(depth+1, buffer)
		}
		node = node.next
	}
}

func (ast *node32) Print(buffer string) {
	ast.print(0, buffer)
}

type element struct {
	node *node32
	down *element
}

/* ${@} bit structure for abstract syntax tree */
type token16 struct {
	pegRule
	begin, end, next int16
}

func (t *token16) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token16) isParentOf(u token16) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token16) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: int32(t.begin), end: int32(t.end), next: int32(t.next)}
}

func (t *token16) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens16 struct {
	tree    []token16
	ordered [][]token16
}

func (t *tokens16) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens16) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens16) Order() [][]token16 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int16, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token16, len(depths)), make([]token16, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = int16(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state16 struct {
	token16
	depths []int16
	leaf   bool
}

func (t *tokens16) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens16) PreOrder() (<-chan state16, [][]token16) {
	s, ordered := make(chan state16, 6), t.Order()
	go func() {
		var states [8]state16
		for i, _ := range states {
			states[i].depths = make([]int16, len(ordered))
		}
		depths, state, depth := make([]int16, len(ordered)), 0, 1
		write := func(t token16, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, int16(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token16 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token16{pegRule: rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token16{pegRule: rulePre_, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token16{pegRule: rule_Suf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens16) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens16) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(buffer[token.begin:token.end]))
	}
}

func (t *tokens16) Add(rule pegRule, begin, end, depth, index int) {
	t.tree[index] = token16{pegRule: rule, begin: int16(begin), end: int16(end), next: int16(depth)}
}

func (t *tokens16) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens16) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i, _ := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next int32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: int32(t.begin), end: int32(t.end), next: int32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = int32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens32) PreOrder() (<-chan state32, [][]token32) {
	s, ordered := make(chan state32, 6), t.Order()
	go func() {
		var states [8]state32
		for i, _ := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, int32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{pegRule: rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre_, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{pegRule: rule_Suf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(buffer[token.begin:token.end]))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth, index int) {
	t.tree[index] = token32{pegRule: rule, begin: int32(begin), end: int32(end), next: int32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i, _ := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

func (t *tokens16) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		for i, v := range tree {
			expanded[i] = v.getToken32()
		}
		return &tokens32{tree: expanded}
	}
	return nil
}

func (t *tokens32) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	return nil
}

type SIPrefixedUnitGrammar struct {
	SIPrefixedUnit

	Buffer string
	buffer []rune
	rules  [8]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	tokenTree
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer string, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer[0:] {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p *SIPrefixedUnitGrammar
}

func (e *parseError) Error() string {
	tokens, error := e.p.tokenTree.Error(), "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.Buffer, positions)
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf("parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n",
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			/*strconv.Quote(*/ e.p.Buffer[begin:end] /*)*/)
	}

	return error
}

func (p *SIPrefixedUnitGrammar) PrintSyntaxTree() {
	p.tokenTree.PrintSyntaxTree(p.Buffer)
}

func (p *SIPrefixedUnitGrammar) Highlighter() {
	p.tokenTree.PrintSyntax()
}

func (p *SIPrefixedUnitGrammar) Execute() {
	buffer, begin, end := p.Buffer, 0, 0
	for token := range p.tokenTree.Tokens() {
		switch token.pegRule {
		case rulePegText:
			begin, end = int(token.begin), int(token.end)
		case ruleAction0:
			p.AddUnitPlusPrefixNode()
		case ruleAction1:
			p.AddUnitPrefix(buffer[begin:end])
		case ruleAction2:
			p.AddUnit(buffer[begin:end])

		}
	}
}

func (p *SIPrefixedUnitGrammar) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != end_symbol {
		p.buffer = append(p.buffer, end_symbol)
	}

	var tree tokenTree = &tokens16{tree: make([]token16, math.MaxInt16)}
	position, depth, tokenIndex, buffer, _rules := 0, 0, 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokenTree = tree
		if matches {
			p.tokenTree.trim(tokenIndex)
			return nil
		}
		return &parseError{p}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin int) {
		if t := tree.Expand(tokenIndex); t != nil {
			tree = t
		}
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 unit_plus_prefix <- <((si_prefix &unit)? unit Action0)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				{
					position2, tokenIndex2, depth2 := position, tokenIndex, depth
					if !_rules[rulesi_prefix]() {
						goto l2
					}
					{
						position4, tokenIndex4, depth4 := position, tokenIndex, depth
						if !_rules[ruleunit]() {
							goto l2
						}
						position, tokenIndex, depth = position4, tokenIndex4, depth4
					}
					goto l3
				l2:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
				}
			l3:
				if !_rules[ruleunit]() {
					goto l0
				}
				if !_rules[ruleAction0]() {
					goto l0
				}
				depth--
				add(ruleunit_plus_prefix, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 si_prefix <- <(<('y' / 'z' / 'a' / 'f' / 'p' / 'n' / 'u' / 'm' / 'c' / 'd' / 'h' / 'k' / 'M' / 'G' / 'T' / 'P' / 'E' / 'Z' / 'Y' / ('d' 'a'))> Action1)> */
		func() bool {
			position5, tokenIndex5, depth5 := position, tokenIndex, depth
			{
				position6 := position
				depth++
				{
					position7 := position
					depth++
					{
						position8, tokenIndex8, depth8 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l9
						}
						position++
						goto l8
					l9:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('z') {
							goto l10
						}
						position++
						goto l8
					l10:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('a') {
							goto l11
						}
						position++
						goto l8
					l11:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('f') {
							goto l12
						}
						position++
						goto l8
					l12:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('p') {
							goto l13
						}
						position++
						goto l8
					l13:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('n') {
							goto l14
						}
						position++
						goto l8
					l14:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('u') {
							goto l15
						}
						position++
						goto l8
					l15:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('m') {
							goto l16
						}
						position++
						goto l8
					l16:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('c') {
							goto l17
						}
						position++
						goto l8
					l17:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('d') {
							goto l18
						}
						position++
						goto l8
					l18:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('h') {
							goto l19
						}
						position++
						goto l8
					l19:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('k') {
							goto l20
						}
						position++
						goto l8
					l20:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('M') {
							goto l21
						}
						position++
						goto l8
					l21:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('G') {
							goto l22
						}
						position++
						goto l8
					l22:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('T') {
							goto l23
						}
						position++
						goto l8
					l23:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('P') {
							goto l24
						}
						position++
						goto l8
					l24:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('E') {
							goto l25
						}
						position++
						goto l8
					l25:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('Z') {
							goto l26
						}
						position++
						goto l8
					l26:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('Y') {
							goto l27
						}
						position++
						goto l8
					l27:
						position, tokenIndex, depth = position8, tokenIndex8, depth8
						if buffer[position] != rune('d') {
							goto l5
						}
						position++
						if buffer[position] != rune('a') {
							goto l5
						}
						position++
					}
				l8:
					depth--
					add(rulePegText, position7)
				}
				if !_rules[ruleAction1]() {
					goto l5
				}
				depth--
				add(rulesi_prefix, position6)
			}
			return true
		l5:
			position, tokenIndex, depth = position5, tokenIndex5, depth5
			return false
		},
		/* 2 unit <- <(<(('r' 'a' 'd' 's') / ('r' 'a' 'd' 'i' 'a' 'n' 's') / ('d' 'e' 'g' 'r' 'e' 'e' 's') / ('H' 'z') / ('r' 'p' 'm') / ('h' / 'H' / 'M' / 'm' / 'l' / 'L' / 'g' / 'V' / 'J' / 'A' / 'C' / 'N' / 's' / '%'))> Action2)> */
		func() bool {
			position28, tokenIndex28, depth28 := position, tokenIndex, depth
			{
				position29 := position
				depth++
				{
					position30 := position
					depth++
					{
						position31, tokenIndex31, depth31 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l32
						}
						position++
						if buffer[position] != rune('a') {
							goto l32
						}
						position++
						if buffer[position] != rune('d') {
							goto l32
						}
						position++
						if buffer[position] != rune('s') {
							goto l32
						}
						position++
						goto l31
					l32:
						position, tokenIndex, depth = position31, tokenIndex31, depth31
						if buffer[position] != rune('r') {
							goto l33
						}
						position++
						if buffer[position] != rune('a') {
							goto l33
						}
						position++
						if buffer[position] != rune('d') {
							goto l33
						}
						position++
						if buffer[position] != rune('i') {
							goto l33
						}
						position++
						if buffer[position] != rune('a') {
							goto l33
						}
						position++
						if buffer[position] != rune('n') {
							goto l33
						}
						position++
						if buffer[position] != rune('s') {
							goto l33
						}
						position++
						goto l31
					l33:
						position, tokenIndex, depth = position31, tokenIndex31, depth31
						if buffer[position] != rune('d') {
							goto l34
						}
						position++
						if buffer[position] != rune('e') {
							goto l34
						}
						position++
						if buffer[position] != rune('g') {
							goto l34
						}
						position++
						if buffer[position] != rune('r') {
							goto l34
						}
						position++
						if buffer[position] != rune('e') {
							goto l34
						}
						position++
						if buffer[position] != rune('e') {
							goto l34
						}
						position++
						if buffer[position] != rune('s') {
							goto l34
						}
						position++
						goto l31
					l34:
						position, tokenIndex, depth = position31, tokenIndex31, depth31
						if buffer[position] != rune('H') {
							goto l35
						}
						position++
						if buffer[position] != rune('z') {
							goto l35
						}
						position++
						goto l31
					l35:
						position, tokenIndex, depth = position31, tokenIndex31, depth31
						if buffer[position] != rune('r') {
							goto l36
						}
						position++
						if buffer[position] != rune('p') {
							goto l36
						}
						position++
						if buffer[position] != rune('m') {
							goto l36
						}
						position++
						goto l31
					l36:
						position, tokenIndex, depth = position31, tokenIndex31, depth31
						{
							position37, tokenIndex37, depth37 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l38
							}
							position++
							goto l37
						l38:
							position, tokenIndex, depth = position37, tokenIndex37, depth37
							if buffer[position] != rune('H') {
								goto l39
							}
							position++
							goto l37
						l39:
							position, tokenIndex, depth = position37, tokenIndex37, depth37
							if buffer[position] != rune('M') {
								goto l40
							}
							position++
							goto l37
						l40:
							position, tokenIndex, depth = position37, tokenIndex37, depth37
							if buffer[position] != rune('m') {
								goto l41
							}
							position++
							goto l37
						l41:
							position, tokenIndex, depth = position37, tokenIndex37, depth37
							if buffer[position] != rune('l') {
								goto l42
							}
							position++
							goto l37
						l42:
							position, tokenIndex, depth = position37, tokenIndex37, depth37
							if buffer[position] != rune('L') {
								goto l43
							}
							position++
							goto l37
						l43:
							position, tokenIndex, depth = position37, tokenIndex37, depth37
							if buffer[position] != rune('g') {
								goto l44
							}
							position++
							goto l37
						l44:
							position, tokenIndex, depth = position37, tokenIndex37, depth37
							if buffer[position] != rune('V') {
								goto l45
							}
							position++
							goto l37
						l45:
							position, tokenIndex, depth = position37, tokenIndex37, depth37
							if buffer[position] != rune('J') {
								goto l46
							}
							position++
							goto l37
						l46:
							position, tokenIndex, depth = position37, tokenIndex37, depth37
							if buffer[position] != rune('A') {
								goto l47
							}
							position++
							goto l37
						l47:
							position, tokenIndex, depth = position37, tokenIndex37, depth37
							if buffer[position] != rune('C') {
								goto l48
							}
							position++
							goto l37
						l48:
							position, tokenIndex, depth = position37, tokenIndex37, depth37
							if buffer[position] != rune('N') {
								goto l49
							}
							position++
							goto l37
						l49:
							position, tokenIndex, depth = position37, tokenIndex37, depth37
							if buffer[position] != rune('s') {
								goto l50
							}
							position++
							goto l37
						l50:
							position, tokenIndex, depth = position37, tokenIndex37, depth37
							if buffer[position] != rune('%') {
								goto l28
							}
							position++
						}
					l37:
					}
				l31:
					depth--
					add(rulePegText, position30)
				}
				if !_rules[ruleAction2]() {
					goto l28
				}
				depth--
				add(ruleunit, position29)
			}
			return true
		l28:
			position, tokenIndex, depth = position28, tokenIndex28, depth28
			return false
		},
		/* 4 Action0 <- <{p.AddUnitPlusPrefixNode()}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		nil,
		/* 6 Action1 <- <{p.AddUnitPrefix(buffer[begin:end])}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 7 Action2 <- <{p.AddUnit(buffer[begin:end])}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
	}
	p.rules = _rules
}
