package dbconf

import (
	"encoding/json"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql" // register mysql driver
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	"github.com/pkg/errors"
)

// DBConfig contains the database driver name and configuration to be passed to Open
type DBConfig struct {
	DriverName     string `json:"driver"`
	DataSourceName string `json:"data_source"`
}

// loadConfigurationFile attempts to load the db configuration file stored at the path
// and returns the configuration. On error, it returns nil.
func loadConfigurationFile(path string) (cfg *DBConfig, err error) {
	if path == "" {
		return nil, errors.Errorf("invalid empty path")
	}

	var body []byte
	body, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read configuration file '%s'", path)
	}

	cfg = &DBConfig{}
	err = json.Unmarshal(body, &cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal configuration '%s'", path)
	}

	if cfg.DataSourceName == "" || cfg.DriverName == "" {
		return nil, errors.Errorf("Invalid db configuration '%s'", path)
	}

	return
}

// DBFromConfig opens a sql.DB from settings in a db config file
func DBFromConfig(path string) (db *sqlx.DB, err error) {
	var dbCfg *DBConfig
	dbCfg, err = loadConfigurationFile(path)
	if err != nil {
		return nil, err
	}

	return sqlx.Open(dbCfg.DriverName, dbCfg.DataSourceName)
}
