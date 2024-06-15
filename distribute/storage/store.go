package main

import (
	"io"
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
	filename := "thefilename"
	f, err := os.Open(pathName + "/" + filename)
	if err != nil {
		return err
	}

	return nil
}
