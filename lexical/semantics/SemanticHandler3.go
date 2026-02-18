package semantics

import . "translator/lexical/bases/model"

type semanticHandler3 struct {
	buffer string
}

func InitSemanticHandler3(buffer string) *semanticHandler3 {
	return &semanticHandler3{
		buffer: buffer,
	}
}

func (instance *semanticHandler3) Run() string {
	AddCodeInCSV(NumConsts, instance.buffer)
	code, _ := FindCodeInCSV(NumConsts, instance.buffer)

	return mark(code, NumConsts)
}
