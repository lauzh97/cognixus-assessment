package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
	data "todo/internal/data"
	google "todo/internal/google"
	service "todo/internal/service"
	pb "todo/proto/todo"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// to read yaml files
func startViper() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
}

func CheckDatabase(ctx context.Context) error {
	query := `SELECT schema_name FROM information_schema.schemata WHERE schema_name='main'`
	row := data.DB.QueryRow(query)

	var result string 
	err := row.Scan(&result);
	if err != nil {
		if err == sql.ErrNoRows {
			// run first time db setup (create schema/tables etc)
			fmt.Println("Running first time setup for db...")
			sqlScript, err := os.ReadFile("postgresql/1_first_time_up.sql")
			if err != nil {
				return err
			}
			query := string(sqlScript)
			if _, err := data.DB.Exec(query); err != nil {
				return err
			}
			fmt.Println("Done")
		} else {
			return err
		}
	}

	return nil
}

func startDB(ctx context.Context)  {
	var err error

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", 
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.name"),
	)
         
    // open database
    data.DB, err = sql.Open("postgres", psqlconn)
    if err != nil {
		log.Fatalln("Failed to open database:", err)
	}

	err = CheckDatabase(ctx)
	for err != nil {
		time.Sleep(2*time.Second)
		fmt.Println(err.Error() + " retrying connection to database...")
		err = CheckDatabase(ctx)
	}
	
	fmt.Println("Serving database on port " + strconv.Itoa(viper.GetInt("database.port")))
}

func startGRPC(ctx context.Context) {
	grpcPort := viper.GetString("server.grpcPort")

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer()
	pb.RegisterTodoServer(s, service.NewTodoServer(ctx))
	log.Println("Serving gRPC on http://0.0.0.0" + grpcPort)
	go func() {
		log.Fatalln(s.Serve(lis))
	}()
}

func startHTTP() {
	grpcPort := viper.GetString("server.grpcPort")
	httpPort := viper.GetString("server.httpPort")

	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0" + grpcPort,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterTodoHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    httpPort,
		Handler: gwmux,
	}

	log.Println("Serving HTTP on http://0.0.0.0" + httpPort)
	go func() {
		log.Fatalln(gwServer.ListenAndServe())
	}()
}

func startFrontend() {
	frontPort := viper.GetString("server.frontPort")

	google.InitializeOAuthGoogle()

    http.HandleFunc("/", google.HandleMain)
    http.HandleFunc("/auth/google/login", google.HandleGoogleLogin)
    http.HandleFunc("/auth/google/callback", google.CallBackFromGoogle)
    http.HandleFunc("/auth/google/authenticated", google.HandleAuthenticated)

	log.Println("Serving Frontend on http://0.0.0.0" + frontPort)
    log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	ctx := context.Background()

	startViper()
	startDB(ctx)
	startGRPC(ctx)
	startHTTP()
	startFrontend()
}