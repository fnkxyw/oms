package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// запись в json с перезаписью данных
func WriteToJSON[T any](path string, val T) error {
	var currentData []T

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return err
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return err
	}

	if len(fileData) > 0 {
		if err := json.Unmarshal(fileData, &currentData); err != nil {
			fmt.Println("Unmurshal error")
			return err
		}
	}

	currentData = append(currentData, val)

	file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Ошибка при открытии файла для записи:", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(currentData); err != nil {
		fmt.Println("Ошибка при кодировании в JSON:", err)
		return err
	}

	return nil
}

func ReadFromJSON(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Open file error:", err)
		return err
	}

	defer file.Close()
	return err
}
