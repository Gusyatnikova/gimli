package impl

import (
	"context"
	"gimli/Internal/gRPC/domain"
	random_string "gimli/Internal/random-string"
	"gimli/Internal/storage"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789_"
	length = 10
)

//ShortenerServiceImpl implements gRPC service
type ShortenerServiceImpl struct {
	db storage.DatabasePtr
}

func (s ShortenerServiceImpl) Begin() {
	if !storage.Running {
		s.db = storage.SetupConnection()
	}
}

func (s ShortenerServiceImpl) End() {
	storage.CloseConnection()
}

func (s ShortenerServiceImpl) Create(
	c context.Context, fullUrl *domain.FullURL) (
	shortPath *domain.ShortPath, err error) {
	//todo: check if URL is Valid

	shortStr := random_string.StringWithCharset(length, charset)
	var shortUrl *storage.ShortURL
	shortUrl, err = storage.InsertPairURL(fullUrl.Url, shortStr)
	if err != nil {
		return nil, err
	}
	shortPath = &domain.ShortPath{}
	*shortPath = storage.ShortURLToShortPath(*shortUrl)
	return
}

func (s ShortenerServiceImpl) Get(
	c context.Context, shortPath *domain.ShortPath) (
	fullURL *domain.FullURL, err error) {
	origin, err := storage.GetOriginByShort(shortPath.Path)
	if err != nil {
		return nil, err
	}
	fullURL = &domain.FullURL{}
	*fullURL = storage.OriginToFull(*origin)
	return fullURL, err
}
