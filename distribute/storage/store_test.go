package main

import (
	"bytes"
	"testing"
)

// TestNewStore 直接进行简单的保存测试
func TestNewStore(t *testing.T) {
	opts := StoreOpts{
		PathTransferFunc: DefaultPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpg bytes"))
	if err := s.writeStream("specialPic", data); err != nil {
		t.Error(err)
	}
}
