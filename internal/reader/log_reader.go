package reader

import "C"
import (
	"bufio"
	"fmt"
	"github.com/sprakhar77/filereader/internal/model"
	"github.com/sprakhar77/filereader/internal/thread"
	"io"
	"os"
	"strings"
	"sync"
)

const KB = 1024
const MB = KB * KB
const GB = MB * MB

// LogReader interface defines the operations that our LogReader provides
type LogReader interface {
	Read(chunkSize int64, threadCount int) error
	MostActiveCookies(date string) []string
}

// logReader implements LogReader interface to concurrently read logs in multiple chunks and push the extracted cookies
// in the outputChan
type logReader struct {
	filePath   string
	wg         sync.WaitGroup
	outputChan chan model.Cookie
	cookieMap  *model.CookieMap
}

// NewLogReader returns a new instance of LogReader
func NewLogReader(filePath string) LogReader {
	return &logReader{
		filePath:   filePath,
		outputChan: make(chan model.Cookie),
		cookieMap:  model.NewCookieMap()}
}

// MostActiveCookies gives a list of most active cookies on the given date
func (r *logReader) MostActiveCookies(date string) []string {
	freqMap := r.cookieMap.Get(date)
	var result []string
	var maxFreq uint64 = 0
	for name, freq := range freqMap {
		if freq > maxFreq {
			maxFreq = freq
			result = nil
		}

		if freq == maxFreq {
			result = append(result, name)
		}
	}

	return result
}

// Read reads the log file concurrently in chunks of given size. The chunk size specifies the size of chunk (in bytes)
// that a thread should read. The thread count specifies the number of threads that should be used for concurrent reading
// operation. For the sake of this example we only use one thread and read the complete file line by line.
// However, this can be easily scaled by reducing the chunk size and increasing the thread count
func (r *logReader) Read(chunkSize int64, threadCount int) error {
	size, err := fileSizeInBytes(r.filePath)
	if err != nil {
		return fmt.Errorf("cannot get info about the file: %w", err)
	}

	var start int64
	var tasks thread.Tasks

	for start < size {
		chunkStart := start
		tasks = append(tasks, thread.NewTask(func() error { return r.readChunk(chunkStart, chunkSize) }))
		start += chunkSize
	}
	pool := thread.NewPool(tasks, threadCount)

	r.wg.Add(1)
	go r.processLogs()

	pool.Run()
	close(r.outputChan)

	r.wg.Wait()

	if errors := tasks.Errors(); len(errors) != 0 {
		return errors[0]
	}

	return nil
}

// readChunk reads a chunk of the file, starting at the offset indicated by start upto the size provided by size
func (r *logReader) readChunk(start, size int64) error {
	file, err := os.Open(r.filePath)
	defer file.Close()


	if err != nil {
		fmt.Errorf("cannot read chunk at start %d and size %d : %w", start, size, err)
	}

	_, err = file.Seek(start, 0)
	if err != nil {
		fmt.Errorf("cannot move file pointer to start %d : %w", start, err)
	}

	var readBytes int64
	reader := bufio.NewReader(file)
	if start != 0 {
		bytes, err := readLine(reader)
		if err != nil {
			return err
		}
		readBytes += int64(len(bytes))
	}

	for readBytes <= size {
		bytes, err := readLine(reader)

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return fmt.Errorf("could not read line: %w", err)
		}

		readBytes += int64(len(bytes))
		line := strings.TrimSpace(string(bytes))
		if len(line) > 0 {
			c, err := model.ToCookie(line)
			if err != nil {
				return fmt.Errorf("could not parse cookie info: %w", err)
			}

			r.outputChan <- c
		}
	}

	return nil
}

// readLine reads the line from the buffer
func readLine(reader *bufio.Reader) ([]byte, error) {
	bytes, err := reader.ReadBytes('\n')
	if err == io.EOF && len(bytes) > 0 {
		return bytes, nil
	}

	if err == io.EOF && len(bytes) == 0 {
		return nil, err
	}

	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to move reader to newline: %w", err)
	}

	return bytes, nil
}

// processLogs process the output that was sent in the outputChan by the threads reading the file.
// It adds the cookies to the cookie map
func (r *logReader) processLogs() {
	lineCount := 0
	for c := range r.outputChan {
		r.cookieMap.Add(c)
		lineCount++
	}

	r.wg.Done()
}

// fileSizeInBytes gets the size of the file without actually reading it
func fileSizeInBytes(filePath string) (int64, error) {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		return 0, err
	}

	fi, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}
