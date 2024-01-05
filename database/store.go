package database

import (
	"os"

	"github.com/joseph-beck/go-redis/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func New() *Store {
	db, err := gorm.Open(sqlite.Open(os.Getenv("SQLITE_DATABASE")))
	if err != nil {
		panic(err)
	}

	return &Store{
		db: db,
	}
}

func (s *Store) Close() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}

	err = db.Close()
	return err
}

func (s *Store) List(i interface{}, t string) error {
	r := s.db.Table(t).Find(i)
	return r.Error
}

func (s *Store) Get(i interface{}, t string) error {
	r := s.db.Table(t).Find(i).First(i)
	return r.Error
}

func (s *Store) Add(i interface{}, t string) error {
	r := s.db.Table(t).Create(i)
	return r.Error
}

func (s *Store) Update(i interface{}, t string) error {
	r := s.db.Table(t).Save(i)
	return r.Error
}

func (s *Store) Delete(i interface{}, t string) error {
	r := s.db.Table(t).Delete(i)
	return r.Error
}

func (s *Store) Contains(i interface{}, t string) bool {
	r := s.db.Table(t).Model(i).First(i)
	return r.Error == nil
}

func (s *Store) AutoMigrate() error {
	err := s.db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	return nil
}
