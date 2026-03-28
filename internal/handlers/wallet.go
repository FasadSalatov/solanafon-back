package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fasad/solanafon-back/internal/config"
	"github.com/fasad/solanafon-back/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// WalletHandler handles /api/wallet/* endpoints
type WalletHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewWalletHandler(db *gorm.DB, cfg *config.Config) *WalletHandler {
	return &WalletHandler{db: db, cfg: cfg}
}

// GetBalance — GET /api/wallet/balance?address=...
func (h *WalletHandler) GetBalance(c *fiber.Ctx) error {
	address := c.Query("address")
	if address == "" {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "address is required"}})
	}

	// Get SOL balance from Solana RPC
	solBalance := h.getSolBalance(address)

	// Get SPL token balances
	tokens := h.getTokenBalances(address)

	totalUsd := solBalance * 148.35 // placeholder SOL price
	for _, t := range tokens {
		if v, ok := t["valueUsd"].(float64); ok {
			totalUsd += v
		}
	}

	return c.JSON(fiber.Map{
		"address":  address,
		"sol":      fmt.Sprintf("%.4f", solBalance),
		"tokens":   tokens,
		"totalUsd": fmt.Sprintf("%.2f", totalUsd),
	})
}

// SendTransaction — POST /api/wallet/send
func (h *WalletHandler) SendTransaction(c *fiber.Ctx) error {
	var input struct {
		SignedTransaction string `json:"signedTransaction"`
	}
	if err := c.BodyParser(&input); err != nil || input.SignedTransaction == "" {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "signedTransaction is required"}})
	}

	// Broadcast to Solana
	sig, err := h.broadcastTransaction(input.SignedTransaction)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fiber.Map{"code": "TX_FAILED", "message": err.Error()}})
	}

	return c.JSON(fiber.Map{"signature": sig, "status": "confirmed"})
}

// GetTransactions — GET /api/wallet/transactions?address=...
func (h *WalletHandler) GetTransactions(c *fiber.Ctx) error {
	address := c.Query("address")
	if address == "" {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "address is required"}})
	}

	sigs := h.getTransactionSignatures(address)
	return c.JSON(fiber.Map{"transactions": sigs, "hasMore": len(sigs) > 0})
}

// GetTokens — GET /api/wallet/tokens
func (h *WalletHandler) GetTokens(c *fiber.Ctx) error {
	tokens := []fiber.Map{
		{"mint": "So11111111111111111111111111111111", "symbol": "SOL", "name": "Solana", "decimals": 9, "icon": "https://raw.githubusercontent.com/solana-labs/token-list/main/assets/mainnet/So11111111111111111111111111111111/logo.png"},
		{"mint": "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", "symbol": "USDC", "name": "USD Coin", "decimals": 6, "icon": "https://raw.githubusercontent.com/solana-labs/token-list/main/assets/mainnet/EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v/logo.png"},
		{"mint": "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB", "symbol": "USDT", "name": "Tether USD", "decimals": 6, "icon": "https://raw.githubusercontent.com/solana-labs/token-list/main/assets/mainnet/Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB/logo.svg"},
	}
	return c.JSON(fiber.Map{"tokens": tokens})
}

// GetPrices — GET /api/wallet/prices
func (h *WalletHandler) GetPrices(c *fiber.Ctx) error {
	// Simplified — in production, use CoinGecko/Jupiter API
	prices := fiber.Map{
		"So11111111111111111111111111111111":              148.35,
		"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v": 1.0,
		"Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB":  1.0,
	}
	return c.JSON(fiber.Map{"prices": prices})
}

// GetTransactionStatus — GET /api/wallet/status?signature=...
func (h *WalletHandler) GetTransactionStatus(c *fiber.Ctx) error {
	sig := c.Query("signature")
	if sig == "" {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "signature is required"}})
	}

	status := h.getSignatureStatus(sig)
	return c.JSON(fiber.Map{"signature": sig, "status": status})
}

// GetManaPoints — GET /api/wallet/mana-points
func (h *WalletHandler) GetManaPoints(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var user models.User
	h.db.First(&user, userID)

	var history []models.ManaTransaction
	h.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(50).Find(&history)

	var income, rewards int
	for _, tx := range history {
		if tx.Amount > 0 {
			income += tx.Amount
			if tx.Type == "reward" {
				rewards++
			}
		}
	}

	historyResult := make([]fiber.Map, 0, len(history))
	for _, tx := range history {
		historyResult = append(historyResult, fiber.Map{
			"id": fmt.Sprintf("mp_tx_%d", tx.ID), "type": tx.Type,
			"amount": tx.Amount, "source": tx.Type, "description": tx.Description,
			"createdAt": tx.CreatedAt,
		})
	}

	return c.JSON(fiber.Map{
		"success":     true,
		"balance":     user.ManaPoints,
		"balanceUsdt": float64(user.ManaPoints) * 0.023,
		"income":      income,
		"rewards":     rewards,
		"history":     historyResult,
	})
}

