package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		Original: key,
	}
}

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

func (k PathKey) filename() string {
	return fmt.Sprintf("%s/%s", k.PathName, k.Original)
}

type StoreOpts struct {
	PathTransferFunc PathTransferFunc
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
	pathKey := s.PathTransferFunc(key)
	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}
	// 打开/创建文件, 路径默认为当前pkg路径下
	fileFullPath := pathKey.filename()
	f, err := os.Create(fileFullPath)
	if err != nil {
		return err
	}
	// 拷贝
	written, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	log.Printf("written (%d) bytes to disk: %s", written, fileFullPath)
	return nil
}
