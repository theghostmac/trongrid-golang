package main

import (
	"fmt"
	"os"

	"github.com/theghostmac/trongrid-golang/config"
)

func main() {
    network := os.Getenv("TRON_NETWORK")
    apiKey := os.Getenv("TRONGRID_API_KEY")

    var networkURL string
    switch network {
        case "mainnet":
            networkURL = config.TRON_MAINNET
    case "shasta":
        networkURL = config.TRON_SHASTA_TESTNET
    default:
        fmt.Println("Invalid TRON network specified. Please use 'mainnet' or 'shasta'.")
        os.Exit(1)
    }

    tronConfig := config.NewTronConfig(networkURL, apiKey)
    fmt.Printf("Configured to use network: %s with API key: %s\n", tronConfig.NetworkURL, apiKey)
}
