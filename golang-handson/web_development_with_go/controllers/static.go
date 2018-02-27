package controllers

import "github.com/YmgchiYt/golang-handson/web_development_with_go/views"

type Static struct {
	Home    *views.View
	Contact *views.View
}

func NewStatic() *Static {
	return &Static{
		Home: views.NewView(
			"bootstrap", "static/home"),
		Contact: views.NewView(
			"bootstrap", "static/contact"),
	}
}
