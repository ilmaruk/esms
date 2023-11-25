package file

import (
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/google/uuid"
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

	filename := path.Join(s.basePath, roster.ID.String()+".rst.json")

	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fh.Close()

	_, err = fh.Write(b)

	return err
}

type JSONRosterFetcher struct {
	basePath string
}

func NewJSONRosterFetcher(basePath string) *JSONRosterFetcher {
	return &JSONRosterFetcher{basePath: basePath}
}

func (s *JSONRosterFetcher) Fetch(id uuid.UUID) (models.Roster, error) {
	var roster models.Roster

	filename := path.Join(s.basePath, id.String()+".rst.json")
	fh, err := os.Open(filename)
	if err != nil {
		return roster, err
	}

	b, err := io.ReadAll(fh)
	if err != nil {
		return roster, err
	}

	err = json.Unmarshal(b, &roster)

	return roster, err
}

type JSONTeamsheetStorer struct {
	basePath string
}

func NewJSONTeamsheetStorer(basePath string) *JSONTeamsheetStorer {
	return &JSONTeamsheetStorer{basePath: basePath}
}

func (s *JSONTeamsheetStorer) Store(teamsheet models.Teamsheet) error {
	b, err := json.MarshalIndent(teamsheet, "", "  ")
	if err != nil {
		return err
	}

	filename := path.Join(s.basePath, teamsheet.ID.String()+".tsh.json")

	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fh.Close()

	_, err = fh.Write(b)

	return err
}
