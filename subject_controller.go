package main

import (
	"github.com/sunranlike/hade/framework/gin"
)

func SubjectAddController(c *gin.Context) error {
	c.IJson("ok, SubjectAddController").ISetStatus(200)
	return nil
}

func SubjectListController(c *gin.Context) error {
	c.IJson("ok, SubjectListController").ISetStatus(200)
	return nil
}

func SubjectDelController(c *gin.Context) error {
	c.IJson("ok, SubjectDelController").ISetStatus(200)
	return nil
}

func SubjectUpdateController(c *gin.Context) error {
	c.IJson("ok, SubjectUpdateController").ISetStatus(200)
	return nil
}

func SubjectGetController(c *gin.Context) error {
	c.IJson("ok, SubjectGetController").ISetStatus(200)
	return nil
}

func SubjectNameController(c *gin.Context) error {
	c.IJson("ok, SubjectNameController").ISetStatus(200)
	return nil
}
