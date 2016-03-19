package tests

import "fmt"
import "os"
import "path"
import "api"

func TestInterfaceBlock() {
	fmt.Println("TestInterfaceBlock")
	wd, _ := os.Getwd()
	interfacesFile := path.Join(wd, "tests", "files", "interface.go")	
	blocks := make([][]string, 0)
	api.FileBlocks(interfacesFile, &blocks)
	checkI(len(blocks), 5)

	checkBlock0(blocks[0])
	checkBlock1(blocks[1])
	checkBlock2(blocks[2])
	checkBlock3(blocks[3])
	checkBlock4(blocks[4])
}

func checkBlock0(block []string) {
	checkI(len(block), 3)
	checkS(block[1], "1")
	checkS(block[2], "package interfaces")
}

func checkBlock1(block []string) {
	checkI(len(block), 4)
	checkS(block[1], "3")
	checkS(block[2], "type First interface {")
	checkS(block[3], "}")
}

func checkBlock2(block []string) {
	checkI(len(block), 3)
	checkS(block[1], "6")
	checkS(block[2], "type Second interface {}")
}

func checkBlock3(block []string) {
	checkI(len(block), 6)
	checkS(block[1], "8")
	checkS(block[2], "type Third interface{")
	checkS(block[3], "")
	checkS(block[4], "")
	checkS(block[5], "}")
}

func checkBlock4(block []string) {
	checkI(len(block), 5)
	checkS(block[1], "13")
	checkS(block[2], "type")
	checkS(block[3], "Fourth interface {")
	checkS(block[4], "}")
}



