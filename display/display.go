package display

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type Display struct {
	writer *tabwriter.Writer
	header RowContents
}

type RowContents string

func New() *Display {
	w := new(Display)
	w.writer = tabwriter.NewWriter(os.Stdout, 0, 1, 1, ' ', 0)
	return w
}

func (d *Display) CreateHeader(column ...string) {
	d.header = d.CreateRowContents(column...)
}

func (d *Display) CreateRowContents(column ...string) RowContents {
	t := strings.Join(column, "\t")
	return RowContents(t)
}

func (d *Display) viewHeader() error {
	_, err := fmt.Fprintln(d.writer, d.header)
	if err != nil {
		return fmt.Errorf("failed to write to stdout: %w", err)
	}

	return nil
}

func (d *Display) View(contents []RowContents) error {
	if d.header != "" {
		err := d.viewHeader()
		if err != nil {
			return err
		}
	}

	for _, row := range contents {
		_, err := fmt.Fprintln(d.writer, row)
		if err != nil {
			return fmt.Errorf("failed to write to stdout: %w", err)
		}
	}

	err := d.writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to write to stdout: %w", err)
	}

	return nil
}
