package client

import (
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
)

// File 呈現了一個欲上傳的檔案資料。
type File struct {
	// Binary 是檔案的二進制。
	Binary []byte `codec:"b" msgpack:"b"`
	// ID 是檔案編號，用於區塊組合。
	ID int `codec:"i" msgpack:"i"`
	// Parts 呈現區塊的分塊進度。索引 0 表示總共區塊數，索引 1 則是本區塊編號。
	// 如果這個切片是空的表示此為實體檔案而非區塊。
	Parts []int `codec:"p" msgpack:"p"`
	// Name 是檔案的原始名稱。
	Name string `codec:"n" msgpack:"n"`

	// source 是這個檔案的源頭，也許是 `string`、`[]byte`、`*os.File`
	source interface{}
	// length 是這個檔案的位元組長度，用以分切。
	length int
	// chunkSize 是切割此檔案的區塊長度。
	chunkSize int
}

// next 會裝載下一個區塊的內容。
func (f *File) next() error {
	// 如果這不是區塊檔案則離開。
	if len(f.Parts) == 0 {
		return nil
	}

	switch v := f.source.(type) {
	// *os.File 表示檔案寫入者。
	case *os.File:
		buf := make([]byte, f.chunkSize)
		_, err := v.ReadAt(buf, int64(f.Parts[0]*f.chunkSize))
		if err != nil {
			return err
		}
		f.Parts[0]++
		f.Binary = buf

	// 位元組表示檔案二進制內容。
	case []byte:
		f.Parts[0]++
		f.Binary = v[f.Parts[0]*f.chunkSize : f.chunkSize]

	// 字串型態表示檔案路徑。
	case string:
		fi, err := os.Open(v)
		if err != nil {
			return err
		}
		buf := make([]byte, f.chunkSize)
		_, err = fi.ReadAt(buf, int64(f.Parts[0]*f.chunkSize))
		if err != nil {
			return err
		}
		f.Parts[0]++
		f.Binary = buf
	}
	return nil
}

// load 會依照這個檔案的來源去讀取內容並裝載到此檔案建構體。
// 當檔案要以區塊方式上傳時，這個方法會盡可能地避免讀取整個檔案的內容。
func (f *File) load(isChunk bool) error {
	switch v := f.source.(type) {
	// *os.File 表示檔案寫入者。
	case *os.File:
		if isChunk {
			s, err := v.Stat()
			if err != nil {
				return err
			}
			f.length = int(s.Size())
			f.Parts = []int{0, int(math.Ceil(float64(f.length) / float64(f.chunkSize)))}
			f.Name = filepath.Base(v.Name())
		} else {
			b, err := ioutil.ReadAll(v)
			if err != nil {
				return err
			}
			f.Binary = b
			f.Name = filepath.Base(v.Name())
		}

	// 位元組表示檔案二進制內容。
	case []byte:
		if isChunk {
			f.length = len(v)
			f.Parts = []int{0, int(math.Ceil(float64(f.length) / float64(f.chunkSize)))}
		} else {
			f.Binary = v
		}

	// 字串型態表示檔案路徑。
	case string:
		if isChunk {
			fi, err := os.Open(v)
			if err != nil {
				return err
			}
			s, err := fi.Stat()
			if err != nil {
				return err
			}
			f.length = int(s.Size())
			f.Parts = []int{0, int(math.Ceil(float64(f.length) / float64(f.chunkSize)))}
			f.Name = filepath.Base(v)
		} else {
			b, err := ioutil.ReadFile(v)
			if err != nil {
				return err
			}
			f.Binary = b
			f.Name = filepath.Base(v)
		}
	}
	return nil
}
