package respond

import (
	"testing"
)

func TestDataErrorHasDetails(t *testing.T) {
	e := NewDataError()

	if e.HasDetails() {
		t.Fatal("returned true when it should be false")
	}

	e.Add("password", "is required")

	if !e.HasDetails() {
		t.Fatal("returned false when it should be true")
	}
}