package auth_service

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func (s *service) GenerateJWT(userRole string, userID string) (string, error) {
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   userID,
		"role": userRole,
		"iat":  time.Now().Unix(),
		"eat":  time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(s.privateKey)
}

func (s *service) ValidateModeratorRoleJWT(jwtToken string) error {
	token, err := s.getToken(jwtToken)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := claims["role"].(string)

	if ok && token.Valid && userRole == RoleModerator {
		return nil
	}

	return errors.New("invalid moderator token provided")
}

func (s *service) ValidateClientRoleJWT(jwtToken string) error {
	token, err := s.getToken(jwtToken)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := claims["role"].(string)

	if ok && token.Valid && (userRole == RoleClient || userRole == RoleModerator) {
		return nil
	}

	return errors.New("invalid author token provided")
}

func (s *service) getToken(token string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.privateKey, nil
	})
	return jwtToken, err
}

func (s *service) GetUserID(jwtToken string) (string, error) {
	token, err := s.getToken(jwtToken)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token provided")
	}

	return claims["id"].(string), nil
}
