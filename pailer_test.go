package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func testData(data chan *ReadData, chunks int) {
	lines := [][]byte{
		[]byte(`AAAAAAAAAAAAAAAAAAAAAAAA`),
		[]byte(`AAAAAAAAAAAAAAAAAAAAAAAA`),
		[]byte(`AAAAAAAAAAAAAAAAAAAAAAAA`),
		[]byte(`AAAAAAAAAAAAAAAAAAAAAAAA`),
		[]byte(`BBBBBBBBBBBBBBBBBBBBBBBB`),
		[]byte(`BBBBBBBBBBBBBBBBBBBBBBBB`),
		[]byte(`BBBBBBBBBBBBBBBBBBBBBBBB`),
		[]byte(`BBBBBBBBBBBBBBBBBBBBBBBB`),
		[]byte(`CCCCCCCCCCCCCCCCCCCCCCCC`),
		[]byte(`CCCCCCCCCCCCCCCCCCCCCCCC`),
		[]byte(`CCCCCCCCCCCCCCCCCCCCCCCC`),
		[]byte(`CCCCCCCCCCCCCCCCCCCCCCCC`),
	}
	for i := 0; i < chunks; i++ {
		rd := &ReadData{}
		rd.Data = bytes.Join(lines, []byte{byte('\n')})
		rd.Data = append(rd.Data, byte('\n'))
		data <- rd
	}
	close(data)
}

func TestPailer(t *testing.T) {
	data := make(chan *ReadData)
	cancel := make(chan bool)
	go testData(data, 10)
	fp, err := ioutil.TempFile("/tmp", "mesos-cli")
	assert.NoError(t, err)
	assert.NoError(t, Pailer(data, cancel, 0, fp))
	assert.NoError(t, fp.Close())
	raw, err := ioutil.ReadFile(fp.Name())
	assert.NoError(t, err)
	//lines := strings.Split(string(raw), "\n")
	lines := bytes.Split(raw, []byte{byte('\n')})
	assert.Equal(t, 121, len(lines))
	assert.Equal(t, []byte(`AAAAAAAAAAAAAAAAAAAAAAAA`), lines[111])
	assert.Equal(t, []byte(`BBBBBBBBBBBBBBBBBBBBBBBB`), lines[115])
	assert.Equal(t, []byte(`CCCCCCCCCCCCCCCCCCCCCCCC`), lines[119])
	fmt.Printf("%s", lines)
}

func TestPailerHuge(t *testing.T) {
	data := make(chan *ReadData)
	cancel := make(chan bool)
	go testData(data, 10000)
	fp, err := ioutil.TempFile("/tmp", "mesos-cli")
	assert.NoError(t, err)
	assert.NoError(t, Pailer(data, cancel, 0, fp))
	assert.NoError(t, fp.Close())
	raw, err := ioutil.ReadFile(fp.Name())
	assert.NoError(t, err)
	lines := bytes.Split(raw, []byte{byte('\n')})
	assert.Equal(t, 120001, len(lines))
	assert.Equal(t, []byte(`AAAAAAAAAAAAAAAAAAAAAAAA`), lines[119991])
	assert.Equal(t, []byte(`BBBBBBBBBBBBBBBBBBBBBBBB`), lines[119995])
	assert.Equal(t, []byte(`CCCCCCCCCCCCCCCCCCCCCCCC`), lines[119999])
}
