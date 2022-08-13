package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type Row struct {
	id       int
	name     string
	dueDate  string
	class    string
	status   string
	priority string
}

func NewShow() *cobra.Command {
	showCmd := &cobra.Command{
		Use:          "show",
		Short:        "Show tasks currently in planner",
		Run:          showAll,
		SilenceUsage: true,
	}

	return showCmd
}

func showAll(cmd *cobra.Command, args []string) {
	ctx, log := CmdContext(cmd)

	cfgFile := cmd.Flag("config").Value.String()
	db, table := DB(ctx, cfgFile)

	// Don't proceed if database is nil.
	if db == nil {
		return
	}

	if _, err := db.Exec(`USE planner;`); err != nil {
		log.Error(err, "Setting table failed.")

		return
	}

	rows, err := db.Query(fmt.Sprintf("SELECT id, name, due_date, class_name, status, priority FROM %s;", table))
	if err != nil {
		log.Error(err, "Query failed.")

		return
	}

	defer rows.Close()

	rs := [][]any{}

	for rows.Next() {
		r := new(Row)

		if err := rows.Scan(&r.id, &r.name, &r.dueDate, &r.class, &r.status, &r.priority); err != nil {
			log.Error(err, "Could not scan row.")

			return
		}

		rs = append(rs, r.List())
	}

	if err := rows.Err(); err != nil {
		log.Error(err, "Error reading rows.")

		return
	}

	header := []any{}
	header = append(header, "ID")
	header = append(header, "Name")
	header = append(header, "Due Date")
	header = append(header, "Class")
	header = append(header, "Status")
	header = append(header, "Priority")

	printTable(header, rs)
}

func printTable(header []any, rows [][]any) {
	tw := table.NewWriter()

	// Print header.
	tw.AppendHeader(header)

	// Print rows.
	for _, r := range rows {
		tw.AppendRow(r)
	}

	fmt.Println(tw.Render())
}

func (r *Row) List() []any {
	return []any{
		r.id,
		r.name,
		r.dueDate,
		r.class,
		r.status,
		r.priority,
	}
}
