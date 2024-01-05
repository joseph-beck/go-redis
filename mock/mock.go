package mock

import (
	"encoding/json"
	"os"

	"github.com/joseph-beck/go-redis/database"
	"github.com/joseph-beck/go-redis/models"
)

func LoadMockData(s *database.Store, p string) error {
	f, err := os.ReadFile(p)
	if err != nil {
		return err
	}

	var d []models.User
	err = json.Unmarshal(f, &d)
	if err != nil {
		return err
	}

	for _, i := range d {
		err := s.Add(&i, "users")
		if err != nil {
			return err
		}
	}

	return nil
}
