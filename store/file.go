package store

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

type CursorStore interface {
	GetValue() (string, error)
	SetValue(v string) error
	Close() error
}

const (
	CursorLength = 200 // this is a rough estimate
)

type historyFileStore struct {
	file *os.File
}

func NewCursorHistoryFileStore(name string) (CursorStore, error) {

	file, err := os.OpenFile(name, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open cursor file: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get cursor file info: %w", err)
	}

	if stat.Size() > CursorLength*1000 {
		err := rollFile(file, CursorLength*100)
		if err != nil {
			return nil, fmt.Errorf("failed to roll file: %w", err)
		}
	}

	return &historyFileStore{file: file}, nil
}

// GetValue reads one character at a time from the end of a file, backtracking
// until it finds a newline character (note, newline characters are different
// depending on the operating system). If the last line of the file is the newline
// character, it is skipped, and looks for the next instance of the newline character.
func (f *historyFileStore) GetValue() (string, error) {
	stat, err := f.file.Stat()
	if err != nil {
		return "", err
	}

	var seekStart int64 = -1
	filesize := stat.Size()
	readChar := make([]byte, 1)
	lastLine := ""
	for {
		if filesize < -seekStart {
			break
		}

		_, err := f.file.Seek(seekStart, io.SeekEnd)
		if err != nil {
			return "", err
		}

		_, err = f.file.Read(readChar)
		if err != nil {
			return "", err
		}

		if (readChar[0] == '\n' || readChar[0] == '\r') && seekStart != -1 {
			break
		}
		lastLine = fmt.Sprintf("%s%s", string(readChar), lastLine)

		seekStart--
	}
	return lastLine, nil
}

func (f *historyFileStore) SetValue(v string) error {
	n, err := fmt.Fprintln(f.file, v)
	if err != nil {
		if n == 0 {
			err = fmt.Errorf("failed to save cursor: %s, with error: %s", v, err)
		} else if n != len(v) {
			err = fmt.Errorf("corrupted state, we saved %d bytes of original cursor: %s, with error: %s", n, v, err)
		} else {
			err = fmt.Errorf("something went wrong when saving cursor: %s, with error: %s", v, err)
		}
		return err
	}
	return nil
}

func (f *historyFileStore) Close() error {
	err := f.file.Close()
	if err != nil {
		return fmt.Errorf("failed to close cursor file. %w", err)
	}

	return nil
}

func rollFile(file *os.File, length int64) error {
	if file == nil {
		return fmt.Errorf("received a nil file")
	}

	saveBytes := make([]byte, length)
	var seekStart int64 = -length
	_, err := file.Seek(seekStart, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	_, err = file.Read(saveBytes)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	bytesReader := bytes.NewReader(saveBytes)
	bufReader := bufio.NewReader(bytesReader)
	bufReader.ReadLine() // trim first line

	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek to beginning of file: %w", err)
	}

	_, err = bufReader.WriteTo(file)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil

}
