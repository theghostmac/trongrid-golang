package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/theghostmac/trongrid-golang/config"
)

// TronClient holds the configuration and HTTP client for making requests to the tron network.
type TronClient struct {
    Config *config.TronConfig
    HttpClient *http.Client
}

// Account represents a TRON network account, including its basic properties.
type Account struct {
    // Address is the account's address in hexadecimal.
    Address string `json:"address"`
    // Balance is the account's balance in Sun (1 TRX = 1,000,000 Sun).
    Balance int64 `json:"balance"`
    // PublicKey is the account's public key/
    PublicKey string `json:"public_key"`
}

// NewTronClient creates a new TronClient with the provided configuration.
func NewTronClient(cfg *config.TronConfig) *TronClient {
    return &TronClient{
        Config: cfg,
        HttpClient: &http.Client{}, // using the defualt HTTP client; customize based on your needs.
    }
}

// GetAccountBalance retrieves the balance of the specified account from the Tron network.
func (c *TronClient) GetAccountBalance(address string) (int64, error) {
    url := fmt.Sprintf("%s/wallet/getaccount", c.Config.NetworkURL)
    requestBody := fmt.Sprintf(`{"address": "%s", "visible": true}`, address)
    req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
    if err != nil {
        return 0, fmt.Errorf("creating request: %w", err)
    }

    req.Header.Add("Content-Type", "application/json")
    if c.Config.APIKey != "" {
        req.Header.Add("TRON-PRO-API-KEY", c.Config.APIKey)
    }

    resp, err := c.HttpClient.Do(req)
    if err != nil {
        return 0, fmt.Errorf("sending request: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return 0, fmt.Errorf("reading response body: %w", err)
    }

    var result struct {
        Balance int64 `json:"balance"`
    }

    if err := json.Unmarshal(body, &result); err != nil {
        return 0, fmt.Errorf("unmarshalling response: %w", err)
    }

    return result.Balance, nil
}

// CreateAccount creates a new account on the TRON network.
func (c *TronClient) CreateAccount() (*Account, error) {
    url := fmt.Sprintf("%s/wallet/createaccount", c.Config.NetworkURL)
    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        return nil, fmt.Errorf("creating request: %w", err)
    }

    req.Header.Add("Content-Type", "application/json")
    if c.Config.APIKey != "" {
        req.Header.Add("TRON-PRO-API-KEY", c.Config.APIKey)
    }

    resp, err := c.HttpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("sending request: %w", err)
    }
    defer resp.Body.Close()

    var account Account
    if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
        return nil, fmt.Errorf("decoding response: %w", err)
    }

    return &account, nil
}

// TransferTRX prepare a transaction to transfer TRX from one account to another.
func (c *TronClient) TransferTRX(fromAddress, toAddress string, amount int64) (string, error) {
    // Generate a transaction, sign it with the private key, and then broadcast it using the broadcast method.

    payload := map[string]interface{} {
        "to_address": toAddress,
        "owner_address": fromAddress,
        "amount": amount,
        "visible": true,
    }
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return "", fmt.Errorf("marshalling payload: %w", err)
    }

    url := fmt.Sprintf("%s/wallet/createtransaction", c.Config.NetworkURL)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
    if err != nil {
        return "", fmt.Errorf("creating request: %w", err)
    }

    req.Header.Add("accept", "application/json")
    req.Header.Add("Content-Type", "application/json")
    if c.Config.APIKey != "" {
        req.Header.Add("TRON-PRO-API-KEY", c.Config.APIKey)
    }

    res, err := c.HttpClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("sending request: %w", err)
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return "", fmt.Errorf("reading response body: %w", err)
    }

    var response struct {
        TxID string `json:"txID"`
    }

    if err := json.Unmarshal(body, &response); err != nil {
        return "", fmt.Errorf("unmarshalling response: %w", err)
    }

    if response.TxID == "" {
        return "", fmt.Errorf("transaction ID is empty, response: %s", string(body))
    }
    
    return response.TxID, nil
}

// BroadcastTransaction sends a signed transaction to the TRON network.
func (c *TronClient) BroadcastTransaction(signedTx string) (string, error) {
    url := fmt.Sprintf("%s/wallet/broadcasttransaction", c.Config.NetworkURL)
    requestBody := fmt.Sprintf(`{"raw_data_hex": "%s"}`, signedTx)
    req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
    if err != nil {
        return "", fmt.Errorf("creating request: %w", err)
    }

    req.Header.Add("Content-Type", "application/json")
    if c.Config.APIKey != "" {
        req.Header.Add("TRON-PRO-API-KEY", c.Config.APIKey)
    }

    resp, err := c.HttpClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("sending request: %w", err)
    }
    defer resp.Body.Close()

    var result struct {
        // Adjust according to the actual API response
        TxID string `json:"txID"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", fmt.Errorf("decoding response: %w", err)
    }

    return result.TxID, nil
}

