package storage

import (
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/akugnerevich/homework-1.git/internal/models"
	"io"
	"os"
)

type ReturnStorage struct {
	Data map[uint]*models.Return
}

func NewReturnStorage() *ReturnStorage {
	return &ReturnStorage{Data: make(map[uint]*models.Return)}
}

func (rs *ReturnStorage) AddReturnToStorage(r *models.Return) error {
	_, ok := rs.Data[r.ID]
	if ok {
		return fmt.Errorf("Order already return ")
	} else {
		rs.Data[r.ID] = r
	}

	return nil
}

func (r *ReturnStorage) DeleteOrderFromStorage(id uint) {
	delete(r.Data, id)
}

// проверка на наличие
func (o *ReturnStorage) IsConsist(id uint) bool {
	_, ok := o.Data[id]
	return ok
}

// считываем с JSON-a
func (o *ReturnStorage) ReadFromJSON(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("Open file erorr: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Read file error: %w", err)
	}

	if len(data) == 0 {
		o.Data = make(map[uint]*models.Return)
		return nil
	}

	var i struct {
		Data map[uint]models.Return `json:"Data"`
	}

	err = json.Unmarshal(data, &i)
	if err != nil {
		return fmt.Errorf("ошибка при декодировании JSON: %w", err)
	}

	o.Data = make(map[uint]*models.Return)
	for returnid, r := range i.Data {
		returnCopy := r
		o.Data[returnid] = &returnCopy
	}

	return nil
}

func (o *ReturnStorage) WritoToJSON() error {
	file, err := os.OpenFile("api/returns.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("OpenFile eror in WriteToJSON", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent(" ", "  ")
	if err := encoder.Encode(o); err != nil {
		fmt.Println("Encoding error in WirteToJSON", err)
		return err
	}
	return nil
}
