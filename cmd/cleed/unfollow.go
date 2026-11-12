package cleed

import "github.com/spf13/cobra"

func (r *Root) initUnfollow() {
	cmd := &cobra.Command{
		Use:   "unfollow [feed]",
		Short: "Unfollow a feed",
		Long: `Unfollow a feed

Removes the feed from every list that contains it. Use --list to remove from one list only.

Examples:
  # Unfollow a feed (removed from all lists)
  cleed unfollow https://example.com/feed.xml

  # Unfollow multiple feeds
  cleed unfollow https://example.com/feed.xml https://example2.com/feed

  # Remove feeds from a specific list only
  cleed unfollow https://example.com/feed.xml --list mylist
`,

		RunE: r.RunUnfollow,
		Args: cobra.MinimumNArgs(1),
	}

	flags := cmd.Flags()
	flags.StringP("list", "L", "", "remove the feed only from this list (default: search all lists)")

	r.Cmd.AddCommand(cmd)
}

func (r *Root) RunUnfollow(cmd *cobra.Command, args []string) error {
	list, err := cmd.Flags().GetString("list")
	if err != nil {
		return err
	}
	return r.feed.Unfollow(args, list)
}
