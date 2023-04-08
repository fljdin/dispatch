package models

import "fmt"

type Connections []Connection

type Connection struct {
	Name string `yaml:"name"`
	URI  string `yaml:"uri"`
}

func (c Connections) GetURIByName(name string) (string, error) {
	for _, connection := range c {
		if connection.Name == name {
			return connection.URI, nil
		}
	}
	return "", fmt.Errorf("connection not found for name %s", name)
}
