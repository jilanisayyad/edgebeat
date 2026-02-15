package controller

import (
	"encoding/json"
	"sync"

	"github.com/jilanisayyad/edgebeat/pkg/utils"
)

type Store struct {
	mu      sync.RWMutex
	payload []byte
	info    *utils.SystemInfo
	hasData bool
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Set(payload []byte) {
	if s == nil {
		return
	}

	copyPayload := make([]byte, len(payload))
	copy(copyPayload, payload)

	var info utils.SystemInfo
	if err := json.Unmarshal(payload, &info); err != nil {
		return
	}

	s.mu.Lock()
	s.payload = copyPayload
	s.info = &info
	s.hasData = true
	s.mu.Unlock()
}

func (s *Store) Get() ([]byte, bool) {
	if s == nil {
		return nil, false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.hasData {
		return nil, false
	}

	copyPayload := make([]byte, len(s.payload))
	copy(copyPayload, s.payload)

	return copyPayload, true
}

// GetCPU returns full system info for CPU metrics access
func (s *Store) GetCPU() (*utils.SystemInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.hasData || s.info == nil {
		return nil, false
	}

	return s.info, true
}

// GetMemory returns full system info for memory metrics access
func (s *Store) GetMemory() (*utils.SystemInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.hasData || s.info == nil {
		return nil, false
	}

	return s.info, true
}

// GetDisk returns full system info for disk metrics access
func (s *Store) GetDisk() (*utils.SystemInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.hasData || s.info == nil {
		return nil, false
	}

	return s.info, true
}

// GetNetwork returns full system info for network metrics access
func (s *Store) GetNetwork() (*utils.SystemInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.hasData || s.info == nil {
		return nil, false
	}

	return s.info, true
}

// GetSystem returns full system info for system metrics access
func (s *Store) GetSystem() (*utils.SystemInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.hasData || s.info == nil {
		return nil, false
	}

	return s.info, true
}

// GetSensors returns full system info for sensor metrics access
func (s *Store) GetSensors() (*utils.SystemInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.hasData || s.info == nil {
		return nil, false
	}

	return s.info, true
}
