package tests

import "fmt"
import "api"

func TestParseSignature() {
	fmt.Println("TestParseSignature")
	s := api.ParseSignature("f()")
	checkS(fmt.Sprint(s), "[[f] [] []]")

	s = api.ParseSignature("f ( )")
	checkS(fmt.Sprint(s), "[[f] [] []]")

	s = api.ParseSignature("testFunc()")
	checkS(fmt.Sprint(s), "[[testFunc] [] []]")

	s = api.ParseSignature("testFunc(a int)")
	checkS(fmt.Sprint(s), "[[testFunc] [int] []]")

	s = api.ParseSignature("testFunc(adsaf string)")
	checkS(fmt.Sprint(s), "[[testFunc] [string] []]")

	s = api.ParseSignature("testFunc(a int, b int)")
	checkS(fmt.Sprint(s), "[[testFunc] [int int] []]")

	s = api.ParseSignature("testFunc(a int, b string, ccc TTTTt)")
	checkS(fmt.Sprint(s), "[[testFunc] [int string TTTTt] []]")

	s = api.ParseSignature("testFunc(a , b int, c ...int)")
	checkS(fmt.Sprint(s), "[[testFunc] [int int ...int] []]")

	s = api.ParseSignature("testFunc(s string, a , b int, s, c ,d bool)")
	checkS(fmt.Sprint(s), "[[testFunc] [string int int bool bool bool] []]")

}

