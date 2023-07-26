package v1

import (
	"fmt"
	"strings"

	"log/slog"

	"github.com/ttacon/chalk"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *HealthStatus) Summary() string {
	if s.Status == nil || s.Health == nil {
		return "Unknown"
	}
	if !s.Status.Connected {
		return "Disconnected"
	}
	if len(s.Health.Conditions) > 0 {
		return fmt.Sprintf("Unhealthy: %s", strings.Join(s.Health.Conditions, ", "))
	}
	if !s.Health.Ready {
		return fmt.Sprintf("Not ready (unknown reason)")
	}
	return "Healthy"
}

type TimestampedLog interface {
	GetTimestamp() *timestamppb.Timestamp
	GetMsg() string
	GetLogLevel() slog.Level
}

var _ TimestampedLog = (*LogEntry)(nil)
var _ TimestampedLog = (*StateTransition)(nil)

func (s *LogEntry) GetLogLevel() slog.Level {
	return slog.Level(s.Level)
}

func (s *StateTransition) GetLogLevel() slog.Level {
	return slog.LevelInfo
}

func (s *StateTransition) GetMsg() string {
	return fmt.Sprintf("State changed to %s",
		stateColor(s.State).Color(s.State.String()))
}

func stateColor(state TaskState) chalk.Color {
	switch state {
	case TaskState_Pending:
		return chalk.Yellow
	case TaskState_Running:
		return chalk.Blue
	case TaskState_Failed, TaskState_Canceled:
		return chalk.Red
	case TaskState_Completed:
		return chalk.Green
	default:
		return chalk.White
	}
}
