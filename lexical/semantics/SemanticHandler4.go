package semantics

import . "translator/lexical/bases/model"

type semanticHandler4 struct {
	buffer string
}

func InitSemanticHandler4(buffer string) *semanticHandler4 {
	return &semanticHandler4{
		buffer: buffer,
	}
}

func (instance *semanticHandler4) Run() string {
	code, _ := FindCodeInCSV(Separators, instance.buffer)

	return mark(code, Separators)
}
