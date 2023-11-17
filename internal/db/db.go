package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/npmaile/wagebot/internal/models"
)

type DataStore interface {
	GetServerConfiguration(guildID string) (models.Server, error)
	GetAllServerConfigs() ([]*models.Server, error)
}

type sqliteStore struct {
	storage *sql.DB
}

func NewSqliteStore(filePath string) (DataStore, error) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to get new sqlite data store: %s", err.Error())
	}
	_, err = db.Exec(`CREATE TABLE if not exists 
	servers(id TEXT, channelPrefix TEXT, rolePrefix TEXT)`)
	if err != nil {
		return nil, fmt.Errorf("unable to get new sqlite data store: %s", err.Error())
	}
	return &sqliteStore{
		storage: db,
	}, nil
}

func (s *sqliteStore) GetServerConfiguration(guildID string) (models.Server, error) {
	row := s.storage.QueryRow(`SELECT
	id, channelPrefix, rolePrefix
	FROM
	servers
	WHERE
	id = ?`, guildID)
	err := row.Err()
	if err != nil {
		return models.Server{}, fmt.Errorf("unable to get server configuration: %s", err.Error())
	}
	var ret models.Server
	err = row.Scan(&ret.ID, &ret.ChannelPrefix, &ret.RolePrefix)
	if err != nil {
		return models.Server{}, fmt.Errorf("unable to get server configuration: %s", err.Error())
	}
	return ret, nil
}

func (s *sqliteStore) GetAllServerConfigs() ([]*models.Server, error) {
	rows, err := s.storage.Query(`SELECT
	id, channelPrefix, rolePrefix
	FROM
	servers`)
	if err != nil {
		return nil, fmt.Errorf("unable to list server configurations: %s", err.Error())
	}
	var ret []*models.Server
	for rows.Next() {
		s := models.Server{}
		err := rows.Scan(&s.ID, &s.ChannelPrefix, &s.RolePrefix)
		if err != nil {
			return nil, fmt.Errorf("unable to scan server configs into struct: %s", err.Error())
		}
		ret = append(ret, &s)
	}
	return ret, nil

}