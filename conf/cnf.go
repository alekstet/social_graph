package conf

import (
	"io/ioutil"

	"github.com/olebedev/config"
)

type Conf struct {
	Host     string
	PortBase string
	User     string
	Password string
	DBName   string
	PortApp  string
}

func Cnf() (*Conf, error) {
	file, err := ioutil.ReadFile("github.com/alekstet/social_graph/conf/cnf.yml")
	if err != nil {
		return nil, err
	}
	yamlString := string(file)

	cfg, err := config.ParseYaml(yamlString)
	if err != nil {
		return nil, err
	}
	Host, err := cfg.String("host")
	if err != nil {
		return nil, err
	}
	PortBase, err := cfg.String("portBase")
	if err != nil {
		return nil, err
	}
	User, err := cfg.String("user")
	if err != nil {
		return nil, err
	}
	Password, err := cfg.String("password")
	if err != nil {
		return nil, err
	}
	DBName, err := cfg.String("dbname")
	if err != nil {
		return nil, err
	}
	PortApp, err := cfg.String("portApp")
	if err != nil {
		return nil, err
	}

	return &Conf{
		Host:     Host,
		PortBase: PortBase,
		User:     User,
		Password: Password,
		DBName:   DBName,
		PortApp:  ":" + PortApp,
	}, nil
}
