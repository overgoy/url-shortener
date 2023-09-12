package server

import (
	"github.com/overgoy/url-shortener/internal/handlers"
	"log"
	"net/http"
)

func Start() {
	http.HandleFunc("/", handlers.HandleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

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
