package semantics

import (
	. "translator/lexical/bases/model"
)

type semanticHandler1 struct {
	buffer string
}

func InitSemanticHandler1(buffer string) *semanticHandler1 {
	return &semanticHandler1{
		buffer: buffer,
	}
}

func (instance *semanticHandler1) Run() string {
	code, _ := FindCodeInCSV(Identificators, instance.buffer)
	if code == -1 {
		AddCodeInCSV(Identificators, instance.buffer)
	}
	code, _ = FindCodeInCSV(Identificators, instance.buffer)
	mark_code := mark(code, Identificators)

	return mark_code
}
