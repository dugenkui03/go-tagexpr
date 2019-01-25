package tagexpr

import (
	"regexp"
	"strings"
)

type selectorExprNode struct {
	exprBackground
	field, name string
	subExprs    []ExprNode
}

func (p *Expr) readSelectorExprNode(expr *string) ExprNode {
	field, name, subSelector, found := findSelector(expr)
	if !found {
		return nil
	}
	operand := &selectorExprNode{
		field: field,
		name:  name,
	}
	operand.subExprs = make([]ExprNode, 0, len(subSelector))
	for _, s := range subSelector {
		grp := newGroupExprNode()
		_, err := p.parseExprNode(&s, grp)
		if err != nil {
			return nil
		}
		sortPriority(grp.RightOperand())
		operand.subExprs = append(operand.subExprs, grp)
	}
	return operand
}

var selectorRegexp = regexp.MustCompile(`^(\([ \t]*[A-Za-z_]+[A-Za-z0-9_\.]*[ \t]*\))?(\$)([\[\+\-\*\/%><\|&!=\^ \t\\]|$)`)

func findSelector(expr *string) (field string, name string, subSelector []string, found bool) {
	raw := *expr
	a := selectorRegexp.FindAllStringSubmatch(raw, -1)
	if len(a) != 1 {
		return
	}
	r := a[0]
	if s0 := r[1]; len(s0) > 0 {
		field = strings.TrimSpace(s0[1 : len(s0)-1])
	}
	name = r[2]
	*expr = (*expr)[len(a[0][0])-len(r[3]):]
	for {
		sub := readPairedSymbol(expr, '[', ']')
		if sub == nil {
			break
		}
		if *sub == "" || (*sub)[0] == '[' {
			*expr = raw
			return "", "", nil, false
		}
		subSelector = append(subSelector, strings.TrimSpace(*sub))
	}
	found = true
	return
}

func (ve *selectorExprNode) Run(currField string, tagExpr *TagExpr) interface{} {
	subFields := make([]interface{}, 0, len(ve.subExprs))
	for _, e := range ve.subExprs {
		subFields = append(subFields, e.Run(currField, tagExpr))
	}
	if ve.field != "" {
		return tagExpr.getValue(ve.field, subFields)
	}
	return tagExpr.getValue(currField, subFields)
}
