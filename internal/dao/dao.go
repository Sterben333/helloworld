package dao

import (
	"context"
	"time"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf/paladin.v2"
	"go-common/library/database/sql"
	"go-common/library/sync/pipeline/fanout"
	xtime "go-common/library/time"
	"helloworld/internal/model"

	"github.com/google/wire"
)

// Provider 声明依赖注入对象
var Provider = wire.NewSet(New, NewDB, NewRedis, NewMC)

//go:generate kratos tool btsgen
// Dao dao interface
type Dao interface {
	Close()
	Ping(ctx context.Context) (err error)
	// bts: -nullcache=&model.Article{ID:-1} -check_null_code=$!=nil&&$.ID==-1
	Article(c context.Context, id int64) (*model.Article, error)
	//新增接口
	AddUser(c context.Context, nickname string, age int32) (user *model.User, err error)
	UpdateUser(c context.Context, uid int64, nickname string, age int32) (row int64, err error)
	GetUser(c context.Context, uid int64) (user *model.User, err error)
	GetUserList(c context.Context) (userlist []*model.User, err error)
}

// dao dao.
type dao struct {
	db         *sql.DB
	redis      *redis.Redis
	mc         *memcache.Memcache
	cache      *fanout.Fanout
	demoExpire int32
}

// New new a dao and return.
// 使用参数接收连接池对象
func New(r *redis.Redis, mc *memcache.Memcache, db *sql.DB) (d Dao, cf func(), err error) {
	return newDao(r, mc, db)
}

// NewDao 根据参数初始化dao
func newDao(r *redis.Redis, mc *memcache.Memcache, db *sql.DB) (d *dao, cf func(), err error) {
	var cfg struct {
		DemoExpire xtime.Duration
	}
	if err = paladin.Get("application.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	d = &dao{
		db:         db, //官方文档直接在这里初始化
		redis:      r,
		mc:         mc,
		cache:      fanout.New("cache"),
		demoExpire: int32(time.Duration(cfg.DemoExpire) / time.Second),
	}
	cf = d.Close
	return
}

// Close close the resource.
func (d *dao) Close() {
	d.cache.Close()
}

// Ping ping the resource.
func (d *dao) Ping(ctx context.Context) (err error) {
	return nil
}
