package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

/* =========================
   MODELS
========================= */

type FileType string

const (
	CSV FileType = "CSV"
)

type JobStatus string

const (
	PENDING    JobStatus = "PENDING"
	PROCESSING JobStatus = "PROCESSING"
	COMPLETED  JobStatus = "COMPLETED"
	FAILED     JobStatus = "FAILED"
)

type Record []string

type Chunk struct {
	ID      int
	Records []Record
}

type ProcessingResult struct {
	Records []Record
}

type ProcessingJob struct {
	ID         string
	InputPath  string
	OutputPath string
	FileType   FileType
	Status     JobStatus
	Error      string
}

/* =========================
   STRATEGY INTERFACES
========================= */

type ChunkReader interface {
	ReadChunk() (*Chunk, error)
	Close() error
}

type ChunkProcessor interface {
	Process(*Chunk) (*ProcessingResult, error)
}

type ChunkWriter interface {
	Write(*ProcessingResult) error
	Close() error
}

/* =========================
   CSV READER
========================= */

type CSVReader struct {
	file      *os.File
	reader    *csv.Reader
	chunkSize int
	chunkID   int
}

func NewCSVReader(path string, chunkSize int) (*CSVReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &CSVReader{
		file:      file,
		reader:    csv.NewReader(file),
		chunkSize: chunkSize,
	}, nil
}

func (r *CSVReader) ReadChunk() (*Chunk, error) {

	records := make([]Record, 0, r.chunkSize)

	for len(records) < r.chunkSize {

		record, err := r.reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		records = append(records, Record(record))
	}

	if len(records) == 0 {
		return nil, io.EOF
	}

	chunk := &Chunk{
		ID:      r.chunkID,
		Records: records,
	}

	r.chunkID++

	return chunk, nil
}

func (r *CSVReader) Close() error {
	return r.file.Close()
}

/* =========================
   CSV PROCESSOR
========================= */

type CSVProcessor struct{}

func NewCSVProcessor() *CSVProcessor {
	return &CSVProcessor{}
}

func (p *CSVProcessor) Process(
	chunk *Chunk,
) (*ProcessingResult, error) {

	result := &ProcessingResult{
		Records: make([]Record, 0),
	}

	for _, record := range chunk.Records {

		// Validation
		if len(record) == 0 {
			return nil, errors.New("invalid record")
		}

		// Transformation
		for i := range record {
			record[i] = strings.TrimSpace(record[i])
		}

		// Filtering
		if len(record) == 1 && record[0] == "" {
			continue
		}

		result.Records = append(
			result.Records,
			record,
		)
	}

	return result, nil
}

/* =========================
   CSV WRITER
========================= */

type CSVWriter struct {
	file   *os.File
	writer *csv.Writer
}

func NewCSVWriter(path string) (*CSVWriter, error) {

	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &CSVWriter{
		file:   file,
		writer: csv.NewWriter(file),
	}, nil
}

func (w *CSVWriter) Write(
	result *ProcessingResult,
) error {

	for _, record := range result.Records {
		w.writer.Write([]string(record))
	}

	w.writer.Flush()

	return w.writer.Error()
}

func (w *CSVWriter) Close() error {
	w.writer.Flush()

	if err := w.writer.Error(); err != nil {
		w.file.Close()
		return err
	}

	return w.file.Close()
}

/* =========================
   FACTORY
========================= */

// Factory creates the appropriate
// Reader + Processor + Writer.
type ComponentFactory struct{}

type Components struct {
	reader    ChunkReader
	processor ChunkProcessor
	writer    ChunkWriter
}

func NewComponentFactory(
	fileType FileType,
	inputPath string,
	outputPath string,
	chunkSize int,
) (*Components, error) {

	switch fileType {

	case CSV:

		reader, err := NewCSVReader(
			inputPath,
			chunkSize,
		)
		if err != nil {
			return nil, err
		}

		writer, err := NewCSVWriter(outputPath)
		if err != nil {
			reader.Close()
			return nil, err
		}

		return &Components{
			reader:    reader,
			processor: NewCSVProcessor(),
			writer:    writer,
		}, nil

	default:
		return nil, errors.New("unsupported file type")
	}
}

