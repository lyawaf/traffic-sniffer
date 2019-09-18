package tshark

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
    "strings"
)

func SeparateSessions(pcapPath string) ([]string, string) {
	dirname, _ := ioutil.TempDir("", "streams")

	num := getNumFromPCAP(pcapPath)

	files := make([]string, num)
	for i := 0; i < num; i++ {
        numS := strconv.Itoa(i)
		sessionPath := dirname + "/" + numS
        r := exec.Command("tshark", "-r", pcapPath,  "-w", sessionPath, "-Y", "tcp.stream==" + numS)
        _, err := r.Output()
        if err != nil {
            fmt.Println("failed to save session", i)
            continue
        }
		files = append(files, sessionPath)
	}
    fmt.Println(files, dirname)
	return files, dirname
}

func getNumFromPCAP(pcapPath string) int {
	result := exec.Command("./countSessions.sh", pcapPath)
	numB, err := result.Output()
    if err != nil {
        fmt.Println("failed to get session num", err)
        return 0
    }
    numS := strings.TrimRight(string(numB), "\n")
	num, err := strconv.Atoi(string(numS))
	return num
}
