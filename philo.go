package main

import (
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	THINKING = iota + 1
	EATING
	STARVING
	EATING_TIME   = 3
	STARVING_TIME = 9
)

func diningPhilosophers() {
	philosophers := make([]*Philosopher, 5)
	for i := range philosophers {
		philosophers[i] = &Philosopher{}
	}

	forks := make([]*Fork, 5)

	for i := 0; i < 5; i++ {
		forks[i] = &Fork{}
	}

	for i := 0; i < 5; i++ {
		philosophers[i].leftFork = forks[i]
		philosophers[i].rightFork = forks[(i+1)%5]
	}

	log.Info().Msg("Philosophers are seated")
	for {
		for id, philosopher := range philosophers {
			go philosopher.init(id)
		}

		log.Info().Msg("Philosophers are eating")

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
	log.Info().Int("id", p.id).Msg("Philosopher is seated")
	p.id = id
	p.think()

	go func() {
		defer table.wg.Done()
		log.Info().Int("id", p.id).Msg("Philosopher is done eating")
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
	log.Info().Int("id", p.id).Msg("Philosopher starts thinking")
	p.state = THINKING
	time.Sleep(time.Duration(rand.Intn(2)*1000) * time.Millisecond)

	go p.eat()
	go p.starving()
}

func (p *Philosopher) starving() {
	log.Info().Int("id", p.id).Msg("Philosopher starts starving")
	p.state = STARVING
	p.starvingCountdown = time.Now().Add(time.Duration(STARVING_TIME*1000) * time.Millisecond)
	for range time.Tick(100 * time.Millisecond) {
		if time.Now().After(p.starvingCountdown) && p.state == STARVING {
			p.die()
		}
	}
}

func (p *Philosopher) eat() {
	log.Info().Int("id", p.id).Msg("Philosopher starts eating")
	p.takeForks()
	p.state = EATING
	time.Sleep(time.Duration(EATING_TIME*1000) * time.Millisecond)
	log.Info().Int("id", p.id).Msg("Philosopher stops eating")
	p.releaseForks()
	go p.think()
}

func (p *Philosopher) die() {
	log.Warn().Int("id", p.id).Msg("Philosopher died")
	os.Exit(1)
}

type Fork struct {
	lock sync.Mutex
}

func (f *Fork) take() {
	f.lock.Lock()
}

func (f *Fork) release() {
	f.lock.Unlock()
}
