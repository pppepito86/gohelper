package main

import "fmt"
import "io/ioutil"
import "path"
import "strings"
import "os"

func main()	{
	blocks := getBlocks()
	for _, block := range blocks {
		if strings.HasPrefix(block[0], "func (") {
			firstLine := ""
			for i := 0; i < len(block) && !strings.Contains(firstLine, "{"); i++ {
				firstLine += block[i]
			}
			if !strings.Contains(firstLine, "{") {
				fmt.Println("!!!!!")
				fmt.Println(firstLine)
				//panic("err")
			}
			fmt.Println(firstLine)
			t, s := parse(firstLine)
			fmt.Println(t, s)
		}
	}
}

func parse(line string) (string, string) {
	return getType(line), getSignature(line)
}

func getType(line string) string {
	line = strings.Split(line, "(")[1]
	line = strings.Split(line, ")")[0]
	if strings.Contains(line, " ") {
		line = strings.Split(line, " ")[1]
	}
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "*") {
		line = line[1:]
	}
	return line
}

func getSignature(line string) string {
	line = strings.SplitN(line, ")", 2)[1]
	line = strings.Split(line, "{")[0]
	line = strings.TrimSpace(line)
	return line
}

func getBlocks() [][]string {
	blocks := make([][]string, 0)
	dir := os.Getenv("pesho")
	if dir == "" {
		dir = "."
	}
	getDirBlocks(dir, &blocks)
	return blocks
}

func getDirBlocks(dirName string, blocks *[][]string) [][]string {
	dirInfo, _ := ioutil.ReadDir(dirName)
	for _, dir := range dirInfo {
		if dir.IsDir() {
			file := path.Join(dirName, dir.Name())
			*blocks = getDirBlocks(file, blocks)
		} else {
			file := path.Join(dirName, dir.Name())
			*blocks = getFileBlocks(file, blocks)
		}
	}
	return *blocks
}

func getFileBlocks(filename string, blocks *[][]string) [][]string {
	b, _ := ioutil.ReadFile(filename)
	s := strings.Split(string(b), "\n")
	block := make([]string, 0)
	diff := 0
	for _, line := range s {
		line = strings.TrimSpace(line)
		block = append(block, line)
		diff = updateDiff(&diff, line)
		if diff == 0 {
			*blocks = append(*blocks, block)
			block = make([]string, 0)
		}
	}
	return *blocks
}

func updateDiff(diff *int, line string) int {
	*diff += strings.Count(line, "(")
	*diff += strings.Count(line, "{")
	*diff -= strings.Count(line, ")")
	*diff -= strings.Count(line, "}")
	return *diff
}
