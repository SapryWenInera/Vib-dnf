package core

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/vanilla-os/vib/api"
)

type AptModule struct {
	Name   string     `json:"name"`
	Type   string     `json:"type"`
	Source api.Source `json:"source"`
}

// BuildAptModule builds a module that installs packages
// using the apt package manager
func BuildAptModule(moduleInterface interface{}, recipe *api.Recipe) (string, error) {
	var module AptModule
	err := mapstructure.Decode(moduleInterface, &module)
	if err != nil {
		return "", err
	}
	if len(module.Source.Packages) > 0 {
		packages := ""
		for _, pkg := range module.Source.Packages {
			packages += pkg + " "
		}

		return fmt.Sprintf("apt install -y %s && apt clean", packages), nil
	}

	if len(module.Source.Paths) > 0 {
		cmd := ""

		for i, path := range module.Source.Paths {
			instPath := filepath.Join(recipe.ParentPath, path+".inst")
			pkgs := ""
			file, err := os.Open(instPath)
			if err != nil {
				return "", err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				pkgs += scanner.Text() + " "
			}

			if err := scanner.Err(); err != nil {
				return "", err
			}

			cmd += fmt.Sprintf("apt install -y %s ", pkgs)

			if i != len(module.Source.Paths)-1 {
				cmd += "&& "
			} else {
				cmd += "&& apt clean"
			}
		}

		return cmd, nil
	}

	return "", errors.New("no packages or paths specified")
}
