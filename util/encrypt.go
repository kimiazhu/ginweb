// Description: util
// Author: ZHU HAIHUA
// Since: 2016-03-15 13:13
package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"math"
	"os"
)

// MD5File 会分段计算文件的MD5值，避免内存消耗太大。
// 如果不指定chunksize，它将使用默认值8KB。
func MD5File(filepath string, chunksize ...int64) string {
	chunk := float64(8192)
	if len(chunksize) > 1 {
		chunk = float64(chunksize[0])
	}
	file, err := os.Open(filepath)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	info, _ := file.Stat()
	size := info.Size()
	blocks := int64(math.Ceil(float64(size) / float64(chunk)))
	hash := md5.New()

	for i := int64(0); i < blocks; i++ {
		blocksize := int(math.Min(chunk, float64(size-int64(float64(i)*chunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash, string(buf))
	}
	return hex.EncodeToString(hash.Sum(nil))
}
