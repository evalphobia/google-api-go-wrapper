package storage

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"

	GCP "cloud.google.com/go/storage"
	"google.golang.org/api/option"
	SDK "google.golang.org/api/storage/v1"

	"github.com/evalphobia/google-api-go-wrapper/config"
	"github.com/evalphobia/google-api-go-wrapper/log"
)

const (
	serviceName = "storage"
)

// Storage repesents Cloud Storage API client.
type Storage struct {
	*GCP.Client
	logger log.Logger
}

// New returns initialized *Storage.
func New(ctx context.Context, conf config.Config) (*Storage, error) {
	if len(conf.Scopes) == 0 {
		conf.Scopes = []string{SDK.CloudPlatformScope}
	}

	httpClient, err := conf.Client()
	if err != nil {
		return nil, err
	}

	svc, err := GCP.NewClient(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	return &Storage{
		Client: svc,
		logger: log.DefaultLogger,
	}, nil
}

// SetLogger sets internal API logger.
func (s *Storage) SetLogger(logger log.Logger) {
	s.logger = logger
}

// UploadByBytes uploads an object from bytes.
func (s *Storage) UploadByBytes(byt []byte, opt ObjectOption) error {
	r := bytes.NewReader(byt)
	return s.Upload(r, opt)
}

// UploadByFile uploads an object from bytes.
func (s *Storage) UploadByFile(filepath string, opt ObjectOption) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	return s.Upload(f, opt)
}

// Upload uploads an object from io.Reader.
func (s *Storage) Upload(r io.Reader, opt ObjectOption) error {
	w := s.getObjectHandle(opt).NewWriter(opt.getOrCreateContext())
	if opt.CacheControl != "" {
		w.CacheControl = opt.CacheControl
	}

	_, err := io.Copy(w, r)
	if err != nil {
		s.Errorf("error on `object.write` operation by Upload; bucket=%s, path=%s, error=%s;", opt.BucketName, opt.Path, err.Error())
		return err
	}
	return w.Close()
}

// Delete deletes an object.
func (s *Storage) Delete(opt ObjectOption) error {
	handler := s.getObjectHandle(opt)
	err := handler.Delete(opt.getOrCreateContext())
	if hasDeleteError(err) {
		s.Errorf("error on `object.delete` operation by Delete; bucket=%s, path=%s, error=%s;", opt.BucketName, opt.Path, err.Error())
		return err
	}
	return nil
}

// Download downloads object data.
func (s *Storage) Download(opt ObjectOption) (data []byte, err error) {
	r, err := s.getObjectHandle(opt).NewReader(opt.getOrCreateContext())
	if err != nil {
		s.Errorf("error on creating reader; bucket=%s, path=%s, error=%s;", opt.BucketName, opt.Path, err.Error())
		return nil, err
	}
	defer r.Close()

	data, err = ioutil.ReadAll(r)
	if err != nil {
		s.Errorf("error on `object.get` operation by Download; bucket=%s, path=%s, error=%s;", opt.BucketName, opt.Path, err.Error())
	}
	return data, err
}

// Rename moves an object from opt.Path to destPath..
func (s *Storage) Rename(destPath string, opt ObjectOption) error {
	destOpt := opt
	destOpt.Path = destPath
	src := s.getObjectHandle(opt)
	dest := s.getObjectHandle(destOpt)

	ctx := opt.getOrCreateContext()
	_, err := dest.CopierFrom(src).Run(ctx)
	if err != nil {
		s.Errorf("error on `object.write` operation by Rename; bucket=%s, src=%s, dest=%s, error=%s;", opt.BucketName, opt.Path, destPath, err.Error())
		return err
	}

	err = src.Delete(ctx)
	if err != nil {
		s.Errorf("error on `object.delete` operation by Delete; bucket=%s, path=%s, error=%s;", opt.BucketName, opt.Path, err.Error())
	}
	return err
}

// Copy copies an object from opt.Path to destPath..
func (s *Storage) Copy(destPath string, opt ObjectOption) error {
	destOpt := opt
	destOpt.Path = destPath
	src := s.getObjectHandle(opt)
	dest := s.getObjectHandle(destOpt)

	ctx := opt.getOrCreateContext()
	_, err := dest.CopierFrom(src).Run(ctx)
	if err != nil {
		s.Errorf("error on `object.write` operation by Copy; bucket=%s, src=%s, dest=%s, error=%s;", opt.BucketName, opt.Path, destPath, err.Error())
	}
	return err
}

// CopyToBucket copies an object from opt.Path to another bucket.
func (s *Storage) CopyToBucket(destBucket, destPath string, opt ObjectOption) error {
	destOpt := opt
	destOpt.BucketName = destBucket
	destOpt.Path = destPath
	src := s.getObjectHandle(opt)
	dest := s.getObjectHandle(destOpt)

	ctx := opt.getOrCreateContext()
	_, err := dest.CopierFrom(src).Run(ctx)
	if err != nil {
		s.Errorf("error on `object.write` operation by Copy; bucket=%s, src=%s, dest=%s, error=%s;", opt.BucketName, opt.Path, destPath, err.Error())
	}
	return err
}

// IsExists checks if an object exists.
func (s *Storage) IsExists(opt ObjectOption) (isExist bool, err error) {
	_, err = s.Attrs(opt)
	switch {
	case isErrObjectNotExist(err):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

// Attrs gets attributes of the object.
func (s *Storage) Attrs(opt ObjectOption) (*GCP.ObjectAttrs, error) {
	handler := s.getObjectHandle(opt)
	a, err := handler.Attrs(opt.getOrCreateContext())
	if err != nil {
		s.Errorf("error on `object.get` operation by Attrs; bucket=%s, path=%s, error=%s;", opt.BucketName, opt.Path, err.Error())
	}
	return a, err
}

func (s *Storage) getObjectHandle(opt ObjectOption) *GCP.ObjectHandle {
	return s.Client.Bucket(opt.BucketName).Object(opt.Path)
}

// Errorf logging error information.
func (s *Storage) Errorf(format string, vv ...interface{}) {
	s.logger.Errorf(serviceName, format, vv...)
}

func hasDeleteError(err error) bool {
	switch {
	case err == nil, isErrObjectNotExist(err):
		return false
	default:
		return true
	}
}

func isErrObjectNotExist(err error) bool {
	return err != nil && err == GCP.ErrObjectNotExist
}
