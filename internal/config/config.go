package config

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Project string `json:"project"`
	DBPath  string `json:"db_path"`
}

func BaseDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".yayo")
}

func Save(project, dbPath string) error {
	cfg := Config{
		Project: project,
		DBPath:  dbPath,
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	os.MkdirAll(BaseDir(), 0755)
	return os.WriteFile(filepath.Join(BaseDir(), "config.json"), data, 0644)
}

func Clean(s string) string {
	return strings.TrimSpace(s)
}

func LoadEnvFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	envs := map[string]string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		envs[parts[0]] = parts[1]
	}

	return envs, nil
}
