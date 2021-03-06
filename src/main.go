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
			for i := 0; i < len(block); i++ {
				firstLine += block[i]
				diff := 0
				ok := false
				for j := 0; j < len(firstLine); j++ {
					if firstLine[j] == '{' && diff == 0 {
						firstLine = firstLine[:j]
						ok = true
						break
					}
					if firstLine[j] == '(' || firstLine[j] == '{' {
						diff++;
					}
					if firstLine[j] == ')' || firstLine[j] == '}' {
						diff--;
					}
				}
				if ok {
					break
				}
			}
			//fmt.Println(firstLine)
		//	t, s := parse(firstLine)
	//		fmt.Println(t, s)
		} else if strings.HasPrefix(block[0], "type") {
			sp := strings.Split(block[0], " ")
			if len(sp) < 3 || !strings.HasPrefix(sp[2], "interface") {
				continue
			}
			if len(block) > 1 {
				fmt.Println(block[1])
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
