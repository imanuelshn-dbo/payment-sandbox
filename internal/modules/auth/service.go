package auth

import (
	"errors"
	"os"
	"payment-sandbox/internal/models"
	apperror "payment-sandbox/pkg/app-error"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db}
}

func (s *Service) Register(req RegisterRequest) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if req.Role == "MERCHANT" && req.MerchantName == "" {
			return apperror.BadRequest("merchant_name required")
		}

		if req.Role == "ADMIN" && req.Name == "" {
			return apperror.BadRequest("name required")
		}

		// cek email
		var count int64
		tx.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
		if count > 0 {
			return apperror.Conflict("email already registered")
		}

		// hash password
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

		user := models.User{
			Email:    req.Email,
			Password: string(hash),
			Role:     req.Role,
		}

		// flow register as merchant
		if req.Role == "MERCHANT" {

			merchant := models.Merchant{
				MerchantName:   req.MerchantName,
				OwnerName:      req.OwnerName,
				Category:       req.Category,
				Address:        req.Address,
				PhoneNumber:    req.PhoneNumber,
				Email:          req.Email,
				Status:         "ACTIVE",
				CommissionRate: req.CommissionRate,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			if err := tx.Create(&merchant).Error; err != nil {
				return err
			}

			user.MerchantID = merchant.ID

			if err := tx.Create(&user).Error; err != nil {
				return err
			}

			// wallet default
			wallet := models.Wallet{
				UserID:  user.ID,
				Balance: 0,
			}

			if err := tx.Create(&wallet).Error; err != nil {
				return err
			}

			return nil
		}

		// flow register as admin
		if req.Role == "ADMIN" {

			admin := models.Admin{
				Name:      req.Name,
				Email:     req.Email,
				MobileNo:  req.MobileNo,
				NIK:       req.NIK,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := tx.Create(&admin).Error; err != nil {
				return err
			}

			user.AdminID = admin.ID

			if err := tx.Create(&user).Error; err != nil {
				return err
			}

			return nil
		}

		return apperror.BadRequest("invalid role")
	})
}

func (s *Service) Login(req LoginRequest) (string, error) {
	var user models.User

	// get user
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id":     user.ID,
		"role":        user.Role,
		"merchant_id": user.MerchantID,
		"admin_id":    user.AdminID,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) RefreshToken(refreshToken string) (TokenPair, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return TokenPair{}, errors.New("invalid refresh token")
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["type"] != "refresh" {
		return TokenPair{}, errors.New("invalid token type")
	}

	userID := int64(claims["user_id"].(float64))
	role := claims["role"].(string)

	return generateToken(userID, role)
}
