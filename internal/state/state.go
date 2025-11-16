package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

type State struct {
	Current string   `json:"current"`
	Configs []Config `json:"configs"`
}

func ConfigExists(name string) (bool, error) {
	state, err := LoadState()
	if err != nil {
		return false, err
	}
	for _, c := range state.Configs {
		if c.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func SaveState(state *State) error {
	stateFile, err := getStateFile()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(stateFile, data, 0o644)
}

func LoadState() (*State, error) {
	stateFile, err := getStateFile()
	if err != nil {
		return nil, err
	}

	if _, err = os.Stat(stateFile); os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Dir(stateFile), 0o755); err != nil {
			return nil, err
		}
		state := &State{
			Current: "",
			Configs: []Config{},
		}
		if err = SaveState(state); err != nil {
			return nil, err
		}
		return state, nil
	}

	data, err := os.ReadFile(stateFile)
	if err != nil {
		return nil, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config"), nil
}

func getStateFile() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "nvmgr", "state.json"), nil
}
