package sync_services

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

// DropboxProvider — реализация для Dropbox
type DropboxProvider struct {
	Client files.Client
}

type Provider interface {
	Upload(string, []byte) error
	Download(string) ([]byte, error)
}

func NewDropboxProvider(token string) *DropboxProvider {
	config := dropbox.Config{
		Token:    token,
		LogLevel: dropbox.LogOff,
	}
	return &DropboxProvider{
		Client: files.New(config),
	}
}

// Upload(path, file) easily uploads/updates/overwrites the file to/on the server
func (d *DropboxProvider) Upload(path string, file []byte) error {
	path = filepath.Join("/", path)

	uploadArg := files.NewUploadArg(path)
	uploadArg.Mode = &files.WriteMode{Tagged: dropbox.Tagged{Tag: "overwrite"}}
	uploadArg.Mute = true

	_, err := d.Client.Upload(uploadArg, bytes.NewReader(file))
	if err != nil {
		return err
	}
	return nil
}

// Download(path string) downloads file, returns errros specified in
func (d *DropboxProvider) Download(path string) ([]byte, error) {
	path = filepath.Join("/", path)
	_, contents, err := d.Client.Download(files.NewDownloadArg(path))
	if err != nil {
		// Parsing error using auth.ParseError to handle not_found response
		var appErr files.DownloadAPIError
		parsedErr := auth.ParseError(err, &appErr)

		// cheking the error on not_found Tag
		if e, ok := parsedErr.(files.DownloadAPIError); ok {
			if e.EndpointError != nil && e.EndpointError.Path != nil && e.EndpointError.Path.Tag == "not_found" {
				return nil, ErrFileDoesntExist
			}
		}

		return nil, parsedErr
	}
	defer func() {
		if err := contents.Close(); err != nil {
			fmt.Println("error closing content")
		}
	}()
	data, err := io.ReadAll(contents)
	if err != nil {
		return nil, fmt.Errorf("failed to read downloaded file: %w", err)
	}
	return data, nil
}
