package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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
				Usage: "alias:path-to-strudel.json to serve (this can be specified multiple times)",
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
	if len(paths) == 0 {
		return fmt.Errorf("no source paths provided")
	}

	// Load all sample packs into memory
	samplePacks := map[string]samplepack{}
	for _, path := range paths {
		parts := strings.Split(path, "<-")
		if len(parts) != 2 {
			return fmt.Errorf("source path [%s] is not in the correct format (alias<-path-to-strudel.json)", path)
		}

		packAlias := parts[0]
		sourcePath := parts[1]

		samplePack, err := readToStrudelSamplePack(sourcePath)
		if err != nil {
			return err
		}

		samplePacks[packAlias] = *samplePack
	}

	router := mux.NewRouter()
	router.HandleFunc(`/favicon.ico`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // nerfed
	})
	router.HandleFunc("/{packAlias}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r) // Get route variables
		packAlias := vars["packAlias"]

		samplePack, ok := samplePacks[packAlias]
		if !ok {
			http.Error(w, fmt.Sprintf("unknown sample pack alias [%s]", packAlias), http.StatusNotFound)
			return
		}

		out, err := samplePack.toData(fmt.Sprintf("http://localhost:%d/%s/", port, packAlias))
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
	router.HandleFunc(`/{packAlias}/{samplePath:.+}`, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r) // Get route variables
		packAlias := vars["packAlias"]

		samplePack, ok := samplePacks[packAlias]
		if !ok {
			http.Error(w, fmt.Sprintf("unknown sample pack alias [%s]", packAlias), http.StatusNotFound)
			return
		}

		samplePath := vars["samplePath"]

		setOpenCORSHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		fmt.Printf("Serving sample [%s : %s]\n", packAlias, samplePath)
		// filepath!
		samplePathFull := fmt.Sprintf(`%s\%s`, samplePack.pathBase, samplePath)
		fmt.Println(samplePathFull)
		http.ServeFile(w, r, samplePathFull)
	})

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("localhost:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Serving from port [%d]\n", port)
	for alias, pack := range samplePacks {
		fmt.Printf(" - /%s = %d samples\n", alias, len(pack.sampleMap))
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
