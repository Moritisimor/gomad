package lexer

import "fmt"

func PrintToken(t Token) {
	switch tok := t.(type) {
	case STRINGLIT:
		fmt.Printf("STRINGLIT(\"%s\")\n", tok.val)

	case BOOLLIT:
		fmt.Printf("BOOLLIT(%t)\n", tok.val)

	case NUMLIT:
		fmt.Printf("NUMLIT(%f)\n", tok.val)

	case LPAREN:
		fmt.Println("LPAREN")

	case RPAREN:
		fmt.Println("RPAREN")

	case UNITLIT:
		fmt.Println("UNIT")

	case SYMBOL:
		fmt.Printf("SYMBOL('%s')\n", tok.val)

	default:
		fmt.Println("UNKNOWN")
	}
}
