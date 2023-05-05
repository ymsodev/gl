package gl

import (
	"fmt"
	"strconv"
	"strings"
)

type GLObject interface {
	glObject() // marker for GL object types

	String() string
}

type GLNil struct{}
type GLSymbol struct{ name string }
type GLBool struct{ val bool }
type GLNumber struct{ val float64 }
type GLString struct{ val string }
type GLError struct{ val error }
type GLFunction struct{ fn func(...GLObject) GLObject }
type GLList struct{ items []GLObject }

func (GLNil) glObject()      {}
func (GLSymbol) glObject()   {}
func (GLBool) glObject()     {}
func (GLNumber) glObject()   {}
func (GLString) glObject()   {}
func (GLError) glObject()    {}
func (GLFunction) glObject() {}
func (GLList) glObject()     {}

func (g GLNil) String() string      { return "nil" }
func (g GLSymbol) String() string   { return fmt.Sprintf("<%s>", g.name) }
func (g GLBool) String() string     { return strconv.FormatBool(g.val) }
func (g GLNumber) String() string   { return strconv.FormatFloat(g.val, 'f', -1, 64) }
func (g GLString) String() string   { return strconv.Quote(g.val) }
func (g GLError) String() string    { return fmt.Sprintf("error: %v", g.val) }
func (g GLFunction) String() string { return "<function>" }
func (g GLList) String() string {
	var b strings.Builder
	b.WriteRune('ʕ')
	for i := range g.items {
		if i != 0 {
			b.WriteRune(' ')
		}
		b.WriteString(g.items[i].String())
	}
	b.WriteRune('ʔ')
	return b.String()
}

func (g GLError) Error() string {
	return g.val.Error()
}
