package application

import (
	"context"
	"time"
)

type Application struct {
	cache      ICache
	fileWriter IFileWriter
	context    context.Context
	ticker     *time.Ticker
}

func NewApplication(cache ICache, filewriter IFileWriter, ctx context.Context) *Application {
	return &Application{cache: cache, fileWriter: filewriter, context: ctx, ticker: time.NewTicker(time.Second * 10)}
}

func (a *Application) Get(key string) (string, error) {
	return a.cache.Get(key)
}

func (a *Application) Set(key, value string) error {
	return a.cache.Set(key, value)
}

func (a *Application) WriteToFile() error {
	for {
		select {
		case <-a.context.Done():
			return nil
		case <-a.ticker.C:
			allValues, err := a.cache.GetAll()
			if err != nil {
				return err
			}
			err = a.fileWriter.Write(allValues)
			if err != nil {
				return err
			}
		}

	}
}

type ICache interface {
	Get(string) (string, error)
	GetAll() (map[string]string, error)
	Set(string, string) error
}

type IFileWriter interface {
	Write(map[string]string) error
}
