package main

import "fmt"

type Post struct {
	Id int
	Content string
	Author string
}

var PostById map[int]*Post
var PostsByAuthor map[string][]*Post

func store(post Post) {
	PostById[post.Id] = &post
	PostsByAuthor[post.Author] = append(PostsByAuthor[post.Author], &post)
}

func main() {
	PostById = make(map[int]*Post, 0)
	PostsByAuthor = make(map[string][]*Post)
	post1 := Post{Id: 1, Content: "Hello World!", Author: "Sau Sheong"}
	post2 := Post{Id: 2, Content: "Bonjour Monde!", Author: "Pierre"}
	post3 := Post{Id: 3, Content: "Hola Mundo!", Author: "Pedro"}
	store(post1)
	store(post2)
	store(post3)
	fmt.Println(PostById[1])
}