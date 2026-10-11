package logger

import (
	"fmt"
	"strconv"
	"strings"
)

// Type selects which Logger implementation a solver uses.
type Type uint8

const (
	TypePrint     Type = 0 // printer.Log
	TypeDashboard Type = 1 // dashboard.Log
)

// String returns the upper-case name of t.
func (t Type) String() string {
	switch t {
	case TypePrint:
		return "PRINT"
	case TypeDashboard:
		return "DASHBOARD"
	}
	return "Type(" + strconv.Itoa(int(t)) + ")"
}

// ParseType parses a type name, ignoring case.
func ParseType(s string) (Type, error) {
	for _, t := range [...]Type{TypePrint, TypeDashboard} {
		if strings.EqualFold(s, t.String()) {
			return t, nil
		}
	}
	return 0, fmt.Errorf("logger: unknown type %q", s)
}
