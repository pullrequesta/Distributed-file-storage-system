package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "myspecialpicture"
	PathKey := CASPathTransformFunc(key)

	expectedOriginalKey := "d9e06924cbe4f7c5f59269e6267f971d02774564"
	expectedPathName := "d9e06/924cb/e4f7c/5f592/69e62/67f97/1d027/74564"

	if PathKey.PathName != expectedPathName {
		t.Errorf("got %s want %s\n", PathKey.PathName, expectedPathName)
	}
	if PathKey.Original != expectedOriginalKey {
		t.Errorf("got %s want %s\n", PathKey.Original, expectedOriginalKey)
	}

}

func TestStore(t *testing.T) {

	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)
	data := bytes.NewReader([]byte("some jpg bytes"))

	if err := s.writeStream("myspecialpicture", data); err != nil {
		t.Error(err)
	}

}
