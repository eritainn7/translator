package tables

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// Структуры для таблиц
type Tables struct {
	ServiceWords map[string]int
	Delimiters   map[string]int
	Operations   map[string]int
	Identifiers  map[string]int
	Constants    map[string]int

	nextIdentifierID int
	nextConstantID   int

	dataPath string
}

// Создание нового экземпляра таблиц
func NewTables(dataPath string) (*Tables, error) {
	t := &Tables{
		ServiceWords:     make(map[string]int),
		Delimiters:       make(map[string]int),
		Operations:       make(map[string]int),
		Identifiers:      make(map[string]int),
		Constants:        make(map[string]int),
		nextIdentifierID: 1,
		nextConstantID:   1,
		dataPath:         dataPath,
	}

	// Загрузка таблиц из CSV
	if err := t.loadServiceWords(); err != nil {
		return nil, fmt.Errorf("error loading service words: %v", err)
	}

	if err := t.loadDelimiters(); err != nil {
		return nil, fmt.Errorf("error loading delimiters: %v", err)
	}

	if err := t.loadOperations(); err != nil {
		return nil, fmt.Errorf("error loading operations: %v", err)
	}

	// Загрузка существующих идентификаторов если есть
	t.loadIdentifiers()

	return t, nil
}

// Загрузка служебных слов
func (t *Tables) loadServiceWords() error {
	file, err := os.Open(t.dataPath + "/service_words.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Пропускаем заголовок
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) >= 2 {
			id, _ := strconv.Atoi(record[0])
			t.ServiceWords[record[1]] = id
		}
	}

	return nil
}

// Загрузка разделителей
func (t *Tables) loadDelimiters() error {
	file, err := os.Open(t.dataPath + "/delimiters.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) >= 2 {
			id, _ := strconv.Atoi(record[0])
			t.Delimiters[record[1]] = id
		}
	}

	return nil
}

// Загрузка операций
func (t *Tables) loadOperations() error {
	file, err := os.Open(t.dataPath + "/operations.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) >= 2 {
			id, _ := strconv.Atoi(record[0])
			t.Operations[record[1]] = id
		}
	}

	return nil
}

// Загрузка идентификаторов из файла
func (t *Tables) loadIdentifiers() {
	file, err := os.Open(t.dataPath + "/identifiers.csv")
	if err != nil {
		// Файл может не существовать при первом запуске
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return
	}

	maxID := 0
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) >= 2 {
			id, _ := strconv.Atoi(record[0])
			t.Identifiers[record[1]] = id
			if id > maxID {
				maxID = id
			}
		}
	}

	t.nextIdentifierID = maxID + 1
}

// Сохранение идентификаторов в файл
func (t *Tables) saveIdentifiers() error {
	file, err := os.Create(t.dataPath + "/identifiers.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Записываем заголовок
	if err := writer.Write([]string{"id", "identifier"}); err != nil {
		return err
	}

	// Записываем все идентификаторы
	for ident, id := range t.Identifiers {
		if err := writer.Write([]string{strconv.Itoa(id), ident}); err != nil {
			return err
		}
	}

	return nil
}

// Сохранение констант в файл
func (t *Tables) saveConstants() error {
	file, err := os.Create(t.dataPath + "/constants.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Записываем заголовок
	if err := writer.Write([]string{"id", "constant"}); err != nil {
		return err
	}

	// Записываем все константы
	for constVal, id := range t.Constants {
		if err := writer.Write([]string{strconv.Itoa(id), constVal}); err != nil {
			return err
		}
	}

	return nil
}

// Получение ID идентификатора
func (t *Tables) GetIdentifierID(name string) int {
	if id, exists := t.Identifiers[name]; exists {
		return id
	}
	id := t.nextIdentifierID
	t.Identifiers[name] = id
	t.nextIdentifierID++

	// Сохраняем в файл при каждом новом идентификаторе
	t.saveIdentifiers()

	return id
}

// Получение ID константы
func (t *Tables) GetConstantID(value string) int {
	if id, exists := t.Constants[value]; exists {
		return id
	}
	id := t.nextConstantID
	t.Constants[value] = id
	t.nextConstantID++

	// Сохраняем в файл при каждой новой константе
	t.saveConstants()

	return id
}

// Получение статистики
func (t *Tables) GetStats() map[string]int {
	return map[string]int{
		"service_words": len(t.ServiceWords),
		"delimiters":    len(t.Delimiters),
		"operations":    len(t.Operations),
		"identifiers":   len(t.Identifiers),
		"constants":     len(t.Constants),
	}
}
