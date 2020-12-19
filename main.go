package main

func main() {
	app := App{}
	app.Initialize("root", "")
	app.Run(":5005")
}