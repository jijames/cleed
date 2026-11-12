package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const (
	configFile = "config.json"
)

type Config struct {
	Version string `json:"version"`

	UserAgent string `json:"userAgent"`
	Timeout   uint   `json:"timeout"`   // in seconds
	BatchSize uint   `json:"batchSize"` // number of feeds to fetch in a batch

	// ConfigPath is an optional override for the directory where lists and
	// other configuration data (other than config.json itself) are stored.
	// When empty the bootstrap config directory is used (see
	// LocalStorage.bootstrapConfigDir).
	ConfigPath string `json:"configPath"`
	// FeedsPath is an optional override for the directory where cached RSS
	// feeds and related cache files are stored. When empty the bootstrap
	// cache directory is used (see LocalStorage.bootstrapCacheDir).
	FeedsPath string `json:"feedsPath"`

	LastRun         time.Time       `json:"lastRun"`
	Styling         uint8           `json:"styling"` // 0: default, 1: enabled, 2: disabled
	Summary         uint8           `json:"summary"` // 0: disabled, 1: enabled
	ColorMap        map[uint8]uint8 `json:"colorMap"`
	HideFutureItems bool            `json:"hideFutureItems"`

	MinifluxToken string `json:"minifluxToken"`
}

func (s *LocalStorage) LoadConfig() (*Config, error) {
	if s.config != nil {
		return s.config, nil
	}
	bootstrapDir, err := s.bootstrapConfigDir()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(filepath.Join(bootstrapDir, configFile))
	if err != nil {
		return nil, err
	}
	s.config = &Config{}
	err = json.Unmarshal(b, s.config)
	if err != nil {
		return nil, err
	}
	if s.config.ColorMap == nil {
		s.config.ColorMap = make(map[uint8]uint8)
	}
	if s.config.BatchSize == 0 {
		s.config.BatchSize = 100
	}
	if s.config.Timeout == 0 {
		s.config.Timeout = 30
	}
	s.sanitizeConfigPaths()
	return s.config, nil
}

func (s *LocalStorage) sanitizeConfigPaths() {
	if s.config == nil {
		return
	}
	changed := false
	if s.config.ConfigPath != "" && !dirExists(s.config.ConfigPath) {
		s.config.ConfigPath = ""
		changed = true
	}
	if s.config.FeedsPath != "" && !dirExists(s.config.FeedsPath) {
		s.config.FeedsPath = ""
		changed = true
	}
	if changed {
		_ = s.SaveConfig()
	}
}

func dirExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}

func (s *LocalStorage) SaveConfig() error {
	if s.config == nil {
		return nil
	}
	bootstrapDir, err := s.bootstrapConfigDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(bootstrapDir, 0755); err != nil {
		return err
	}
	b, err := json.Marshal(s.config)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(bootstrapDir, configFile), b, 0600)
}
