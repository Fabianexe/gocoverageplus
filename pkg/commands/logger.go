package commands

import (
	"log/slog"
	"os"

	"github.com/spf13/pflag"
)

type multiFlag int

var _ pflag.Value = (*multiFlag)(nil)

func (m *multiFlag) String() string {
	return "verbose"
}

func (m *multiFlag) Set(_ string) error {
	*m++

	return nil
}

func (m *multiFlag) Type() string {
	return "bool"
}

var verbose multiFlag

func initLogger() {
	var programLevel = new(slog.LevelVar)
	switch verbose {
	case 0:
		programLevel.Set(slog.LevelError)
	case 1:
		programLevel.Set(slog.LevelInfo)
	default:
		programLevel.Set(slog.LevelDebug)
	}
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: programLevel,
	})
	slog.SetDefault(slog.New(h))
}
