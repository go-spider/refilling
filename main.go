package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"
)

func JetBrianDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home +"\\AppData\\Roaming\\JetBrains\\"
	}
	return os.Getenv("HOME")+"/.config/JetBrains/"
}

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
	flag.StringVar(&filename, "f", "", "f不为空在当前目录生成evaluation.key，否则批量替换key(仅适用于Windows&Linux)")
	flag.Parse()
	if filename != ""{
		err := wk.writeKey(filename)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		fmt.Println(filename)
	}else{
		IntelliPath := filepath.FromSlash(JetBrianDir())
		fileInfoList,err := ioutil.ReadDir(IntelliPath)
		if err != nil {
			log.Fatal(err)
		}
		reg := regexp.MustCompile(`.?(IntelliJIdea|GoLand|CLion|PyCharm|DataGrip|RubyMine|AppCode|PhpStorm|WebStorm|Rider|idea).*`)
		for i := range fileInfoList {
			if fileInfoList[i].IsDir(){
				dirName := fileInfoList[i].Name()
				if reg.Match([]byte(dirName)){
					evalPath := filepath.Join(IntelliPath+dirName,"eval")
					keyList,err := ioutil.ReadDir(evalPath)
					if err != nil {
						log.Fatal(err)
					}
					for j :=range keyList{
						if !keyList[j].IsDir(){
							keyName := keyList[j].Name()
							if reg.Match([]byte(keyName)){
								keyPath := filepath.Join(evalPath,keyName)
								fmt.Println(keyPath)
								err := wk.writeKey(keyPath)
								if err != nil {
									fmt.Printf("Error: %s\n", err)
								}
							}
						}
					}
				}
			}
		}
	}
}

