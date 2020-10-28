package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
	"io"
)

type selectOptions struct {
	File     string
	Parser   string
	Selector string
	Reader   io.Reader
	Writer   io.Writer
}

func runSelectCommand(opts selectOptions, cmd *cobra.Command) error {
	parser, err := getParser(opts.File, opts.Parser)
	if err != nil {
		return err
	}
	rootNode, err := getRootNode(getRootNodeOpts{
		File:   opts.File,
		Parser: parser,
		Reader: opts.Reader,
	}, cmd)
	if err != nil {
		return err
	}

	var res *dasel.Node
	if opts.Selector == "." {
		res = rootNode
	} else {
		res, err = rootNode.Query(opts.Selector)
		if err != nil {
			return fmt.Errorf("could not query node: %w", err)
		}
	}

	if opts.Writer == nil {
		opts.Writer = cmd.OutOrStdout()
	}

	_, _ = fmt.Fprintf(opts.Writer, "%v\n", res.InterfaceValue())

	return nil
}

func selectCommand() *cobra.Command {
	var fileFlag, selectorFlag, parserFlag string

	cmd := &cobra.Command{
		Use:   "select -f <file> -p <json,yaml> -s <selector>",
		Short: "Select properties from the given file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if selectorFlag == "" && len(args) > 0 {
				selectorFlag = args[0]
				args = args[1:]
			}
			return runSelectCommand(selectOptions{
				File:     fileFlag,
				Parser:   parserFlag,
				Selector: selectorFlag,
			}, cmd)
		},
	}

	cmd.Flags().StringVarP(&fileFlag, "file", "f", "", "The file to query.")
	cmd.Flags().StringVarP(&selectorFlag, "selector", "s", "", "The selector to use when querying the data structure.")
	cmd.Flags().StringVarP(&parserFlag, "parser", "p", "", "The parser to use with the given file.")

	_ = cmd.MarkFlagFilename("file")

	return cmd
}
