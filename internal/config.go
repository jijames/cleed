package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/radulucut/cleed/internal/utils"
)

func (f *TerminalFeed) SetVersion(version string) {
	f.version = version
}

func (f *TerminalFeed) SetDefaultExploreRepository(repository string) {
	f.defaultExploreRepository = repository
}

func (f *TerminalFeed) DisplayConfig() error {
	config, err := f.storage.LoadConfig()
	if err != nil {
		return utils.NewInternalError("failed to load config: " + err.Error())
	}
	f.printer.Println("User-Agent:", config.UserAgent)
	f.printer.Println("Timeout:", config.Timeout)
	f.printer.Println("Batch size:", config.BatchSize)
	styling := "default"
	if config.Styling == 0 {
		styling = "enabled"
	} else if config.Styling == 1 {
		styling = "disabled"
	}
	f.printer.Println("Styling:", styling)
	f.printer.Print("Color map:")
	for k, v := range config.ColorMap {
		f.printer.Printf(" %d:%d", k, v)
	}
	f.printer.Println()
	summary := "disabled"
	if config.Summary == 1 {
		summary = "enabled"
	}
	f.printer.Println("Summary:", summary)
	futureItems := "show"
	if config.HideFutureItems {
		futureItems = "hide"
	}
	f.printer.Println("Future items:", futureItems)
	if config.ConfigPath != "" {
		f.printer.Println("Config path:", config.ConfigPath)
	} else {
		f.printer.Println("Config path:")
	}
	if config.FeedsPath != "" {
		f.printer.Println("Feeds path:", config.FeedsPath)
	} else {
		f.printer.Println("Feeds path:")
	}
	if config.MinifluxToken != "" {
		f.printer.Println("Miniflux token:", "******"+config.MinifluxToken[len(config.MinifluxToken)-6:])
	} else {
		f.printer.Println("Miniflux token:")
	}
	return nil
}

// SetConfigPath updates the override for the configuration directory (lists,
// etc.). Pass an empty string to clear the override and fall back to the
// default location (e.g. %AppData%/cleed on Windows). Note: config.json
// itself always lives in the bootstrap location.
func (f *TerminalFeed) SetConfigPath(p string) error {
	config, err := f.storage.LoadConfig()
	if err != nil {
		return utils.NewInternalError("failed to load config: " + err.Error())
	}
	config.ConfigPath = p
	if err := f.storage.SaveConfig(); err != nil {
		return utils.NewInternalError("failed to save config: " + err.Error())
	}
	if p == "" {
		f.printer.Println("config path was cleared (using default)")
	} else {
		f.printer.Println("config path was updated")
	}
	return nil
}

// SetFeedsPath updates the override for the cached RSS feeds directory.
// Pass an empty string to clear the override and fall back to the default
// location (e.g. %LocalAppData%/cleed on Windows).
func (f *TerminalFeed) SetFeedsPath(p string) error {
	config, err := f.storage.LoadConfig()
	if err != nil {
		return utils.NewInternalError("failed to load config: " + err.Error())
	}
	config.FeedsPath = p
	if err := f.storage.SaveConfig(); err != nil {
		return utils.NewInternalError("failed to save config: " + err.Error())
	}
	if p == "" {
		f.printer.Println("feeds path was cleared (using default)")
	} else {
		f.printer.Println("feeds path was updated")
	}
	return nil
}

func (f *TerminalFeed) SetTimeout(timeout uint) error {
	config, err := f.storage.LoadConfig()
	if err != nil {
		return utils.NewInternalError("failed to load config: " + err.Error())
	}
	config.Timeout = timeout
	err = f.storage.SaveConfig()
	if err != nil {
		return utils.NewInternalError("failed to save config: " + err.Error())
	}
	f.printer.Println("timeout was updated")
	return nil
}

func (f *TerminalFeed) SetBatchSize(batchSize uint) error {
	config, err := f.storage.LoadConfig()
	if err != nil {
		return utils.NewInternalError("failed to load config: " + err.Error())
	}
	config.BatchSize = batchSize
	err = f.storage.SaveConfig()
	if err != nil {
		return utils.NewInternalError("failed to save config: " + err.Error())
	}
	f.printer.Println("batch size was updated")
	return nil
}

