package app

import (
	"fmt"

	"scan-eth/internal/scan/config"
	"scan-eth/internal/scan/services"
	"scan-eth/pkg/log"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

func NewScanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scan",
		Short:   "Scan ethereum blockchain",
		Version: "1.0.0",
		RunE:    ScanCommand,
	}

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "etc/scan/dev.yaml", "config file path")

	return cmd
}

func ScanCommand(cmd *cobra.Command, args []string) error {
	c, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("load config error: %s", err)
	}

	log.InitLogger(&c.Log)
	log.Infof("config: %+v", c)

	s, err := services.New(c)
	if err != nil {
		return err
	}
	defer s.Stop()

	return s.Run()
}
