package pailer

import (
	"bytes"
	"fmt"
	"github.com/mesos/mesos-go/agent"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func testData(data chan *agent.Response_ReadFile, chunks int) {
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
		raw := bytes.Join(lines, []byte{byte(rune('\n'))})
		raw = append(raw, byte(rune('\n')))
		data <- &agent.Response_ReadFile{Data: raw}
	}
	close(data)
}

func TestPailer(t *testing.T) {
	data := make(chan *agent.Response_ReadFile)
	cancel := make(chan bool)
	go testData(data, 10)
	fp, err := ioutil.TempFile("/tmp", "mesos-cli")
	assert.NoError(t, err)
	assert.NoError(t, Pailer(data, cancel, 0, fp))
	assert.NoError(t, fp.Close())
	raw, err := ioutil.ReadFile(fp.Name())
	assert.NoError(t, err)
	lines := bytes.Split(raw, []byte{byte(rune('\n'))})
	assert.Equal(t, 121, len(lines))
	assert.Equal(t, []byte(`AAAAAAAAAAAAAAAAAAAAAAAA`), lines[111])
	assert.Equal(t, []byte(`BBBBBBBBBBBBBBBBBBBBBBBB`), lines[115])
	assert.Equal(t, []byte(`CCCCCCCCCCCCCCCCCCCCCCCC`), lines[119])
	for _, line := range lines {
		fmt.Println(string(line))
	}
}

func TestPailerHuge(t *testing.T) {
	data := make(chan *agent.Response_ReadFile)
	cancel := make(chan bool)
	go testData(data, 10000)
	fp, err := ioutil.TempFile("/tmp", "mesos-cli")
	assert.NoError(t, err)
	assert.NoError(t, Pailer(data, cancel, 0, fp))
	assert.NoError(t, fp.Close())
	raw, err := ioutil.ReadFile(fp.Name())
	assert.NoError(t, err)
	lines := bytes.Split(raw, []byte{byte(rune('\n'))})
	assert.Equal(t, 120001, len(lines))
	assert.Equal(t, []byte(`AAAAAAAAAAAAAAAAAAAAAAAA`), lines[119991])
	assert.Equal(t, []byte(`BBBBBBBBBBBBBBBBBBBBBBBB`), lines[119995])
	assert.Equal(t, []byte(`CCCCCCCCCCCCCCCCCCCCCCCC`), lines[119999])
}
