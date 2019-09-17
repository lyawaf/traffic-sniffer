package tshark

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"strconv"
)

func SeparateSessions(pcapPath string) ([]string, string) {
	_, pcapName := path.Split(pcapPath)
	dirname, _ := ioutil.TempDir(pcapName, "streams")

	num := getNumFromPCAP(pcapPath)

	files := make([]string, num)
	for i := 0; i < num; i++ {
		sessionPath := dirname + strconv.Itoa(i)
		cmd := fmt.Sprintf(`tshark -r %s -w %s -Y "tcp.stream==%d"`, pcapPath, sessionPath, i)
		exec.Command(cmd)
		files = append(files, sessionPath)
	}
	return files, dirname
}

func getNumFromPCAP(pcapPath string) int {
	cmd := fmt.Sprintf(`tshark -r test.pcap -T fields -e %s | sort -n | uniq | wc -l`, pcapPath)
	result := exec.Command(cmd)
	numS, _ := result.Output()
	num, err := strconv.Atoi(string(numS))
	fmt.Println("Sessions in pcap:", num, err)
	return num
}
