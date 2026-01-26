package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/hastefuI/ffrelayctl/api"
)

const (
	FormatText = "text"
	FormatJSON = "json"
)

type CombinedMask struct {
	Type string      `json:"type"`
	Mask interface{} `json:"mask"`
}

func ValidFormats() []string {
	return []string{FormatText, FormatJSON}
}

func IsValidFormat(format string) bool {
	for _, f := range ValidFormats() {
		if f == format {
			return true
		}
	}
	return false
}

func Print(format string, v interface{}) error {
	return Fprint(os.Stdout, format, v)
}

func Fprint(w io.Writer, format string, v interface{}) error {
	switch format {
	case FormatJSON:
		return printJSON(w, v)
	case FormatText:
		return printText(w, v)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

func printJSON(w io.Writer, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("error formatting output: %v", err)
	}
	fmt.Fprintln(w, string(data))
	return nil
}

func printText(w io.Writer, v interface{}) error {
	switch data := v.(type) {
	case []api.User:
		return printUsers(w, data)
	case []api.Profile:
		return printProfiles(w, data)
	case api.Profile:
		return printProfiles(w, []api.Profile{data})
	case []CombinedMask:
		return printCombinedMasks(w, data)
	case []api.RelayAddress:
		return printRelayAddresses(w, data)
	case api.RelayAddress:
		return printRelayAddresses(w, []api.RelayAddress{data})
	case []api.DomainAddress:
		return printDomainAddresses(w, data)
	case api.DomainAddress:
		return printDomainAddresses(w, []api.DomainAddress{data})
	case []api.RelayNumber:
		return printRelayNumbers(w, data)
	case api.RelayNumber:
		return printRelayNumbers(w, []api.RelayNumber{data})
	case []api.InboundContact:
		return printInboundContacts(w, data)
	case api.InboundContact:
		return printInboundContacts(w, []api.InboundContact{data})
	case []api.RealPhone:
		return printRealPhones(w, data)
	case api.RealPhone:
		return printRealPhones(w, []api.RealPhone{data})
	case api.RelayNumberSuggestions:
		return printRelayNumberSuggestions(w, data)
	case []api.PhoneNumberOption:
		return printPhoneNumberOptions(w, data)
	default:
		return printJSON(w, v)
	}
}

func printUsers(w io.Writer, users []api.User) error {
	if len(users) == 0 {
		fmt.Fprintln(w, "No users found.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "EMAIL")
	for _, u := range users {
		fmt.Fprintln(tw, u.Email)
	}
	return tw.Flush()
}

func printProfiles(w io.Writer, profiles []api.Profile) error {
	if len(profiles) == 0 {
		fmt.Fprintln(w, "No profiles found.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tSUBDOMAIN\tPREMIUM\tPHONE\tFORWARDED\tBLOCKED\tREPLIED")
	for _, p := range profiles {
		subdomain := "-"
		if p.Subdomain != nil {
			subdomain = *p.Subdomain
		}
		fmt.Fprintf(tw, "%d\t%s\t%t\t%t\t%d\t%d\t%d\n",
			p.ID,
			subdomain,
			p.HasPremium,
			p.HasPhone,
			p.EmailsForwarded,
			p.EmailsBlocked,
			p.EmailsReplied,
		)
	}
	return tw.Flush()
}

func printCombinedMasks(w io.Writer, masks []CombinedMask) error {
	if len(masks) == 0 {
		fmt.Fprintln(w, "No masks found.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tTYPE\tADDRESS\tENABLED\tDESCRIPTION\tFORWARDED\tBLOCKED")
	for _, m := range masks {
		switch mask := m.Mask.(type) {
		case api.RelayAddress:
			desc := truncate(mask.Description, 30)
			fmt.Fprintf(tw, "%d\t%s\t%s\t%t\t%s\t%d\t%d\n",
				mask.ID,
				m.Type,
				mask.FullAddress,
				mask.Enabled,
				desc,
				mask.NumForwarded,
				mask.NumBlocked,
			)
		case api.DomainAddress:
			desc := truncate(mask.Description, 30)
			fmt.Fprintf(tw, "%d\t%s\t%s\t%t\t%s\t%d\t%d\n",
				mask.ID,
				m.Type,
				mask.FullAddress,
				mask.Enabled,
				desc,
				mask.NumForwarded,
				mask.NumBlocked,
			)
		}
	}
	return tw.Flush()
}

func printRelayAddresses(w io.Writer, addresses []api.RelayAddress) error {
	if len(addresses) == 0 {
		fmt.Fprintln(w, "No random masks found.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tADDRESS\tENABLED\tDESCRIPTION\tFORWARDED\tBLOCKED")
	for _, a := range addresses {
		desc := truncate(a.Description, 30)
		fmt.Fprintf(tw, "%d\t%s\t%t\t%s\t%d\t%d\n",
			a.ID,
			a.FullAddress,
			a.Enabled,
			desc,
			a.NumForwarded,
			a.NumBlocked,
		)
	}
	return tw.Flush()
}

func printDomainAddresses(w io.Writer, addresses []api.DomainAddress) error {
	if len(addresses) == 0 {
		fmt.Fprintln(w, "No custom domain masks found.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tADDRESS\tENABLED\tDESCRIPTION\tFORWARDED\tBLOCKED")
	for _, a := range addresses {
		desc := truncate(a.Description, 30)
		fmt.Fprintf(tw, "%d\t%s\t%t\t%s\t%d\t%d\n",
			a.ID,
			a.FullAddress,
			a.Enabled,
			desc,
			a.NumForwarded,
			a.NumBlocked,
		)
	}
	return tw.Flush()
}

func printRelayNumbers(w io.Writer, numbers []api.RelayNumber) error {
	if len(numbers) == 0 {
		fmt.Fprintln(w, "No phone masks found.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tNUMBER\tENABLED\tLOCATION\tTEXTS LEFT\tMINS LEFT")
	for _, n := range numbers {
		fmt.Fprintf(tw, "%d\t%s\t%t\t%s\t%d\t%d\n",
			n.ID,
			n.Number,
			n.Enabled,
			n.Location,
			n.RemainingText,
			n.RemainingMin,
		)
	}
	return tw.Flush()
}

func printInboundContacts(w io.Writer, contacts []api.InboundContact) error {
	if len(contacts) == 0 {
		fmt.Fprintln(w, "No inbound contacts found.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tNUMBER\tBLOCKED\tCALLS\tTEXTS\tLAST CONTACT")
	for _, c := range contacts {
		lastDate := c.LastInboundDate
		fmt.Fprintf(tw, "%d\t%s\t%t\t%d\t%d\t%s\n",
			c.ID,
			c.InboundNumber,
			c.Blocked,
			c.NumCalls,
			c.NumTexts,
			lastDate,
		)
	}
	return tw.Flush()
}

func printRealPhones(w io.Writer, phones []api.RealPhone) error {
	if len(phones) == 0 {
		fmt.Fprintln(w, "No forwarding numbers found.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tNUMBER\tVERIFIED\tCOUNTRY")
	for _, p := range phones {
		fmt.Fprintf(tw, "%d\t%s\t%t\t%s\n",
			p.ID,
			p.Number,
			p.Verified,
			p.CountryCode,
		)
	}
	return tw.Flush()
}

func printRelayNumberSuggestions(w io.Writer, suggestions api.RelayNumberSuggestions) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	if suggestions.RealNum != nil {
		fmt.Fprintf(w, "Real Number: %s\n\n", *suggestions.RealNum)
	}

	printSuggestionCategory := func(title string, options []api.PhoneNumberOption) {
		if len(options) > 0 {
			fmt.Fprintf(tw, "%s:\n", title)
			fmt.Fprintln(tw, "  NUMBER\tLOCATION\tREGION")
			for _, opt := range options {
				locality := "-"
				if opt.Locality != nil {
					locality = *opt.Locality
				}
				fmt.Fprintf(tw, "  %s\t%s\t%s\n",
					opt.PhoneNumber,
					locality,
					opt.Region,
				)
			}
			fmt.Fprintln(tw)
		}
	}

	printSuggestionCategory("Same Prefix Options", suggestions.SamePrefixOptions)
	printSuggestionCategory("Same Area Options", suggestions.SameAreaOptions)
	printSuggestionCategory("Other Areas Options", suggestions.OtherAreasOptions)
	printSuggestionCategory("Random Options", suggestions.RandomOptions)

	return tw.Flush()
}

func printPhoneNumberOptions(w io.Writer, options []api.PhoneNumberOption) error {
	if len(options) == 0 {
		fmt.Fprintln(w, "No phone numbers found.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NUMBER\tLOCATION\tREGION\tCOUNTRY")
	for _, opt := range options {
		locality := "-"
		if opt.Locality != nil {
			locality = *opt.Locality
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			opt.PhoneNumber,
			locality,
			opt.Region,
			opt.ISOCountry,
		)
	}
	return tw.Flush()
}

func truncate(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
