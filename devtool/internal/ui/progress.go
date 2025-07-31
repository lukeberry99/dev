package ui

import (
	"fmt"
	"time"

	"github.com/schollz/progressbar/v3"
)

type ProgressReporter struct {
	bar *progressbar.ProgressBar
}

func NewProgressReporter(max int, description string) *ProgressReporter {
	bar := progressbar.NewOptions(max,
		progressbar.OptionSetDescription(description),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
	)

	return &ProgressReporter{
		bar: bar,
	}
}

func (p *ProgressReporter) Increment() {
	p.bar.Add(1)
}

func (p *ProgressReporter) SetDescription(desc string) {
	p.bar.Describe(desc)
}

func (p *ProgressReporter) Finish() {
	p.bar.Finish()
}
