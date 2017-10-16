package static

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// Engine is an interface like Gin Engine
type Engine interface {
	Static(relativePath, root string) gin.IRoutes
	StaticFile(relativePath, filepath string) gin.IRoutes
}

// Factory .
type Factory struct {
	engine Engine
	pathes map[string]*File
}

// File .
type File struct {
	path string
	URL  string
}

// Install init Static url struct
func Install(router Engine) *Factory {
	return &Factory{
		engine: router,
		pathes: map[string]*File{},
	}
}

// TODO fix
func getFileHash(path string) string {
	file, err := os.Open("utf8.txt")
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return ""
	}

	return string(hash.Sum(nil))
}

// Static patch Static from gin
func (factory *Factory) Static(relativePath, root string) *Factory {
	relativePath = "/" + strings.Trim(relativePath, "/")
	root = "/" + strings.Trim(root, "/")

	filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		factory.StaticFile(relativePath+"/"+path, root+"/"+path)
		return nil
	})
	return factory
}

// StaticFile patch Static from gin
func (factory *Factory) StaticFile(relativePath, filepath string) *Factory {
	factory.pathes[relativePath] = &File{
		path: filepath,
	}
	return factory
}

// URL get Static url with query
func (factory *Factory) URL(path string) string {
	file, ok := factory.pathes[path]
	if !ok {
		return path
	}
	if file.URL == "" {
		file.URL = fmt.Sprintf("%s?hash=%s", path, getFileHash(file.path))
	}
	return file.URL
}

// URL get Static url with query
func URL(path string) string {
	return ""
}
