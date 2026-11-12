package storage

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	cacheInfoFile = "cache_info"
	feedPrefix    = "feed_"
)

type CacheInfoItem struct {
	LastFetch  time.Time
	FetchAfter time.Time
	ETag       string
	URL        string
}

func (s *LocalStorage) LoadCacheInfo() (map[string]*CacheInfoItem, error) {
	cacheinfo := make(map[string]*CacheInfoItem)
	path, err := s.JoinCacheDir(cacheInfoFile)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cacheinfo, nil
		}
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		item, err := parseCacheInfoItem(scanner.Text())
		if err != nil {
			return nil, err
		}
		cacheinfo[item.URL] = item
	}
	return cacheinfo, scanner.Err()
}

func (s *LocalStorage) SaveCacheInfo(cacheinfo map[string]*CacheInfoItem) error {
	path, err := s.JoinCacheDir(cacheInfoFile)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, item := range cacheinfo {
		_, err = f.Write(getCacheInfoItemLine(item))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *LocalStorage) SaveFeedCache(r io.Reader, name string) error {
	path, err := s.JoinCacheDir(FeedCacheFileName(name))
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}

func (s *LocalStorage) OpenFeedCache(name string) (io.ReadCloser, error) {
	path, err := s.JoinCacheDir(FeedCacheFileName(name))
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err == nil {
		return f, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}
	legacyPath, lerr := s.JoinCacheDir(legacyFeedCacheFileName(name))
	if lerr != nil {
		return nil, lerr
	}
	return os.Open(legacyPath)
}

func (s *LocalStorage) RemoveFeedCaches(names []string) error {
	cacheinfo, err := s.LoadCacheInfo()
	if err != nil {
		return err
	}
	for i := range names {
		delete(cacheinfo, names[i])
	}
	err = s.SaveCacheInfo(cacheinfo)
	if err != nil {
		return err
	}
	for i := range names {
		path, err := s.JoinCacheDir(FeedCacheFileName(names[i]))
		if err != nil {
			continue
		}
		os.Remove(path)
		legacyPath, err := s.JoinCacheDir(legacyFeedCacheFileName(names[i]))
		if err != nil {
			continue
		}
		os.Remove(legacyPath)
	}
	return nil
}

func FeedCacheFileName(name string) string {
	hash := sha256.Sum256([]byte(name))
	return feedPrefix + hex.EncodeToString(hash[:])
}

func legacyFeedCacheFileName(name string) string {
	return feedPrefix + url.QueryEscape(name)
}

func getCacheInfoItemLine(item *CacheInfoItem) []byte {
	return []byte(fmt.Sprintf("%s %d %s %d\n",
		item.URL,
		item.LastFetch.Unix(),
		url.QueryEscape(item.ETag),
		item.FetchAfter.Unix()),
	)
}

func parseCacheInfoItem(line string) (*CacheInfoItem, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid cache info line: %s", line)
	}
	lastCheck, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}
	etag, err := url.QueryUnescape(parts[2])
	if err != nil {
		return nil, err
	}
	var fetchAfter int64
	if len(parts) == 4 {
		fetchAfter, _ = strconv.ParseInt(parts[3], 10, 64)
	}
	return &CacheInfoItem{
		LastFetch:  time.Unix(lastCheck, 0),
		ETag:       etag,
		URL:        parts[0],
		FetchAfter: time.Unix(fetchAfter, 0),
	}, nil
}
