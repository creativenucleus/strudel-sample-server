package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "port",
				Value: 5432,
				Usage: "port to serve from",
			},
			&cli.StringSliceFlag{
				Name:  "sources",
				Value: []string{"./"},
				Usage: "strudel.jsons to serve",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			port := cmd.Int("port")
			sources := cmd.StringSlice("sources")

			err := serve(port, sources)
			if err != nil {
				fmt.Println("Error:", err)
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func serve(port int, paths []string) error {
	fmt.Printf("Serving from port [%d]\n", port)

	if len(paths) == 0 {
		return fmt.Errorf("no source paths provided")
	}

	if len(paths) > 1 {
		return fmt.Errorf("Currently only one source path is supported")
	}

	sampleMap := make(strudelSampleMap)
	for _, path := range paths {
		fmt.Printf(" <- directory [%s]\n", path)

		var err error
		sampleMap, err = addToStrudelSampleMap(sampleMap, path)
		if err != nil {
			return err
		}
	}
	sampleMap["_base"] = []string{fmt.Sprintf("http://localhost:%d/", port)}

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		out, err := json.Marshal(sampleMap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		setOpenCORSHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		_, err = w.Write(out)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	router.HandleFunc(`/favicon.ico`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // nerfed
	})
	router.HandleFunc(`/{samplePath:.+}`, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r) // Get route variables
		samplePath := vars["samplePath"]
		fmt.Printf("Serving sample [%s]\n", samplePath)

		setOpenCORSHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		samplePathFull := fmt.Sprintf(`%s\%s`, filepath.Dir(paths[0]), samplePath)
		fmt.Println(samplePathFull)
		http.ServeFile(w, r, samplePathFull)
	})

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("localhost:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func setOpenCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
}
