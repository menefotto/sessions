package session

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"sync"
)

type Store interface {
	Put(key string, val Session) error
	Get(key string) (Session, error)
}

func NewMapStore() Store {
	return &MapStore{new(sync.RWMutex), make(map[string]Session)}
}

type MapStore struct {
	mtx *sync.RWMutex
	mem map[string]Session
}

func (s *MapStore) Put(key string, val Session) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.mem[key] = val

	return nil
}

func (s *MapStore) Get(key string) (Session, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	val, ok := s.mem[key]
	if !ok {
		return &NewSession(""), nil
	}

	return val, nil
}

func NewFileSystemStore() Store {
	return &FileSystemStore{}
}

type FileSystemStore struct {
	in  *bytes.Buffer
	dec *gob.Decoder
	enc *gob.Encoder
}

func (s *FileSystemStore) Put(key string, val Session) error {
	s.enc = gob.NewEncoder(bytes.NewBuffer([]byte))
	if err := s.enc.Encode(val); err != nil {
		return err
	}

	err := ioutil.WriteFile("_session"+key, s.in.Bytes, 0666)
	if err != nil {
		return err
	}
	defer s.in.Reset()

	return nil
}

func (s *FileSystemStore) Get(key string) (Session, error) {
	var ss Session

	data, err := ioutil.ReadFile("_session" + key)
	if err != nil {
		return ss, err
	}

	s.enc = gob.NewDecoder(bytes.NewBuffer(data))
	err := s.dec.Decode(&ss)
	if err != nil {
		return ss, err
	}

	return ss, nil
}
