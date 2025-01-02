package main

import "golang.org/x/sync/errgroup"

func main() {
	var eg errgroup.Group

	if err := eg.Wait(); err != nil {

	}
}
