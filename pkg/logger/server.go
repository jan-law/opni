package logger

import (
	"io"
	"log/slog"
	"regexp"
	"sync"

	controlv1 "github.com/rancher/opni/pkg/apis/control/v1"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/proto"
)

type LogServer struct {
	controlv1.UnsafeLogServer
	LogServerOptions

	clientsMu sync.RWMutex
	clients   map[string]controlv1.LogClient
}

type LogServerOptions struct{}

type LogServerOption func(*LogServerOptions)

func (o *LogServerOptions) apply(opts ...LogServerOption) {
	for _, op := range opts {
		op(o)
	}
}

func NewLogServer(opts ...LogServerOption) *LogServer {
	options := &LogServerOptions{}
	options.apply(opts...)

	return &LogServer{
		clients: make(map[string]controlv1.LogClient),
	}
}

func (ls *LogServer) AddClient(name string, client controlv1.LogClient) {
	ls.clientsMu.Lock()
	defer ls.clientsMu.Unlock()
	ls.clients[name] = client
}

func (ls *LogServer) RemoveClient(name string) {
	ls.clientsMu.Lock()
	defer ls.clientsMu.Unlock()
	delete(ls.clients, name)
}

// todo point protoHandler to write to LogServer, which will write to all clients
// who are following logs
func (ls *LogServer) Write(p []byte) (int, error) {
	len := len(p)
	ls.clientsMu.Lock()
	defer ls.clientsMu.Unlock()
	// for i, c := range w.clients {
	// }

	return len, nil
}

func (ls *LogServer) StreamLogs(req *controlv1.LogStreamRequest, stream controlv1.Log_StreamLogsServer) error {
	f, err := OpenLogFile()
	if err != nil {
		return err
	}

	defer f.Close()

	for {
		msg, err := getLogMessage(req, f)
		if err != nil && err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		if msg == nil {
			continue
		}

		if err := stream.Send(msg); err != nil {
			return err
		}
	}
}
func getLogMessage(req *controlv1.LogStreamRequest, f afero.File) (*controlv1.StructuredLogRecord, error) {
	// todo optimize filtering logic
	since := req.Since.AsTime()
	until := req.Until.AsTime()
	minLevel := req.Filters.Level
	nameFilters := req.Filters.NamePattern

	sizeBuf := make([]byte, 4)
	record := &controlv1.StructuredLogRecord{}
	_, err := io.ReadFull(f, sizeBuf)
	if err != nil {
		return nil, err
	}

	size := int32(sizeBuf[0]) |
		int32(sizeBuf[1])<<8 |
		int32(sizeBuf[2])<<16 |
		int32(sizeBuf[3])<<24

	recordBytes := make([]byte, size)
	_, err = io.ReadFull(f, recordBytes)
	if err != nil {
		return nil, err
	}

	if err := proto.Unmarshal(recordBytes, record); err != nil {
		return nil, err
	}

	if ParseLevel(record.Level) < slog.Level(minLevel) {
		return nil, nil
	}

	time := record.Time.AsTime()
	if time.Before(since) || time.After(until) {
		return nil, nil
	}

	if nameFilters != nil && !matchesNameFilter(nameFilters, record.Name) {
		return nil, nil
	}

	return record, nil
}

func matchesNameFilter(patterns []string, name string) bool {
	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, name)
		if matched {
			return true
		}
	}
	return false
}
