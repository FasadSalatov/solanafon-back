package models

// ManaPointTariff — pricing tiers for Mana Points purchase
type ManaPointTariff struct {
	ID        string  `gorm:"primarykey" json:"id"`            // mp_250, mp_500, etc.
	MPAmount  int     `gorm:"not null" json:"mp"`
	PriceUSDT float64 `gorm:"not null" json:"priceUsdt"`
	Currency  string  `gorm:"default:USD" json:"currency"`
	IsPopular bool    `gorm:"default:false" json:"isPopular"`
	IsActive  bool    `gorm:"default:true" json:"isActive"`
}

// WalletNetwork — supported blockchain networks
type WalletNetwork struct {
	ID        string `gorm:"primarykey" json:"id"` // solana, bitcoin
	Name      string `gorm:"not null" json:"name"`
	IconColor string `json:"iconColor"`
	IsDefault bool   `gorm:"default:false" json:"isDefault"`
}
