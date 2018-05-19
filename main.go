package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	ctf "github.com/ohookins/contentful-go"
)

func main() {
	filename := flag.String("filename", "", "Filename of XML export of Wordpress blog")
	space := flag.String("space", "", "Space ID on Contentful")
	token := flag.String("token", "", "Personal CMA Token")
	flag.Parse()

	if *filename == "" {
		fmt.Println("Please supply a filename input")
		return
	}

	if *space == "" {
		fmt.Println("Please supply a Space ID")
		return
	}

	if *token == "" {
		fmt.Println("Please supply a CMA Token")
		return
	}
	cma := ctf.NewCMA(*token)

	file, err := ioutil.ReadFile(*filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := parseDoc(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = deleteCategories(cma, *space)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	err = createCategories(cma, body.Channel.WPCategory, *space)
	if err != nil {
		fmt.Println(err)
	}

	err = deleteTags(cma, *space)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	err = createTags(cma, body.Channel.WPTag, *space)
	if err != nil {
		fmt.Println(err)
	}
}
