package benchmarks

import (
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	r "gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/returnStorage"
	"testing"
)

func BenchmarkAddReturnToStorage(b *testing.B) {
	rs := r.NewReturnStorage()
	testReturn := &models.Return{ID: 1}

	for i := 0; i < b.N; i++ {
		err := rs.AddReturnToStorage(testReturn)
		if err != nil {
			b.Errorf("error adding return to storage: %v", err)
		}
	}
}

func BenchmarkDeleteReturnFromStorage(b *testing.B) {
	rs := r.NewReturnStorage()
	testReturn := &models.Return{ID: 1}
	err := rs.AddReturnToStorage(testReturn)
	if err != nil {
		b.Fatalf("error setting up delete benchmark: %v", err)
	}

	for i := 0; i < b.N; i++ {
		rs.DeleteReturnFromStorage(1)
	}
}

func BenchmarkReadFromJSONR(b *testing.B) {
	rs := r.NewReturnStorage()
	err := rs.WriteToJSON()
	if err != nil {
		b.Fatalf("error writing to JSON before reading: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := rs.ReadFromJSON()
		if err != nil {
			b.Errorf("error reading from JSON: %v", err)
		}
	}
}

func BenchmarkWriteToJSONR(b *testing.B) {
	rs := r.NewReturnStorage()
	err := rs.AddReturnToStorage(&models.Return{ID: 1})
	if err != nil {
		b.Fatalf("error adding return before writing to JSON: %v", err)
	}

	for i := 0; i < b.N; i++ {
		err := rs.WriteToJSON()
		if err != nil {
			b.Errorf("error writing to JSON: %v", err)
		}
	}
}

func BenchmarkIsConsistR(b *testing.B) {
	rs := r.NewReturnStorage()
	err := rs.AddReturnToStorage(&models.Return{ID: 1})
	if err != nil {
		b.Fatalf("error adding return for consist benchmark: %v", err)
	}

	for i := 0; i < b.N; i++ {
		exists := rs.IsConsist(1)
		if !exists {
			b.Errorf("error: return with ID 1 not found")
		}
	}
}

func BenchmarkGetReturn(b *testing.B) {
	rs := r.NewReturnStorage()
	err := rs.AddReturnToStorage(&models.Return{ID: 1})
	if err != nil {
		b.Fatalf("error adding return for get benchmark: %v", err)
	}

	for i := 0; i < b.N; i++ {
		rs.GetReturn(1)
	}
}

func BenchmarkGetReturnIDs(b *testing.B) {
	rs := r.NewReturnStorage()
	err := rs.AddReturnToStorage(&models.Return{ID: 1})
	if err != nil {
		b.Fatalf("error adding first return: %v", err)
	}
	err = rs.AddReturnToStorage(&models.Return{ID: 2})
	if err != nil {
		b.Fatalf("error adding second return: %v", err)
	}

	for i := 0; i < b.N; i++ {
		ids := rs.GetReturnIDs()
		if len(ids) == 0 {
			b.Errorf("error: no return IDs found")
		}
	}
}
