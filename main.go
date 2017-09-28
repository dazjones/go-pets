package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
)

type Impl struct {
	DB *gorm.DB
}

func (i *Impl) InitDB() {
	var err error
	i.DB, err = gorm.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}
	i.DB.LogMode(true)
}

func (i *Impl) InitSchema() {
	i.DB.CreateTable(&Pet{})
	i.DB.AutoMigrate(&Pet{})
}

type Pet struct {
  ID uint `gorm:"primary_key" json:"id"`
  Type string `gorm:"not null" json:"type"`
  Name string `gotm:"not null" json:"name"`
}


func main() {
	i := Impl{}
	i.InitDB()
	i.InitSchema()
	defer i.DB.Close()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/pets", i.GetAllPets),
		rest.Post("/pets", i.PostPet),
		rest.Get("/", i.Home),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))

}

func (i *Impl) Home(w rest.ResponseWriter, r *rest.Request) {
    w.WriteJson("go-pets is up and running!")
}

func (i *Impl) GetAllPets(w rest.ResponseWriter, r *rest.Request) {
	pets := []Pet{}
	i.DB.Find(&pets)
	w.WriteJson(&pets)
}

func (i *Impl) PostPet(w rest.ResponseWriter, r *rest.Request) {
	pet := Pet{}
	if err := r.DecodeJsonPayload(&pet); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := i.DB.Save(&pet).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&pet)
}

func (i *Impl) Status(w rest.ResponseWriter, r *rest.Request) {
    w.WriteJson(map[string]string{"status": "OK"})
}
