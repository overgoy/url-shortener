package main

import (
	app "github.com/overgoy/url-shortener/internal/app"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	http.HandleFunc("/", app.HandleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
	//log := logrus.New()
	//log.SetOutput(os.Stdout)
	//log.SetLevel(logrus.InfoLevel)
	//
	////controller = controller.NewBaseController(*log)
	//r := chi.NewRouter()
	////r.Mount("/", controller.Route())
	//
	//log.Info("Server started")
	//err := http.ListenAndServe(":8080", r)
	//
	//if err != nil {
	//	fmt.Printf("server: #{err}")
	//}

}
