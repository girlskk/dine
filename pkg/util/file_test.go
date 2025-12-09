package util

import "testing"

func TestGetFileNameAndExt(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		fileName, ext := GetFileNameAndExt("/path/to/test.txt")
		if fileName != "test" || ext != ".txt" {
			t.Errorf("expected (test, .txt), got (%s, %s)", fileName, ext)
		}
	})

	t.Run("test2", func(t *testing.T) {
		fileName, ext := GetFileNameAndExt("path/to/test")
		if fileName != "test" || ext != "" {
			t.Errorf("expected (test, ''), got (%s, %s)", fileName, ext)
		}
	})
}
