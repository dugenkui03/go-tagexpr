package tagexpr

import (
	"reflect"
	"testing"
)

func TestVMFunc(t *testing.T) {
	g := &struct {
		_ int
		h string `tagexpr:"$"`
		s []string
		m map[string][]string
	}{
		h: "haha",
		s: []string{"1"},
		m: map[string][]string{"0": {"2"}},
	}
	d := "ddd"
	e := new(int)
	*e = 3
	var cases = []struct {
		tagName   string
		structure interface{}
		tests     map[string]interface{}
	}{
		{
			tagName: "tagexpr",
			structure: &struct {
				A int            `tagexpr:"$>0&&$<10"`
				b string         `tagexpr:"{is:$=='test'}{msg:sprintf('want: test, but got: %s',$)}"`
				c float32        `tagexpr:"(A)$+$"`
				d *string        `tagexpr:"$"`
				e **int          `tagexpr:"$"`
				f *[3]int        `tagexpr:"{x:len($)}{y:len()}"`
				g string         `tagexpr:"{x:regexp('g\\d{3}$',$)}{y:regexp('g\\d{3}$')}"`
				h []string       `tagexpr:"{x:$[1]}{y:$[10]}"`
				i map[string]int `tagexpr:"{x:$['a']}{y:$[0]}"`
			}{
				A: 5.0,
				b: "x",
				c: 1,
				d: &d,
				e: &e,
				f: new([3]int),
				g: "g123",
				h: []string{"", "hehe"},
				i: map[string]int{"a": 7},
			},
			tests: map[string]interface{}{
				"A.$":   true,
				"b.is":  false,
				"b.msg": "want: test, but got: x",
				"c.$":   6.0,
				"d.$":   d,
				"e.$":   float64(*e),
				"f.x":   float64(3),
				"f.y":   float64(3),
				"g.x":   true,
				"g.y":   true,
				"h.x":   "hehe",
				"h.y":   nil,
				"i.x":   7.0,
				"i.y":   nil,
			},
		},
		{
			tagName: "tagexpr",
			structure: &struct {
				A int    `tagexpr:"$>0&&$<10"`
				b string `tagexpr:"{is:$=='test'}{msg:sprintf('want: test, but got: %s',$)}"`
				c struct {
					_ int
					d bool `tagexpr:"$"`
				}
				e *struct {
					_ int
					f bool `tagexpr:"$"`
				}
				g **struct {
					_ int
					h string `tagexpr:"$"`
					s []string
					m map[string][]string
				} `tagexpr:"$['h']"`
				i string `tagexpr:"(g.s)$[0]+(g.m)$['0'][0]==$"`
			}{
				A: 5.0,
				b: "x",
				c: struct {
					_ int
					d bool `tagexpr:"$"`
				}{d: true},
				e: &struct {
					_ int
					f bool `tagexpr:"$"`
				}{f: true},
				g: &g,
				i: "12",
			},
			tests: map[string]interface{}{
				"A.$":   true,
				"b.is":  false,
				"b.msg": "want: test, but got: x",
				"c.d.$": true,
				"e.f.$": true,
				"g.h.$": "haha",
				"i.$":   true,
			},
		},
	}
	for i, c := range cases {
		vm := New(c.tagName)
		// vm.WarmUp(c.structure)
		tagExpr, err := vm.Run(c.structure)
		if err != nil {
			t.Fatal(err)
		}
		for selector, value := range c.tests {
			val := tagExpr.Eval(selector)
			if !reflect.DeepEqual(val, value) {
				t.Fatalf("Eval NO: %d, selector: %q, got: %v, want: %v", i, selector, val, value)
			}
		}
		tagExpr.Range(func(selector string, eval func() interface{}) {
			t.Logf("selector: %s", selector)
			value := c.tests[selector]
			val := eval()
			if !reflect.DeepEqual(val, value) {
				t.Fatalf("Range NO: %d, selector: %q, got: %v, want: %v", i, selector, val, value)
			}
		})
	}
}
