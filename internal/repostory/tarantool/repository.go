package tarantool

import (
	"context"
	"github.com/tarantool/go-tarantool/v2"
	"vk_tarantool_project/internal/domain"
)

const (
	userSpaceName    = "users"
	userSpaceIndexPK = "pk"
)

type Tarantool struct {
	conn *tarantool.Connection
}

// New create new repository for data manipulating in Tarantool DBMS
func New(conn *tarantool.Connection) *Tarantool {
	return &Tarantool{
		conn: conn,
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
