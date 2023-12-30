package bq

import (
	"context"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/m-mizutani/devourer/pkg/domain/interfaces"
	"github.com/m-mizutani/devourer/pkg/domain/model"
	"github.com/m-mizutani/devourer/pkg/utils"
	"github.com/m-mizutani/goerr"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type implClient struct {
	projectID string
	datasetID string
	client    *bigquery.Client
	dataSet   *bigquery.Dataset
}

type flowLog struct {
	ID          string    `bigquery:"id"`
	Protocol    string    `bigquery:"protocol"`
	SrcAddr     string    `bigquery:"src_addr"`
	DstAddr     string    `bigquery:"dst_addr"`
	SrcPort     int       `bigquery:"src_port"`
	DstPort     int       `bigquery:"dst_port"`
	FirstSeenAt time.Time `bigquery:"first_seen_at"`
	LastSeenAt  time.Time `bigquery:"last_seen_at"`
	SrcBytes    int64     `bigquery:"src_bytes"`
	DstBytes    int64     `bigquery:"dst_bytes"`
	SrcPackets  int64     `bigquery:"src_packets"`
	DstPackets  int64     `bigquery:"dst_packets"`
	Status      string    `bigquery:"status"`
}

const (
	tblFlowLogs = "flow_logs"
)

func New(ctx context.Context, projectID, datasetID string, opts ...option.ClientOption) (interfaces.Dumper, error) {
	bqClient, err := bigquery.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, err
	}

	dataSet := bqClient.Dataset(datasetID)

	tables := []struct {
		name   string
		schema any
	}{
		{
			name:   tblFlowLogs,
			schema: flowLog{},
		},
	}

	for _, t := range tables {
		table := dataSet.Table(t.name)
		schema, err := bigquery.InferSchema(t.schema)
		if err != nil {
			return nil, goerr.Wrap(err, "failed to infer schema").With("table", t.name)
		}

		meta := &bigquery.TableMetadata{
			Schema: schema,
			TimePartitioning: &bigquery.TimePartitioning{
				Type:  bigquery.DayPartitioningType,
				Field: "first_seen_at",
			},
		}
		if err := table.Create(ctx, meta); err != nil {
			if gerr, ok := err.(*googleapi.Error); !ok || gerr.Code != 409 {
				return nil, goerr.Wrap(err, "failed to create table").With("table", t.name)
			}
		}
	}

	return &implClient{
		projectID: projectID,
		datasetID: datasetID,
		client:    bqClient,
		dataSet:   dataSet,
	}, nil
}

func (x *implClient) Dump(ctx context.Context, record *model.Record) error {
	if len(record.FlowLogs) > 0 {
		rows := make([]flowLog, len(record.FlowLogs))
		for i, flow := range record.FlowLogs {
			rows[i] = flowLog{
				ID:          flow.ID.String(),
				Protocol:    flow.Protocol,
				SrcAddr:     flow.Src.Addr.String(),
				DstAddr:     flow.Dst.Addr.String(),
				SrcPort:     int(flow.Src.Port),
				DstPort:     int(flow.Dst.Port),
				FirstSeenAt: flow.FirstSeenAt,
				LastSeenAt:  flow.LastSeenAt,
				SrcBytes:    int64(flow.SrcStat.Bytes),
				DstBytes:    int64(flow.DstStat.Bytes),
				SrcPackets:  int64(flow.SrcStat.Packets),
				DstPackets:  int64(flow.DstStat.Packets),
				Status:      flow.Status,
			}
		}

		insert := x.dataSet.Table(tblFlowLogs).Inserter()
		if err := insert.Put(ctx, rows); err != nil {
			return goerr.Wrap(err, "failed to insert row of new flows")
		}
	}

	return nil
}

func (x *implClient) Close() {
	utils.SafeClose(x.client)
}
