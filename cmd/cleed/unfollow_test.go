package cleed

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/radulucut/cleed/internal"
	_storage "github.com/radulucut/cleed/internal/storage"
	"github.com/radulucut/cleed/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_Unfollow_Default(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeMock := mocks.NewMockTime(ctrl)
	timeMock.EXPECT().Now().Return(defaultCurrentTime).AnyTimes()

	out := new(bytes.Buffer)
	printer := internal.NewPrinter(nil, out, out)
	storage := _storage.NewLocalStorage("cleed_test", timeMock)
	defer localStorageCleanup(t, storage)

	listsDir := path.Join(testConfigDir(t), "lists")
	err := os.MkdirAll(listsDir, 0700)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(path.Join(listsDir, "default"),
		[]byte(fmt.Sprintf("%d %s\n%d %s\n",
			defaultCurrentTime.Unix(), "https://example.com",
			defaultCurrentTime.Unix(), "https://test.com",
		),
		), 0600)
	if err != nil {
		t.Fatal(err)
	}

	feed := internal.NewTerminalFeed(timeMock, printer, storage)

	root, err := NewRoot("0.1.0", timeMock, printer, storage, feed)
	assert.NoError(t, err)

	os.Args = []string{"cleed", "unfollow", "https://example.com", "https://test.com", "https://not-in-list.com"}

	err = root.Cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, `https://example.com was unfollowed
https://test.com was unfollowed
https://not-in-list.com was not found
`, out.String())

	files, err := os.ReadDir(listsDir)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, files, 1)
	file := files[0]
	assert.Equal(t, "default", file.Name())
	b, err := os.ReadFile(path.Join(listsDir, file.Name()))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "", string(b))
}

func Test_Unfollow_Custom_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeMock := mocks.NewMockTime(ctrl)
	timeMock.EXPECT().Now().Return(defaultCurrentTime).AnyTimes()

	out := new(bytes.Buffer)
	printer := internal.NewPrinter(nil, out, out)
	storage := _storage.NewLocalStorage("cleed_test", timeMock)
	defer localStorageCleanup(t, storage)

	listsDir := path.Join(testConfigDir(t), "lists")
	err := os.MkdirAll(listsDir, 0700)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(path.Join(listsDir, "test"),
		[]byte(fmt.Sprintf("%d %s\n%d %s\n",
			defaultCurrentTime.Unix(), "https://example.com",
			defaultCurrentTime.Unix(), "https://test.com",
		),
		), 0600)
	if err != nil {
		t.Fatal(err)
	}

	feed := internal.NewTerminalFeed(timeMock, printer, storage)

	root, err := NewRoot("0.1.0", timeMock, printer, storage, feed)
	assert.NoError(t, err)

	os.Args = []string{"cleed", "unfollow", "https://example.com", "--list", "test"}

	err = root.Cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, `https://example.com was unfollowed
`, out.String())

	files, err := os.ReadDir(listsDir)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, files, 1)
	file := files[0]
	assert.Equal(t, "test", file.Name())
	b, err := os.ReadFile(path.Join(listsDir, file.Name()))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t,
		fmt.Sprintf("%d %s\n", defaultCurrentTime.Unix(), "https://test.com"),
		string(b),
	)
}

func Test_Unfollow_No_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeMock := mocks.NewMockTime(ctrl)
	timeMock.EXPECT().Now().Return(defaultCurrentTime).AnyTimes()

	out := new(bytes.Buffer)
	printer := internal.NewPrinter(nil, out, out)
	storage := _storage.NewLocalStorage("cleed_test", timeMock)
	defer localStorageCleanup(t, storage)

	feed := internal.NewTerminalFeed(timeMock, printer, storage)

	root, err := NewRoot("0.1.0", timeMock, printer, storage, feed)
	assert.NoError(t, err)

	os.Args = []string{"cleed", "unfollow", "https://example.com"}

	err = root.Cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com was not found\n", out.String())
}

func Test_Unfollow_All_Lists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeMock := mocks.NewMockTime(ctrl)
	timeMock.EXPECT().Now().Return(defaultCurrentTime).AnyTimes()

	out := new(bytes.Buffer)
	printer := internal.NewPrinter(nil, out, out)
	storage := _storage.NewLocalStorage("cleed_test", timeMock)
	defer localStorageCleanup(t, storage)

	listsDir := path.Join(testConfigDir(t), "lists")
	err := os.MkdirAll(listsDir, 0700)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(path.Join(listsDir, "default"),
		[]byte(fmt.Sprintf("%d %s\n",
			defaultCurrentTime.Unix(), "https://example.com",
		),
		), 0600)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(path.Join(listsDir, "test"),
		[]byte(fmt.Sprintf("%d %s\n",
			defaultCurrentTime.Unix(), "https://example.com",
		),
		), 0600)
	if err != nil {
		t.Fatal(err)
	}

	feed := internal.NewTerminalFeed(timeMock, printer, storage)

	root, err := NewRoot("0.1.0", timeMock, printer, storage, feed)
	assert.NoError(t, err)

	os.Args = []string{"cleed", "unfollow", "https://example.com"}

	err = root.Cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com was unfollowed\n", out.String())

	for _, name := range []string{"default", "test"} {
		b, err := os.ReadFile(path.Join(listsDir, name))
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "", string(b), "list %s should be empty", name)
	}
}

