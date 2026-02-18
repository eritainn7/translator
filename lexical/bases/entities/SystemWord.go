package entities

type systemWordBase struct {
	entity
}

func InitSystemWordBase() *systemWordBase {
	instance := systemWordBase{entity: entity{mark: "W"}}

	return &instance
}

func (base *systemWordBase) FindCode() (string, error) {
	code, err := base.findCodeInCSV(endpoints["SystemWordsBase"])

	if err != nil {
		return "error", err
	}

	code = base.mark + code
	return code, nil

}
