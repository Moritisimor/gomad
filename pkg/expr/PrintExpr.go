package expr

import "fmt"

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
			if i != len(expr.Val) - 1 {
				fmt.Print(", ")
			}
		}

		fmt.Println(")")
	}
}
