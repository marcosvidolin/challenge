package producer

import (
	"encoding/csv"
	"io"
	"os"
)

// csvReader is a wrapper for reading a CSV file in chunks,
// tracking the file being read and the current chunk of
// records being processed.
type csvReader struct {
	file  *os.File
	chunk *Chunk
}

// NewCsvReader initializes a new csvReader for the given filename.
// Returns an error if the file does not exist or cannot be opened.
func NewCsvReader(filename string) (*csvReader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return &csvReader{
		file: f,
		chunk: &Chunk{
			Filename: filename,
		},
	}, nil
}

// Read reads a chunk of records at time
// based on the give chunk size
func (r *csvReader) Read(chunkSize int, chunkCh chan<- Chunk) error {
	reader := csv.NewReader(r.file)

	chunk := Chunk{
		Filename: r.chunk.Filename,
		Records:  []Record{},
	}

	for line := 0; ; line++ {
		record, err := reader.Read()
		if err == io.EOF {
			if len(chunk.Records) > 0 {
				chunkCh <- chunk
			}
			close(chunkCh)
			return io.EOF
		}
		if err != nil {
			// TODO: just log the issue
			line++
			continue
		}

		if line == 0 {
			continue
		}

		chunk.Records = append(chunk.Records, record)

		if len(chunk.Records) == chunkSize {
			chunkCh <- chunk
			chunk = Chunk{
				Filename: r.chunk.Filename,
				Records:  []Record{},
			}
		}

		line++
	}
}

// Close closes the file. Close will return error
// if its already been closed
func (r *csvReader) Close() error {
	return r.file.Close()
}

// Chunk represents a portion of the file, containing
// successfully processed records and any failures
// encountered during processing
type Chunk struct {
	Filename string
	Records  []Record
}

// Record represents a single CSV record as a slice of string values
type Record []string
