package config

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"

	"github.com/contractGuru/pkg/logger"
)

// Config struct stores configuration values which will be read from env file.
type Config struct {
	port       string
	dbUsername string
	dbPassword string
	dbHost     string
	dbPort     string
	dbName     string
	dbOpt      string
	dbOptVal   string
	dbLoc      string
	dbLocVal   string
	cert       string
	key        string
}

// GetConfig creates a pointer to Config using env variables as default values.
func GetConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.dbUsername, "dbusername", os.Getenv("DB_USERNAME"), "DB username")
	flag.StringVar(&cfg.dbPassword, "dbpassword", os.Getenv("DB_PASSWORD"), "DB password")
	flag.StringVar(&cfg.dbPort, "dbport", os.Getenv("DB_PORT"), "DB port")
	flag.StringVar(&cfg.dbHost, "dbhost", os.Getenv("DB_HOST"), "DB host")
	flag.StringVar(&cfg.dbName, "dbname", os.Getenv("DB_NAME"), "DB name")
	flag.StringVar(&cfg.dbOpt, "dbopt", os.Getenv("DB_OPTIONS"), "DB options")
	flag.StringVar(&cfg.dbOptVal, "dboptval", os.Getenv("DB_OPT_VALUE"), "DB option value")
	flag.StringVar(&cfg.dbLoc, "dbloc", os.Getenv("DB_LOC_TIME"), "DB loc time")
	flag.StringVar(&cfg.dbLocVal, "dblocval", os.Getenv("DB_LOC_TIME_OPT"), "DB loc time value")
	flag.StringVar(&cfg.port, "port", os.Getenv("PORT"), "server port")
	flag.StringVar(&cfg.cert, "cert", os.Getenv("CERT"), "certificate")
	flag.StringVar(&cfg.key, "key", os.Getenv("KEY"), "key")

	flag.Parse()

	return cfg
}

// GetDBConnStr returns the config instance of DB connection string.
func (c *Config) GetDBConnStr() string {
	return c.getDBConnStr(c.dbHost, c.dbName)
}

func (c *Config) getDBConnStr(dbhost, dbname string) string {
	dbSettings := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s=%s&%s=%s",
		c.dbUsername,
		c.dbPassword,
		dbhost,
		c.dbPort,
		dbname,
		c.dbOpt,
		c.dbOptVal,
		c.dbLoc,
		c.dbLocVal,
	)

	return dbSettings
}

// GetAPIPort returns the API port as string.
func (c *Config) GetAPIPort() string {
	listenAt := fmt.Sprintf(":%s", c.port)

	return listenAt
}

// GetCert returns certificate.
func (c *Config) GetCert() (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(c.cert, c.key)

	if err != nil {
		logger.Error.Println(err.Error())
		return nil, err
	}

	return &cert, nil
}
