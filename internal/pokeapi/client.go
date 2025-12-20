package pokeapi

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/C0Mon/pokedexcli/internal/pokecache"
)

type Client struct {
	httpClient http.Client
	c          pokecache.Cache
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	c := Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		c: pokecache.NewCache(cacheInterval),
	}
	return c
}

func (c *Client) GetData(url string) ([]byte, error) {
	// Check if data in cache
	data, found := c.c.Get(url)
	if found {
		return data, nil
	}

	// Call data through client
	res, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()

	// Read into correct format
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// Add response to cache
	c.c.Add(url, body)
	return body, nil
}
