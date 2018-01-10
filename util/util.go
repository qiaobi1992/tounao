package util

import (
	"log"
	"github.com/huichen/sego"
	"strings"
	"bytes"
	"os/exec"
)

var segmenter sego.Segmenter

func init() {
	segmenter.LoadDictionary("data/dictionary.txt")
}

func Check(e error) (bool) {
	if e != nil {
		log.Fatalln(e)
		return false
	}
	return true
}

func Split(str string) []string {
	segments := segmenter.Segment([]byte(str))
	return sego.SegmentsToSlice(segments, true)
}


func RunWithAdb(args ...string) {

	log.Printf("adb %s", strings.Join(args, " "))

	var buffer bytes.Buffer
	cmd := exec.Command("adb", args...)
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer
	err := cmd.Run()
	if cmd.Process != nil {
		cmd.Process.Kill()
	}
	if err != nil {
		log.Fatalf("adb %s: %v", strings.Join(args, " "), err)
	}
}
