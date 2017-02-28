// Copyright (C) 2017 Anonymous
// RWTXliExJCquO54R+qP94i4V+X8bQegE6L9EjhKIH23ePweJG8u7dqDK
//
// This program is free software: you can redistribute it and/or modify it under
// the terms of the GNU Affero General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option) any
// later version. See COPYING for the full text of the License.

package rants

import (
	"encoding/xml"
	"github.com/russross/blackfriday"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
)

type State struct {
	License template.HTML
	Index   template.HTML
	ctx     context.Context
	tmpl    *template.Template
}

func newState() *State {
	license, err := ioutil.ReadFile("templates/license-snippet.txt")
	if err != nil {
		panic(err)
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		panic(err)
	}

	md, err := ioutil.ReadFile("content/index.md")
	index := blackfriday.MarkdownBasic(md)

	return &State{
		License: template.HTML(license),
		Index:   template.HTML(index),
		tmpl:    tmpl,
	}
}

func (c *State) error(err error) {
	log.Errorf(c.ctx, "%v", err)
}

func (c *State) index(w http.ResponseWriter, r *http.Request) {
	c.ctx = appengine.NewContext(r)
	w.Header().Set("Content-Type", "text/html")
	if err := c.tmpl.Execute(w, c); err != nil {
		c.error(err)
	}
}

func (c *State) api(w http.ResponseWriter, r *http.Request) {
	c.ctx = appengine.NewContext(r)
	w.Header().Set("Content-Type", "text/xml")
	io.WriteString(w, xml.Header)
}

func init() {
	c := newState()
	http.HandleFunc("/", c.index)
	http.HandleFunc("/v1/", c.api)
}
