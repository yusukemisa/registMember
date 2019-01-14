package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yusukemisa/registMember/src/infra"
)

var indexTemplate = template.Must(template.ParseFiles("template/index.html"))

type templateParams struct {
	Notice string
	Name   string
	Posts  []*Post
}

type Post struct {
	Author  string
	Message string
	Posted  time.Time
}

type PostHandler struct {
	Cli *infra.Client
}

func (h PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = h.handleGet(w, r)
	case "POST":
		err = h.handlePost(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h PostHandler) handleGet(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	params := templateParams{}
	posts, err := h.getPosts(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = fmt.Sprintf("failed to get posts:%s", err.Error())
		return indexTemplate.Execute(w, params)
	}
	params.Posts = posts
	return indexTemplate.Execute(w, params)

}

func (h PostHandler) handlePost(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	params := templateParams{}
	name := r.FormValue("name")
	message := r.FormValue("message")

	if name == "" || message == "" {
		w.WriteHeader(http.StatusBadRequest)
		params.Notice = "No name or message provided"
		return indexTemplate.Execute(w, params)
	}

	post := &Post{
		Author:  name,
		Message: message,
		Posted:  time.Now(),
	}

	key := h.Cli.Datastore.NameKey("Post", uuid.New().String(), nil)
	if _, err := h.Cli.Datastore.Put(ctx, key, post); err != nil {
		return err
	}

	posts, err := h.getPosts(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = fmt.Sprintf("failed to get posts:%s", err.Error())
		return indexTemplate.Execute(w, params)
	}
	params.Posts = posts
	params.Notice = fmt.Sprintf("Thank you for your submission, %s!", name)

	return indexTemplate.Execute(w, params)

}

func (h PostHandler) getPosts(ctx context.Context) ([]*Post, error) {
	var posts []*Post
	query := h.Cli.Datastore.NewQuery("Post").Order("Posted")
	if _, err := h.Cli.Datastore.GetAll(ctx, query, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}
