package pailer

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/mesos/mesos-go/api/v1/lib/agent"
	"github.com/mesos/mesos-go/api/v1/lib/agent/calls"
	"github.com/mesos/mesos-go/api/v1/lib/httpcli"
	"io"
	"sync"
	"time"
)

const PollInterval = 100 * time.Millisecond

var (
	ErrMaxExceeded   = errors.New("max exceeded")
	ErrEndPagination = errors.New("no more items to paginate")
)

/*
FilePaginator paginates requests to the HTTP operator
endpoint sending the response as fileData on the data channel.
If max is a negative number paginate forever or until
a message is received at f.cancel. Otherwise find the
length of the file and paginate until we reach that position.
*/
type FilePaginator struct {
	data   chan *agent.Response_ReadFile // Results
	cancel chan bool                     // Cancel the pagination
	offset uint64                        // Current offset
	Max    uint64                        // Maximum offset
	Follow bool                          // Begin at the end
	Path   string                        // Path to file
}

func (f *FilePaginator) init(caller *httpcli.Client) error {
	caller.Send
	resp, err := caller.CallAgent(calls.ReadFileWithLength(f.Path, uint64(0), uint64(0)))
	if err != nil {
		return err
	}
	fd := resp.GetReadFile()
	// If we are tailing output we start
	// at the end of the file. Since it
	// is impossible to tell how many
	// of the previous bytes equate to
	// a single line we just start at the
	// end and wait for more data.
	if f.Follow {
		f.offset = fd.Size_
	} else {
		// Signal the end of this file once
		// we reach the total offset.
		f.Max = fd.Size_
	}
	return nil
}

func (f *FilePaginator) Next(caller *httpcli.Client) error {
	select {
	case <-f.cancel:
		return ErrEndPagination
	default:
	}
	resp, err := caller.CallAgent(calls.ReadFile(f.Path, uint64(50000), f.offset))
	if err != nil {
		return err
	}
	fd := resp.GetReadFile()
	f.offset += fd.Size_
	//f.offset += fd.Length()
	f.data <- fd
	// If max == 0 this is a tail operation
	if f.Max > 0 {
		if f.Max == f.offset {
			// We've reached the end of the file
			return ErrMaxExceeded
		}
	}
	// If this is not a tail operation and no
	// data was returned break pagination.
	if !f.Follow && fd.Size() == 0 {
		return ErrEndPagination
	}
	time.Sleep(PollInterval)
	return nil
}

func (f *FilePaginator) Close() { close(f.data) }

// Pailer reads until n lines
func Pailer(data <-chan *agent.Response_ReadFile, cancel chan bool, n int, w io.Writer) error {
	writer := bufio.NewWriter(w)
	var (
		buf   bytes.Buffer
		err   error
		lines int
	)
loop:
	for fd := range data {
		_, err = buf.Write(fd.Data)
		if err != nil {
			break loop
		}
		for {
			line, err := buf.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					buf.Write(line)
					break
				} else {
					break loop
				}
			}
			_, err = writer.Write(line)
			if err != nil {
				break loop
			}
			writer.Flush()
			lines++
			if n > 0 {
				if lines >= n {
					cancel <- true
					break loop
				}
			}
		}
		if buf.Len() > 0 {
			writer.Write(buf.Bytes())
			writer.Flush()
		}
	}
	return err
}

func Monitor(caller *httpcli.Client, w io.Writer, lines int, pag *FilePaginator) (err error) {
	if err = pag.init(caller); err != nil {
		return err
	}
	pag.data = make(chan *agent.Response_ReadFile)
	pag.cancel = make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for err == nil {
			err = pag.Next(caller)
		}
		pag.Close()
	}()
	go func() {
		defer wg.Done()
		err = Pailer(pag.data, pag.cancel, lines, w)
	}()
	wg.Wait()
	return err
}
