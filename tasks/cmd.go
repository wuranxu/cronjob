package tasks

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

var StreamData = Stream{}

type Stream struct {
	sync.Map
}

func (s *Stream) Send(jobId string, output *string) {
	ch, ok := s.Load(jobId)
	if !ok {
		ch = make(chan *string, 10)
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

func Command(ctx context.Context, jobId uint, cmd string, websocket bool) error {
	c := exec.CommandContext(ctx, "bash", "-c", cmd)
	//c := exec.CommandContext(ctx, "cmd", "/C", cmd)
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	go func() {
		reader := bufio.NewReader(stdout)
		for {
			select {
			case <-ctx.Done():
				if ctx.Err() != nil {
					fmt.Printf("程序出现错误: %q", ctx.Err())
				} else {
					fmt.Println("程序被终止")
				}
				StreamData.Send(fmt.Sprintf("%d", jobId), nil)
				return
			default:
				readString, err := reader.ReadString('\n')
				if err != nil || err == io.EOF {
					break
				}
				if websocket {
					StreamData.Send(fmt.Sprintf("%d", jobId), &readString)
				}
				fmt.Print(readString)
			}
		}
	}()
	return c.Run()
}
