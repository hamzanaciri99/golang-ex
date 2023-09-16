package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/hamzanaciri99/golang-ex/util"
)

type Chapter struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func (c Chapter) String() string {
	return c.Title
}

func ParseChaptersFromJSON(data []byte) map[string]Chapter {
	chapters := []Chapter{}
	chapterMap := make(map[string]Chapter, 0)
	err := json.Unmarshal(data, &chapters)
	util.CheckError(err)

	for _, chapter := range chapters {
		chapterMap[chapter.ID] = chapter
	}

	return chapterMap
}

// func NextChapters(chapters map[string]Chapter, current Chapter) ([]Chapter, error) {
// 	next := []Chapter{}

// 	for _, option := range current.Options {
// 		nextChapter, ok := chapters[option.Arc]
// 		if !ok {
// 			return nil, fmt.Errorf("arc with id %s not found", option.Arc)
// 		}
// 		next = append(next, nextChapter)
// 	}
// 	return next, nil
// }

const FIRST_CHAPTER = "intro"

type ChapterHandler struct {
	chapter Chapter
}

func (h ChapterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Title   string
		Content string
		Options string
	}{
		h.chapter.Title,
		strings.Join(h.chapter.Story, " "),
		"<div class=\"end\">THE END</div>",
	}

	if len(h.chapter.Options) > 0 {
		params.Options = ""
		for _, option := range h.chapter.Options {
			params.Options += fmt.Sprintf(`
				<a class="option" href="%s">%s</a>
			`, option.Arc, option.Text)
		}
	}

	html, err := os.ReadFile("story/template.html")
	util.CheckError(err)

	t, _ := template.New("t").Parse(string(html))
	t.Execute(w, params)
}

func Handle(chapter Chapter) http.Handler {
	return ChapterHandler{chapter}
}

var (
	runInCLI bool
)

func flagsInit() {
	flag.BoolVar(&runInCLI, "run_in_cli", false,
		`if set to true, a CLI version of the program will be run,
		otherwise it will be run on the browser!`)
}

func init() {
	flagsInit()
	flag.Parse()
}

func runServer(chapters map[string]Chapter) {
	mux := http.NewServeMux()
	mux.Handle("/", Handle(chapters[FIRST_CHAPTER]))

	for _, chapter := range chapters {
		mux.Handle("/"+chapter.ID, Handle(chapter))
	}

	fmt.Printf("listening to localhost:8080")
	http.ListenAndServe("localhost:8080", mux)
}

func runCLI(chapters map[string]Chapter) {
	current := chapters[FIRST_CHAPTER]

	for {
		fmt.Printf("Title: %s\nStory: %s\n------------------------\n",
			current.Title, strings.Join(current.Story, " "))

		if len(current.Options) == 0 {
			break
		}

		for i, option := range current.Options {
			fmt.Printf("%d) %s\n", i, option.Text)
		}

		var index int
		_, err := fmt.Scanf("%d", &index)
		util.CheckError(err)

		if index < 0 || index > len(current.Options) {
			fmt.Println("Wrong choice, please choose again!")
			continue
		}

		current = chapters[current.Options[index].Arc]
	}
}

func main() {
	json, err := os.ReadFile("story/story.json")
	util.CheckError(err)

	chapters := ParseChaptersFromJSON(json)

	if runInCLI {
		runCLI(chapters)
	} else {
		runServer(chapters)
	}

}