// GetTariffs — GET /api/wallet/mana-points/tariffs
func (h *WalletHandler) GetTariffs(c *fiber.Ctx) error {
	var tariffs []models.ManaPointTariff
	h.db.Where("is_active = true").Order("mp_amount ASC").Find(&tariffs)

	if len(tariffs) == 0 {
		// Seed default tariffs
		defaults := []models.ManaPointTariff{
			{ID: "mp_250", MPAmount: 250, PriceUSDT: 5.79, Currency: "USD"},
			{ID: "mp_500", MPAmount: 500, PriceUSDT: 11.90, Currency: "USD", IsPopular: true},
			{ID: "mp_1000", MPAmount: 1000, PriceUSDT: 22.99, Currency: "USD"},
			{ID: "mp_2500", MPAmount: 2500, PriceUSDT: 57.98, Currency: "USD"},
		}
		for _, t := range defaults {
			t.IsActive = true
			h.db.Create(&t)
		}
		tariffs = defaults
	}

	result := make([]fiber.Map, len(tariffs))
	for i, t := range tariffs {
		result[i] = fiber.Map{
			"id": t.ID, "mp": t.MPAmount, "priceUsdt": t.PriceUSDT,
			"currency": t.Currency, "isPopular": t.IsPopular,
		}
	}
	return c.JSON(fiber.Map{"success": true, "tariffs": result})
}

// PurchaseMana — POST /api/wallet/mana-points/purchase
func (h *WalletHandler) PurchaseMana(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var input struct {
		TariffID      string `json:"tariffId"`
		PaymentMethod string `json:"paymentMethod"`
	}
	if err := c.BodyParser(&input); err != nil || input.TariffID == "" {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "tariffId is required"}})
	}

	var tariff models.ManaPointTariff
	if err := h.db.First(&tariff, "id = ?", input.TariffID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fiber.Map{"code": "NOT_FOUND", "message": "Tariff not found"}})
	}

	var user models.User
	h.db.First(&user, userID)
	user.ManaPoints += tariff.MPAmount
	h.db.Save(&user)

	h.db.Create(&models.ManaTransaction{
		UserID: userID, Amount: tariff.MPAmount, Type: "purchase",
		Description: fmt.Sprintf("Purchased %d MP", tariff.MPAmount),
	})

	return c.JSON(fiber.Map{
		"success": true, "transactionId": fmt.Sprintf("tx_%d", time.Now().UnixMilli()),
		"newBalance": user.ManaPoints, "amountPurchased": tariff.MPAmount, "amountCharged": tariff.PriceUSDT,
	})
}

// GiftMana — POST /api/wallet/mana-points/gift
func (h *WalletHandler) GiftMana(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var input struct {
		RecipientAddress string `json:"recipientAddress"`
		Amount           int    `json:"amount"`
	}
	if err := c.BodyParser(&input); err != nil || input.Amount <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "VALIDATION_ERROR", "message": "Invalid gift amount"}})
	}

	var user models.User
	h.db.First(&user, userID)
	if user.ManaPoints < input.Amount {
		return c.Status(400).JSON(fiber.Map{"error": fiber.Map{"code": "INSUFFICIENT_BALANCE", "message": "Not enough Mana Points"}})
	}

	user.ManaPoints -= input.Amount
	h.db.Save(&user)

	h.db.Create(&models.ManaTransaction{
		UserID: userID, Amount: -input.Amount, Type: "gift_sent",
		Description: fmt.Sprintf("Gift sent: %d MP", input.Amount),
	})

	return c.JSON(fiber.Map{
		"success": true, "transactionId": fmt.Sprintf("tx_gift_%d", time.Now().UnixMilli()),
		"newBalance": user.ManaPoints, "amountSent": input.Amount,
	})
}

