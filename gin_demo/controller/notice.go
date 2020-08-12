package controller

import (
	"github.com/gin-gonic/gin"
	"learn_tools/gin_demo/model"
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

	if e.Code == model.ModelNoErr {
		c.SuccessWithData(g, n)
	} else {
		log.Errorf("Get get data  error(%v)", err)
		c.FailureForServer(g, " failure.")
	}
}

func (c *NoticeController) GetAll(g *gin.Context) {
	n, e := model.GetAll()
	if e.Code == model.ModelNoErr {
		c.SuccessWithData(g, n)
	} else {
		log.Errorf("Get get all data  error(%v)", e)
		c.FailureForServer(g, "get data failure.")
	}
}
