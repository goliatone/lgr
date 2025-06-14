package cmd

import (
	"fmt"
	"strings"

	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(qrCommand)
	qrCommand.Flags().BoolP("invert", "i", false, "Inver QR colors")

}

var qrCommand = &cobra.Command{
	Use:     "qr [text]",
	Aliases: []string{"q"},
	Short:   "Generate QR code from given text",
	Example: `
	lgr qr "https://example.com"
	lgr qr --invert "Large QR code"
	lgr qr -i "Extra small QR code"
	`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		invert, _ := cmd.Flags().GetBool("invert")
		return printQR(args, invert)
	},
}

func printQR(args []string, invert bool) error {
	text := strings.Join(args, "")

	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		return fmt.Errorf("failed to generate QR code: %w", err)
	}

	qrString := qr.ToSmallString(invert)
	fmt.Println(qrString)
	return nil
}
