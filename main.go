package main

import (
	"fmt"
	"image/jpeg"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/nfnt/resize"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	THINKING = iota + 1
	EATING
	STARVING
	EATING_TIME   = 3
	STARVING_TIME = 9
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	arg1 := os.Args[1]

	if arg1 == "resize" {
		resizeImages()
	} else if arg1 == "dining" {
		fmt.Println("je dine, tu dines, il dine, nous dinons, vous dinez, ils dinent")
		diningPhilosophers()
	} else {
		fmt.Println("Invalid argument")
		os.Exit(1)
	}
}

func diningPhilosophers() {
	philosophers := make([]*Philosopher, 5)
	for i := range philosophers {
		philosophers[i] = &Philosopher{}
	}

	forks := make([]*Fork, 5)

	for i := 0; i < 5; i++ {
		forks[i] = &Fork{}
		forks[i].init()
	}

	for i := 0; i < 5; i++ {
		philosophers[i].leftFork = forks[i]
		philosophers[i].rightFork = forks[(i+1)%5]
	}

	fmt.Println("Philosophers are seated")
	for {
		for id, philosopher := range philosophers {
			go philosopher.init(id)
		}

		fmt.Println("Philosophers are eating")

		table.wg.Wait()
	}
}

type Table struct {
	wg sync.WaitGroup
}

var table = &Table{}

type Philosopher struct {
	id                int
	leftFork          *Fork
	rightFork         *Fork
	state             int
	starvingCountdown time.Time
}

func (p *Philosopher) init(id int) {
	fmt.Println("Philosopher ", id, " is seated")
	p.id = id
	p.think()

	go func() {
		defer table.wg.Done()
		fmt.Println("Philosopher ", id, " is done eating")
	}()
}

func (p *Philosopher) takeForks() {
	p.leftFork.take()
	p.rightFork.take()
}

func (p *Philosopher) releaseForks() {
	p.leftFork.release()
	p.rightFork.release()
}

func (p *Philosopher) think() {
	table.wg.Add(1)
	fmt.Println("Philosopher ", p.id, " starts thinking")
	p.state = THINKING
	time.Sleep(time.Duration(rand.Intn(2)*1000) * time.Millisecond)

	go p.eat()
	go p.starving()
}

func (p *Philosopher) starving() {
	fmt.Println("Philosopher ", p.id, " starts starving")
	p.state = STARVING
	p.starvingCountdown = time.Now().Add(time.Duration(STARVING_TIME*1000) * time.Millisecond)
	for range time.Tick(100 * time.Millisecond) {
		if time.Now().After(p.starvingCountdown) && p.state == STARVING {
			p.die()
		}
	}
}

func (p *Philosopher) eat() {
	fmt.Println("Philosopher ", p.id, " starts eating")
	p.takeForks()
	p.state = EATING
	time.Sleep(time.Duration(EATING_TIME*1000) * time.Millisecond)
	fmt.Println("Philosopher ", p.id, " stops eating")
	p.releaseForks()
	go p.think()
}

func (p *Philosopher) die() {
	fmt.Println("Philosopher", p.id, "died")
	os.Exit(1)
}

type Fork struct {
	lock sync.Mutex
}

func (f *Fork) init() {
	// f.lock = sync.Mutex{}
}

func (f *Fork) take() {
	f.lock.Lock()
}

func (f *Fork) release() {
	f.lock.Unlock()
}

func resizeImages() {
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
