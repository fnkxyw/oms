package returnStorage

import (
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"reflect"
	"sort"
	"testing"
)

func TestNewReturnStorage(t *testing.T) {
	want := &ReturnStorage{
		Data: make(map[uint]*models.Return),
		path: "api/return.json",
	}
	if got := NewReturnStorage(); !reflect.DeepEqual(got, want) {
		t.Errorf("NewReturnStorage() = %v, want %v", got, want)
	}
}

func TestReturnStorage_AddReturnToStorage(t *testing.T) {
	rs := NewReturnStorage()
	r := &models.Return{
		ID:     1,
		UserID: 123,
	}

	err := rs.AddReturnToStorage(r)
	if err != nil {
		t.Errorf("AddReturnToStorage() error = %v, wantErr %v", err, false)
	}

	if got, exists := rs.Data[1]; !exists || got != r {
		t.Errorf("AddReturnToStorage() failed, expected return ID %v to be added", r.ID)
	}
}

func TestReturnStorage_DeleteReturnFromStorage(t *testing.T) {
	rs := NewReturnStorage()
	r := &models.Return{
		ID:     1,
		UserID: 123,
	}

	err := rs.AddReturnToStorage(r)
	if err != nil {
		t.Errorf("DeleteReturnFromStorage() failed , AddReturnToStorage")
	}
	rs.DeleteReturnFromStorage(1)

	if _, exists := rs.Data[1]; exists {
		t.Errorf("DeleteReturnFromStorage() failed, expected return ID 1 to be deleted")
	}
}

func TestReturnStorage_IsConsist(t *testing.T) {
	rs := NewReturnStorage()
	r := &models.Return{
		ID:     1,
		UserID: 123,
	}

	err := rs.AddReturnToStorage(r)
	if err != nil {
		t.Errorf("IsConsist() failed , AddReturnToStorage")
	}

	if !rs.IsConsist(1) {
		t.Errorf("IsConsist() failed, expected return with ID 1 to exist")
	}

	if rs.IsConsist(2) {
		t.Errorf("IsConsist() failed, expected return with ID 2 to not exist")
	}
}

func TestReturnStorage_GetReturn(t *testing.T) {
	rs := NewReturnStorage()
	r := &models.Return{
		ID:     1,
		UserID: 123,
	}

	err := rs.AddReturnToStorage(r)
	if err != nil {
		t.Errorf("GetReturn() failed , AddReturnToStorage")
	}

	got, exists := rs.GetReturn(1)
	if !exists || !reflect.DeepEqual(got, r) {
		t.Errorf("GetReturn() got = %v, want %v", got, r)
	}

	if _, exists = rs.GetReturn(2); exists {
		t.Errorf("GetReturn() expected return with ID 2 to not exist")
	}
}

func TestReturnStorage_GetReturnIDs(t *testing.T) {
	rs := NewReturnStorage()
	err := rs.AddReturnToStorage(&models.Return{ID: 1, UserID: 123})
	if err != nil {
		t.Errorf("GetReturnIDs() failed, AddReturnToStorage")
	}
	err = rs.AddReturnToStorage(&models.Return{ID: 2, UserID: 456})
	if err != nil {
		t.Errorf("GetReturnIDs() failed, AddReturnToStorage")
	}

	got := rs.GetReturnIDs()
	sort.Slice(got, func(i, j int) bool {
		return got[i] < got[j]
	})
	want := []uint{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetReturnIDs() = %v, want %v", got, want)
	}

}