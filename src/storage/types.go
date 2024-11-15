package storage

import "github.com/scope3-dio/clients/scope3"

type Fetcher interface {
	FetchProperty([]scope3.PropertyQuery) ([]scope3.PropertyResponse, error)
}
