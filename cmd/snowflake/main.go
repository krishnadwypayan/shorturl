package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/krishnadwypayan/shorturl/internal/routes"
	"github.com/krishnadwypayan/shorturl/internal/snowflake"
)

func main() {
	printBanner()
	generator := snowflake.NewGenerator(1)
	r := gin.Default()
	routes.RegisterSnowflakeRoutes(r, generator)
	r.Run(":8080")
}

func printBanner() {
	fmt.Print(`
███████╗███╗   ██╗ ██████╗ ██╗    ██╗███████╗██╗      █████╗ ██╗  ██╗███████╗
██╔════╝████╗  ██║██╔═══██╗██║    ██║██╔════╝██║     ██╔══██╗██║ ██╔╝██╔════╝
███████╗██╔██╗ ██║██║   ██║██║ █╗ ██║█████╗  ██║     ███████║█████╔╝ █████╗  
╚════██║██║╚██╗██║██║   ██║██║███╗██║██╔══╝  ██║     ██╔══██║██╔═██╗ ██╔══╝  
███████║██║ ╚████║╚██████╔╝╚███╔███╔╝██║     ███████╗██║  ██║██║  ██╗███████╗
╚══════╝╚═╝  ╚═══╝ ╚═════╝  ╚══╝╚══╝ ╚═╝     ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝                      
`)
}
