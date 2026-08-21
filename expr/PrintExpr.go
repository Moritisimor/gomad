package expr

import (
	"fmt"
	"strings"
)

func PrintExpr(e Expression) {
	switch expr := e.(type) {
	case Lambda:
		fmt.Print("<LAMBDA>")

	case Boolean:
		fmt.Printf("Boolean(%t)", expr.Val)

	case String:
		fmt.Printf("String(\"%s\")", expr.Val)

	case Number:
		fmt.Printf("Number(%f)", expr.Val)

	case Symbol:
		fmt.Printf("Symbol('%s')", expr.Val)

	case Unit:
		fmt.Println("Unit")

	case List:
		fmt.Print("List(")
		for i, e := range expr.Val {
			PrintExpr(e)
			if i != len(expr.Val)-1 {
				fmt.Print(", ")
			}
		}

		fmt.Print(")")
	}
}

func SprintExpr(e Expression) string {
	switch expr := e.(type) {
	case Lambda:
		return "<LAMBDA>"

	case Boolean:
		return fmt.Sprintf("Boolean(%t)", expr.Val)

	case String:
		return fmt.Sprintf("String(\"%s\")", expr.Val)

	case Number:
		return fmt.Sprintf("Number(%f)", expr.Val)

	case Symbol:
		return fmt.Sprintf("Symbol('%s')", expr.Val)

	case Unit:
		return "Unit"

	case List:
		acc := strings.Builder{}
		acc.WriteString("List(")
		for i, e := range expr.Val {
			acc.WriteString(SprintExpr(e))
			if i != len(expr.Val)-1 {
				acc.WriteString(", ")
			}
		}

		acc.WriteByte(')')
		return acc.String()

	default:
		return "Unknown expression"
	}
}
