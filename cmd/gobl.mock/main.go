package main

import (
	"encoding/json"
	"fmt"
	"os"

	mock "github.com/pablocaeg/gobl.mock"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	date    = ""
)

func main() {
	if err := root().Execute(); err != nil {
		os.Exit(1)
	}
}

func root() *cobra.Command {
	var (
		regime     string
		addon      string
		lines      int
		credit     bool
		simplified bool
		seed       int64
		hasSeed    bool
		output     string
	)

	cmd := &cobra.Command{
		Use:     "gobl.mock",
		Short:   "Generate realistic, valid GOBL invoices",
		Version: versionString(),
		RunE: func(_ *cobra.Command, _ []string) error {
			opts := []mock.Option{
				mock.WithRegime(l10n.TaxCountryCode(regime)),
				mock.WithLines(lines),
			}

			if addon != "" {
				opts = append(opts, mock.WithAddon(cbc.Key(addon)))
			}
			if credit {
				opts = append(opts, mock.WithCredit())
			}
			if simplified {
				opts = append(opts, mock.WithSimplified())
			}
			if hasSeed {
				opts = append(opts, mock.WithSeed(seed))
			}

			env, err := mock.Envelope(opts...)
			if err != nil {
				return err
			}

			data, err := json.MarshalIndent(env, "", "  ")
			if err != nil {
				return fmt.Errorf("marshaling JSON: %w", err)
			}

			if output != "" {
				if err := os.WriteFile(output, data, 0644); err != nil {
					return fmt.Errorf("writing file: %w", err)
				}
				fmt.Fprintf(os.Stderr, "Written to %s\n", output)
				return nil
			}

			fmt.Println(string(data))
			return nil
		},
	}

	cmd.Flags().StringVar(&regime, "regime", "ES", "Tax regime country code (ES, DE, MX)")
	cmd.Flags().StringVar(&addon, "addon", "", "Addon to apply (e.g. es-facturae-v3, de-xrechnung-v3)")
	cmd.Flags().IntVar(&lines, "lines", 3, "Number of line items")
	cmd.Flags().BoolVar(&credit, "credit", false, "Generate a credit note")
	cmd.Flags().BoolVar(&simplified, "simplified", false, "Generate a simplified invoice")
	cmd.Flags().Int64Var(&seed, "seed", 0, "Random seed for reproducible output")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path")

	// Track if --seed was explicitly set.
	cmd.PreRun = func(_ *cobra.Command, _ []string) {
		hasSeed = cmd.Flags().Changed("seed")
	}

	return cmd
}

func versionString() string {
	if date != "" {
		return fmt.Sprintf("%s (%s)", version, date)
	}
	return version
}
