package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/url"
	"sync"
	"time"
)

const PollInterval = 100 * time.Millisecond

/*
FilePaginator paginates requests to a /file/read endpoint
sending the response as fileData on the data channel.
If max is a negative number paginate forever or until
a message is received at f.cancel. Otherwise find the
length of the file and paginate until we reach that position.
*/
type FilePaginator struct {
	data   chan *fileData // Results
	cancel chan bool      // Cancel the pagination
	offset int            // Current offset
	max    int            // Maximum offset
	tail   bool           // Begin at the end
	path   string         // Path to file
}

func (f *FilePaginator) init(c *Client) error {
	u := &url.URL{
		Path: "/files/read",
		RawQuery: url.Values{
			"path": []string{f.path},
		}.Encode(),
	}
	fd := &fileData{}
	if err := c.Get(u, fd); err != nil {
		return err
	}
	// If we are tailing output we start
	// at the end of the file. Since it
	// is impossible to tell how many
	// of the previous bytes equate to
	// a single line we just start at the
	// end and wait for more data.
	if f.tail {
		f.offset = fd.Offset
	} else {
		// Signal the end of this file once
		// we reach the total offset.
		f.max = fd.Offset
	}
	return nil
}

func (f *FilePaginator) Next(c *Client) error {
	select {
	case <-f.cancel:
		return ErrEndPagination
	default:
	}
	u := &url.URL{
		Path: "/files/read",
		RawQuery: url.Values{
			"offset": []string{fmt.Sprintf("%d", f.offset)},
			"path":   []string{f.path},
		}.Encode(),
	}
	fd := &fileData{}
	if err := c.Get(u, fd); err != nil {
		return err
	}
	f.offset += fd.Length()
	f.data <- fd
	// If max == 0 this is a tail operation
	if f.max > 0 {
		if f.max == f.offset {
			// We've reached the end of the file
			return ErrMaxExceeded
		}
	}
	// If this is not a tail operation and no
	// data was returned break pagination.
	if !f.tail && fd.Length() == 0 {
		return ErrEndPagination
	}
	time.Sleep(PollInterval)
	return nil
}

func (f *FilePaginator) Close() { close(f.data) }

// Pailer reads until n lines
// TODO: Add support for binary data
func Pailer(data <-chan *fileData, cancel chan bool, n int, w io.Writer) error {
	writer := bufio.NewWriter(w)
	var (
		buf   bytes.Buffer
		err   error
		lines int
	)
loop:
	for fd := range data {
		_, err = buf.WriteString(fd.Data)
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

func Monitor(c *Client, w io.Writer, lines int, pag *FilePaginator) (err error) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		err = c.PaginateFile(pag)
	}()
	go func() {
		defer wg.Done()
		err = Pailer(pag.data, pag.cancel, lines, w)
	}()
	wg.Wait()
	return err
}
