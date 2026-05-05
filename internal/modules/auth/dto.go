package auth

type RegisterRequest struct {
	Role     string `json:"role" binding:"required,oneof=MERCHANT ADMIN"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`

	// merchant
	MerchantName   string  `json:"merchant_name"`
	OwnerName      string  `json:"owner_name"`
	Category       string  `json:"category"`
	Address        string  `json:"address"`
	PhoneNumber    string  `json:"phone_number"`
	CommissionRate float64 `json:"commission_rate"`

	// admin
	Name     string `json:"name"`
	MobileNo string `json:"mobile_no"`
	NIK      string `json:"nik"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
