package db

import (
    "context"
    "log"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "github.com/mercimat/wikiapp/web"
)

type MongoDB struct {
    Collection  *mongo.Collection
    Ctx         context.Context
}

func NewMongoDB(server string, database string, collection string) *MongoDB {
    clientOptions := options.Client().ApplyURI(server)
    ctx := context.TODO()

    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    return &MongoDB{
        Collection: client.Database(database).Collection(collection),
        Ctx: ctx,
    }
}

func (m MongoDB) SavePage(p *web.Page) error {
    var page web.Page
    filter := bson.D{{"title", p.Title}}
    update := bson.D{{"$set", bson.D{{"body", p.Body}}}}
    opts := options.FindOneAndUpdate().SetUpsert(true)
    err := m.Collection.FindOneAndUpdate(m.Ctx, filter, update, opts).Decode(&page)
    if err != nil && err != mongo.ErrNoDocuments {
        return err
    }
    return nil
}

func (m MongoDB) GetPage(title string) (*web.Page, error) {
    var p web.Page
    filter := bson.D{{"title", title}}
    err := m.Collection.FindOne(m.Ctx, filter).Decode(&p)
    return &p, err
}

