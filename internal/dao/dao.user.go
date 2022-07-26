package dao

import (
	"context"
	"fmt"
	"go-common/library/database/sql"
	"go-common/library/log"
	"helloworld/internal/model"
	"time"
)

//AddUser 添加用户
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

	return user, nil
}

//DeleteUser 更新用户信息
func (d *dao) DeleteUser(c context.Context, uid int64) (row int64, err error) {
	deleteSql := fmt.Sprintf("DELETE FROM `users` WHERE uid=?;")

	res, err := d.db.Exec(c, deleteSql, uid)
	if err != nil {
		log.Error("db.Exec(%s) error(%v)", deleteSql, err)
		return 0, err
	}

	row, err = res.RowsAffected()
	return row, nil
}

//UpdateUser 更新用户信息
func (d *dao) UpdateUser(c context.Context, uid int64, nickname string, age int32) (row int64, err error) {
	querySql := fmt.Sprintf("UPDATE `users` SET nickname=?,age=?,uptime=? WHERE uid=?;")

	timenow := time.Now().Unix()
	res, err := d.db.Exec(c, querySql, nickname, age, timenow, uid)
	if err != nil {
		log.Error("db.Exec(%s) error(%v)", querySql, err)
		return 0, err
	}

	row, err = res.RowsAffected()
	return row, nil
}

//GetUser 查询用户
func (d *dao) GetUser(c context.Context, uid int64) (user *model.User, err error) {
	querySql := fmt.Sprintf("SELECT * FROM `users` WHERE uid=?;")

	user = new(model.User)
	err = d.db.QueryRow(c, querySql, uid).Scan(&user.Uid, &user.Nickname, &user.Age, &user.Uptime, &user.Addtime)
	if err != nil && err != sql.ErrNoRows {
		log.Error("d.QueryRow error(%v)", err)
		return
	}
	return user, nil
}

//GetUserList 查询用户列表
func (d *dao) GetUserList(c context.Context) (userlist []*model.User, err error) {
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
