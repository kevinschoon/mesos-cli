package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func testData(data chan *fileData, chunks int) {
	lines := []string{
		"AAAAAAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAAAAA",
		"BBBBBBBBBBBBBBBBBBBBBBBB",
		"BBBBBBBBBBBBBBBBBBBBBBBB",
		"BBBBBBBBBBBBBBBBBBBBBBBB",
		"BBBBBBBBBBBBBBBBBBBBBBBB",
		"CCCCCCCCCCCCCCCCCCCCCCCC",
		"CCCCCCCCCCCCCCCCCCCCCCCC",
		"CCCCCCCCCCCCCCCCCCCCCCCC",
		"CCCCCCCCCCCCCCCCCCCCCCCC",
	}
	for i := 0; i < chunks; i++ {
		data <- &fileData{Data: strings.Join(lines, "\n") + "\n"}
	}
	close(data)
}

func TestPailer(t *testing.T) {
	data := make(chan *fileData)
	cancel := make(chan bool)
	go testData(data, 10)
	fp, err := ioutil.TempFile("/tmp", "mesos-cli")
	assert.NoError(t, err)
	assert.NoError(t, Pailer(data, cancel, 0, fp))
	assert.NoError(t, fp.Close())
	raw, err := ioutil.ReadFile(fp.Name())
	assert.NoError(t, err)
	lines := strings.Split(string(raw), "\n")
	assert.Equal(t, 121, len(lines))
	assert.Equal(t, "AAAAAAAAAAAAAAAAAAAAAAAA", lines[111])
	assert.Equal(t, "BBBBBBBBBBBBBBBBBBBBBBBB", lines[115])
	assert.Equal(t, "CCCCCCCCCCCCCCCCCCCCCCCC", lines[119])
	fmt.Println(lines)
}

func TestPailerHuge(t *testing.T) {
	data := make(chan *fileData)
	cancel := make(chan bool)
	go testData(data, 10000)
	fp, err := ioutil.TempFile("/tmp", "mesos-cli")
	assert.NoError(t, err)
	assert.NoError(t, Pailer(data, cancel, 0, fp))
	assert.NoError(t, fp.Close())
	raw, err := ioutil.ReadFile(fp.Name())
	assert.NoError(t, err)
	lines := strings.Split(string(raw), "\n")
	assert.Equal(t, 120001, len(lines))
	assert.Equal(t, "AAAAAAAAAAAAAAAAAAAAAAAA", lines[119991])
	assert.Equal(t, "BBBBBBBBBBBBBBBBBBBBBBBB", lines[119995])
	assert.Equal(t, "CCCCCCCCCCCCCCCCCCCCCCCC", lines[119999])
}
