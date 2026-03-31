package shared

import (
	"encoding/csv"
	"fmt"
	"lexer/models"
	"os"
)

func SaveLexemesToCSV(filename string, lexemes []models.Lexeme) error {
	// Создаем файл
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	// Создаем CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Записываем заголовок (на русском или английском - по желанию)
	header := []string{"Код лексемы", "Вид лексемы", "Значение лексемы"}
	// Альтернативный заголовок на английском:
	// header := []string{"Lexeme Code", "Lexeme Type", "Lexeme Value"}

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("ошибка записи заголовка: %v", err)
	}

	// Преобразуем тип лексемы из кода в человекочитаемый вид
	typeNames := map[byte]string{
		'W': "Служебное слово",
		'I': "Идентификатор",
		'C': "Константа",
		'R': "Разделитель",
		'O': "Операция",
		'E': "Конец файла",
		'F': "Ошибка",
	}

	// Записываем каждую лексему
	for _, lex := range lexemes {
		// Получаем человекочитаемое название типа
		typeName, exists := typeNames[lex.Type]
		if !exists {
			typeName = fmt.Sprintf("Неизвестный тип (%c)", lex.Type)
		}

		// Формируем полный код лексемы (например: "W73" или "I1")
		lexemeCode := fmt.Sprintf("%c%d", lex.Type, lex.Index)

		// Создаем запись
		record := []string{
			lexemeCode,
			typeName,
			lex.Value,
		}

		// Записываем в CSV
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("ошибка записи лексемы: %v", err)
		}
	}

	return nil
}
