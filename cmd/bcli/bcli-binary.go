package main
import (
	"os"
	//"time"
	"encoding/binary"
	//"math/rand"
	"fmt"
	"strconv"
	"bytes"
)

func writeFileToBinary(binary_strings string, to_file_path string) {
	if len(binary_strings) <= 0 {
		return
	}
	
	file, err := os.Create(to_file_path)
	defer file.Close()
	if err != nil {
		fmt.Println("ERROR!")
	}

	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//fmt.Println(r)
	
	for i := 0; i < len(binary_strings); i+=2 {
		value, _ :=  strconv.ParseInt(binary_strings[i:i+2], 16, 32);
		var bin_buf bytes.Buffer
		m := uint8(value)
		binary.Write(&bin_buf,binary.BigEndian, m)
		//b :=bin_buf.Bytes()
		//l := len(b)
		//fmt.Println(l)
		writeNextBytes(file, bin_buf.Bytes())

	}
}

func writeNextBytes(file *os.File, bytes []byte) {

	_, err := file.Write(bytes)

	if err != nil {
		fmt.Println("Error when try writeNextBytes!")
	}
}
