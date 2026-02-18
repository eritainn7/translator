package semantics

import . "translator/lexical/bases/model"

type semanticHandler2 struct {
	buffer string
}

func InitSemanticHandler2(buffer string) *semanticHandler2 {
	return &semanticHandler2{
		buffer: buffer,
	}
}

func (instance *semanticHandler2) Run() string {
	code, _ := FindCodeInCSV(SystemWord, instance.buffer)
	if code == -1 {
		return InitSemanticHandler1(instance.buffer).Run()
	}

	mark_code := mark(code, SystemWord)
	return mark_code
}
