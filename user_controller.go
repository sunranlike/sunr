package main

import "coredemo/framework"

func UserLoginController(c *framework.Context) error {
	c.Json("ok, UserLoginController").SetStatus(200)
	return nil
}
