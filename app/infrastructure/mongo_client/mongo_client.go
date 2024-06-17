package mongo_client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"

	"cnores-skeleton-golang-app/app/infrastructure/constant"
	utils_context "cnores-skeleton-golang-app/app/shared/utils/context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DatabaseInterface[R any] interface {
	Collection(name string) CollectionInterface[R]
	Client() MongoClientInterface
}

type CollectionInterface[R any] interface {
	FindOne(interface{}, interface{}) (R, error)
	Find(interface{}, interface{}, ...*options.FindOptions) ([]R, error)
	CountDocuments(ctx interface{}, filter interface{}) (int64, error)
	InsertOne(ctx interface{}, param interface{}) (string, error)
	UpdateOne(ctx interface{}, param interface{}, update interface{}) (int, error)
	DeleteOne(ctx interface{}, filter interface{}) (int64, error)
	UpsertOne(ctx, filter, update interface{}) (string, int64, int64, error)
	UpdateMany(ctx interface{}, filter interface{}, update interface{}) (int, int, error)
}

type SingleResultInterface interface {
	Decode(v interface{}) error
}

type MongoClientInterface interface {
	Disconnect() error
	Connect() error
	StartSession() (mongo.Session, error)
	UseSession(ctx context.Context, fn func(sessCtx mongo.SessionContext) error) error
}

type mongoClient struct {
	cl *mongo.Client
}

type mongoDatabase[R any] struct {
	db *mongo.Database
}

type mongoCollection[R any] struct {
	coll *mongo.Collection
}

func NewClient(uri string, isTSL bool) (MongoClientInterface, error) {
	if isTSL {
		caFilePath := "/app/global-bundle.cer"
		tlsConfig, err := getCustomTLSConfig(caFilePath)
		if err != nil {
			return nil, err
		}
		client, err := mongo.NewClient(options.Client().ApplyURI(uri).SetTLSConfig(tlsConfig))
		return &mongoClient{cl: client}, err
	} else {
		client, err := mongo.NewClient(options.Client().ApplyURI(uri))
		return &mongoClient{cl: client}, err
	}
}

func getCustomTLSConfig(caFile string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)
	certs, err := ioutil.ReadFile(caFile)

	if err != nil {
		return tlsConfig, err
	}

	tlsConfig.RootCAs = x509.NewCertPool()
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

	if !ok {
		return tlsConfig, errors.New("Failed parsing pem file")
	}

	return tlsConfig, nil
}

func (mc *mongoClient) UseSession(ctx context.Context, fn func(sessCtx mongo.SessionContext) error) error {
	return mc.cl.UseSession(ctx, fn)
}

func database[R any](mc *mongoClient, dbName string) (DatabaseInterface[R], error) {
	db := mc.cl.Database(dbName)
	ctx := context.Background()
	log := utils_context.GetLogFromContext(ctx, constant.InfrastructureLayer, "mongo_client.Database")

	log.Info(fmt.Sprintf("checking database .... %s", dbName))
	e := mc.cl.Ping(context.Background(), readpref.Primary())
	log.Info(fmt.Sprintf("checking database .... %v", e))
	if e != nil {
		return nil, e
	}
	log.Info("mongo database connected successfully")
	return &mongoDatabase[R]{db: db}, nil
}

func (mc *mongoClient) StartSession() (mongo.Session, error) {
	return mc.cl.StartSession()
}

func (mc *mongoClient) Connect() error {
	// mongo client does not use context on connect method. There is a ticket
	// with a request to deprecate this functionality and another one with
	// explanation why it could be useful in synchronous requests.
	// https://jira.mongodb.org/browse/GODRIVER-1031
	// https://jira.mongodb.org/browse/GODRIVER-979
	return mc.cl.Connect(nil)
}

func (mc *mongoClient) Disconnect() error {
	return mc.cl.Disconnect(context.Background())
}

func NewDatabase[R any](dbName string, client MongoClientInterface) DatabaseInterface[R] {
	db, _ := database[R](client.(*mongoClient), dbName)
	return db
}