/* =========================
   SERVICE
========================= */

type FileProcessingService struct {
	workerCount int
	bufferSize  int

	mu   sync.RWMutex
	jobs map[string]*ProcessingJob
}

func NewFileProcessingService(
	workerCount int,
	bufferSize int,
) *FileProcessingService {

	return &FileProcessingService{
		workerCount: workerCount,
		bufferSize:  bufferSize,
		jobs:        make(map[string]*ProcessingJob),
	}
}

/* =========================
   CREATE JOB
========================= */

func (s *FileProcessingService) CreateJob(
	inputPath string,
	outputPath string,
	fileType FileType,
) (*ProcessingJob, error) {

	job := &ProcessingJob{
		ID:         fmt.Sprintf("job-%d", time.Now().UnixNano()),
		InputPath:  inputPath,
		OutputPath: outputPath,
		FileType:   fileType,
		Status:     PENDING,
	}

	s.mu.Lock()
	s.jobs[job.ID] = job
	s.mu.Unlock()

	go s.processJob(job)

	return job, nil
}

/*
	=========================
	  WORKER POOL

=========================
*/
func (s *FileProcessingService) processJob(job *ProcessingJob) {

	s.updateStatus(job.ID, PROCESSING, "")

	// Create reader, processor and writer
	components, err := NewComponentFactory(
		job.FileType,
		job.InputPath,
		job.OutputPath,
		1000,
	)
	if err != nil {
		s.failJob(job.ID, err)
		return
	}

	defer components.reader.Close()
	defer components.writer.Close()

	jobs := make(chan *Chunk, 10)

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < s.workerCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for chunk := range jobs {

				result, err := processWithRetry(
					components.processor,
					chunk,
					3,
				)

				if err != nil {
					// Handle failed chunk
					continue
				}

				// Write immediately
				if err := components.writer.Write(result); err != nil {
					return
				}
			}
		}()
	}

	// Read and distribute chunks
	for {
		chunk, err := components.reader.ReadChunk()

		if err == io.EOF {
			break
		}

		if err != nil {
			s.failJob(job.ID, err)
			break
		}

		jobs <- chunk
	}

	// No more chunks
	close(jobs)

	// Wait for all workers
	wg.Wait()

	s.updateStatus(job.ID, COMPLETED, "")
}

/* =========================
   RETRY
========================= */

func processWithRetry(
	processor ChunkProcessor,
	chunk *Chunk,
	maxRetries int,
) (*ProcessingResult, error) {

	var err error

	for attempt := 0; attempt <= maxRetries; attempt++ {

		result, processErr := processor.Process(chunk)

		if processErr == nil {
			return result, nil
		}

		err = processErr

		if attempt < maxRetries {
			time.Sleep(
				time.Duration(1<<attempt) *
					100 *
					time.Millisecond,
			)
		}
	}

	return nil, err
}

/* =========================
   STATUS
========================= */

func (s *FileProcessingService) GetJobStatus(
	jobID string,
) (*ProcessingJob, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	job, ok := s.jobs[jobID]

	if !ok {
		return nil, errors.New("job not found")
	}

	return job, nil
}

func (s *FileProcessingService) updateStatus(
	jobID string,
	status JobStatus,
	errMsg string,
) {

	s.mu.Lock()
	defer s.mu.Unlock()

	job := s.jobs[jobID]

	job.Status = status
	job.Error = errMsg
}

func (s *FileProcessingService) failJob(
	jobID string,
	err error,
) {

	s.updateStatus(
		jobID,
		FAILED,
		err.Error(),
	)
}

/* =========================
   MAIN
========================= */

func main() {

	/*
		4 workers
		10 chunks max waiting in memory
	*/

	service := NewFileProcessingService(
		4,
		10,
	)

	job, err := service.CreateJob(
		"input.csv",
		"output.csv",
		CSV,
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Job:", job.ID)

	time.Sleep(2 * time.Second)

	status, _ := service.GetJobStatus(job.ID)

	fmt.Println("Status:", status.Status)
}
