package report

import (
	"encoding/csv"
	"fmt"
	"io"
	"text/tabwriter"
)

type BitResult struct {
	Bit                      int
	Class                    string
	Accuracy, Drops, Invalid float64
}

func BitTable(out io.Writer, rows []BitResult) {
	w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "Bit\tClass\tMean accuracy\tDrop >10 points\tNaN/Inf")
	for _, row := range rows {
		fmt.Fprintf(w, "%d\t%s\t%.2f%%\t%.2f%%\t%.2f%%\n", row.Bit, row.Class, row.Accuracy*100, row.Drops*100, row.Invalid*100)
	}
	w.Flush()
}

func BitCSV(out io.Writer, rows []BitResult) error {
	w := csv.NewWriter(out)
	if err := w.Write([]string{"bit_index", "bit_class", "mean_accuracy_percent", "drop_over_10_points_percent", "nan_inf_percent"}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write([]string{fmt.Sprint(row.Bit), row.Class, fmt.Sprintf("%.6f", row.Accuracy*100), fmt.Sprintf("%.6f", row.Drops*100), fmt.Sprintf("%.6f", row.Invalid*100)}); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}
