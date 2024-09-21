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
		rs.AddReturnToStorage(testReturn)
	}
}

func BenchmarkDeleteReturnFromStorage(b *testing.B) {
	rs := r.NewReturnStorage()
	testReturn := &models.Return{ID: 1}
	rs.AddReturnToStorage(testReturn)

	for i := 0; i < b.N; i++ {
		rs.DeleteReturnFromStorage(1)
	}
}

func BenchmarkReadFromJSONR(b *testing.B) {
	rs := r.NewReturnStorage()
	rs.WriteToJSON()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rs.ReadFromJSON()
	}
}

func BenchmarkWriteToJSONR(b *testing.B) {
	rs := r.NewReturnStorage()
	rs.AddReturnToStorage(&models.Return{ID: 1})

	for i := 0; i < b.N; i++ {
		rs.WriteToJSON()
	}
}

func BenchmarkIsConsistR(b *testing.B) {
	rs := r.NewReturnStorage()
	rs.AddReturnToStorage(&models.Return{ID: 1})

	for i := 0; i < b.N; i++ {
		rs.IsConsist(1)
	}
}

func BenchmarkGetReturn(b *testing.B) {
	rs := r.NewReturnStorage()
	rs.AddReturnToStorage(&models.Return{ID: 1})

	for i := 0; i < b.N; i++ {
		rs.GetReturn(1)
	}
}

func BenchmarkGetReturnIDs(b *testing.B) {
	rs := r.NewReturnStorage()
	rs.AddReturnToStorage(&models.Return{ID: 1})
	rs.AddReturnToStorage(&models.Return{ID: 2})

	for i := 0; i < b.N; i++ {
		rs.GetReturnIDs()
	}
}
