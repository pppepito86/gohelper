package tests

import "fmt"
import "api"

func TestTrimComment() {
	fmt.Println("TestTrimComment")
	s, b := api.TrimComment("", false)
	checkS(s, "")
	checkB(b, false)
	s, b = api.TrimComment("aaa//bbb", false)
	checkS(s, "aaa")
	checkB(b, false)
	s, b = api.TrimComment("aaa//bbb/a/*sdf*/dsf", false)
	checkS(s, "aaa")
	checkB(b, false)
	s, b = api.TrimComment("aaa/*bbb/a//*sdfdsf", false)
	checkS(s, "aaa")
	checkB(b, true);
	s, b = api.TrimComment("aaa/*bbb/a//*sdf*/dsf", false)
	checkS(s, "aaadsf")
	checkB(b, false);
	s, b = api.TrimComment("aaa/*bbb/a//*sdf*/dsf//dfasdf/*", false)
	checkS(s, "aaadsf")
	checkB(b, false);

	s, b = api.TrimComment("aaa//bbb", true)
	checkS(s, "")
	checkB(b, true)
	s, b = api.TrimComment("aaa//bbb/a/*sdf*/dsf", true)
	checkS(s, "dsf")
	checkB(b, false)
	s, b = api.TrimComment("aaa/*bbb/a//*sdfdsf", true)
	checkS(s, "")
	checkB(b, true);
	s, b = api.TrimComment("aaa/*bbb/a//*sdf*/dsf", true)
	checkS(s, "dsf")
	checkB(b, false);
	s, b = api.TrimComment("aaa/*bbb/a//*sdf*/dsf//dfasdf/*", true)
	checkS(s, "dsf")
	checkB(b, false);
}

