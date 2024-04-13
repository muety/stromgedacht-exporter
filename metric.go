package main

import (
	"fmt"
	"strings"
)

// Hand-crafted Prometheus metrics
// Since we're only using very simple counters in this application,
// we don't actually need the official client SDK as a dependency
// Inspired by https://github.com/muety/wakapi/tree/master/models/metrics

type Metrics []Metric

func (m Metrics) Print() (output string) {
	printedMetrics := make(map[string]bool)
	for _, m := range m {
		if _, ok := printedMetrics[m.Key()]; !ok {
			output += fmt.Sprintf("%s\n", m.Header())
			printedMetrics[m.Key()] = true
		}
		output += fmt.Sprintf("%s\n", m.Print())
	}

	return output
}

func (m Metrics) Len() int {
	return len(m)
}

func (m Metrics) Less(i, j int) bool {
	return strings.Compare(m[i].Key(), m[j].Key()) < 0
}

func (m Metrics) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type Metric interface {
	Key() string
	Header() string
	Print() string
}

type GaugeMetric struct {
	Name   string
	Value  int64
	Desc   string
	Labels Labels
}

func (c GaugeMetric) Key() string {
	return c.Name
}

func (c GaugeMetric) Print() string {
	return fmt.Sprintf("%s%s %d", c.Name, c.Labels.Print(), c.Value)
}

func (c GaugeMetric) Header() string {
	return fmt.Sprintf("# HELP %s %s\n# TYPE %s gauge", c.Name, c.Desc, c.Name)
}

type Labels []Label

type Label struct {
	Key   string
	Value string
}

func (l Labels) Print() string {
	printedLabels := make([]string, len(l))
	for i, e := range l {
		printedLabels[i] = e.Print()
	}
	if len(l) == 0 {
		return ""
	}
	return fmt.Sprintf("{%s}", strings.Join(printedLabels, ","))
}

func (l Label) Print() string {
	return fmt.Sprintf("%s=\"%s\"", l.Key, l.Value)
}
