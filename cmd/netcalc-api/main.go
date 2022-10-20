package main

func main() {

	config := &Configuration{
		HttpPort: 3000,
	}
	app := NewApp(config)
	app.Run()
}