func (f *TerminalFeed) SetUserAgent(agent string) error {
	config, err := f.storage.LoadConfig()
	if err != nil {
		return utils.NewInternalError("failed to load config: " + err.Error())
	}
	config.UserAgent = agent
	err = f.storage.SaveConfig()
	if err != nil {
		return utils.NewInternalError("failed to save config: " + err.Error())
	}
	f.printer.Println("User-Agent was updated")
	return nil
}

func (f *TerminalFeed) SetStyling(v uint8) error {
	config, err := f.storage.LoadConfig()
	if err != nil {
		return utils.NewInternalError("failed to load config: " + err.Error())
	}
	if v > 2 {
		return utils.NewInternalError("invalid value for styling")
	}
	config.Styling = v
	err = f.storage.SaveConfig()
	if err != nil {
		return utils.NewInternalError("failed to save config: " + err.Error())
	}
	f.printer.Println("styling was updated")
	return nil
}

func (f *TerminalFeed) SetSummary(v uint8) error {
	config, err := f.storage.LoadConfig()
	if err != nil {
		return utils.NewInternalError("failed to load config: " + err.Error())
	}
	if v > 1 {
		return utils.NewInternalError("invalid value for summary")
	}
	config.Summary = v
	err = f.storage.SaveConfig()
	if err != nil {
		return utils.NewInternalError("failed to save config: " + err.Error())
	}
	f.printer.Println("summary was updated")
	return nil
}

func (f *TerminalFeed) UpdateColorMap(mappings string) error {
	config, err := f.storage.LoadConfig()
	if err != nil {
		return utils.NewInternalError("failed to load config: " + err.Error())
	}
	if mappings == "" {
		config.ColorMap = make(map[uint8]uint8)
	} else {
		colors := strings.Split(mappings, ",")
		for i := range colors {
			parts := strings.Split(colors[i], ":")
			if len(parts) == 0 {
				return utils.NewInternalError("failed to parse color mapping: " + colors[i])
			}
			left, err := strconv.Atoi(parts[0])
			if err != nil {
				return utils.NewInternalError("failed to parse color mapping: " + parts[0])
			}
			if len(parts) == 1 || parts[1] == "" {
				delete(config.ColorMap, uint8(left))
			} else {
				right, err := strconv.Atoi(parts[1])
				if err != nil {
					return utils.NewInternalError("failed to parse color mapping: " + parts[1])
				}
				config.ColorMap[uint8(left)] = uint8(right)
			}
		}
	}
	err = f.storage.SaveConfig()
	if err != nil {
		return utils.NewInternalError("failed to save config: " + err.Error())
	}
	f.printer.Println("color map updated")
	return nil
}

func (f *TerminalFeed) UpdateFutureItems(value uint8) error {
	config, err := f.storage.LoadConfig()
	if err != nil {
		return utils.NewInternalError("failed to load config: " + err.Error())
	}
	if value == 0 {
		config.HideFutureItems = true
	} else if value == 1 {
		config.HideFutureItems = false
	} else {
		return utils.NewInternalError("invalid value for future items")
	}
	err = f.storage.SaveConfig()
	if err != nil {
		return utils.NewInternalError("failed to save config: " + err.Error())
	}
	f.printer.Println("future items was updated")
	return nil
}

func (f *TerminalFeed) DisplayColorRange() {
	styling := f.printer.GetStyling()
	f.printer.SetStyling(true)
	for i := 0; i < 256; i++ {
		f.printer.Print(f.printer.ColorForeground(fmt.Sprintf("%d ", i), uint8(i)))
	}
	f.printer.Println()
	f.printer.SetStyling(styling)
}

func (f *TerminalFeed) SetMinifluxToken(token string) error {
	config, err := f.storage.LoadConfig()
	if err != nil {
		return utils.NewInternalError("failed to load config: " + err.Error())
	}
	config.MinifluxToken = token
	err = f.storage.SaveConfig()
	if err != nil {
		return utils.NewInternalError("failed to save config: " + err.Error())
	}
	f.printer.Println("miniflux token was updated")
	return nil
}
