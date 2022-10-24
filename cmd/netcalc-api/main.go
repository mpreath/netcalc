package main

func main() {

	config := LoadConfig()
	app := NewApp(config)
	app.Run()
}
