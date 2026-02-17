package lexical

import "fmt"

type state string

func (s state) String() (*string, error) {
	switch s {
	case "S":
	case "F":
	case "Z":
	case "0":
	case "1":
	case "2":
	case "3":
	case "4":
	case "5":
	case "6":
	case "7":
	case "8":
	case "9":
	default:
		return nil, fmt.Errorf("Wrong value")
	}

	sym := string(s)
	return &sym, nil
}
