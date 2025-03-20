package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// JWT 签名密钥，用于加密和验证 Token
// 服务器会用这个密钥 签名 Token，并在验证 Token 时使用相同的密钥来检查其有效性。
var jwtSecret = []byte("Metafora-TodoList-JWT")

// Claims 结构体 是 JWT 负载（Payload），用于存储用户信息和 JWT 相关声明。
// 用于签发用户 token
type Claims struct {
	// 用户 ID（类型 uint）。
	Id uint `json:"id"`
	// 用户名（类型 string）。
	UserName  string `json:"username"`
	Authority int    `json:"authority"`
	// jwt.RegisteredClaims：ExpiresAt：Token 过期时间。 Issuer：签发者（例如 Metafora）。 其他字段如 Subject、IssuedAt 等可以根据需求使用。
	jwt.RegisteredClaims
}

// GenerateToken 签发用户 token
func GenerateToken(id uint, username string, authority int) (string, error) {
	nowTime := time.Now()
	// 设置 Token 过期时间（24 小时后过期）。
	expireTime := nowTime.Add(24 * time.Hour)

	claims := Claims{
		Id:        id,
		UserName:  username,
		Authority: authority,
		// ExpiresAt: jwt.NewNumericDate(expireTime) 将 time.Time 转换为 *jwt.NumericDate（符合 jwt/v4 的要求）。
		// Issuer: 说明 Token 是由 "Metafora" 这个应用生成的
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "Metafora",
		},
	}

	// 创建一个新的 JWT Token。SigningMethodHS256：使用 HMAC-SHA256 作为签名算法。 claims：将 Claims 结构体作为 JWT 的负载（Payload）。
	// SignedString(jwtSecret): 使用 jwtSecret 进行 HMAC-SHA256 签名。
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
}

// ParseToken 解析和验证 JWT 令牌（token），如果 JWT 合法且未过期，则返回解析出的 Claims（声明信息），否则返回错误。
func ParseToken(tokenString string) (*Claims, error) { // tokenString string：JWT 令牌的字符串，通常是客户端传来的 Token。
	// 解析 tokenString 并将 payload 解析到 Claims 结构体。
	// &Claims{}：用于存储解析后的 JWT 负载（Payload）。
	// callbackFunction：用于提供签名密钥（SecretKey）。
	// 服务器会使用 jwtSecret 重新计算签名，并与 JWT 的 Signature 部分进行对比，如果匹配，则证明 Token 未被篡改。
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil { // 确保 Token 解析成功。
		// 尝试将 Claims 结构体转换回 *Claims 类型，以便访问 JWT 里的用户信息。
		// tokenClaims.Valid 检查 Token 是否有效（是否过期、签名是否正确等）。
		// 如果解析成功，返回 claims（JWT 的 Payload 信息）。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	// 如果 JWT 无效或解析失败，返回 nil 和 err 进行错误处理。
	return nil, err
}

// EmailClaims 是用于发送邮箱的 token 负载
// 用于签发邮箱 token
type EmailClaims struct {
	UserID        uint   `json:"user_id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	OperationType uint   `json:"operation_type"`
	jwt.RegisteredClaims
}

// GenerateEmailToken 签发邮箱 token
func GenerateEmailToken(userId, operation uint, email, password string) (string, error) {
	nowTime := time.Now()
	// 设置 Token 过期时间（24 小时后过期）。
	expireTime := nowTime.Add(24 * time.Hour)

	claims := EmailClaims{
		UserID:        userId,
		Email:         email,
		Password:      password,
		OperationType: operation,
		// ExpiresAt: jwt.NewNumericDate(expireTime) 将 time.Time 转换为 *jwt.NumericDate（符合 jwt/v4 的要求）。
		// Issuer: 说明 Token 是由 "Metafora" 这个应用生成的
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "Metafora",
		},
	}

	// 创建一个新的 JWT Token。SigningMethodHS256：使用 HMAC-SHA256 作为签名算法。 claims：将 Claims 结构体作为 JWT 的负载（Payload）。
	// SignedString(jwtSecret): 使用 jwtSecret 进行 HMAC-SHA256 签名。
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
}

// ParseEmailToken 解析和验证 JWT 令牌（邮箱token），如果 JWT 合法且未过期，则返回解析出的 Claims（声明信息），否则返回错误。
func ParseEmailToken(tokenString string) (*EmailClaims, error) { // tokenString string：JWT 令牌的字符串，通常是客户端传来的 Token。
	// 解析 tokenString 并将 payload 解析到 Claims 结构体。
	// &Claims{}：用于存储解析后的 JWT 负载（Payload）。
	// callbackFunction：用于提供签名密钥（SecretKey）。
	// 服务器会使用 jwtSecret 重新计算签名，并与 JWT 的 Signature 部分进行对比，如果匹配，则证明 Token 未被篡改。
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &EmailClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil { // 确保 Token 解析成功。
		// 尝试将 Claims 结构体转换回 *Claims 类型，以便访问 JWT 里的用户信息。
		// tokenClaims.Valid 检查 Token 是否有效（是否过期、签名是否正确等）。
		// 如果解析成功，返回 claims（JWT 的 Payload 信息）。
		if claims, ok := tokenClaims.Claims.(*EmailClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	// 如果 JWT 无效或解析失败，返回 nil 和 err 进行错误处理。
	return nil, err
}
