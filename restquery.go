package brest

import (
	"fmt"
	"net/http"
	"strings"
)

// RestQuery structure
type RestQuery struct {
	Request     *http.Request
	Action      Action
	Resource    string
	Key         string
	ContentType string
	Accept      string
	Content     interface{}
	Offset      int
	Limit       int
	Fields      []*Field
	Relations   []*Relation
	Sorts       []*Sort
	Filter      *Filter
	Debug       bool
}

func (q *RestQuery) String() string {
	var str string
	if q.Action == Get {
		if q.Key == "" {
			str = fmt.Sprintf("<action=%v resource=%v offset=%v limit=%v fields=%v relations=%v sorts=%v filter=%v>", q.Action, q.Resource, q.Offset, q.Limit, q.Fields, q.Relations, q.Sorts, q.Filter)
		} else {
			str = fmt.Sprintf("<action=%v resource=%v key=%v fields=%v relations=%v>", q.Action, q.Resource, q.Key, q.Fields, q.Relations)
		}
	} else if q.Action == Delete {
		str = fmt.Sprintf("<action=%v resource=%v key=%v>", q.Action, q.Resource, q.Key)
	} else {
		str = fmt.Sprintf("<action=%v resource=%v key=%v content-type=%v>", q.Action, q.Resource, q.Key, q.ContentType)
	}
	return str
}

// Field structure
type Field struct {
	Name string
}

func (f *Field) String() string {
	return f.Name
}

// Relation structure
type Relation struct {
	Name string
}

func (r *Relation) String() string {
	return r.Name
}

// Sort structure
type Sort struct {
	Name string
	Asc  bool
}

func (s *Sort) String() string {
	if s.Asc {
		return fmt.Sprintf("asc(%v)", s.Name)
	}
	return fmt.Sprintf("desc(%v)", s.Name)
}

// Filter structure
type Filter struct {
	Op      Op          // operation
	Attr    string      // attribute name
	Value   interface{} // attribute value
	Filters []*Filter   // sub filters for 'and' and 'or' operations
}

func (f *Filter) String() string {
	if f.Op == "" {
		return "None"
	}
	if f.Op == And || f.Op == Or {
		var sb strings.Builder
		for _, filter := range f.Filters {
			sb.WriteRune(' ')
			sb.WriteString(filter.String())
			sb.WriteRune(' ')
		}
		return fmt.Sprintf("%v (%v)", f.Op, sb)
	}
	return fmt.Sprintf("%v %v %v", f.Attr, f.Op, f.Value)
}
