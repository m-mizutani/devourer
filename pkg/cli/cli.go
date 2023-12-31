package cli

import (
	"github.com/m-mizutani/devourer/pkg/cli/config"
	"github.com/m-mizutani/devourer/pkg/domain/types"
	"github.com/m-mizutani/devourer/pkg/utils"
	"github.com/urfave/cli/v2"
)

func Run(args []string) error {
	var (
		logger config.Logger

		closer func()
	)

	app := cli.App{
		Name:    "devourer",
		Flags:   mergeFlags([]cli.Flag{}, &logger),
		Version: types.AppVersion,
		Commands: []*cli.Command{
			cmdCapture(),
		},
		Before: func(ctx *cli.Context) error {
			f, err := logger.Configure()
			if err != nil {
				return err
			}
			closer = f
			return nil
		},
		After: func(ctx *cli.Context) error {
			if closer != nil {
				closer()
			}
			return nil
		},
	}

	if err := app.Run(args); err != nil {
		utils.Logger().Error("Failed to run devourer", utils.ErrLog(err))
		return err
	}

	return nil
}
