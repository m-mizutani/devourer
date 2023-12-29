package cli

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/m-mizutani/devourer/pkg/cli/config"
	"github.com/m-mizutani/devourer/pkg/domain/interfaces"
	"github.com/m-mizutani/devourer/pkg/domain/logic"
	"github.com/m-mizutani/devourer/pkg/infra"
	"github.com/m-mizutani/devourer/pkg/infra/capture"
	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v2"
)

func cmdCapture() *cli.Command {
	var (
		iface     string
		output    string
		writeFile string
		bigquery  config.BigQuery
	)

	return &cli.Command{
		Name:    "capture",
		Usage:   "Capture packets from network interface",
		Aliases: []string{"c"},
		Flags: mergeFlags([]cli.Flag{
			&cli.StringFlag{
				Name:        "interface",
				Category:    "capture",
				Aliases:     []string{"i"},
				EnvVars:     []string{"DEVOURER_INTERFACE"},
				Usage:       "Network interface to capture packets",
				Destination: &iface,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "output",
				Category:    "capture",
				Aliases:     []string{"o"},
				EnvVars:     []string{"DEVOURER_OUTPUT"},
				Usage:       "Output destination (stdout, file, bigquery)",
				Destination: &output,
				Value:       "stdout",
			},
			&cli.StringFlag{
				Name:        "write-file",
				Category:    "capture",
				Aliases:     []string{"w"},
				EnvVars:     []string{"DEVOURER_WRITE_FILE"},
				Usage:       "Write packets to file. This option works only with output=file",
				Destination: &writeFile,
			},
		}, &bigquery),
		Action: func(c *cli.Context) error {
			// configure dumper
			var dumper interfaces.Dumper
			switch output {
			case "file":
				fd, err := os.Create(writeFile)
				if err != nil {
					return goerr.Wrap(err, "Failed to create file")
				}
				dumper = &jsonDumper{
					encoder: json.NewEncoder(os.Stdout),
					closer:  fd,
				}

			case "stdout":
				dumper = &jsonDumper{encoder: json.NewEncoder(os.Stdout)}

			case "bigquery":
				v, err := bigquery.Configure(c.Context)
				if err != nil {
					return err
				}
				dumper = v
			}

			// configure device
			device, err := capture.NewDevice(iface)
			if err != nil {
				return err
			}

			clients := infra.New(
				infra.WithDumper(dumper),
				infra.WithCapture(device),
			)

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
			ticker := time.NewTicker(1 * time.Second)

			ctx := c.Context
			engine := logic.NewEngine()

			for {
				select {
				case <-ctx.Done():
					return nil

				case pkt := <-clients.Capture().Read():
					out, err := engine.InputPacket(pkt)
					if err != nil {
						return err
					}
					if out != nil {
						if err := clients.Dumper().Dump(ctx, out); err != nil {
							return err
						}
					}

				case <-ticker.C:
					out, err := engine.Tick(time.Now())
					if err != nil {
						return err
					}
					if out != nil {
						if err := clients.Dumper().Dump(ctx, out); err != nil {
							return err
						}
					}

				case <-sigCh:
					out := engine.Flush()
					if out != nil {
						if err := clients.Dumper().Dump(ctx, out); err != nil {
							return err
						}
					}
					return nil
				}
			}
		},
	}
}
