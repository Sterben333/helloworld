package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
	"helloworld/internal/model"
	"time"
)

// AddUser 添加用户
func (d *dao) AddUser(c context.Context, nickname string, age int32) (user *model.User, err error) {
	querySql := fmt.Sprintf("INSERT INTO `users`(uid,nickname,age,uptime,addtime) VALUES(null,?,?,?,?);")

	timenow := time.Now().Unix()
	res, err := d.db.Exec(c, querySql, nickname, age, timenow, timenow)
	if err != nil {
		log.Error("db.Exec(%s) error(%v)", querySql, err)
		return nil, err
	}
	user = new(model.User)
	user.Uid, _ = res.LastInsertId()
	user.Nickname = nickname
	user.Age = age
	user.Addtime = int32(timenow)
	user.Uptime = int32(timenow)

	//--------------同时添加到redis中---------------//
	bytes, _ := json.Marshal(&user)                        // 将user结构体转换成切片
	_, err = d.redis.Do(c, "SET", user.Uid, string(bytes)) //将[]byte转换成string，并进行添加到redis中，key为uid。value是user
	// 错误处理
	if err != nil && err != redis.ErrNoReply {
		log.Error("d.redis.Do error(%v)", err)
		return
	}
	d.redis.Do(c, "EXPIRE", user.Uid, 20) //设置过期时间20s

	return user, nil
}

//DeleteUser 删除用户信息
func (d *dao) DeleteUser(c context.Context, uid int64) (row int64, err error) {
	//这里考虑一下多线程可能造成到问题，先删redis还是先删除mysql到问题

	fmt.Println("删除开始")
	//我采用了延时双删，先删redis，再删mysql，再删redis
	flag, _ := redis.Int(d.redis.Do(c, "EXISTS", uid))
	fmt.Println(flag)
	if flag != 0 {
		d.redis.Do(c, "DEL", uid)
	}
	deleteSql := fmt.Sprintf("DELETE FROM `users` WHERE uid=?;")
	res, err := d.db.Exec(c, deleteSql, uid)
	if err != nil {
		log.Error("db.Exec(%s) error(%v)", deleteSql, err)
		return 0, err
	}
	time.Sleep(2) //睡2s再删除redis
	flag, _ = redis.Int(d.redis.Do(c, "EXISTS", uid))
	if flag != 0 {
		d.redis.Do(c, "DEL", uid)
	}

	row, err = res.RowsAffected()
	return row, nil
}

//UpdateUser 更新用户信息
func (d *dao) UpdateUser(c context.Context, uid int64, nickname string, age int32) (row int64, err error) {
	querySql := fmt.Sprintf("UPDATE `users` SET nickname=?,age=?,uptime=? WHERE uid=?;")

	//为了防止读到脏数据，对redis不做更新处理而做删除处理
	flag, _ := redis.Int(d.redis.Do(c, "EXISTS", uid))
	if flag != 0 {
		d.redis.Do(c, "DEL", uid)
	}
	timenow := time.Now().Unix()
	res, err := d.db.Exec(c, querySql, nickname, age, timenow, uid)
	if err != nil {
		log.Error("db.Exec(%s) error(%v)", querySql, err)
		return 0, err
	}

	time.Sleep(2) //睡2s再删除redis
	flag, _ = redis.Int(d.redis.Do(c, "EXISTS", uid))
	if flag != 0 {
		d.redis.Do(c, "DEL", uid)
	}

	row, err = res.RowsAffected()
	return row, nil
}

//GetUser 查询用户
func (d *dao) GetUser(c context.Context, uid int64) (user *model.User, err error) {
	user = new(model.User)

	// 首先查看redis缓存，如果有直接返回
	flag, _ := redis.Int(d.redis.Do(c, "EXISTS", uid))
	if flag != 0 {
		out, _ := redis.String(d.redis.Do(c, "GET", uid))
		if err != nil && err != redis.ErrNoReply {
			log.Error("d.redis.Do error(%v)", err)
			return
		}
		fmt.Println("redis中存在值，直接返回")
		json.Unmarshal([]byte(out), &user) //将字符串变成[]byte，再变成user结构体
		return user, nil
	}

	fmt.Println("redis中没有不存在该值，将要去mysql中寻找")

	querySql := fmt.Sprintf("SELECT * FROM `users` WHERE uid=?;")

	err = d.db.QueryRow(c, querySql, uid).Scan(&user.Uid, &user.Nickname, &user.Age, &user.Uptime, &user.Addtime)
	if err != nil && err != sql.ErrNoRows {
		log.Error("d.QueryRow error(%v)", err)
		return
	}
	//从mysql中寻找到了值，同时我也要加入到redis中，方便下次查找
	bytes, _ := json.Marshal(&user)
	_, err = d.redis.Do(c, "SET", user.Uid, string(bytes)) //将user结构体转换成[]byte，再转换成string
	if err != nil && err != redis.ErrNoReply {
		log.Error("d.redis.Do error(%v)", err)
		return
	}
	d.redis.Do(c, "EXPIRE", user.Uid, 20) //设置过期时间

	return user, nil
}

//GetUserList 查询用户列表
func (d *dao) GetUserList(c context.Context) (userlist []*model.User, err error) {
	//对于用户列表，我觉得太大了，全放到redis中不是很好，而且这个操作不是很频繁，所以不做redis对处理
	querySql := fmt.Sprintf("SELECT * FROM `users`;")
	rows, err := d.db.Query(c, querySql)
	if err != nil {
		log.Error("query  error(%v)", err)
		return
	}
	defer rows.Close()

	userlist = make([]*model.User, 0)
	for rows.Next() {
		user := new(model.User)

		if err = rows.Scan(&user.Uid, &user.Nickname, &user.Age, &user.Uptime, &user.Addtime); err != nil {
			log.Error("scan demo log error(%v)", err)
			return
		}
		userlist = append(userlist, user)
	}
	return userlist, nil
}
