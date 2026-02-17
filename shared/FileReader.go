package shared

import (
	"fmt"
	"io"
	"os"
)

type fileReader struct {
	path_to_file string
}

func InitFileReader(path_to_file string) fileReader {
	instance_fileReader := fileReader{path_to_file}

	return instance_fileReader
}

func (instance_fileReader fileReader) GetTextFile() string {
	file, err := os.Open(instance_fileReader.path_to_file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	data := make([]byte, 64)
	result := ""
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		result += string(data[:n])
	}

	return result
}
