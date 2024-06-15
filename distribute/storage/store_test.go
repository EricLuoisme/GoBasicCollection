package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	// call
	key := "anotherPic"
	pathName := CASPathTransformFunc(key)
	targetPathName := "0f9f2/84da2/864d3/3a2be/14b51/0f6f0/fd538/3991d"
	// check
	if pathName != targetPathName {
		t.Errorf("have %s want %s", pathName, targetPathName)
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
	data := bytes.NewReader([]byte("some jpg bytes"))
	// check
	if err := s.writeStream("specialPic", data); err != nil {
		t.Error(err)
	}
}
