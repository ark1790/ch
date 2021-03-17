package cmd

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ark1790/ch/eventstore/api"
	"github.com/ark1790/ch/eventstore/belt"
	"github.com/ark1790/ch/eventstore/belt/localqueue"
	elasticrepo "github.com/ark1790/ch/eventstore/repo/elastic"
	elastic "github.com/olivere/elastic/v7"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start API server",
	Long:  `Start the API server`,

	Run: func(cmd *cobra.Command, args []string) {
		c, err := elastic.NewClient(
			elastic.SetSniff(false),
			elastic.SetURL(viper.GetString("ELASTIC_URL")),
		)
		if err != nil {
			panic(err)
		}

		eRepo := elasticrepo.NewEventRepo(c)

		q := localqueue.NewLocalQueue()
		beltSize := viper.GetInt("BELT_SIZE")

		for beltSize > 0 {
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {

				log.Println("Starting worker", beltSize)
				wg.Done()
				w := belt.NewWorker(q, eRepo)
				w.Start()
			}()
			beltSize--
			wg.Wait()

		}

		srvr := api.NewServer(*eRepo, q)

		port := viper.GetString("PORT")
		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT)

		go func() {
			log.Println("Listening on " + port)
			if err := srvr.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		}()

		<-stop

		log.Println("Shutting down server...")

		srvr.GracefulStop()

	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
