package api

import "io/ioutil"
import "path"
import "strings"
import "strconv"
import "regexp"

func ParseSignature(signature string) [][]string {
	signature = strings.TrimSpace(signature)
	s := make([][]string, 3)
	s[0] = []string{getName(signature)}
	s[1] = getParams(signature)
	s[2] = []string{}
	return s
}

func getName(signature string) string {
	name := strings.Split(signature, "(")[0]
	name = strings.TrimSpace(name)
	return name
}

func getParams(signature string) []string {
	paramsString := getParamsString(signature)
	if paramsString == "" {
		return []string{}
	}
	paramsSplit := strings.Split(paramsString, ",")
	lastParam := parseParam(paramsSplit[len(paramsSplit)-1], "_")
	params := make([]string, len(paramsSplit))
	params[len(params)-1] = lastParam
	for i := len(params) - 2; i >=0; i-- {
		lastParam = parseParam(paramsSplit[i], lastParam)
		params[i] = lastParam
	}
	return params
}

func parseParam(param string, lastType string) string {
	param = strings.TrimSpace(param)
	p := strings.Split(param, " ")
	if len(p) == 1 {
		return lastType
	} else if len(p) == 2 {
		return strings.TrimSpace(p[1])
	} else if len(p) > 2 {
		return "error"
	}
	return ""
}

func getParamsString(signature string) string {
	params := strings.SplitN(signature, "(", 2)[1]
	params = strings.TrimSpace(params)
	diff := 1
	for i, c := range params {
		if c == '('{
			diff++
		}
		if c == ')'{
			diff--
		}
		if diff == 0 {
			return strings.TrimSpace(params[:i])
		}
	}
	return "error"
}

func FindInterface(interfaceName, dir string) [][]string {
	bs := Blocks(dir)
	result := make([][]string, 0)
	for _, b := range bs {
		if strings.HasPrefix(b[2], "type " + interfaceName + " interface") {
			result = append(result, b)
		} else if b[2] == "type" && len(b) >= 4 && strings.HasPrefix(b[3], interfaceName + " interface") {
			result = append(result, b)
		}
	}
	return result
}

func Blocks(dir string) [][]string {
	blocks := make([][]string, 0)
	DirBlocks(dir, &blocks)
	return blocks
}

func DirBlocks(dirName string, blocks *[][]string) [][]string {
	dirInfo, _ := ioutil.ReadDir(dirName)
	for _, dir := range dirInfo {
		if dir.IsDir() {
			file := path.Join(dirName, dir.Name())
			*blocks = DirBlocks(file, blocks)
		} else {
			if !strings.HasSuffix(dir.Name(), ".go") {
				continue
			}
			file := path.Join(dirName, dir.Name())
			*blocks = FileBlocks(file, blocks)
		}
	}
	return *blocks
}

func FileBlocks(filename string, blocks *[][]string) [][]string {
	b, _ := ioutil.ReadFile(filename)
	s := strings.Split(string(b), "\n")
	block := createBlock(filename, 1)
	diff := 0
	isCommentStarted := false
	for idx, line := range s {
		line = RemoveMultipleSpaces(line)
		line = strings.TrimSpace(line)
		line, isCommentStarted = TrimComment(line, isCommentStarted)
		if diff == 0 && line == "" {
			block = createBlock(filename, idx+2)
			continue
		}
		block = append(block, line)
		diff = updateDiff(&diff, line)
		if diff == 0 && line != "type" {
			*blocks = append(*blocks, block)
			block = createBlock(filename, idx+2)
		}
	}
	return *blocks
}

func RemoveMultipleSpaces(line string) string {
	re := regexp.MustCompile("[ \t\r]+")
	return re.ReplaceAllString(line, " ")
}

func createBlock(fileName string, lineNumber int) []string {
	block := make([]string, 0)
	block = append(block, fileName)
	block = append(block, strconv.Itoa(lineNumber))
	return block
}

func TrimComment(line string, isCommentStarted bool) (string, bool) {
	if isCommentStarted {
		if strings.Contains(line, "*/") {
			line = strings.SplitN(line, "*/", 2)[1]
		} else {
			return "", true
		}
	}
	for i:=1; i < len(line); i++ {
		if line[i-1:i+1] == "//" {
			return line[:i-1], false
		}
		if line[i-1:i+1] == "/*" {
			s, b := TrimComment(line[i+1:], true)
			return line[:i-1]+s, b
		}
	}
	return line, false
}


func updateDiff(diff *int, line string) int {
	*diff += strings.Count(line, "(")
	*diff += strings.Count(line, "{")
	*diff -= strings.Count(line, ")")
	*diff -= strings.Count(line, "}")
	return *diff
}