// GetNetworks — GET /api/wallet/networks
func (h *WalletHandler) GetNetworks(c *fiber.Ctx) error {
	var networks []models.WalletNetwork
	h.db.Find(&networks)

	if len(networks) == 0 {
		defaults := []models.WalletNetwork{
			{ID: "solana", Name: "Solana", IconColor: "#9945FF", IsDefault: true},
			{ID: "bitcoin", Name: "Bitcoin", IconColor: "#F7931A"},
		}
		for _, n := range defaults {
			h.db.Create(&n)
		}
		networks = defaults
	}

	result := make([]fiber.Map, len(networks))
	for i, n := range networks {
		result[i] = fiber.Map{
			"id": n.ID, "name": n.Name, "iconColor": n.IconColor, "isDefault": n.IsDefault,
		}
	}
	return c.JSON(fiber.Map{"success": true, "networks": result})
}

// --- Solana RPC helpers ---

func (h *WalletHandler) solanaRPC(method string, params interface{}) (json.RawMessage, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0", "id": 1, "method": method, "params": params,
	})
	resp, err := http.Post(h.cfg.SolanaRPCURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)

	var result struct {
		Result json.RawMessage `json:"result"`
		Error  interface{}     `json:"error"`
	}
	json.Unmarshal(data, &result)
	return result.Result, nil
}

func (h *WalletHandler) getSolBalance(address string) float64 {
	result, err := h.solanaRPC("getBalance", []interface{}{address})
	if err != nil {
		return 0
	}
	var balResp struct {
		Value uint64 `json:"value"`
	}
	json.Unmarshal(result, &balResp)
	return float64(balResp.Value) / 1e9
}

func (h *WalletHandler) getTokenBalances(address string) []fiber.Map {
	result, err := h.solanaRPC("getTokenAccountsByOwner", []interface{}{
		address,
		map[string]string{"programId": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"},
		map[string]string{"encoding": "jsonParsed"},
	})
	if err != nil {
		return nil
	}

	var parsed struct {
		Value []struct {
			Account struct {
				Data struct {
					Parsed struct {
						Info struct {
							Mint        string `json:"mint"`
							TokenAmount struct {
								UIAmountString string  `json:"uiAmountString"`
								Decimals       int     `json:"decimals"`
								UIAmount       float64 `json:"uiAmount"`
							} `json:"tokenAmount"`
						} `json:"info"`
					} `json:"parsed"`
				} `json:"data"`
			} `json:"account"`
		} `json:"value"`
	}
	json.Unmarshal(result, &parsed)

	tokens := make([]fiber.Map, 0)
	for _, v := range parsed.Value {
		info := v.Account.Data.Parsed.Info
		if info.TokenAmount.UIAmount > 0 {
			tokens = append(tokens, fiber.Map{
				"mint": info.Mint, "balance": info.TokenAmount.UIAmountString,
				"decimals": info.TokenAmount.Decimals,
			})
		}
	}
	return tokens
}

func (h *WalletHandler) broadcastTransaction(signedTx string) (string, error) {
	result, err := h.solanaRPC("sendTransaction", []interface{}{signedTx, map[string]string{"encoding": "base64"}})
	if err != nil {
		return "", err
	}
	var sig string
	json.Unmarshal(result, &sig)
	return sig, nil
}

func (h *WalletHandler) getTransactionSignatures(address string) []fiber.Map {
	result, err := h.solanaRPC("getSignaturesForAddress", []interface{}{address, map[string]int{"limit": 20}})
	if err != nil {
		return nil
	}
	var sigs []struct {
		Signature string `json:"signature"`
		Slot      int    `json:"slot"`
		BlockTime int64  `json:"blockTime"`
	}
	json.Unmarshal(result, &sigs)

	txs := make([]fiber.Map, 0, len(sigs))
	for _, s := range sigs {
		txs = append(txs, fiber.Map{
			"signature": s.Signature, "slot": s.Slot,
			"timestamp": s.BlockTime * 1000, "status": "finalized",
		})
	}
	return txs
}

func (h *WalletHandler) getSignatureStatus(sig string) string {
	result, err := h.solanaRPC("getSignatureStatuses", []interface{}{[]string{sig}})
	if err != nil {
		return "pending"
	}
	var statuses struct {
		Value []struct {
			ConfirmationStatus string `json:"confirmationStatus"`
		} `json:"value"`
	}
	json.Unmarshal(result, &statuses)
	if len(statuses.Value) > 0 && statuses.Value[0].ConfirmationStatus != "" {
		return statuses.Value[0].ConfirmationStatus
	}
	return "pending"
}

// Suppress unused import warning
var _ = strings.TrimSpace
