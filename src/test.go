package main

import "tests"
import "fmt"

func main() {
	tests.TestTrimComment()
	tests.TestInterfaceBlock()
	tests.TestFindInterface()
	tests.TestParseSignature()
	fmt.Println("Tests PASSED")
}
