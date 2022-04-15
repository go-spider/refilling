package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	var filename string
	flag.StringVar(&filename, "f", "evaluation.key", "默认,evaluation.key")
	flag.Parse()
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
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 666)
	if err != nil {
		fmt.Println(err)
	}
	n, err := f.Write(writeBuffer)
	if err == nil && n < len(writeBuffer) {
		err = io.ErrShortWrite
	}
	if err != nil {
		fmt.Println(err)
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	if err != nil {
		fmt.Println(err)
	}
}
