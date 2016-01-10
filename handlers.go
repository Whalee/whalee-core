package main

import (
	"fmt"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func ForceDeploy(githubRepo string) {
	//Force the dheployment inside the dockers
}
