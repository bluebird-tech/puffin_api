package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	"log"
	"testing"
)

func TestPostEvent(t *testing.T) {
	i := Impl{}
	i.InitDB()
	i.InitSchema()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/events", i.PostEvent),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	recorded := test.RunRequest(t, api.MakeHandler(),
		test.MakeSimpleRequest("POST", "http://1.2.3.4/events", &map[string]string{"measurement": "stress"}))
	recorded.CodeIs(201)
	recorded.ContentTypeIsJson()
}

func TestPostInstallation(t *testing.T) {
	i := Impl{}
	i.InitDB()
	i.InitSchema()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/installations", i.PostInstallation),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	recorded := test.RunRequest(t, api.MakeHandler(),
		test.MakeSimpleRequest("POST", "http://1.2.3.4/installations", &map[string]string{"deviceType": "android"}))
	recorded.CodeIs(201)
	recorded.ContentTypeIsJson()
}

func TestGetHealth(t *testing.T) {
	i := Impl{}
	i.InitDB()
	i.InitSchema()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/health", i.GetHealth),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	recorded := test.RunRequest(t, api.MakeHandler(),
		test.MakeSimpleRequest("GET", "http://1.2.3.4/health", nil))
	recorded.CodeIs(200)
	recorded.BodyIs("{\n  \"status\": \"ok\"\n}")
}
