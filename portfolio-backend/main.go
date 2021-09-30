package main

import (
  "log"
  "net/http"
  "context"
  "time"

  "github.com/gin-gonic/gin"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "go.mongodb.org/mongo-driver/mongo/readpref"
)

type message struct {
  Message string `bson:"message" json:"message"`
}

func main() {
  log.Println("Connecting to MongoDB")
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo-0.mongo:27017,mongo-1.mongo:27017,mongo-2.mongo:27017/?replicaSet=rs0"))
  if err != nil {
    log.Fatalf("Could not connect to database: %s", err)
  }

  err = client.Connect(ctx)
  if err != nil {
      log.Fatalf("Could not connect to database: %s", err)
  }

  defer func() {
    if err := client.Disconnect(ctx); err != nil {
      log.Fatalf("Could not disconnect from database: %s", err)
    } 
  }()

  router := gin.Default()
  router.GET("/api/health", getHealthHandlerFunction(client))
  router.GET("/api/message", getGetMessagesHandlerFunction(client))
  router.POST("/api/message", getPostMessageHandlerFunction(client))
  router.PUT("/api/message", getPutMessageHandlerFunction(client))

  router.Run(":8080")
}

func getHealthHandlerFunction(client *mongo.Client) func(*gin.Context) {
  return func (c* gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    err := client.Ping(ctx, readpref.Primary())
    if err != nil {
      log.Printf("Could not ping database: %s", err)
      c.AbortWithStatus(http.StatusInternalServerError)
    }
  }
}

func getPostMessageHandlerFunction(client *mongo.Client) func(*gin.Context) {
  return func (c *gin.Context) {
    var m map[string]interface{}
    if err := c.BindJSON(&m); err != nil {
      c.JSON(http.StatusInternalServerError, message{Message: "Could not parse the body of the request. Check that the body contains valid JSON"})
    }
    log.Printf("Successfully unmarshalled post body. Writing message '%s' to the database", m)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := client.Database("portfolio").Collection("messages")
    _, insertErr := collection.InsertOne(ctx, m) // TODO: Should probably have some kind of validation on the stuff we put in the database
    if insertErr != nil {
      errorMessage := "Could not insert the posted document"
      log.Printf("%s: %s", errorMessage, insertErr)
      c.JSON(http.StatusInternalServerError, message{Message: errorMessage})
    }

    log.Printf("Wrote message '%s' to the database", m)
  }
}

func getPutMessageHandlerFunction(client *mongo.Client) func(*gin.Context) {
  return func (c *gin.Context) {
    var m map[string]interface{}
    if err := c.BindJSON(&m); err != nil {
      c.JSON(http.StatusInternalServerError, message{Message: "Could not parse the body of the request. Check that the body contains valid JSON"})
    }
    log.Printf("Successfully unmarshalled put body. Upserting message '%s' to the database", m)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // valuesToUpdate := map[string]
    // for key,value := range m {
      // valuesToUpdate[key] = value
      // // valuesToUpdate = append(valuesToUpdate, bson.E{key, value})
    // }
    // if err != nil {
      // errorMessage := "Could not update the document" // TODO: Improve error message to make debugging easier
      // log.Printf("%s: %s", errorMessage, err)
      // c.JSON(http.StatusInternalServerError, message{Message: errorMessage})
    // }

    collection := client.Database("portfolio").Collection("messages")
    opts := options.Update().SetUpsert(true)
    filter := bson.M{"email": m["email"]}
    update := bson.M{"$set": m}
    _, err := collection.UpdateOne(ctx, filter, update, opts) // TODO: Should probably have some kind of validation on the stuff we put in the database
    if err != nil {
      errorMessage := "Could not update the document"
      log.Printf("%s: %s", errorMessage, err)
      c.JSON(http.StatusInternalServerError, message{Message: errorMessage})
    }

    log.Printf("Wrote message '%s' to the database", m)
  }
}

func getGetMessagesHandlerFunction(client *mongo.Client) func(*gin.Context) {
  return func(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    queryParams := c.Request.URL.Query()
    // TODO: Add query parameter validation
    filter := bson.M{"email": queryParams["email"][0]}
    log.Printf("%s", filter)

    collection := client.Database("portfolio").Collection("messages")
    cur, err := collection.Find(ctx, filter)
    if err != nil {
      errorMessage := "Could not execute search request against database"
      log.Printf("%s: %s", errorMessage, err)
      c.JSON(http.StatusInternalServerError, errorMessage)
    }

    defer cur.Close(ctx)

    var messages []*map[string]interface{}
    for cur.Next(ctx) {
      var m map[string]interface{}
      err := cur.Decode(&m)
      if err != nil { log.Fatal(err) }
      messages = append(messages, &m)
    }

    c.JSON(http.StatusOK, messages)
  }
}
