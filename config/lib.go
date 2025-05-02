package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
)

type Session struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type SessionProfile struct {
	Name        string  `json:"name"`
	Sessions []*Session `json:"sessions"`
}

type Config struct {
	BookmarkRoots   []string          `json:"BookmarkRoots"`
	SessionProfiles []*SessionProfile `json:"SessionProfiles"`
}

func (c *Config) SaveConfig() error {
	b, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}

	configPath, _, err := getConfigPath()
	if err != nil {
		return err
	}

	_, err = os.Lstat(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		f, err := os.Create(configPath)
		if err != nil {
			return err
		}

		_, err = f.Write(b)
		if err != nil {
			return err
		}

		return f.Close()
	}

	return os.WriteFile(configPath, b, 0666)
}

func (c *Config) Sanitize() {
	if len(c.BookmarkRoots) <= 0 {
		return
	}

	c.BookmarkRoots = slices.DeleteFunc(
		c.BookmarkRoots,
		func(root string) bool {
			expanded := os.ExpandEnv(root)
			_, err := os.Lstat(expanded)
			return err != nil
		},
	)
	for _, profile := range c.SessionProfiles {
		profile.Sessions = slices.DeleteFunc(
			profile.Sessions,
			func(s *Session) bool {
				expanded := os.ExpandEnv(s.Path)
				_, err := os.Lstat(expanded)
				return err != nil
			},
		)
	}
}

func LoadConfig() (*Config, error) {
	configPath, home, err := getConfigPath()

	_, err = os.Lstat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return initConfig(configPath, home)
		}
	}

	b, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var conf Config

	err = json.Unmarshal(b, &conf)
	if err != nil {
		return nil, err
	}

	conf.Sanitize()

	return &conf, nil
}

func getConfigPath() (string, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}

	return filepath.Join(home, ".shutil.json"), home, nil 
}

func initConfig(configPath string, home string) (*Config, error) {
	f, err := os.Create(configPath)
	if err != nil {
		return nil, err
	}

	conf := Config{
		[]string{ home },
		[]*SessionProfile{},
	}

	blob, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return nil, err
	}

	_, err = f.Write(blob)
	if err != nil {
		return nil, err
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
