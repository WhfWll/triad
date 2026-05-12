package file

import (
	"encoding/csv"
	"errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os"
	"unicode/utf8"
)

var (
	ErrWriteToReader  = errors.New("write to a csv file with reading mode")
	ErrReadFromWriter = errors.New("read from a csv file with writing mode")
)

type File struct {
	osf *os.File
	w   *csv.Writer
	r   *csv.Reader
}

// CreateCSV create a new csv file for writing.
// A UTF-8 BOM header will be put at the very beginning, it will avoid encoding problems.
// If the file already exists you can append record to it.
func CreateCSV(path string) (*File, error) {
	create := false
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path does not exist
		create = true
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	if create {
		_, err = f.WriteString("\xEF\xBB\xBF") //写入UTF-8 BOM
		if err != nil {
			return nil, err
		}
	}

	return &File{f, csv.NewWriter(f), nil}, nil
}

// OpenCSV return a *File for reading,
// if the csv file does not exist an error will be returned.
func OpenCSV(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &File{f, nil, csv.NewReader(f)}, nil
}

// Close should be called when you finish csv reading or writing.
func (f *File) Close() error {
	if f.w != nil {
		err := f.Flush()
		if err != nil {
			return err
		}
	}
	return f.osf.Close()
}

// Write write one line record to the cache.
// the record will not save to disk until Flush be called.
func (f *File) Write(record []string) error {
	if f.w == nil {
		return ErrWriteToReader
	}
	return f.w.Write(record)
}

// WriteAll writes multiple records to csv file and then calls Flush,
// returning any error from the Flush.
func (f *File) WriteAll(records [][]string) error {
	if f.w == nil {
		return ErrWriteToReader
	}
	return f.w.WriteAll(records)
}

// Flush will sync the data from cache to disk file.
func (f *File) Flush() error {
	if f.w == nil {
		return ErrWriteToReader
	}
	f.w.Flush()
	return nil
}

// Read next record line.
func (f *File) Read() (record []string, err error) {
	if f.r == nil {
		return nil, ErrReadFromWriter
	}
	return f.r.Read()
}

// ReadAll returns all records.
func (f *File) ReadAll() (records [][]string, err error) {
	if f.r == nil {
		return nil, ErrReadFromWriter
	}
	return f.r.ReadAll()
}

// ReadCsv csv from file
func ReadCsv(filepath string) ([][]string, error) {
	f, err := OpenCSV(filepath)
	if err != nil {
		return nil, errors.New("OpenCSV fails with err:" + err.Error())
	}
	defer f.Close()
	rows, err := f.ReadAll()
	if err != nil {
		return rows, errors.New("ReadAll fails with err:" + err.Error())
	}
	result := DealCellData(rows)
	return result, nil
}

// ReadCsvFromFormFile csv from form file
func ReadCsvFromFormFile(file io.Reader) ([][]string, error) {
	r := csv.NewReader(file)
	if r == nil {
		return nil, ErrReadFromWriter
	}
	rows, err := r.ReadAll()
	if err != nil {
		return rows, errors.New("ReadAll fails with err:" + err.Error())
	}
	result := DealCellData(rows)
	return result, nil
}

// WriteCsv csv from [][]string
func WriteCsv(data [][]string, filepath string) error {
	f, err := CreateCSV(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	err = f.WriteAll(data)
	if err != nil {
		return err
	}
	return nil
}

// DealCellData GBK to utf8
func DealCellData(row [][]string) [][]string {
	for i := 0; i < len(row); i++ {
		var resRow []string
		for _, v := range row[i] {
			if !utf8.Valid([]byte(v)) {
				data, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(v))
				resRow = append(resRow, string(data))
			} else {
				resRow = append(resRow, v)
			}
		}
		row[i] = resRow
	}
	return row
}
