package semantics

type semanticHandler1 struct {
	current_sequence_leksems string
}

func InitSemanticHandler1(sequence_leksems string) *semanticHandler1 {
	return &semanticHandler1{
		current_sequence_leksems: sequence_leksems,
	}
}

func (instance *semanticHandler1) Run() string {

	return instance.current_sequence_leksems
}
