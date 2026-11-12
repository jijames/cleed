package storage

import (
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/radulucut/cleed/internal/utils"
)

const (
	dataPathOverrideName = "CLEED_DATA_PATH"

	listsDir   = "lists"
	exploreDir = "explore"
)

// defaultLocalConfigDir and defaultLocalCacheDir are used when set from main and
// CLEED_DATA_PATH is unset.
var (
	defaultLocalConfigDir string
	defaultLocalCacheDir  string
)

// SetDefaultLocalDirs sets the default config and cache directories when
// CLEED_DATA_PATH is unset. Pass empty strings to clear.
func SetDefaultLocalDirs(configDir, cacheDir string) {
	defaultLocalConfigDir = configDir
	defaultLocalCacheDir = cacheDir
}

type LocalStorage struct {
	name string
	time utils.Time

	config *Config
}

func NewLocalStorage(
	name string,
	time utils.Time,
) *LocalStorage {
	return &LocalStorage{
		name: name,
		time: time,
	}
}

func (s *LocalStorage) Init(version string) error {
	bootstrapDir, err := s.bootstrapConfigDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(bootstrapDir, 0755); err != nil {
		return err
	}
	_, err = s.LoadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			s.config = &Config{
				Version:   version,
				Timeout:   30,
				BatchSize: 100,
				ColorMap:  make(map[uint8]uint8),
			}
			if err := s.SaveConfig(); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	listsPath, err := s.JoinConfigDir(listsDir)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(listsPath, 0755); err != nil {
		return err
	}
	cacheDir, err := s.JoinCacheDir("")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}
	return s.Migrate()
}

func (s *LocalStorage) Migrate() error {
	// handle migration here
	return nil
}

func (s *LocalStorage) ClearAll() error {
	configDir, err := s.JoinConfigDir("")
	if err != nil {
		return err
	}
	if err := os.RemoveAll(configDir); err != nil {
		return err
	}
	cacheDir, err := s.JoinCacheDir("")
	if err != nil {
		return err
	}
	if err := os.RemoveAll(cacheDir); err != nil {
		return err
	}
	bootstrapDir, err := s.bootstrapConfigDir()
	if err != nil {
		return err
	}
	if err := os.RemoveAll(bootstrapDir); err != nil {
		return err
	}
	return nil
}

// bootstrapConfigDir returns the directory that always holds config.json.
// Resolution order:
//  1. CLEED_DATA_PATH env override -> <CLEED_DATA_PATH>/config
//  2. defaultLocalConfigDir (set from main) or OS user config directory
func (s *LocalStorage) bootstrapConfigDir() (string, error) {
	if p := os.Getenv(dataPathOverrideName); p != "" {
		return filepath.Join(p, "config"), nil
	}
	if defaultLocalConfigDir != "" {
		return resolveLocalDir(defaultLocalConfigDir)
	}
	appBase, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appBase, s.name), nil
}

// bootstrapCacheDir returns the default cache directory used when no
// `feedsPath` override is configured.
//  1. CLEED_DATA_PATH env override -> <CLEED_DATA_PATH>/cache
//  2. defaultLocalCacheDir (set from main) or OS user cache directory
func (s *LocalStorage) bootstrapCacheDir() (string, error) {
	if p := os.Getenv(dataPathOverrideName); p != "" {
		return filepath.Join(p, "cache"), nil
	}
	if defaultLocalCacheDir != "" {
		return resolveLocalDir(defaultLocalCacheDir)
	}
	appBase, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appBase, s.name), nil
}

func resolveLocalDir(p string) (string, error) {
	p = filepath.Clean(p)
	if filepath.IsAbs(p) {
		return p, nil
	}
	return filepath.Abs(p)
}

// JoinConfigDir returns a path under the configuration directory. If the user
// configured a custom `configPath` it is honoured, otherwise the bootstrap
// config dir is used. Note: config.json itself is always read from / written
// to the bootstrap dir (see LoadConfig / SaveConfig).
func (s *LocalStorage) JoinConfigDir(file string) (string, error) {
	if s.config != nil && s.config.ConfigPath != "" && defaultLocalConfigDir == "" {
		return filepath.Join(s.config.ConfigPath, file), nil
	}
	base, err := s.bootstrapConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, file), nil
}

// JoinCacheDir returns a path under the cache (RSS feed) directory. The
// `feedsPath` setting overrides the bootstrap cache dir when set.
func (s *LocalStorage) JoinCacheDir(file string) (string, error) {
	if s.config != nil && s.config.FeedsPath != "" && defaultLocalCacheDir == "" {
		return filepath.Join(s.config.FeedsPath, file), nil
	}
	base, err := s.bootstrapCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, file), nil
}

func (s *LocalStorage) joinListsDir(file string) (string, error) {
	base, err := s.JoinConfigDir(listsDir)
	if err != nil {
		return "", err
	}
	return path.Join(base, url.QueryEscape(file)), nil
}

func (s *LocalStorage) JoinExploreDir(file string) (string, error) {
	base, err := s.JoinCacheDir(exploreDir)
	if err != nil {
		return "", err
	}
	return path.Join(base, url.QueryEscape(file)), nil
}
