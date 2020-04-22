package handlers

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/gin-gonic/gin"
	"github.com/jogit82/google-oauth-gin/structs"
	"github.com/nishanths/go-xkcd/v2"
)

// MovieSearchFormHandler handles form loading
// func MovieSearchFormHandler(c *gin.Context) {
// 	c.HTML(http.StatusOK, "searchmovies", nil)
// }

// SearchMoviesHandler handles the movie/shows/people search using the movie DB called themoviedb.org which has a Go specific library!
func SearchMoviesHandler(c *gin.Context) {
	log.Println("Hello SearchMoviesHandler")
	tmdbClient, err := tmdb.Init("2d46d415e612dfe8f47fdf1d121bd545")

	if err != nil {
		fmt.Println(err)
	}

	options := make(map[string]string)
	// options["language"] = "pt-BR"

	// Multi Search
	// https://developers.themoviedb.org/3/search/multi-search
	search, err := tmdbClient.GetSearchMulti(c.PostForm("query"), options)

	if err != nil {
		log.Fatal(err)
	}
	var movie structs.Movie
	var show structs.Show
	var person structs.Person
	movies := make([]structs.Movie, 5)
	shows := make([]structs.Show, 5)
	people := make([]structs.Person, 5)

	// Iterate
	for _, v := range search.Results {
		if v.MediaType == "movie" {
			// fmt.Println("Movie Title: ", v.Title)
			movie.Name = v.Title
			movie.Overview = v.Overview
			movie.Popularity = v.Popularity
			movie.VoteCount = v.VoteCount
			movies = append(movies, movie)
		} else if v.MediaType == "tv" {
			show.Name = v.Name
			show.Overview = v.Overview
			show.Popularity = v.Popularity
			show.VoteCount = v.VoteCount
			shows = append(shows, show)
		} else if v.MediaType == "person" {
			person.Name = v.Name
			person.Popularity = v.Popularity
			people = append(people, person)
		}
	}

	c.HTML(http.StatusOK, "movieresults", gin.H{
		"query":  c.PostForm("query"),
		"movies": movies,
		"shows":  shows,
		"people": people,
	})
}

// SearchHandler handles searches by any user
func SearchHandler(c *gin.Context) {
	client := xkcd.NewClient()

	comic, err := client.Get(context.Background(), rand.Intn(2300))
	// comic, err := client.Latest(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("search query >>>>>")
	// log.Println(c.PostForm("query"))
	c.HTML(http.StatusOK, "searchresults", gin.H{
		"query":      c.PostForm("query"),
		"name":       comic.Title,
		"image":      comic.ImageURL,
		"number":     comic.Number,
		"year":       comic.Year,
		"transcript": comic.Transcript,
	})
	// c.JSON(http.StatusOK, gin.H{
	// 	"query": c.PostForm("query"),
	// })
}

// RandomHandler handles random pick
func RandomHandler(c *gin.Context) {
	client := xkcd.NewClient()

	comic, err := client.Get(context.Background(), rand.Intn(2300))
	// comic, err := client.Latest(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("search query >>>>>")
	// log.Println(c.PostForm("query"))
	c.HTML(http.StatusOK, "searchresults", gin.H{
		"query":      c.PostForm("query"),
		"title":      comic.Title,
		"image":      comic.ImageURL,
		"number":     comic.Number,
		"year":       comic.Year,
		"transcript": comic.Transcript,
	})
	// c.JSON(http.StatusOK, gin.H{
	// 	"query": c.PostForm("query"),
	// })
}
