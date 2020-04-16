package main

import "github.com/zserge/webview"

func main() {
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Beehive")
	w.SetSize(1024, 768, webview.HintNone)
	w.Navigate("http://localhost:8181")
	w.Run()
}
