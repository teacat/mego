package client

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
}
