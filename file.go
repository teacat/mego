package mego

import "os"

const (
	// B 為 Byte。
	B = 1
	// KB 為 Kilobyte。
	KB = 1024 * B
	// MB 為 Megabyte。
	MB = 1024 * KB
	// GB 為 Gigabyte。
	GB = 1024 * MB
	// TB 為 Terabyte。
	TB = 1024 * GB
)

// File 呈現了一個接收的檔案資料。
type File struct {
	// Name 為此檔案的原生名稱。
	Name string
	// Size 是這個檔案的總位元組大小。
	Size int
	// Extension 是這個檔案的副檔名。
	Extension string
	// Path 為此檔案上傳後的本地路徑。
	Path string
	// Keys 為此檔案的鍵值組，可供開發者存放自訂資料。
	Keys map[string]interface{}
}

// Remove 會移除這個檔案。
func (f *File) Remove() error {
	return os.Remove(f.Path)
}

// Move 會移動接收到的檔案到指定路徑。
func (f *File) Move(dest string) error {
	return os.Rename(f.Path, dest)
}
