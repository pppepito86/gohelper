package tests

import "strconv"

func checkS(result, expected string) {
	if result != expected {
		panic("Test Failed. Expected <" + expected+">, but got <" + result + ">")
	}
}

func checkB(result, expected bool) {
	if result != expected {
		panic("Test Failed. Expected <" + strconv.FormatBool(expected)+">, but got <" + strconv.FormatBool(result) + ">")
	}
}

func checkI(result, expected int) {
	if result != expected {
		panic("Test Failed. Expected <" + strconv.Itoa(expected)+">, but got <" + strconv.Itoa(result) + ">")
	}
}
