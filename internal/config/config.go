// Package config handles loading and saving secrets-cli configuration files.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the global secrets configuration (.secrets/config.yaml)
type Config struct {
	Version string `yaml:"version"`
	Owner   string `yaml:"owner"`
}

// VaultConfig represents a vault's configuration (vault.yaml)
type VaultConfig struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description,omitempty"`
	Members     []string `yaml:"members"`
	CreatedAt   string   `yaml:"created_at"`
	UpdatedAt   string   `yaml:"updated_at,omitempty"`
}

// LoadConfig loads the global config from .secrets/config.yaml
func LoadConfig(secretsDir string) (*Config, error) {
	path := filepath.Join(secretsDir, "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// SaveConfig saves the global config to .secrets/config.yaml
func SaveConfig(secretsDir string, cfg *Config) error {
	path := filepath.Join(secretsDir, "config.yaml")
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// LoadVaultConfig loads a vault's configuration
func LoadVaultConfig(vaultDir string) (*VaultConfig, error) {
	path := filepath.Join(vaultDir, "vault.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault config: %w", err)
	}

	var cfg VaultConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse vault config: %w", err)
	}

	return &cfg, nil
}

// SaveVaultConfig saves a vault's configuration
func SaveVaultConfig(vaultDir string, cfg *VaultConfig) error {
	path := filepath.Join(vaultDir, "vault.yaml")
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to serialize vault config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write vault config: %w", err)
	}

	return nil
}

// VaultExists checks if a vault exists
func VaultExists(secretsDir, vaultName string) bool {
	vaultDir := filepath.Join(secretsDir, "vaults", vaultName)
	info, err := os.Stat(vaultDir)
	return err == nil && info.IsDir()
}

// GetVaultDir returns the path to a vault directory
func GetVaultDir(secretsDir, vaultName string) string {
	return filepath.Join(secretsDir, "vaults", vaultName)
}

// GetKeysDir returns the path to the keys directory
func GetKeysDir(secretsDir string) string {
	return filepath.Join(secretsDir, "keys")
}

// ListVaults returns a list of all vault names
func ListVaults(secretsDir string) ([]string, error) {
	vaultsDir := filepath.Join(secretsDir, "vaults")
	entries, err := os.ReadDir(vaultsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to list vaults: %w", err)
	}

	var vaults []string
	for _, entry := range entries {
		if entry.IsDir() {
			vaults = append(vaults, entry.Name())
		}
	}

	return vaults, nil
}
