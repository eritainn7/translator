package entities

type operationsBase struct {
	entity
}

func InitOperationsBase() *operationsBase {
	instance := operationsBase{entity: entity{mark: "O"}}

	return &instance
}

func (base *operationsBase) FindCode() (string, error) {
	code, err := base.findCodeInCSV(endpoints["OperationsBase"])

	if err != nil {
		return "error", err
	}

	code = base.mark + code
	return code, nil

}
