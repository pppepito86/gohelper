package main

import "fmt"
import "io/ioutil"
import "path"
import "strings"
import "os"
import "strconv"

func main()	{
	interfaceName := os.Args[1]
	dir := os.Getenv("pesho")
	if dir == "" {
		dir = "."
	}
	
	fmt.Println("Searching for interface: " + interfaceName + " in " + dir)
	blocks := getBlocks(dir)
	for _, block := range blocks {
		if strings.HasPrefix(strings.TrimSpace(block[2]), "type " + interfaceName + " interface") {
			fmt.Println(block[0] + ": " + block[1])
			for i:=2; i < len(block); i++ {
				fmt.Println("   " + block[i])
			}
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
	line = strings.TrimSpace(line)
	line = strings.SplitN(line, ")", 2)[1]
	line = strings.TrimSpace(line)
	//line = strings.Split(line, "{")[0]
	line = strings.TrimSpace(line)
	name := strings.Split(line, "(")[0]
	name = strings.TrimSpace(name)
	//fmt.Println("name: ", name)
	line = strings.SplitN(line, "(", 2)[1]
	//fmt.Println("AAAA", line)
	line = strings.TrimSpace(line)
	diff := 1
	params := ""
	for i:=0; i < len(line); i++ {
		if line[i] == '(' || line[i] == '{' {
			diff++;
		}
		if line[i] == ')' || line[i] == '}' {
			diff--;
		}
		
		if line[i] == ')' && diff == 0 {
			params = line[:i]
			line = line[i:]
			break;
		}
	}

	params = parseParams(params)
	//fmt.Println("params: ", params)
	line = strings.SplitN(line, ")", 2)[1]
	line = strings.TrimSpace(line)
	return name + "(" + params + ")"
}

func parseParams(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return line
	}
	l := strings.Split(line, ",")
	lastParam := parseParam(l[len(l)-1], "_")
	result := lastParam
	for i := len(l)-2; i >= 0; i-- {
		lastParam = parseParam(l[i], lastParam)
		result = lastParam + "," + result
	}
	return result
}

func parseParam(param string, lastType string) string {
	param = strings.TrimSpace(param)
	p := strings.Split(param, " ")
	if len(p) == 1 {
		return lastType
	} else if len(p) == 2 {
		return strings.TrimSpace(p[1])
	} else if len(p) > 2 {
		return "_"
	}
	return ""
}

func getBlocks(dir string) [][]string {
	blocks := make([][]string, 0)
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
			if !strings.HasSuffix(dir.Name(), ".go") {
				continue
			}
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
	block = append(block, filename)
	block = append(block, "1")
	diff := 0
	for idx, line := range s {
		line = strings.TrimSpace(line)
		block = append(block, line)
		diff = updateDiff(&diff, line)
		if diff == 0 {
			*blocks = append(*blocks, block)
			block = make([]string, 0)
			block = append(block, filename)
			block = append(block, strconv.Itoa(idx+2))
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
