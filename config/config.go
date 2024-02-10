package config

// TronConfig holds the configuration for connecting to the Tron network.
type TronConfig struct {
    NetworkURL string
    APIKey string
}

// Netowrk constants.
const (
    TRON_MAINNET = "https://api.trongrid.io"
    TRON_SHASTA_TESTNET = "https://api.shasta.trongrid.io"
)

// NewTronConfig creates a new TronConfig with specified network and API key.
func NewTronConfig(network string, apiKey string) *TronConfig {
    return &TronConfig{
        NetworkURL: network,
        APIKey: apiKey,
    }
}

