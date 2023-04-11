package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/nfnt/resize"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// init logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Suppression des fichiers dans imgResize pour vraiment générer les images
	deleteResizeFiles()

	// Ouverture du dossier img
	d, err := os.Open("./img")
	if err != nil {
		log.Fatal().Err(err).Msg("error opening current directory")
		panic(err)
	}
	// Récupération des fichiers du dossier img
	files, err := d.Readdir(-1)
	if err != nil {
		log.Fatal().Err(err).Msg("error reading current directory")
		panic(err)
	}

	// Sans go routine
	elapsedWithout := time.Now()
	for i, file := range files {
		resizeImage(i, file)
	}

	// Suppression des fichiers dans imgResize pour vraiment générer les images
	deleteResizeFiles()

	// Avec go routine
	elapsedWith := time.Now()
	var wg sync.WaitGroup
	for i, file := range files {
		wg.Add(1)
		go func(i int, file os.FileInfo) {
			defer wg.Done()
			resizeImage(i, file)
		}(i, file)
	}
	wg.Wait()

	log.Info().Msg(fmt.Sprintf("Without go routine: %s", time.Since(elapsedWithout)))
	log.Info().Msg(fmt.Sprintf("With go routine: %s", time.Since(elapsedWith)))
	log.Info().Msg(fmt.Sprintf("Difference: %s", time.Since(elapsedWith)-time.Since(elapsedWithout)))
	log.Info().Msg(fmt.Sprintf("Difference: %f", float64(time.Since(elapsedWith)-time.Since(elapsedWithout))/float64(time.Since(elapsedWithout))*100) + "%")
}

func deleteResizeFiles() {
	// delete all files in imgResize
	f, err := os.ReadDir("./imgResize")
	if err != nil {
		log.Fatal().Err(err).Msg("error reading imgResize directory")
		panic(err)
	}

	for _, file := range f {
		err = os.Remove("./imgResize/" + file.Name())
		if err != nil {
			log.Fatal().Err(err).Msg("error deleting file")
			panic(err)
		}
	}
}

func resizeImage(i int, file os.FileInfo) {
	// check if file is a directory
	if !file.IsDir() {
		// get file extension
		extension := filepath.Ext(file.Name())
		if extension != ".jpg" {
			log.Error().Msg("File is not a jpg")
		} else {
			// open file
			file, err := os.Open("./img/" + file.Name())
			if err != nil {
				log.Fatal().Err(err).Msg("error opening file")
			}

			// decode jpeg into image.Image
			img, err := jpeg.Decode(file)
			if err != nil {
				log.Fatal().Err(err).Msg("error decoding file")
			}
			file.Close()

			// resize to width 1000 using Lanczos resampling
			// and preserve aspect ratio
			m := resize.Resize(100, 0, img, resize.Lanczos3)

			out, err := os.Create("imgResize/" + fmt.Sprint(i) + ".jpg")
			if err != nil {
				log.Fatal().Err(err).Msg("error creating file")
			}
			defer out.Close()

			// write new image to file
			err = jpeg.Encode(out, m, nil)
			if err != nil {
				log.Fatal().Err(err).Msg("error writing file")
			}
		}
	}
}