func (md *mongoDatabase[R]) Collection(colName string) CollectionInterface[R] {
	collection := md.db.Collection(colName)
	return &mongoCollection[R]{coll: collection}
}

func (md *mongoDatabase[R]) Client() MongoClientInterface {
	client := md.db.Client()
	return &mongoClient{cl: client}
}

func (mc *mongoCollection[R]) FindOne(ctx interface{}, filter interface{}) (R, error) {
	var singleResult *mongo.SingleResult
	var findContext context.Context
	if reflect.TypeOf(ctx).String() == "*context.emptyCtx" || reflect.TypeOf(ctx).String() == "*context.valueCtx" {
		findContext = (ctx).(context.Context)
	} else {
		findContext = (ctx).(mongo.SessionContext)
	}
	log := utils_context.GetLogFromContext(findContext, constant.InfrastructureLayer, "database_client")
	singleResult = mc.coll.FindOne(findContext, filter)
	var res R
	e := singleResult.Decode(&res)
	if e != nil {
		log.Error(fmt.Sprintf("Error on find one with message %s", e.Error()))
		return res, e
	}
	log.Info(fmt.Sprintf("Found one result sucessfully"))
	return res, nil
}

func (mc *mongoCollection[R]) CountDocuments(ctx interface{}, filter interface{}) (int64, error) {

	var findContext context.Context
	if reflect.TypeOf(ctx).String() == "*context.emptyCtx" || reflect.TypeOf(ctx).String() == "*context.valueCtx" {
		findContext = (ctx).(context.Context)
	} else {
		findContext = (ctx).(mongo.SessionContext)
	}
	return mc.coll.CountDocuments(findContext, filter)

}

func (mc *mongoCollection[R]) Find(ctx interface{}, filter interface{}, findOptions ...*options.FindOptions) ([]R, error) {

	var err error
	var resList []R
	var result *mongo.Cursor
	if reflect.TypeOf(ctx).String() == "*context.emptyCtx" || reflect.TypeOf(ctx).String() == "*context.valueCtx" {
		sessionContext := (ctx).(context.Context)
		if len(findOptions) > 0 {
			result, err = mc.coll.Find(sessionContext, filter, findOptions[0])
		} else {
			result, err = mc.coll.Find(sessionContext, filter)
		}
		if err != nil {
			return nil, err
		}

		for result.Next(sessionContext) {
			var singleResult R
			singleErr := result.Decode(&singleResult)
			if singleErr != nil {
				fmt.Println("cursor.Next() error:", err)
				return nil, singleErr

			} else {
				resList = append(resList, singleResult)
			}
		}
	} else {
		sessionContext := (ctx).(mongo.SessionContext)
		result, err = mc.coll.Find(sessionContext, filter)
		var singleResult R
		result, err = mc.coll.Find(sessionContext, filter)
		for result.Next(sessionContext) {
			singleErr := result.Decode(&singleResult)
			if singleErr != nil {
				fmt.Println("cursor.Next() error:", err)
				return nil, singleErr

			} else {
				resList = append(resList, singleResult)
			}
		}
	}

	if err != nil {
		return resList, err
	}

	return resList, err
}

func (mc *mongoCollection[R]) UpdateOne(ctx interface{}, param interface{}, update interface{}) (int, error) {
	var id *mongo.UpdateResult
	var err error
	if reflect.TypeOf(ctx).String() == "*context.emptyCtx" || reflect.TypeOf(ctx).String() == "*context.valueCtx" || reflect.TypeOf(ctx).String() == "*context.cancelCtx" {
		sessionContext := (ctx).(context.Context)
		id, err = mc.coll.UpdateOne(sessionContext, param, update)
	} else {
		sessionContext := (ctx).(mongo.SessionContext)
		id, err = mc.coll.UpdateOne(sessionContext, param, update)
	}
	if err != nil {
		return 0, err
	}
	resInserted := id.ModifiedCount
	return int(resInserted), err
}

