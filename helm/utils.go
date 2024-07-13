package main

import (
	"context"
	"dagger/helm/internal/dagger"

	"helm.sh/helm/v3/pkg/chart"
	"sigs.k8s.io/yaml"
)

func getChartMetadata(ctx context.Context, c *dagger.Directory) (*chart.Metadata, error) {
	chartYaml, err := c.File("Chart.yaml").Contents(ctx)
	if err != nil {
		return nil, err
	}

	meta := new(chart.Metadata)
	err = yaml.Unmarshal([]byte(chartYaml), meta)
	if err != nil {
		return nil, err
	}

	return meta, nil
}
