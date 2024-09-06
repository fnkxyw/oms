package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

// запись в JSON
func WriteToJSON[T any](path string, val T) error {

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("OpenFile eror in WriteToJSON", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent(" ", "  ")
	if err := encoder.Encode(val); err != nil {
		fmt.Println("Encoding error in WirteToJSON", err)
		return err
	}

	return nil
}
