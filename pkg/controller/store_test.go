package controller

import (
	"encoding/json"
	"testing"

	"github.com/jilanisayyad/edgebeat/pkg/utils"
)

func TestStoreSetGet(t *testing.T) {
	store := NewStore()
	info := utils.SystemInfo{Timestamp: "2026-02-15T00:00:00Z"}
	payload, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	store.Set(payload)
	got, ok := store.Get()
	if !ok {
		t.Fatal("expected payload in store")
	}
	if string(got) != string(payload) {
		t.Fatalf("payload mismatch: %s", string(got))
	}
}

func TestStoreSetInvalidJSON(t *testing.T) {
	store := NewStore()
	store.Set([]byte("not-json"))
	_, ok := store.Get()
	if ok {
		t.Fatal("expected no payload after invalid JSON")
	}
}

func TestStoreGetNil(t *testing.T) {
	var store *Store
	if _, ok := store.Get(); ok {
		t.Fatal("expected no payload from nil store")
	}
}

func TestStoreAccessors(t *testing.T) {
	store := NewStore()
	info := utils.SystemInfo{
		Timestamp: "2026-02-15T00:00:00Z",
		CPU:       utils.CPUStats{TotalPercent: 42.5},
	}
	payload, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	store.Set(payload)

	if data, ok := store.GetCPU(); !ok || data.Timestamp == "" || data.CPU.TotalPercent != 42.5 {
		t.Fatalf("GetCPU = %+v, ok=%v", data, ok)
	}
	if data, ok := store.GetMemory(); !ok || data.Timestamp == "" {
		t.Fatalf("GetMemory = %+v, ok=%v", data, ok)
	}
	if data, ok := store.GetDisk(); !ok || data.Timestamp == "" {
		t.Fatalf("GetDisk = %+v, ok=%v", data, ok)
	}
	if data, ok := store.GetNetwork(); !ok || data.Timestamp == "" {
		t.Fatalf("GetNetwork = %+v, ok=%v", data, ok)
	}
	if data, ok := store.GetSystem(); !ok || data.Timestamp == "" {
		t.Fatalf("GetSystem = %+v, ok=%v", data, ok)
	}
	if data, ok := store.GetSensors(); !ok || data.Timestamp == "" {
		t.Fatalf("GetSensors = %+v, ok=%v", data, ok)
	}
}
