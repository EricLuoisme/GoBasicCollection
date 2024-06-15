package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	// call
	key := "anotherPic"
	pathKey := CASPathTransformFunc(key)
	targetPathOriginal := "0f9f284da2864d33a2be14b510f6f0fd5383991d"
	targetPathName := "0f9f284d/a2864d33/a2be14b5/10f6f0fd/5383991d"
	// check
	if pathKey.PathName != targetPathName || pathKey.Original != targetPathOriginal {
		t.Errorf("have %s want %s", pathKey.PathName, targetPathName)
		t.Errorf("have %s want %s", pathKey.Original, targetPathOriginal)
	}
}

// TestNewStore 直接进行简单的保存测试
func TestNewStore(t *testing.T) {
	// call
	opts := StoreOpts{
		//PathTransferFunc: DefaultPathTransformFunc,
		PathTransferFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	data := bytes.NewReader([]byte("some jpg bytes again"))
	// check
	if err := s.writeStream("specialPic", data); err != nil {
		t.Error(err)
	}
}
