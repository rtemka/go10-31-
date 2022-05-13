package memDb

import "GoNews/pkg/storage"

var FakeData = []storage.Post{
	{Id: 1, Title: "mem db post 1", Content: "Lorem ipsum"},
	{Id: 2, Title: "mem db post 2", Content: "Lorem ipsum"},
}

var FakePost = storage.Post{Id: 99, Title: "mem db post 99", Content: "Lorem ipsum"}

// MemDb заглушка для реализации БД в памяти
type MemDb struct{}

func New() *MemDb {
	return &MemDb{}
}

func (db *MemDb) Posts() ([]storage.Post, error) {
	return FakeData, nil
}

func (db *MemDb) AddPost(_ storage.Post) error {
	return nil
}

func (db *MemDb) UpdatePost(_ storage.Post) error {
	return nil
}

func (db *MemDb) DeletePost(_ storage.Post) error {
	return nil
}

func (db *MemDb) Close() { return }
