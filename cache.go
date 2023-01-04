package unixodbc

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type cachedStatement struct {
	stmt       *stmt
	accessTime time.Time
}

type psCache struct {
	cache map[string]*cachedStatement
}

func (c *psCache) Get(query string, set func() (*stmt, error)) (*stmt, error) {
	key := c.key(query)
	if cachedStmt, found := c.cache[key]; found {
		cachedStmt.accessTime = time.Now()
		return cachedStmt.stmt, nil
	}
	s, err := set()
	if err != nil {
		return s, err
	}

	c.prune()
	c.cache[key] = &cachedStatement{
		stmt:       s,
		accessTime: time.Now(),
	}
	return s, nil
}
func (c *psCache) key(query string) string {
	h := sha256.New()
	h.Write([]byte(query))
	return hex.EncodeToString(h.Sum(nil))
}

func (c *psCache) prune() {

}
