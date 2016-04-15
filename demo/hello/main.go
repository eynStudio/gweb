package main

import (
	"log"

	"github.com/eynstudio/gweb"
)

func main() {
	log.Println("Hello Start...")
	app := gweb.NewAppWithCfg(&gweb.Cfg{Port: 80})

	app.Root.AddNode(NewHome())
	app.Root.AddNode(NewApi())
	app.Start()

}
