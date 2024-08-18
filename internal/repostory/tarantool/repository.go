package tarantool

import (
	"context"
	"github.com/tarantool/go-tarantool/v2"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
	"vk_tarantool_project/internal/domain"
)

const (
	userSpaceName      = "users"
	userSpaceIndexPK   = "primary"
	dataSpaceName      = "data"
	dataSpaceIndexPK   = "primary"
	dataSpaceIndexHash = "hash_key"
)

type Tarantool struct {
	conn       *tarantool.Connection
	reqTimeout time.Duration
	mu         sync.Mutex
}

// New create new repository for data manipulating in Tarantool DBMS
func New(conn *tarantool.Connection, reqTimeout time.Duration) *Tarantool {
	return &Tarantool{
		conn:       conn,
		reqTimeout: reqTimeout,
	}
}

// GetUserByName get *domain.UserInfo by string username
func (t *Tarantool) GetUserByName(ctx context.Context, username string) (*domain.UserInfo, error) {
	data, err := t.conn.Do(
		tarantool.NewSelectRequest(userSpaceName).
			Index(userSpaceIndexPK).Limit(1).
			Iterator(tarantool.IterEq).
			Key([]interface{}{username}),
	).Get()

	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, domain.ErrUserNotFound
	}

	if len(data) > 1 {
		return nil, domain.ErrMultipleUser
	}

	var usr domain.UserInfo
	usr.Username = data[0].([]interface{})[0].(string)
	usr.Password = data[0].([]interface{})[1].(string)

	return &usr, nil
}

// WriteData write KV data of format string:any
// WriteData use Upsert method for insert or update records in Tarantool
// The method starts updating each field in a separate goroutine
// If one of the requests to update a field returns an error, then the entire request fails with an error
func (t *Tarantool) WriteData(ctx context.Context, data domain.Data) error {

	// Error group for errors handling in goroutine
	g, gctx := errgroup.WithContext(ctx)

	// Start insert or update each field in the request
	for key, val := range data.Data {
		key, val := key, val
		g.Go(func() error {
			// Timeout limits the time it takes to complete a request
			reqTimeoutCtx, cancel := context.WithTimeout(gctx, t.reqTimeout)
			defer cancel()

			if _, err := t.conn.Do(
				tarantool.
					NewUpsertRequest(dataSpaceName).
					Context(reqTimeoutCtx).
					Tuple([]interface{}{key, val}).
					Operations(
						tarantool.NewOperations().Assign(1, val),
					),
			).Get(); err != nil {
				return err
			}

			return nil
		})
	}

	// Check error from goroutine
	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (t *Tarantool) ReadData(ctx context.Context, keys domain.DataKeys) (map[interface{}]interface{}, error) {

	resultData := make(map[interface{}]interface{})

	// Error group for errors handling in goroutine
	g, gctx := errgroup.WithContext(ctx)

	// Start select each field in the request
	for _, key := range keys.Keys {
		key := key
		g.Go(func() error {

			// Timeout limits the time it takes to complete a request
			reqTimeoutCtx, cancel := context.WithTimeout(gctx, t.reqTimeout)
			defer cancel()

			// Select query with key data
			data, err := t.conn.Do(
				tarantool.NewSelectRequest(dataSpaceName).
					Context(reqTimeoutCtx).
					Limit(1).
					Index(dataSpaceIndexHash).
					Iterator(tarantool.IterEq).
					Key([]interface{}{key}),
			).Get()

			if err != nil {
				return err
			}

			// If no data found, return without error
			// and don't update key in resultData
			if len(data) == 0 {
				return nil
			}

			// Update key in resultData if data found
			t.mu.Lock()
			resultData[key] = data[0].([]interface{})[1]
			t.mu.Unlock()

			return nil
		})
	}

	// Check error from goroutine
	if err := g.Wait(); err != nil {
		return resultData, err
	}

	return resultData, nil
}
