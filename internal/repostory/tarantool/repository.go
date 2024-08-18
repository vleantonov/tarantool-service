package tarantool

import (
	"context"
	"github.com/tarantool/go-tarantool/v2"
	"golang.org/x/sync/errgroup"
	"time"
	"vk_tarantool_project/internal/domain"
)

const (
	userSpaceName    = "users"
	userSpaceIndexPK = "pk"
	dataSpaceName    = "data"
	dataSpaceIndexPK = "pk"
)

type Tarantool struct {
	conn       *tarantool.Connection
	reqTimeout time.Duration
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
					Tuple([]interface{}{key, val}),
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
