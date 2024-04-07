package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

type Storage struct {
}

func NewFileStorage() Storage {
	return Storage{}
}

func (s Storage) SaveFile(ctx context.Context, file *multipart.File, root, ext, name string) (string, error) {
	path := fmt.Sprintf("%s/%s%s", root, name, ext)

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, *file)
	if err != nil {
		return "", err
	}

	return path, nil
}
