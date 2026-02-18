package semantics

import . "translator/lexical/bases/model"

type semanticHandler6 struct {
	buffer string
}

func InitSemanticHandler6(buffer string) *semanticHandler6 {
	return &semanticHandler6{
		buffer: buffer,
	}
}

func (instance *semanticHandler6) Run() string {
	code, _ := FindCodeInCSV(Operations, instance.buffer)

	return mark(code, Operations)
}
