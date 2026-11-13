package goose

import (
	"testing"
)

func TestVersion(t *testing.T) {
	var err error

	db := testDatabase(t)
	if err = Version(db, ""); err != nil {
		t.Errorf("не ожидалась ошибка вывода версии: %v", err)
	}
}

func TestVersionClosedDb(t *testing.T) {
	var err error

	db := testDatabase(t)
	if err = db.Close(); err != nil {
		t.Fatalf("ошибка закрытия базы данных: %v", err)
	}
	if err = Version(db, ""); err == nil {
		t.Fatal("ожидалась ошибка вывода версии")
	}
}
