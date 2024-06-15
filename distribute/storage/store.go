package main

import (
	"io"
	"log"
	"os"
)

type PathTransferFunc func(string) string

type StoreOpts struct {
	PathTransferFunc PathTransferFunc
}

var DefaultPathTransformFunc = func(key string) string {
	return key
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
	if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
		return err
	}
	// 打开文件
	fullPath := pathName + "/" + "thefilename"
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	// 拷贝
	written, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	log.Printf("written (%d) bytes to disk: %s", written, fullPath)
	return nil
}
