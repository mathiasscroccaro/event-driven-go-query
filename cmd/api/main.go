package main

import (
	"github.com/mathiasscroccaro/event-driven-go-query/internal/use_case"
)

func main() {
	use_case.NewCLIAPIUseCase().Execute()
}
