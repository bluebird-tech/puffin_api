package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/yvasiyarov/go-metrics"
	"github.com/yvasiyarov/gorelic"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	i := Impl{}
	i.InitDB()
	i.InitSchema()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&NewRelicMiddleware{
		License: os.Getenv("NEWRELIC_LICENCE_KEY"),
		Name:    "Puffin Event API",
		Verbose: true,
	})
	router, err := rest.MakeRouter(
		rest.Post("/events", i.PostEvent),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), api.MakeHandler()))
}

type NewRelicMiddleware struct {
	License string
	Name    string
	Verbose bool
	agent   *gorelic.Agent
}

func (mw *NewRelicMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {

	mw.agent = gorelic.NewAgent()
	mw.agent.NewrelicLicense = mw.License
	mw.agent.HTTPTimer = metrics.NewTimer()
	mw.agent.Verbose = mw.Verbose
	mw.agent.NewrelicName = mw.Name
	mw.agent.CollectHTTPStat = true
	mw.agent.Run()

	return func(writer rest.ResponseWriter, request *rest.Request) {

		handler(writer, request)

		// the timer middleware keeps track of the time
		startTime := request.Env["START_TIME"].(*time.Time)
		mw.agent.HTTPTimer.UpdateSince(*startTime)
	}
}

type Event struct {
	Id          int64     `json:"id"`
	Measurement string    `sql:"size:1024" json:"measurement"`
	Fields      string    `sql:"size:1024" json:"fields"`
	Tags        string    `sql:"size:1024" json:"tags"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   time.Time `json:"-"`
}

type Impl struct {
	DB gorm.DB
}

func (i *Impl) InitDB() {
	var err error
	connection := os.Getenv("DATABASE_URL")
	i.DB, err = gorm.Open("postgres", connection)
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}
	i.DB.LogMode(true)
}

func (i *Impl) InitSchema() {
	i.DB.AutoMigrate(&Event{})
}

func (i *Impl) PostEvent(w rest.ResponseWriter, r *rest.Request) {
	event := Event{}
	if err := r.DecodeJsonPayload(&event); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := i.DB.Save(&event).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&event)
}

func helloHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello!"))
	})
}
