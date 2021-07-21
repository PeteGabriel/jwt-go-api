package main

import "jwtGoApi/internal/api"

func main(){
	app := api.New()
	app.Start()
}