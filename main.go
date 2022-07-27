package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	// set header
	w.Header().Set("Content-Type", "application/json")
	// return all movies in json format
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// get all variables from mux in params
	params := mux.Vars(r)
	// Iterate over the range of movies using index, item
	for index, item := range movies {
		// when id match
		if item.ID == params["id"] {
			// movie at index will be replaced by all movies after index i.e. index+1
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	// return remaining movies
	json.NewEncoder(w).Encode(movies)
}

func getMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// get all variables from mux in params
	params := mux.Vars(r)
	//loop through all movies
	for _, movie := range movies {
		if movie.ID == params["id"] {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// create a new variable
	var movie Movie
	// decode movie data in json body to movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	// create a random number for movieId
	movie.ID = strconv.Itoa(rand.Intn(10000000))
	movies = append(movies, movie)

	// return created movie
	json.NewEncoder(w).Encode(movie)
}

// Logic : first delete from array & craete a new one with updated data
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, movie := range movies {
		if movie.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			// return updated movies
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	movies = append(
		movies, Movie{
			ID:       "1",
			Isbn:     "438227",
			Title:    "Matrix",
			Director: &Director{Firstname: "Moli", Lastname: "Maralina"},
		},
	)
	movies = append(
		movies, Movie{
			ID:       "2",
			Isbn:     "438347",
			Title:    "Edge of tomorrow",
			Director: &Director{Firstname: "Marati", Lastname: "Muina"},
		},
	)
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovieById).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// start the server
	fmt.Print("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
