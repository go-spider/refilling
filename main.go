package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

type WriteFile struct {
	KeyContent []byte
}


func (w *WriteFile)writeKey(filename string) error{
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 666)
	if err != nil {
		return err
	}
	n, err := f.Write(w.KeyContent)
	if err == nil && n < len(w.KeyContent) {
		err = io.ErrShortWrite
	}
	if err != nil {
		return err
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var writeBuffer = make([]byte,8)
	timeMill:=-((time.Now().UnixNano() / int64(time.Millisecond))+1)
	writeBuffer[0] = byte(uint64(timeMill) >> 56)
	writeBuffer[1] = byte(uint64(timeMill) >> 48)
	writeBuffer[2] = byte(uint64(timeMill) >> 40)
	writeBuffer[3] = byte(uint64(timeMill) >> 32)
	writeBuffer[4] = byte(uint64(timeMill) >> 24)
	writeBuffer[5] = byte(uint64(timeMill) >> 16)
	writeBuffer[6] = byte(uint64(timeMill) >> 8)
	writeBuffer[7] = byte(uint64(timeMill) >> 0)
	wk := WriteFile{KeyContent: writeBuffer}
	var filename string
	flag.StringVar(&filename, "f", "", "f不为空在当前目录生成evaluation.key，否则按config文件路径批量替换key")
	flag.Parse()
	if filename!=""{
		err := wk.writeKey(filename)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		fmt.Println(filename)
	}else{
		f, err := os.Open("config")
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		defer func() {
			_ = f.Close()
		}()
		br := bufio.NewReader(f)
		for {
			line, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			err := wk.writeKey(string(line))
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			fmt.Println(string(line))
		}
	}
	}

