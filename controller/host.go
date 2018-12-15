package controller

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"

	"../model"
)

type Host struct {
	temp *template.Template
}

func (h Host) RegisterRoutes() {
	http.HandleFunc("/host", h.servepage)
}

func (h Host) servepage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t := h.temp.Lookup("host.html")
		if t != nil {
			err := t.Execute(w, nil)
			if err != nil {
				log.Println(err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()

		if err != nil {
			log.Println(err)
		}

		f := r.Form

		data := model.HostType{
			Port1:    f.Get("port1"),
			Port2:    f.Get("port2"),
			Image1:   f.Get("image1"),
			Image2:   f.Get("image2"),
			Email:    f.Get("email"),
			Password: f.Get("password"),
		}

		//log.Println(data)

		// executing terminal logic
		go func() {
			cmd := exec.Command("session", data.Port1, data.Image1)
			cmd.Stdout = os.Stdout
			if err := cmd.Run(); err != nil {
				log.Fatalln(err)
			} else {
				log.Printf("Running %s on port %d", data.Image1, data.Port1)
			}
		}()

		go func() {
			cmd := exec.Command("session", data.Port2, data.Image2)
			cmd.Stdout = os.Stdout
			if err := cmd.Run(); err != nil {
				log.Fatalln(err)
			} else {
				log.Printf("Running %s on port %d", data.Image2, data.Port2)
			}
		}()

		t := h.temp.Lookup("session.html")
		if t != nil {
			err := t.Execute(w, data)
			if err != nil {
				log.Println(err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	}
}
