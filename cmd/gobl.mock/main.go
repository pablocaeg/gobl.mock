// Package main provides the gobl.mock CLI tool.
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
		invType    string
		scenario   string
		lines      int
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
			if invType != "" {
				opts = append(opts, mock.WithType(cbc.Key(invType)))
			}
			if scenario != "" {
				opts = append(opts, mock.WithScenario(cbc.Key(scenario)))
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

	cmd.Flags().StringVar(&regime, "regime", "ES", "Tax regime country code")
	cmd.Flags().StringVar(&addon, "addon", "", "Addon to apply (e.g. es-facturae-v3)")
	cmd.Flags().StringVar(&invType, "type", "", "Invoice type: standard, credit-note, corrective, debit-note, proforma")
	cmd.Flags().StringVar(&scenario, "scenario", "", "Scenario: hotel, freelance, reverse-charge, restaurant, ecommerce")
	cmd.Flags().IntVar(&lines, "lines", 3, "Number of line items")
	cmd.Flags().BoolVar(&simplified, "simplified", false, "Generate a simplified invoice")
	cmd.Flags().Int64Var(&seed, "seed", 0, "Random seed for reproducible output")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path")

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
