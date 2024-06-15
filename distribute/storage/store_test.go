package main

import "testing"

func TestNewStore(t *testing.T) {
	opts := StoreOpts{
		PathTransferFunc: DefaultPathTransformFunc,
	}
	s := NewStore(opts)

}
