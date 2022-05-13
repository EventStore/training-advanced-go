package main

import (
	"context"
	"net/http"
	"time"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/EventStore/training-introduction-go/application"
	"github.com/EventStore/training-introduction-go/controllers"
	"github.com/EventStore/training-introduction-go/domain/doctorday"
	"github.com/EventStore/training-introduction-go/eventsourcing"
	"github.com/EventStore/training-introduction-go/infrastructure"
	"github.com/EventStore/training-introduction-go/infrastructure/inmemory"
	"github.com/EventStore/training-introduction-go/infrastructure/mongodb"
	"github.com/EventStore/training-introduction-go/infrastructure/projections"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	EsdbConnectionString  = "esdb://localhost:2113?tls=false"
	MongoConnectionString = "mongodb://localhost"
)

func main() {
	esdbClient, err := createESDBClient()
	if err != nil {
		panic(err)
	}

	mongoClient, err := createMongoClient()
	if err != nil {
		panic(err)
	}

	typeMapper := eventsourcing.NewTypeMapper()
	serde := infrastructure.NewEsEventSerde(typeMapper)
	eventStore := infrastructure.NewEsEventStore(esdbClient, "scheduling", serde)

	doctorday.RegisterTypes(typeMapper)

	dispatcher := getDispatcher(eventStore)
	mongoDatabase := mongoClient.Database("projections")
	availableSlotsRepo := mongodb.NewAvailableSlotsRepository(mongoDatabase)
	commandStore := infrastructure.NewEsCommandStore(eventStore, esdbClient, serde, dispatcher)

	dayArchiver := application.NewDayArchiverProcessManager(
		inmemory.NewColdStorage(),
		mongodb.NewArchivableDayRepository(mongoDatabase),
		commandStore,
		eventStore,
		-180*24*time.Hour,
	)

	subManager := projections.NewSubscriptionManager(
		esdbClient,
		infrastructure.NewEsCheckpointStore(esdbClient, "DaySubscription", serde),
		serde,
		"$all",
		projections.NewProjector(application.NewAvailableSlotsProjection(availableSlotsRepo)),
		projections.NewProjector(dayArchiver))

	err = subManager.Start(context.TODO())
	if err != nil {
		panic(err)
	}

	err = commandStore.Start()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", hello)

	s := controllers.NewSlotsController(dispatcher, availableSlotsRepo.AvailableSlotsRepository, eventStore)
	s.Register("/api/", e)

	e.Logger.Fatal(e.Start(":5001"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, Training!")
}

func createESDBClient() (*esdb.Client, error) {
	settings, err := esdb.ParseConnectionString(EsdbConnectionString)
	if err != nil {
		return nil, err
	}

	db, err := esdb.NewClient(settings)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createMongoClient() (*mongo.Client, error) {
	return mongo.Connect(context.TODO(), options.Client().ApplyURI(MongoConnectionString))
}

func getDispatcher(eventStore infrastructure.EventStore) *infrastructure.Dispatcher {
	aggregateStore := infrastructure.NewEsAggregateStore(eventStore, 5)
	dayRepository := doctorday.NewEventStoreDayRepository(aggregateStore)
	handlers := doctorday.NewHandlers(dayRepository)
	cmdHandlerMap := infrastructure.NewCommandHandlerMap(handlers)
	dispatcher := infrastructure.NewDispatcher(cmdHandlerMap)
	return &dispatcher
}
