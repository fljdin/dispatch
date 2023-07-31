package tasks

import (
	"fmt"
	"net/url"
	"strconv"
)

type Connection struct {
	Name     string `yaml:"name"`
	URI      string `yaml:"uri"`
	Service  string `yaml:"service"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Dbname   string `yaml:"dbname"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (c Connection) CombinedURI() string {
	uri := "postgresql://?"
	params := url.Values{}

	if c.Service != "" {
		params.Add("service", c.Service)
	}

	if c.Host != "" {
		params.Add("host", c.Host)
	}

	if c.Port != 0 {
		params.Add("port", strconv.Itoa(c.Port))
	}

	if c.Dbname != "" {
		params.Add("dbname", c.Dbname)
	}

	if c.User != "" {
		params.Add("user", c.User)
	}

	if c.Password != "" {
		params.Add("password", c.Password)
	}

	uri += params.Encode()

	return uri
}

type Connections []Connection

func (c Connections) GetURIByName(name string) (string, error) {
	for _, connection := range c {
		if connection.Name != name {
			continue
		}
		if connection.URI == "" {
			return connection.CombinedURI(), nil
		}
		return connection.URI, nil
	}
	return "", fmt.Errorf("connection not found for name %s", name)
}
