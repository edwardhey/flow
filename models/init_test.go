package models

import (
	"testing"
)

func Test_DB_Instance(t *testing.T) {
	if db == nil {
		t.Fatalf("db is nil")
	}
}
