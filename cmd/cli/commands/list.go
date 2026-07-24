package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	dnsresolve "github.com/francomano/proxydoctor/core/checks/dns_resolve"
	portscan "github.com/francomano/proxydoctor/core/checks/port_scan"
	publicip "github.com/francomano/proxydoctor/core/checks/public_ip"
	tlscert "github.com/francomano/proxydoctor/core/checks/tls_cert"
	"github.com/francomano/proxydoctor/core/engine"
)

var listChecksCmd = &cobra.Command{
	Use:   "list-checks",
	Short: "List all available checks",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create registry with all checks
		registry := engine.NewCheckRegistry()
		registry.Register(publicip.NewPublicIPCheck())
		registry.Register(dnsresolve.NewDNSResolveCheck())
		registry.Register(tlscert.NewTLSCertCheck())
		registry.Register(portscan.NewPortScanCheck())

		checks := registry.ListChecks()

		fmt.Println("📋 Available Checks")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()

		for _, checker := range checks {
			fmt.Printf("ID: %s\n", checker.ID())
			fmt.Printf("Name: %s\n", checker.Name())
			fmt.Printf("Category: %s\n", checker.Category())
			fmt.Printf("Description: %s\n", checker.Description())
			if deps := checker.DependsOn(); len(deps) > 0 {
				fmt.Printf("Depends on: %v\n", deps)
			}
			fmt.Println()
		}

		fmt.Printf("Total: %d checks available\n", len(checks))
		return nil
	},
}
