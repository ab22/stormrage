package main

import "log"

func main() {
	log.Println("Starting server...")

	s, err := NewServer()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Listening...")

	err = s.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