func Test_Unfollow_Clean_Cache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeMock := mocks.NewMockTime(ctrl)
	timeMock.EXPECT().Now().Return(defaultCurrentTime).AnyTimes()

	out := new(bytes.Buffer)
	printer := internal.NewPrinter(nil, out, out)
	storage := _storage.NewLocalStorage("cleed_test", timeMock)
	defer localStorageCleanup(t, storage)

	listsDir := path.Join(testConfigDir(t), "lists")
	err := os.MkdirAll(listsDir, 0700)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(path.Join(listsDir, "default"),
		[]byte(fmt.Sprintf("%d %s\n%d %s\n",
			defaultCurrentTime.Unix(), "https://example.com",
			defaultCurrentTime.Unix(), "https://test.com",
		),
		), 0600)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(path.Join(listsDir, "test"),
		[]byte(fmt.Sprintf("%d %s\n",
			defaultCurrentTime.Unix(), "https://test.com",
		),
		), 0600)
	if err != nil {
		t.Fatal(err)
	}

	cacheDir := testCacheDir(t)
	err = os.MkdirAll(cacheDir, 0700)
	if err != nil {
		t.Fatal(err)
	}

	err = storage.SaveCacheInfo(map[string]*_storage.CacheInfoItem{
		"https://example.com": {
			URL:        "https://example.com",
			LastFetch:  defaultCurrentTime,
			ETag:       "etag",
			FetchAfter: time.Unix(0, 0),
		},
		"https://test.com": {
			URL:        "https://test.com",
			LastFetch:  defaultCurrentTime,
			ETag:       "etag",
			FetchAfter: time.Unix(0, 0),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = storage.SaveFeedCache(bytes.NewBuffer([]byte("example")), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}

	err = storage.SaveFeedCache(bytes.NewBuffer([]byte("test")), "https://test.com")
	if err != nil {
		t.Fatal(err)
	}

	feed := internal.NewTerminalFeed(timeMock, printer, storage)

	root, err := NewRoot("0.1.0", timeMock, printer, storage, feed)
	assert.NoError(t, err)

	os.Args = []string{"cleed", "unfollow", "https://example.com", "https://test.com"}

	err = root.Cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, `https://example.com was unfollowed
https://test.com was unfollowed
`, out.String())

	files, err := os.ReadDir(cacheDir)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, files, 1)

	cacheInfo, err := storage.LoadCacheInfo()
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, cacheInfo, 0)

	assert.NoFileExists(t, path.Join(cacheDir, _storage.FeedCacheFileName("https://example.com")))
	assert.NoFileExists(t, path.Join(cacheDir, _storage.FeedCacheFileName("https://test.com")))
}

func Test_Unfollow_Clean_Cache_Single_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeMock := mocks.NewMockTime(ctrl)
	timeMock.EXPECT().Now().Return(defaultCurrentTime).AnyTimes()

	out := new(bytes.Buffer)
	printer := internal.NewPrinter(nil, out, out)
	storage := _storage.NewLocalStorage("cleed_test", timeMock)
	defer localStorageCleanup(t, storage)

	listsDir := path.Join(testConfigDir(t), "lists")
	err := os.MkdirAll(listsDir, 0700)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(path.Join(listsDir, "default"),
		[]byte(fmt.Sprintf("%d %s\n%d %s\n",
			defaultCurrentTime.Unix(), "https://example.com",
			defaultCurrentTime.Unix(), "https://test.com",
		),
		), 0600)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(path.Join(listsDir, "test"),
		[]byte(fmt.Sprintf("%d %s\n",
			defaultCurrentTime.Unix(), "https://test.com",
		),
		), 0600)
	if err != nil {
		t.Fatal(err)
	}

	cacheDir := testCacheDir(t)
	err = os.MkdirAll(cacheDir, 0700)
	if err != nil {
		t.Fatal(err)
	}

	err = storage.SaveCacheInfo(map[string]*_storage.CacheInfoItem{
		"https://example.com": {
			URL:        "https://example.com",
			LastFetch:  defaultCurrentTime,
			ETag:       "etag",
			FetchAfter: time.Unix(0, 0),
		},
		"https://test.com": {
			URL:        "https://test.com",
			LastFetch:  defaultCurrentTime,
			ETag:       "etag",
			FetchAfter: time.Unix(0, 0),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = storage.SaveFeedCache(bytes.NewBuffer([]byte("example")), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}

	err = storage.SaveFeedCache(bytes.NewBuffer([]byte("test")), "https://test.com")
	if err != nil {
		t.Fatal(err)
	}

	feed := internal.NewTerminalFeed(timeMock, printer, storage)

	root, err := NewRoot("0.1.0", timeMock, printer, storage, feed)
	assert.NoError(t, err)

	os.Args = []string{"cleed", "unfollow", "https://example.com", "https://test.com", "--list", "default"}

	err = root.Cmd.Execute()
	assert.NoError(t, err)

	cacheInfo, err := storage.LoadCacheInfo()
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, cacheInfo, 1)
	assert.Equal(t, &_storage.CacheInfoItem{
		URL:        "https://test.com",
		LastFetch:  time.Unix(defaultCurrentTime.Unix(), 0),
		ETag:       "etag",
		FetchAfter: time.Unix(0, 0),
	}, cacheInfo["https://test.com"])

	assert.NoFileExists(t, path.Join(cacheDir, _storage.FeedCacheFileName("https://example.com")))

	b, err := os.ReadFile(path.Join(cacheDir, _storage.FeedCacheFileName("https://test.com")))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "test", string(b))
}
