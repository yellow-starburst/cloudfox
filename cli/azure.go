package cli

import (
	"log"

	"github.com/BishopFox/cloudfox/azure"
	"github.com/spf13/cobra"
)

var (
	AzTenantID        string
	AzSubscription    string
	AzRGName          string
	AzOutputFormat    string
	AzOutputDirectory string
	AzVerbosity       int
	AzWrapTable       bool
	AzMergedTable     bool

	AzCommands = &cobra.Command{
		Use:     "azure",
		Aliases: []string{"az"},
		Long:    `See "Available Commands" for Azure Modules below`,
		Short:   "See \"Available Commands\" for Azure Modules below",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	AzWhoamiListRGsAlso bool
	AzWhoamiCommand     = &cobra.Command{
		Use:     "whoami",
		Aliases: []string{},
		Short:   "Display available Azure CLI sessions",
		Long: `
Display Available Azure CLI Sessions:
./cloudfox az whoami`,
		Run: func(cmd *cobra.Command, args []string) {
			err := azure.AzWhoamiCommand(AzOutputDirectory, cmd.Root().Version, AzWrapTable, AzVerbosity, AzWhoamiListRGsAlso)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	AzInventoryCommand = &cobra.Command{
		Use:     "inventory",
		Aliases: []string{"inv"},
		Short:   "Display an inventory table of all resources per location",
		Long: `
Enumerate inventory for a specific tenant:
./cloudfox az inventory --tenant TENANT_ID

Enumerate inventory for a specific subscription:
./cloudfox az inventory --subscription SUBSCRIPTION_ID
`,
		Run: func(cmd *cobra.Command, args []string) {
			err := azure.AzInventoryCommand(AzTenantID, AzSubscription, AzOutputDirectory, cmd.Root().Version, AzVerbosity, AzWrapTable, AzMergedTable)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	AzRBACCommand = &cobra.Command{
		Use:     "rbac",
		Aliases: []string{},
		Short:   "Display role assignments for Azure principals",
		Long: `
Enumerate role assignments for a specific tenant:
./cloudfox az rbac --tenant TENANT_ID

Enumerate role assignments for a specific subscription:
./cloudfox az rbac --subscription SUBSCRIPTION_ID
`,
		Run: func(cmd *cobra.Command, args []string) {
			err := azure.AzRBACCommand(AzTenantID, AzSubscription, AzOutputFormat, AzOutputDirectory, cmd.Root().Version, AzVerbosity, AzWrapTable, AzMergedTable)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	AzVMsCommand = &cobra.Command{
		Use:     "vms",
		Aliases: []string{"vms", "virtualmachines"},
		Short:   "Enumerates Azure Compute virtual machines",
		Long: `
Enumerate VMs for a specific tenant:
./cloudfox az vms --tenant TENANT_ID

Enumerate VMs for a specific subscription:
./cloudfox az vms --subscription SUBSCRIPTION_ID`,
		Run: func(cmd *cobra.Command, args []string) {
			err := azure.AzVMsCommand(AzTenantID, AzSubscription, AzOutputFormat, AzOutputDirectory, cmd.Root().Version, AzVerbosity, AzWrapTable, AzMergedTable)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	AzStorageCommand = &cobra.Command{
		Use:     "storage",
		Aliases: []string{},
		Short:   "Enumerates azure storage accounts",
		Long: `
Enumerate storage accounts for a specific tenant:
./cloudfox az storage --tenant TENANT_ID

Enumerate storage accounts for a specific subscription:
./cloudfox az storage --subscription SUBSCRIPTION_ID
`,
		Run: func(cmd *cobra.Command, args []string) {
			err := azure.AzStorageCommand(AzTenantID, AzSubscription, AzOutputFormat, AzOutputDirectory, cmd.Root().Version, AzVerbosity, AzWrapTable, AzMergedTable)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	AzRunbooksCommand = &cobra.Command{
		Use:     "runbooks",
		Aliases: []string{},
		Short:   "Downloads all run-books in the account",
		Long: `
Downloads all run-books for a specific tenant:
./cloudfox az runbooks --tenant TENANT_ID

Downloads all run-books for a specific subscription:
./cloudfox az runbooks --subscription SUBSCRIPTION_ID
`,
		Run: func(cmd *cobra.Command, args []string) {
			err := azure.AzRunbooksCommand(AzTenantID, AzSubscription, AzOutputDirectory, cmd.Root().Version, AzVerbosity, AzWrapTable, AzMergedTable)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	AzAllVulnChecksCommand = &cobra.Command{
		Use:     "all-vulnchecks",
		Aliases: []string{},
		Short:   "Runs all vulnerability checks",
		Long:    `Executes all vulnerability checks to collect and display information from all supported Azure services.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Run the runbooks command as part of the all-vulnchecks command
			AzRunbooksCommand.Run(cmd, args)
		},
	}
	AzAllEnumChecksCommand = &cobra.Command{
		Use:     "all-enumchecks",
		Aliases: []string{},
		Short:   "Runs all enumeration checks",
		Long:    `Executes all enumeration checks to collect and display information from all supported Azure services.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Run the enumeration commands
			AzInventoryCommand.Run(cmd, args)
			AzRBACCommand.Run(cmd, args)
			AzVMsCommand.Run(cmd, args)
			AzStorageCommand.Run(cmd, args)
		},
	}
	AzLogger = internal.NewLogger()
)

func init() {
	AzWhoamiCommand.Flags().BoolVarP(&AzWhoamiListRGsAlso, "list-rgs", "l", false, "Drill down to the resource group level")

	// Global flags
	AzCommands.PersistentFlags().StringVarP(&AzOutputFormat, "output", "o", "all", "[\"table\" | \"csv\" | \"all\" ]")
	AzCommands.PersistentFlags().StringVar(&AzOutputDirectory, "outdir", defaultOutputDir, "Output Directory ")
	AzCommands.PersistentFlags().IntVarP(&AzVerbosity, "verbosity", "v", 2, "1 = Print control messages only\n2 = Print control messages, module output\n3 = Print control messages, module output, and loot file output\n")
	AzCommands.PersistentFlags().StringVarP(&AzTenantID, "tenant", "t", "", "Tenant name")
	AzCommands.PersistentFlags().StringVarP(&AzSubscription, "subscription", "s", "", "Subscription ID or Name")
	AzCommands.PersistentFlags().StringVarP(&AzRGName, "resource-group", "g", "", "Resource Group name")
	AzCommands.PersistentFlags().BoolVarP(&AzWrapTable, "wrap", "w", false, "Wrap table to fit in terminal (complicates grepping)")
	AzCommands.PersistentFlags().BoolVarP(&AzMergedTable, "merged-table", "m", false, "Writes a single table for all subscriptions in the tenant. Default writes a table per subscription.")

	AzCommands.AddCommand(
		AzWhoamiCommand,
		AzRBACCommand,
		AzVMsCommand,
		AzStorageCommand,
		AzInventoryCommand,
		AzRunbooksCommand,    // Add the new runbooks command
		AzAllVulnChecksCommand, // Add the new all-vulnchecks command
		AzAllEnumChecksCommand, // Add the new all-enumchecks command
	)
}
