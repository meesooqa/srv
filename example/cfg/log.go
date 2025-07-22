package cfg

import "log/slog"

var logLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
}

// Log - log parameters
type Log struct {
	RawLevelCode    string `yaml:"level"`
	RawLevel        slog.Level
	RawOutputFormat string `yaml:"output_format"`
	RawWriteToFile  bool   `yaml:"write_to_file"`
	RawPath         string `yaml:"path"`
}

func (cfg *Log) Path() string {
	return cfg.RawPath
}

func (cfg *Log) OutputFormat() string {
	return cfg.RawOutputFormat
}

func (cfg *Log) Level() slog.Level {
	level, ok := logLevelMap[cfg.RawLevelCode]
	if ok {
		cfg.RawLevel = level
	} else {
		cfg.RawLevel = slog.LevelInfo
	}
	return cfg.RawLevel
}

func (cfg *Log) IsWriteToFile() bool {
	return cfg.RawWriteToFile
}
