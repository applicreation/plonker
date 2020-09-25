package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Dir string
	Connection
	Tables []Table
}

type File struct {
	Name string
	Path string
}

type Connection struct {
	Engine   string `yaml:"engine"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Table struct {
	Name       string   `yaml:"name"`
	PrimaryKey string   `yaml:"primaryKey"`
	Columns    []string `yaml:"columns"`
	Range      `yaml:"range"`
	Order      `yaml:"order"`
	Relations  []Relation `yaml:"relations"`
}

type Range struct {
	Records    int `yaml:"records"`
	Percentage int `yaml:"percentage"`
	Timeframe  `yaml:"timeframe"`
}

type Order struct {
	Column    string `yaml:"column"`
	Direction string `yaml:"direction"`
}

type Relation struct {
	Key    string `yaml:"key"`
	Table  string `yaml:"table"`
	Column string `yaml:"column"`
}

type Timeframe struct {
	Period string `yaml:"period"`
	Column string `yaml:"column"`
}

func (c *Config) Load() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting pwd: %s", err)
		return
	}

	files, err := dirFiles(pwd + "/" + c.Dir)
	if err != nil {
		log.Printf("Error getting files: %s", err)
		return
	}

	for _, file := range files {
		content, err := ioutil.ReadFile(file.Path)
		if err != nil {
			return
		}

		text := string(content)

		if file.Name == "connection.yaml" {
			var connection Connection

			err = yaml.Unmarshal([]byte(text), &connection)
			if err != nil {
				return
			}

			c.Connection = connection
		} else {
			var table Table

			err = yaml.Unmarshal([]byte(text), &table)
			if err != nil {
				return
			}

			c.Tables = append(c.Tables, table)
		}
	}
}

func dirFiles(dir string) ([]File, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("Configuration path must be a directory: %s", dir)
	}

	var files []File
	err = nil
	for err != io.EOF {
		var fis []os.FileInfo
		fis, err = f.Readdir(128)
		if err != nil && err != io.EOF {
			return nil, err
		}

		for _, fi := range fis {
			if fi.IsDir() {
				continue
			}

			name := fi.Name()
			if !strings.HasSuffix(name, ".yaml") {
				continue
			}

			path := filepath.Join(dir, name)

			files = append(files, File{
				Name: name,
				Path: path,
			})
		}
	}

	return files, nil
}
