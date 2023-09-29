package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

type Tutorial struct {
	ID       int
	Title    string
	Author   Author
	Comments []Comment
}

type Author struct {
	Name      string
	Tutorials []int
}

type Comment struct {
	Body string
}

func populate() []Tutorial {
	author := &Author{Name: "Dandan", Tutorials: []int{1}}
	tutorial := Tutorial{
		ID:     1,
		Title:  "Little Gophers tutorial",
		Author: *author,
		Comments: []Comment{
			{Body: "Headbutt!"},
			{Body: "Pokemon"},
		},
	}

	var tutorials []Tutorial

	tutorials = append(tutorials, tutorial)

	return tutorials
}

func main() {
	var commentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Comment",
			Fields: graphql.Fields{
				"Body": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	var authorType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Author",
			Fields: graphql.Fields{
				"Name": &graphql.Field{
					Type: graphql.String,
				},
				"Tutorials": &graphql.Field{
					// Declare a new list/array then the type of the list
					Type: graphql.NewList(graphql.Int),
				},
			},
		},
	)

	var tutorialType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Tutorial",
			Fields: graphql.Fields{
				"ID": &graphql.Field{
					Type: graphql.Int,
				},
				"Title": &graphql.Field{
					Type: graphql.String,
				},
				"Author": &graphql.Field{
					// Can reference other types created
					Type: authorType,
				},
				"Comments": &graphql.Field{
					// Can reference other types created
					Type: graphql.NewList(commentType),
				},
			},
		},
	)

	tutorials := populate()

	// Schema
	fields := graphql.Fields{
		"Tutorial": &graphql.Field{
			Type:        tutorialType,
			Description: "Get tutorial by ID",
			Args: graphql.FieldConfigArgument{
				"ID": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["ID"].(int)
				if ok {
					for _, tutorial := range tutorials {
						if int(tutorial.ID) == id {
							return tutorial, nil
						}
					}
				}
				return nil, nil
			},
		},
		"List": &graphql.Field{
			Type:        graphql.NewList(tutorialType),
			Description: "Get tutorial list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return tutorials, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	query := `
		{
			List {
				ID
				Title
				Comments {
					Body
				}
				Author {
					Name
					Tutorials
				}
			}
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON) // {"data":{"hello":"world"}}
}
