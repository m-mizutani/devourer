package cli

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/m-mizutani/devourer/pkg/cli/config"
	"github.com/m-mizutani/devourer/pkg/domain/interfaces"
	"github.com/m-mizutani/devourer/pkg/domain/logic"
	"github.com/m-mizutani/devourer/pkg/infra"
	"github.com/m-mizutani/devourer/pkg/infra/capture"
	"github.com/m-mizutani/devourer/pkg/utils"
	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v2"
)

func cmdCapture() *cli.Command {
	var (
		iface        string
		output       string
		writeFile    string
		bigquery     config.BigQuery
		statInterval time.Duration
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
			&cli.DurationFlag{
				Name:        "stat-interval",
				Category:    "capture",
				Aliases:     []string{"s"},
				EnvVars:     []string{"DEVOURER_STAT_INTERVAL"},
				Usage:       "Show statistics in every interval",
				Destination: &statInterval,
			},
		}, &bigquery),
		Action: func(c *cli.Context) error {
			// configure dumper
			var dumper interfaces.Dumper
			switch output {
			case "file":
				fd, err := os.Create(filepath.Clean(writeFile))
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
			expireTicker := time.NewTicker(1 * time.Second)
			var statTicker <-chan time.Time
			if statInterval > 0 {
				statTicker = time.NewTicker(statInterval).C
			} else {
				statTicker = make(chan time.Time)
			}

			ctx := c.Context
			engine := logic.NewEngine()

			utils.Logger().Info("Starting capture...",
				slog.Any("interface", iface),
				slog.Any("output", output),
				slog.Any("stat_interval", statInterval.String()),
			)

			var packetCount int64
			var sizeCount int64
			lastTime := time.Now()

			for {
				select {
				case <-ctx.Done():
					return nil

				case pkt := <-clients.Capture().Read():
					packetCount++
					sizeCount += int64(pkt.Metadata().Length)

					out, err := engine.InputPacket(pkt)
					if err != nil {
						return err
					}
					if out != nil {
						if err := clients.Dumper().Dump(ctx, out); err != nil {
							return err
						}
					}

				case <-statTicker:
					d := time.Since(lastTime)
					utils.Logger().Info("Statistics",
						slog.String("pps", fmt.Sprintf("%.2f", float64(packetCount)/d.Seconds())),
						slog.String("bps", convertToBps(float64(sizeCount)/d.Seconds())),
						slog.Int("flow count", engine.FlowCount()),
					)
					packetCount = 0
					sizeCount = 0
					lastTime = time.Now()

				case <-expireTicker.C:
					out, err := engine.Tick(time.Now())
					if err != nil {
						return err
					}
					if err := clients.Dumper().Dump(ctx, out); err != nil {
						return err
					}

				case <-sigCh:
					out := engine.Flush()
					utils.Logger().Info("Caught signal, flushing flow logs...",
						slog.Int("flow_logs", len(out.FlowLogs)),
					)
					if err := clients.Dumper().Dump(ctx, out); err != nil {
						return err
					}

					return nil
				}
			}
		},
	}
}
