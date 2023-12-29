package cli

import (
	"context"
	"encoding/json"
	"io"

	"github.com/m-mizutani/devourer/pkg/domain/model"
	"github.com/m-mizutani/devourer/pkg/utils"
	"github.com/urfave/cli/v2"
)

type flagConfig interface {
	Flags() []cli.Flag
}

func mergeFlags(base []cli.Flag, configs ...flagConfig) []cli.Flag {
	ret := base[:]
	for _, config := range configs {
		ret = append(ret, config.Flags()...)
	}

	return ret
}

type jsonDumper struct {
	encoder *json.Encoder
	closer  io.Closer
}

func (s *jsonDumper) Dump(ctx context.Context, record *model.Record) error {
	for _, flow := range record.NewFlows {
		if err := s.encoder.Encode(flow); err != nil {
			return err
		}
	}
	for _, flow := range record.ClosedFlows {
		if err := s.encoder.Encode(flow); err != nil {
			return err
		}
	}
	return nil
}

func (s *jsonDumper) Close() {
	if s.closer != nil {
		utils.SafeClose(s.closer)
	}
}
