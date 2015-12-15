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
	//api.Use(&rest.AccessLogJsonMiddleware{})
	router, err := rest.MakeRouter(
		rest.Post("/events", i.PostEvent),
		rest.Post("/installations", i.PostInstallation),
		rest.Get("/health", i.GetHealth),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+os.Getenv("PORT"), api.MakeHandler()))
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
	Id                          int64      `json:"id"`
	DeviceId                    string     `sql:"size:40" json:"deviceId"`
	Identifier                  string     `sql:"size:6" json:"identifier"`
	AppVersion                  string     `sql:"size:20" json:"AppVersion"`
	MeasuredAt                  time.Time  `json:"measuredAt"`
	Measurement                 string     `sql:"size:1024" json:"measurement"`
	Value                       float64    `json:"value"`
	Fields                      string     `sql:"size:1024" json:"fields"`
	Tags                        string     `sql:"size:1024" json:"tags"`
	Note                        string     `sql:"size:1024" json:"note"`
	Tenant                      string     `sql:"size:30" json:"tenant"`
	CreatedAt                   time.Time  `json:"createdAt"`
	UpdatedAt                   time.Time  `json:"updatedAt"`
	DeletedAt                   time.Time  `json:"-"`
	NotificationIntervalMinutes int64      `json:"notificationIntervalMinutes"`
	NotificationShownAt         *time.Time `json:"notificationShownAt"`
	NotificationDismissedAt     *time.Time `json:"notificationDismissedAt"`
	NotificationAcceptedAt      *time.Time `json:"notificationAcceptedAt"`
	StartAt                     *time.Time `json:"startAt"`
	EndAt                       *time.Time `json:"EndAt"`
}

type Installation struct {
	Id                  int64     `json:"id"`
	DeviceType          string    `sql:"size:40" json:"deviceType"`
	DeviceManufacturer  string    `sql:"size:40" json:"deviceManufacturer"`
	DeviceModel         string    `sql:"size:40" json:"deviceModel"`
	DeviceProduct       string    `sql:"size:40" json:"deviceProduct"`
	DeviceSdk           string    `sql:"size:40" json:"deviceSdk"`
	DeviceSystemLocale  string    `sql:"size:40" json:"deviceSystemLocale"`
	DeviceAdvertisingId string    `sql:"size:40" json:"deviceAdverstisingId"`
	DeviceMobileCarrier string    `sql:"size:40" json:"deviceMobileCarrier"`
	AppVersion          string    `sql:"size:40" json:"appVersion"`
	AppVersionCode      string    `sql:"size:40" json:"appVersionCode"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
	DeletedAt           time.Time `json:"-"`
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

	if err = i.DB.DB().Ping(); err != nil {
		log.Fatalf("Unable to verify connection to database: '%v'", err)
	}

	i.DB.LogMode(true)
}

func (i *Impl) InitSchema() {
	i.DB.AutoMigrate(&Event{})
	i.DB.AutoMigrate(&Installation{})
}

func (i *Impl) PostEvent(w rest.ResponseWriter, r *rest.Request) {
	event := Event{}
	if err := r.DecodeJsonPayload(&event); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := i.DB.Create(&event).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.WriteJson(&event)
}

func (i *Impl) PostInstallation(w rest.ResponseWriter, r *rest.Request) {
	installation := Installation{}
	if err := r.DecodeJsonPayload(&installation); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := i.DB.Save(&installation).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.WriteJson(&installation)
}

func (i *Impl) GetHealth(w rest.ResponseWriter, r *rest.Request) {
	if err := i.DB.DB().Ping(); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(map[string]string{"status": "ok"})
}
