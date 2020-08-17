package model

import (
	"learn_tools/gin_demo/global"
	"time"
)

type Notice struct {
	Id         int    `xorm:"not null pk autoincr INT(11)" json:"id,omitempty"`
	Title      string `xorm:"not null default '' comment('标题') VARCHAR(50)" json:"title,omitempty"`
	Content    string `xorm:"not null comment('内容') TEXT" json:"content,omitempty"`
	Url        string `xorm:"default '' comment('超链接') VARCHAR(150)" json:"url,omitempty"`
	CreateTime int64  `xorm:"not null default 0 comment('创建时间') BIGINT(20)" json:"create_time,omitempty"`
	UpdateTime int64  `xorm:"default 0 comment('更新时间') BIGINT(20)" json:"update_time,omitempty"`
	Del        bool   `xorm:"not null default 0 comment('0未删除，1已删除') Bool" json:"del,omitempty"`
}

func Add(notice *Notice) Err {
	exit, err := global.My.DB.Exist(&Notice{Title: notice.Title})
	if err != nil {
		return Err{ModelErr, err}
	}
	if exit {
		return Err{ModelErrExisted, nil}
	}
	t := time.Now().Unix()

	notice.CreateTime = t
	notice.UpdateTime = t
	_, err = global.My.DB.Insert(notice)
	if err != nil {
		return Err{ModelErr, err}
	}
	return Err{ModelNoErr, nil}
}

func Del(id int) Err {
	//n, err := global.My.DB.Id(id).Delete(&Notice{})
	//if err != nil || n <= 0 {
	//	return Err{ModelErrDel, err}
	//}
	//return Err{ModelNoErr, nil}
	n, err := global.My.DB.Id(notice.Id).Update(notice)
	if err != nil || n <= 0 {
		return Err{ModelErrUpdate, err}
	}
	return Err{ModelNoErr, nil}
}

func Update(notice *Notice) Err {
	n, err := global.My.DB.Id(notice.Id).Update(notice)
	if err != nil || n <= 0 {
		return Err{ModelErrUpdate, err}
	}
	return Err{ModelNoErr, nil}
}

func Get(id int) (notice Notice, err Err) {
	ok, e := global.My.DB.Id(id).Get(&notice)
	if !ok || e != nil {
		err.Code = ModelErrFind
		err.Err = e
		return
	}
	err.Code = ModelNoErr
	err.Err = nil
	return
}

func GetAll() (notices []Notice, err Err) {
	e := global.My.DB.Where("del = 0").Find(&notices)
	if e != nil {
		err.Code = ModelErr
		err.Err = e
		return
	}
	err.Code = ModelNoErr
	err.Err = nil
	return
}
