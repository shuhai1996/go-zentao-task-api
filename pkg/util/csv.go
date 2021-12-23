package util

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"os"
)

// ReadCsv 通过 channel 遍历csv
func ReadCsv(fileName string, c chan interface{}) {
	str, _ := os.Getwd() //获取当前目录
	fs, err := os.Open(str + "/" + fileName)
	if err != nil {
		log.Fatalf("can not open the file, err is %+v\n", err)
	}
	defer fs.Close()
	r := csv.NewReader(fs)
	//针对大文件，一行一行的读取文件
	for {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			log.Fatalf("can not read, err is %+v", err)
		}
		if err == io.EOF {
			break
		}
		c <- row //通过chan c 传递
	}
	close(c) //遍历结束关闭通道
}

func ReadAll(fileName string) [][]string {
	// 针对小文件，也可以一次性读取所有的文件
	fs1, _ := os.Open(fileName)
	r1 := csv.NewReader(fs1)
	content, er := r1.ReadAll()
	if er != nil {
		log.Fatalf("can not readall, err is %+v", er)
	}
	return content
}

func LineCounter(fileName string) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}
	str, _ := os.Getwd() //获取当前目录
	fs, err := os.Open(str + "/" + fileName)
	if err != nil {
		log.Fatalf("can not open the file, err is %+v\n", err)
	}
	defer fs.Close()
	fileReader := bufio.NewReader(fs)
	for {
		c, err := fileReader.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func ReadCsv2(fileName string) [][]string {
	str, _ := os.Getwd() //获取当前目录
	fs, err := os.Open(str + "/" + fileName)
	var content [][]string
	if err != nil {
		log.Fatalf("can not open the file, err is %+v\n", err)
	}
	defer fs.Close()
	r := csv.NewReader(fs)
	//针对大文件，一行一行的读取文件
	for {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			log.Fatalf("can not read, err is %+v", err)
		}
		if err == io.EOF {
			break
		}
		content = append(content, row)
	}
	return content
}
