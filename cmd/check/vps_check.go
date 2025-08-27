package check

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/TheoBrigitte/kimsufi-notifier/cmd/flag"
	"github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi"
)

// runnerVPS handles VPS plan checking
func runnerVPS(cmd *cobra.Command, k *kimsufi.Service, planCode string, datacenters []string) error {
	countryCode := cmd.Flag(flag.CountryFlagName).Value.String()

	// Get VPS availability data
	vpsAvailabilities, err := k.GetVPSAvailabilities(planCode, countryCode, "")
	if err != nil {
		return fmt.Errorf("failed to get VPS availability for %s: %w", planCode, err)
	}

	// Display VPS availability
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "planCode\tdatacenter\tstatus\tlinuxStatus\twindowsStatus")
	fmt.Fprintln(w, "--------\t----------\t------\t-----------\t-------------")

	nothingAvailable := true
	for _, dc := range vpsAvailabilities.Datacenters {
		// Filter by requested datacenters if specified
		if len(datacenters) > 0 {
			found := false
			for _, requestedDC := range datacenters {
				if strings.EqualFold(dc.Datacenter, requestedDC) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Check if this datacenter has any availability
		isAvailable := dc.Status == "available" || dc.LinuxStatus == "available" || dc.WindowsStatus == "available"
		if isAvailable {
			nothingAvailable = false
		}

		// Display datacenter availability
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			planCode,
			dc.Datacenter,
			dc.Status,
			dc.LinuxStatus,
			dc.WindowsStatus)
	}
	w.Flush()

	if nothingAvailable {
		message := datacenterAvailableMessageFormatter(datacenters)
		log.Printf("%s is not available in %s\n", planCode, message)
		os.Exit(1)
	}

	return nil
}
