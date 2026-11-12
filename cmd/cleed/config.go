package cleed

import (
	"github.com/spf13/cobra"
)

func (r *Root) initConfig() {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Display or change configuration",
		Long: `Display or change configuration

Examples:
  # Display configuration
  cleed config

  # Set the user agent
  cleed config --user-agent="My User Agent"

  # Disable styling
  cleed config --styling=2

  # Map color 0 to 230 and color 1 to 213
  cleed config --map-colors=0:230,1:213

  # Remove color mapping for color 0
  cleed config --map-colors=0:

  # Clear all color mappings
  cleed config --map-colors=

  # Display color range. Useful for finding colors to map
  cleed config --color-range

  # Enable run summary
  cleed config --summary=1

  # Set the miniflux token
  cleed config --miniflux-token="your_token_here"

  # Override the config directory (lists, etc.). Pass "" to use the default.
  cleed config --config-path="D:/cleed-data/config"

  # Override the cached RSS feeds directory. Pass "" to use the default.
  cleed config --feeds-path="D:/cleed-data/feeds"
`,
		RunE: r.RunConfig,
	}

	flags := cmd.Flags()
	flags.Uint8("styling", 0, "disable or enable styling (0: default, 1: enable, 2: disable)")
	flags.Uint8("summary", 0, "disable or enable summary (0: disable, 1: enable)")
	flags.String("map-colors", "", "map colors to other colors, e.g. 0:230,1:213. Use --color-range to check available colors")
	flags.Bool("color-range", false, "display color range. Useful for finding colors to map")
	flags.String("user-agent", "", "set the user agent. Setting the value to '-' will not send the user agent")
	flags.Uint("batch-size", 100, "set the batch (queue) size for fetching feeds")
	flags.Uint("timeout", 30, "set the timeout in seconds for fetching feeds")
	flags.Uint8("future-items", 1, "show or hide future items (0: hide, 1: show)")
	flags.String("miniflux-token", "", "set the miniflux token")
	flags.String("config-path", "", "override the directory where lists and other config files are stored (config.json itself stays in the OS user-config dir). Pass an empty string to use the default")
	flags.String("feeds-path", "", "override the directory where cached RSS feeds are stored. Pass an empty string to use the default")

	r.Cmd.AddCommand(cmd)
}

func (r *Root) RunConfig(cmd *cobra.Command, args []string) error {
	if cmd.Flag("styling").Changed {
		styling, err := cmd.Flags().GetUint8("styling")
		if err != nil {
			return err
		}
		return r.feed.SetStyling(styling)
	}
	if cmd.Flag("summary").Changed {
		summary, err := cmd.Flags().GetUint8("summary")
		if err != nil {
			return err
		}
		return r.feed.SetSummary(summary)
	}
	if cmd.Flag("map-colors").Changed {
		return r.feed.UpdateColorMap(cmd.Flag("map-colors").Value.String())
	}
	if cmd.Flag("color-range").Changed {
		r.feed.DisplayColorRange()
		return nil
	}
	if cmd.Flag("user-agent").Changed {
		return r.feed.SetUserAgent(cmd.Flag("user-agent").Value.String())
	}
	if cmd.Flag("batch-size").Changed {
		batchSize, err := cmd.Flags().GetUint("batch-size")
		if err != nil {
			return err
		}
		return r.feed.SetBatchSize(batchSize)
	}
	if cmd.Flag("timeout").Changed {
		timeout, err := cmd.Flags().GetUint("timeout")
		if err != nil {
			return err
		}
		return r.feed.SetTimeout(timeout)
	}
	if cmd.Flag("future-items").Changed {
		value, err := cmd.Flags().GetUint8("future-items")
		if err != nil {
			return err
		}
		return r.feed.UpdateFutureItems(value)
	}
	if cmd.Flag("miniflux-token").Changed {
		token := cmd.Flag("miniflux-token").Value.String()
		return r.feed.SetMinifluxToken(token)
	}
	if cmd.Flag("config-path").Changed {
		return r.feed.SetConfigPath(cmd.Flag("config-path").Value.String())
	}
	if cmd.Flag("feeds-path").Changed {
		return r.feed.SetFeedsPath(cmd.Flag("feeds-path").Value.String())
	}
	return r.feed.DisplayConfig()
}
