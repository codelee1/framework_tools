package controller

import (
	"gin_demo/model"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type NoticeController struct {
	BaseController
}

func (c *NoticeController) Add(g *gin.Context) {
	var m model.Notice
	err := g.BindJSON(&m)
	if err != nil {
		log.Errorf("Add bind json  error(%v)", err)
		c.FailureForClient(g, "json change failure.")
		return
	}
	e := model.Add(&m)
	log.Debug("add data ",m)
	if e.Code == model.ModelNoErr {
		c.Success(g)
	} else {
		log.Errorf("Add error(%v)", err)
		c.FailureForServer(g, "add data failure.")
	}
}

func (c *NoticeController) Del(g *gin.Context) {
	sid := g.Param("id")
	if sid == "" {
		c.FailureForClient(g, "id is null")
		return
	}
	id, err := strconv.Atoi(sid)
	if err != nil {
		log.Errorf("Add id change  error(%v)", err)
		c.FailureForClient(g, "id change failure.")
	}
	e := model.Del(id)
	log.Debug("del data ",id)
	if e.Code == model.ModelNoErr {
		c.Success(g)
	} else {
		log.Errorf("Del  error(%v)", err)
		c.FailureForServer(g, "add data failure.")
	}
}

func (c *NoticeController) Update(g *gin.Context) {
	sid := g.Param("id")
	if sid == "" {
		c.FailureForClient(g, "id is null")
		return
	}
	id, err := strconv.Atoi(sid)
	if err != nil {
		log.Errorf("Update id change  error(%v)", err)
		c.FailureForClient(g, "id change failure.")
	}
	var m model.Notice
	err = g.BindJSON(&m)
	if err != nil {
		log.Errorf("Update bind json  error(%v)", err)
		c.FailureForClient(g, "json change failure.")
	}
	m.Id = id
	m.UpdateTime = time.Now().Unix()
	e := model.Update(&m)
	log.Debug("update data ",m)
	if e.Code == model.ModelNoErr {
		c.Success(g)
	} else {
		log.Errorf("Update  error(%v)", err)
		c.FailureForServer(g, "update data failure.")
	}
}
func (c *NoticeController) Get(g *gin.Context) {
	sid := g.Param("id")
	if sid == "" {
		c.FailureForClient(g, "id is null")
		return
	}
	id, err := strconv.Atoi(sid)
	if err != nil {
		log.Errorf("Get id change error(%v)", err)
		c.FailureForClient(g, "id change failure.")
	}
	n, e := model.Get(id)
	log.Debug("get get data ",n)

	if e.Code == model.ModelNoErr {
		c.SuccessWithData(g, n)
	} else {
		log.Errorf("get data  error(%v)", err)
		c.FailureForServer(g, " failure.")
	}
}

func (c *NoticeController) GetAll(g *gin.Context) {
	n, e := model.GetAll()
	log.Debug("get all data ",n)
	if e.Code == model.ModelNoErr {
		c.SuccessWithData(g, n)
	} else {
		log.Errorf("get all data  error(%v)", e)
		c.FailureForServer(g, "get data failure.")
	}
}
