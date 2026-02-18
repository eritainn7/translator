package semantics

import (
	"fmt"
	. "translator/lexical/bases/model"
)

type semanticHandler7 struct {
	buffer string
}

func InitSemanticHandler7(buffer string) *semanticHandler7 {
	return &semanticHandler7{
		buffer: buffer,
	}
}

func (instance *semanticHandler7) Run(current_symbol string) (string, error) {
	instance.buffer += current_symbol
	code, _ := FindCodeInCSV(Operations, instance.buffer)
	if code == -1 {
		err := fmt.Errorf("Not exist in operations")
		return "", err
	}

	return mark(code, Operations), nil
}
