package utils

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

var tableColors = []tablewriter.Colors{
	{tablewriter.Normal, tablewriter.FgGreenColor},
	{tablewriter.Normal, tablewriter.FgRedColor},
}

// Get a tablewriter color depending of a given ID
func GetTableRowColor(color string) tablewriter.Colors {
	if color == "green" {
		return tableColors[0]
	}
	return tableColors[1]
}

// Generate a writer for tablewriter
func GetTableWriter(header []string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader(header)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	return table
}
