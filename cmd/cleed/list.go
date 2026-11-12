package cleed

import (
	"github.com/spf13/cobra"
)

func (r *Root) initList() {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show all lists, feeds in a list or manage lists",
		Long: `Show all lists, feeds in a list or manage lists

Examples:
  # Show all lists
  cleed list

  # Show all feeds across every list
  cleed list --feeds

  # Show all feeds in a list
  cleed list mylist

  # Rename a list
  cleed list mylist --rename newlist

  # Merge a list. Move all feeds from anotherlist to mylist and remove anotherlist
  cleed list mylist --merge anotherlist

  # Remove a list
  cleed list mylist --remove

  # Import feeds from a file
  cleed list mylist --import-from-file feeds.txt

  # Import feeds from an OPML file into a list
  cleed list mylist --import-from-opml feeds.opml

  # Import feeds from an OPML file into multiple lists
  cleed list --import-from-opml feeds.opml

  # Export feeds to a file
  cleed list mylist --export-to-file feeds.txt

  # Export feeds from a list to a file
  cleed list mylist --export-to-opml feeds.opml

  # Export all feeds to an OPML file grouped by lists
  cleed list --export-to-opml feeds.opml

  # Export only cached feeds to an OPML file
  cleed list --export-to-opml feeds.opml -C
`,

		RunE: r.RunList,
		Args: cobra.MaximumNArgs(1),
	}

	flags := cmd.Flags()
	flags.Bool("feeds", false, "show all feeds from every list")
	flags.String("rename", "", "rename a list")
	flags.String("merge", "", "merge a list")
	flags.Bool("remove", false, "remove a list")
	flags.String("import-from-file", "", "import feeds from a file. Newline separated URLs")
	flags.String("import-from-opml", "", "import feeds from an OPML file")
	flags.String("export-to-file", "", "export feeds to a file. Newline separated URLs")
	flags.String("export-to-opml", "", "export feeds to an OPML file")
	flags.BoolP("cached-only", "C", false, "include only cached feeds (opml export)")

	r.Cmd.AddCommand(cmd)
}

func (r *Root) RunList(cmd *cobra.Command, args []string) error {
	list := ""
	if len(args) > 0 {
		list = args[0]
	}
	importFromOPML := cmd.Flag("import-from-opml").Value.String()
	if importFromOPML != "" {
		return r.feed.ImportFromOPML(importFromOPML, list)
	}
	exportToOPML := cmd.Flag("export-to-opml").Value.String()
	if exportToOPML != "" {
		cachedOnly, err := cmd.Flags().GetBool("cached-only")
		if err != nil {
			return err
		}
		return r.feed.ExportToOPML(exportToOPML, list, cachedOnly)
	}
	if list == "" {
		showFeeds, err := cmd.Flags().GetBool("feeds")
		if err != nil {
			return err
		}
		if showFeeds {
			return r.feed.ListAllFeeds()
		}
		return r.feed.Lists()
	}
	rename := cmd.Flag("rename").Value.String()
	if rename != "" {
		return r.feed.RenameList(list, rename)
	}
	merge := cmd.Flag("merge").Value.String()
	if merge != "" {
		return r.feed.MergeLists(list, merge)
	}
	if cmd.Flag("remove").Changed {
		return r.feed.RemoveList(list)
	}
	importFromFile := cmd.Flag("import-from-file").Value.String()
	if importFromFile != "" {
		return r.feed.ImportFromFile(importFromFile, list)
	}
	exportToFile := cmd.Flag("export-to-file").Value.String()
	if exportToFile != "" {
		return r.feed.ExportToFile(exportToFile, list)
	}
	return r.feed.ListFeeds(list)
}
