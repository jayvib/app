package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

type Post struct {
	XMLName xml.Name `xml:"post"`
	Id string `xml:"id,attr"`
	Content string `xml:"content"`
	Author Author `xml:"author"`
	Xml string `xml:",innerxml"`
	Comments []Comment `xml:"comments>comment"`
}

type Author struct {
	Id string `xml:"id,attr"`
	Name string `xml:",chardata"`
}

type Comment struct {
	Id string `xml:"id,attr"`
	Content string `xml:"content"`
	Author Author `xml:"author"`
}

func main() {
	file, err := os.Open("post.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var post Post
	err = xml.NewDecoder(file).Decode(&post)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", post)

}
