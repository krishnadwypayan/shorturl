package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/krishnadwypayan/shorturl/internal/routes"
)

func main() {
	printBanner()
	r := gin.Default()
	routes.RegisterShortURLRoutes(r)
	r.Run(":8081")
}

func printBanner() {
	fmt.Print(`
███████╗██╗  ██╗ ██████╗ ██████╗ ████████╗██╗   ██╗██████╗ ██╗     
██╔════╝██║  ██║██╔═══██╗██╔══██╗╚══██╔══╝██║   ██║██╔══██╗██║     
███████╗███████║██║   ██║██████╔╝   ██║   ██║   ██║██████╔╝██║     
╚════██║██╔══██║██║   ██║██╔══██╗   ██║   ██║   ██║██╔══██╗██║     
███████║██║  ██║╚██████╔╝██║  ██║   ██║   ╚██████╔╝██║  ██║███████╗
╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚══════╝
`)
}
