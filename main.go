package main

import (
	"context"
	"fmt"
	"log"
    "time"
	"net/http"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

type Note struct {
	Title string             `json:"title,omitempty" bson:"title,omitempty"`
	Note  string             `json:"note,omitempty" bson:"note,omitempty"`
}

var router *gin.Engine

func main() {

    router = gin.Default()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(
       "mongodb+srv://topher:topher@cluster0.ygkqf.mongodb.net/NotesGo?retryWrites=true&w=majority",
    ))
    if err != nil { log.Fatal(err) }

    err = client.Ping(ctx, readpref.Primary())
    if err != nil {
        log.Fatal(err)
    }

	var quickstartDatabase = client.Database("NotesGo")
	var collection = quickstartDatabase.Collection("GoNotes")

    databases, err := client.ListDatabaseNames(ctx, bson.M{})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(databases)

	// insert new note route and handler
    router.POST("/api/new", func (c *gin.Context) {
		var tempTitle = c.PostForm("title")
		var tempNoteBody = c.PostForm("note")
        var tempNote = Note{Title:tempTitle,Note:tempNoteBody}
        insertResult, err := collection.InsertOne(context.TODO(), tempNote) //inserting
        if err != nil {
            log.Fatal(err)
            c.JSON(http.StatusOK)
        }else{
			c.JSON(http.StatusOK)
		}
        fmt.Println("Inserted a single document: ", insertResult)
    })

	// finding a note by title route and handler
	router.GET("/api/find/:title", func (c *gin.Context) {
		var tempTitle = c.Request.URL.Query()["title"][0]
		fmt.Println(tempTitle)
		filter := bson.D{{"title", tempTitle}}
		var result Note

	    err = collection.FindOne(context.TODO(), filter).Decode(&result) //finding
	    if err != nil {
	        log.Fatal(err)
	    }
		fmt.Println("found a single document: ", result)
		c.JSON(http.StatusOK, gin.H{"title":result.Title,"note": result.Note})

    })

	// finding all notes route and handler
	router.POST("/api/all", func (c *gin.Context) {
		// Pass these options to the Find method
		findOptions := options.Find()
		// Here's an array in which you can store the decoded documents
		var results []*Note

		// Passing bson.D{{}} as the filter matches all documents in the collection
		cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
		if err != nil {
		    log.Fatal(err)
		}

		// Finding multiple documents returns a cursor
		// Iterating through the cursor allows us to decode documents one at a time
		for cur.Next(context.TODO()) {

		    // create a value into which the single document can be decoded
		    var elem Note
		    err := cur.Decode(&elem)
		    if err != nil {
		        log.Fatal(err)
		    }

		    results = append(results, &elem)
		}

		if err := cur.Err(); err != nil {
		    log.Fatal(err)
		}

		// Close the cursor once finished
		cur.Close(context.TODO())

		fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
		c.JSON(http.StatusOK, results)


	})

	// deletes a note by title route and handler
	router.POST("api/delete", func (c *gin.Context){
		// delete
		var tempTitle = c.PostForm("title")
		var tempNote = Note{Title:tempTitle}
	    deleteResult, err := collection.DeleteMany(context.TODO(), tempNote)
	    if err != nil {
	        log.Fatal(err)
			c.JSON(http.StatusOK, gin.H{"err": true})
	    }else{
			c.JSON(http.StatusOK, gin.H{"err": false})
		}
	    fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	})

	router.Run()
}
