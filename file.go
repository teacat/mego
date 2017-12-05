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

// ChunkProcessor 是一個區塊處理介面，這令你可以自行撰寫區塊處理函式，
// 例如將每個區塊上傳至 Amazon S3 而不是在本地拼裝。
type ChunkProcessor interface {
	// Process 會在接收到每個區塊的時候被呼叫，並接收區塊建構體。
	// 當這個函式回傳了一個錯誤會中斷檔案的上傳。
	Process(*Chunk) error
	// Done 會在所有區塊都處理完畢時呼叫，會接收一個最後上傳的區塊建構體。
	// 需要回傳一個有資料的檔案結構體。
	Done(*Chunk) *File
}

// DefaultChunkProcessor 是預設的區塊處理建構體，用以實作區塊處理介面。
type DefaultChunkProcessor struct {
}

// Process 會處理每個區塊。
func (p *DefaultChunkProcessor) Process(chunk *Chunk) error {
	return nil
}

// Done 會在最終處理區塊的收尾。
func (p *DefaultChunkProcessor) Done(last *Chunk) *File {
	return nil
}