func (mc *mongoCollection[R]) UpdateMany(ctx interface{}, filter interface{}, update interface{}) (int, int, error) {
	var updateResult *mongo.UpdateResult
	var err error

	log := utils_context.GetLogFromContext(ctx.(context.Context), constant.InfrastructureLayer, "database_client")

	opts := options.Update().SetUpsert(false)

	if reflect.TypeOf(ctx).String() == "*context.emptyCtx" || reflect.TypeOf(ctx).String() == "*context.valueCtx" || reflect.TypeOf(ctx).String() == "*context.cancelCtx" {
		sessionContext := (ctx).(context.Context)
		updateResult, err = mc.coll.UpdateMany(sessionContext, filter, update, opts)
	} else {
		sessionContext := (ctx).(mongo.SessionContext)
		updateResult, err = mc.coll.UpdateMany(sessionContext, filter, update, opts)
	}
	if err != nil {
		log.Error("Error in Update Many, err: %s", err.Error())
		return 0, 0, err
	}
	matchedCount := updateResult.MatchedCount
	modifiedCount := updateResult.ModifiedCount

	return int(matchedCount), int(modifiedCount), err
}

func (mc *mongoCollection[R]) InsertOne(ctx interface{}, document interface{}) (string, error) {

	var id *mongo.InsertOneResult
	var err error
	if reflect.TypeOf(ctx).String() == "*context.emptyCtx" || reflect.TypeOf(ctx).String() == "*context.valueCtx" {
		sessionContext := (ctx).(context.Context)
		id, err = mc.coll.InsertOne(sessionContext, document)
	} else {
		sessionContext := (ctx).(mongo.SessionContext)
		id, err = mc.coll.InsertOne(sessionContext, document)
	}
	if err != nil {
		return "", err
	}
	resInserted := id.InsertedID.(primitive.ObjectID).Hex()
	return resInserted, err

}

func (mc *mongoCollection[R]) DeleteOne(ctx interface{}, filter interface{}) (int64, error) {
	var count *mongo.DeleteResult
	var err error
	if reflect.TypeOf(ctx).String() == "*context.emptyCtx" || reflect.TypeOf(ctx).String() == "*context.valueCtx" {
		sessionContext := (ctx).(context.Context)
		count, err = mc.coll.DeleteOne(sessionContext, filter)
	} else {
		sessionContext := (ctx).(mongo.SessionContext)
		count, err = mc.coll.DeleteOne(sessionContext, filter)
	}
	return count.DeletedCount, err
}

func (mc *mongoCollection[R]) UpsertOne(ctx, filter, update interface{}) (string, int64, int64, error) {
	log := utils_context.GetLogFromContext(ctx.(context.Context), constant.InfrastructureLayer, "mongo_client.UpsertOne")

	log.Info("type context %s", reflect.TypeOf(ctx).String())
	var updateResult *mongo.UpdateResult
	var err error

	var sessionContext context.Context
	opts := options.Update().SetUpsert(true)
	ctxType := reflect.TypeOf(ctx).String()
	if ctxType == "*context.emptyCtx" ||
		ctxType == "*context.valueCtx" || ctxType == "*context.cancelCtx" {
		sessionContext = (ctx).(context.Context)
	} else {
		sessionContext = (ctx).(mongo.SessionContext)
	}

	updateResult, err = mc.coll.UpdateOne(sessionContext, filter, update, opts)
	if err != nil {
		log.Error("Error in called to UpdateOne: [%s]", err.Error())
		return "", 0, 0, err
	}

	matchedCount := updateResult.MatchedCount
	modifiedCount := updateResult.ModifiedCount
	upsertedCount := updateResult.UpsertedCount
	upsertedId := updateResult.UpsertedID

	var resUpserted string
	if upsertedId != nil {
		resUpserted = upsertedId.(primitive.ObjectID).String()
	}

	if matchedCount != 0 {
		log.Info("Matched and replaced an existing document in UpdateOne: [%d]", matchedCount)
	}
	if upsertedCount != 0 {
		log.Info("The new document was inserted successfully with ID: [%v]", upsertedId)
	}

	return resUpserted, modifiedCount, upsertedCount, nil
}
