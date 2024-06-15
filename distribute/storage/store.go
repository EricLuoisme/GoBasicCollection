package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

// CASPathTransformFunc 将Key进行hash, 再对hashString进行分段
// 作用可以看作是相较于直接使用Key作为存储路径, 使用Hash可以平均的分布
// 文件到具体的存储位置, e.g. -> 0f9f2/84da2/864d3/3a2be/14b51
func CASPathTransformFunc(key string) PathKey {
	// [20]byte -> []byte -> [:]
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 8
	sliceLen := len(hashStr) / blockSize

	path := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i+1)*blockSize
		path[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(path, "/"),
		Original: hashStr,
	}
}

type PathTransferFunc func(string) PathKey

type PathKey struct {
	PathName string
	Original string
}

type StoreOpts struct {
	PathTransferFunc PathTransferFunc
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		Original: key,
	}
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) writeStream(key string, r io.Reader) error {
	// 获取Path名称并进行os级别创建
	pathName := s.PathTransferFunc(key)
	if err := os.MkdirAll(pathName.PathName, os.ModePerm); err != nil {
		return err
	}
	// 添加Buf进行读写
	buf := new(bytes.Buffer)
	io.Copy(buf, r)
	// 名称也进行hash -> 使用md5
	filenameBs := md5.Sum(buf.Bytes())
	filename := hex.EncodeToString(filenameBs[:])
	fullPath := pathName.PathName + "/" + filename
	// 打开/创建文件, 路径默认为当前pkg路径下
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	// 拷贝
	written, err := io.Copy(f, buf)
	if err != nil {
		return err
	}
	log.Printf("written (%d) bytes to disk: %s", written, fullPath)
	return nil
}
