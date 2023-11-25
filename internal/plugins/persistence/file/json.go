package file

import (
	"encoding/json"
	"os"
	"path"

	"github.com/ilmaruk/esms/internal/models"
)

type JSONRosterStorer struct {
	basePath string
}

func NewJSONRosterStorer(basePath string) *JSONRosterStorer {
	return &JSONRosterStorer{basePath: basePath}
}

func (s *JSONRosterStorer) Store(roster models.Roster) error {
	b, err := json.MarshalIndent(roster, "", "  ")
	if err != nil {
		return err
	}

	filename := path.Join(s.basePath, roster.ID.String()+".json")

	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fh.Close()

	_, err = fh.Write(b)

	return err
}
