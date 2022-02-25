package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/ksuid"
	"github.com/spf13/cobra"
)

var (
	generateCount       = 1
	generateResource    = "example"
	generateEnvironment = ksuid.Production
)

// RootCmd is the initial entrypoint where all commands are mounted.
var RootCmd = &cobra.Command{
	Use:   "ksuid <command>",
	Short: "utility to parse and generate ksuid",
}

// GenerateCommand is executed to generate one or more ksuid.
var GenerateCommand = &cobra.Command{
	Use:     "generate [options]",
	Aliases: []string{"gen", "g"},
	Short:   "generate one or more ksuid",

	Run: func(cmd *cobra.Command, args []string) {
		ksuid.SetEnvironment(generateEnvironment)

		for n := 0; n < generateCount; n++ {
			id := ksuid.Generate(generateResource)

			fmt.Println(id.String())
		}
	},
}

// ParseCommand is executed to parse ksuids given as a command line argument.
var ParseCommand = &cobra.Command{
	Use:     "parse <ksuid [ksuid...]>",
	Aliases: []string{"p"},
	Short:   "parse ksuids given on the command line",

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("at least one ksuid required to parse")
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			id, err := ksuid.Parse(arg)
			if err != nil {
				fmt.Printf("ID:    %s\nError: %s\n\n", arg, err)
				continue
			}

			fmt.Printf(
				"ID:          %s\nResource:    %s\nEnvironment: %s\nTimestamp:   %s\n",
				arg, id.Resource, id.Environment, time.Unix(int64(id.Timestamp), 0).Format(time.RFC3339),
			)

			iid := id.InstanceID

			switch iid.SchemeData {
			case 'H':
				fmt.Printf(
					"Machine ID:  %s\nProcess ID:  %d\n",
					net.HardwareAddr(iid.BytesData[:6]), binary.BigEndian.Uint16(iid.BytesData[6:]),
				)

			case 'D':
				fmt.Printf(
					"Docker ID:   %x\n",
					iid.Bytes(),
				)

			case 'R':
				fmt.Printf(
					"Random ID:   %x\n",
					iid.Bytes(),
				)

			default:
				fmt.Printf(
					"Node ID:     %x\n",
					iid.Bytes(),
				)
			}

			fmt.Printf("Sequence ID: %d\n\n", id.SequenceID)
		}
	},
}

func init() {
	GenerateCommand.Flags().IntVarP(&generateCount, "count", "n", 1, "number of ksuid to generate")
	GenerateCommand.Flags().StringVarP(&generateResource, "resource", "r", "example", "resource prefix")
	GenerateCommand.Flags().StringVarP(&generateEnvironment, "environment", "e", ksuid.Production, "environment prefix")

	RootCmd.AddCommand(GenerateCommand, ParseCommand)
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
