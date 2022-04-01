package bidirectional

import "testing"

func TestBiMap_Insert(t *testing.T) {
	m := NewBiMap[string, string]()

	k, v := "key", "value"
	m.Insert(k, v)

	vv, ok := m.Get(k)
	if !ok {
		t.Fatal("item should exist")
	}

	if vv != v {
		t.Fatal("wrong value")
	}
}

func TestBiMap_Delete(t *testing.T) {
	m := NewBiMap[string, string]()

	k, v := "key", "value"
	m.Insert(k, v)

	vv, ok := m.Get(k)
	if !ok {
		t.Fatal("item should exist")
	}

	if vv != v {
		t.Fatal("wrong value")
	}

	m.Delete(k)
	if _, ok := m.Get(k); ok {
		t.Fatal("item should exist")
	}
}

func TestBiMap_DeleteInverse(t *testing.T) {
	m := NewBiMap[string, string]()

	k, v := "key", "value"
	m.Insert(k, v)

	kk, ok := m.Inverse(v)
	if !ok {
		t.Fatal("item should exist")
	}

	if kk != k {
		t.Fatal("wrong value")
	}

	m.DeleteInverse(v)
	if _, ok := m.Get(k); ok {
		t.Fatal("item should be removed")
	}
}
