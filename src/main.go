package main

func main() {
	app := App{}
	app.Initialize(DbUser, DbPass, DbHost, DbName)
	app.Run("localhost:8080")
}
