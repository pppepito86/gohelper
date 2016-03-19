package tests

import "fmt"
import "os"
import "path"
import "api"

func TestFindInterface() {
	fmt.Println("TestFindInterface")
	wd, _ := os.Getwd()
	interfacesDir := path.Join(wd, "tests", "files")	
	checkInterface1(interfacesDir)
	checkInterface2(interfacesDir)
}

func checkInterface1(dir string) {
	blocks := api.FindInterface("Third", dir)
	checkI(len(blocks), 1)
	checkI(len(blocks[0]), 6)
	checkS(blocks[0][1], "8")
}

func checkInterface2(dir string) {
	blocks := api.FindInterface("Fourth", dir)
	checkI(len(blocks), 1)
	checkI(len(blocks[0]), 5)
	checkS(blocks[0][1], "13")
}

