package config

import (
	"context"

	"github.com/m-mizutani/devourer/pkg/domain/interfaces"
	"github.com/m-mizutani/devourer/pkg/domain/types"
	"github.com/m-mizutani/devourer/pkg/infra/bq"
	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v2"
	"google.golang.org/api/option"
)

type BigQuery struct {
	projectID string
	datasetID string
	saKeyData string `masq:"secret"`
	saKeyFile string
}

func (x *BigQuery) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "bq-project-id",
			Usage:       "BigQuery project ID",
			Destination: &x.projectID,
			EnvVars:     []string{"DEVOURER_BIGQUERY_PROJECT_ID"},
		},
		&cli.StringFlag{
			Name:        "bq-dataset-id",
			Usage:       "BigQuery dataset ID",
			Destination: &x.datasetID,
			EnvVars:     []string{"DEVOURER_BIGQUERY_DATASET_ID"},
		},
		&cli.StringFlag{
			Name:        "bq-sa-key-data",
			Usage:       "BigQuery service account key data",
			Destination: &x.saKeyData,
			EnvVars:     []string{"DEVOURER_BIGQUERY_SA_KEY_DATA"},
		},
		&cli.StringFlag{
			Name:        "bq-sa-key-file",
			Usage:       "BigQuery service account key file",
			Destination: &x.saKeyFile,
			EnvVars:     []string{"DEVOURER_BIGQUERY_SA_KEY_FILE"},
		},
	}
}

func (x *BigQuery) Configure(ctx context.Context) (interfaces.Dumper, error) {
	if x.projectID == "" {
		return nil, goerr.Wrap(types.ErrInvalidOption, "BigQuery project ID is empty")
	}
	if x.datasetID == "" {
		return nil, goerr.Wrap(types.ErrInvalidOption, "BigQuery dataset ID is empty")
	}

	var options []option.ClientOption
	if x.saKeyData != "" {
		options = append(options, option.WithCredentialsJSON([]byte(x.saKeyData)))
	}
	if x.saKeyFile != "" {
		options = append(options, option.WithCredentialsFile(x.saKeyFile))
	}

	return bq.New(ctx,
		x.projectID,
		x.datasetID,
		options...,
	)
}
