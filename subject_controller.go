package main

import "coredemo/framework"

func SubjectAddController(c *framework.Context) error {
	c.Json("ok, SubjectAddController").SetStatus(200)
	return nil
}

func SubjectListController(c *framework.Context) error {
	c.Json("ok, SubjectListController").SetStatus(200)
	return nil
}

func SubjectDelController(c *framework.Context) error {
	c.Json("ok, SubjectDelController").SetStatus(200)
	return nil
}

func SubjectUpdateController(c *framework.Context) error {
	c.Json("ok, SubjectUpdateController").SetStatus(200)
	return nil
}

func SubjectGetController(c *framework.Context) error {
	c.Json("ok, SubjectGetController").SetStatus(200)
	return nil
}

func SubjectNameController(c *framework.Context) error {
	c.Json("ok, SubjectNameController").SetStatus(200)
	return nil
}
