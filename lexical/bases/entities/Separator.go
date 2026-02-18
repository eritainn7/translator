package entities

type separatorsBase struct {
	entity
}

func InitSeparatorsBase() *separatorsBase {
	instance := separatorsBase{entity: entity{mark: "R"}}

	return &instance
}

func (base *separatorsBase) FindCode() (string, error) {
	code, err := base.findCodeInCSV(endpoints["SeparatorsBase"])

	if err != nil {
		return "error", err
	}

	code = base.mark + code
	return code, nil
}
