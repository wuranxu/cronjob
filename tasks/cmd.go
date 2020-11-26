package tasks

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
)

func Command(ctx context.Context, cmd string) error {
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
				return
			default:
				readString, err := reader.ReadString('\n')
				if err != nil || err == io.EOF {
					break
				}
				fmt.Print(readString)
			}
		}
	}()
	return c.Run()
}
