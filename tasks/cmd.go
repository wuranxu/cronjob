package tasks

import (
	"bufio"
	"context"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os/exec"
	"sync"
)

var StreamData = Stream{}

type Stream struct {
	sync.Map
}

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func (s *Stream) Send(jobId string, output *string) {
	ch, ok := s.Load(jobId)
	if !ok {
		ch = make(chan *string)
		s.Store(jobId, ch)
	}
	strings := ch.(chan *string)
	strings <- output
}

func (s *Stream) Close(jobId string) {
	load, ok := s.Load(jobId)
	if !ok {
		return
	}
	strings := load.(chan *string)
	close(strings)
	s.Delete(jobId)
}

func (s *Stream) Read(jobId string) chan *string {
	ch, ok := s.Load(jobId)
	if !ok {
		return nil
	}
	return ch.(chan *string)
}

func reader(ctx context.Context, wg *sync.WaitGroup, jobId uint, stdout io.Reader, websocket bool) {
	reader := bufio.NewReader(stdout)
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			StreamData.Send(fmt.Sprintf("%d", jobId), nil)
			return
		default:
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			byte2String := ConvertByte2String([]byte(readString), "GB18030")
			if websocket {
				fmt.Println(byte2String)
				StreamData.Send(fmt.Sprintf("%d", jobId), &byte2String)
			}
		}
	}
}

func SyncCommand(ctx context.Context, cmd string) error {
	c := exec.CommandContext(ctx, "cmd", "/C", cmd)
	output, err := c.CombinedOutput()
	fmt.Println(string(output))
	return err
}

func Command(ctx context.Context, jobId uint, cmd string, websocket bool) error {
	c := exec.CommandContext(ctx, "cmd", "/C", cmd)
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go reader(ctx, &wg, jobId, stdout, websocket)
	go reader(ctx, &wg, jobId, stderr, websocket)
	err = c.Start()
	wg.Wait()
	StreamData.Send(fmt.Sprintf("%d", jobId), nil)
	return err
}

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
