package api

import (
	"net/http"
	"time"

	db "bank_system/database"
	"bank_system/util"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
	}
}

func (server *Server) createUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := server.database.Where("username = ?", req.Username).Or("email = ?", req.Email).First(&req).Error; err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username or email already in use"})
	}
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	user := db.User{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.Fullname,
		Email:          req.Email,
		CreatedAt:      time.Now(),
	}
	if err := server.database.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	res := newUserResponse(user)
	c.JSON(http.StatusOK, res)

}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
}

type loginUserResponse struct {
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  User      `json:"user"`
}

// response內不能含有password否則會回傳洩漏密碼

func (server *Server) loginUser(c *gin.Context) {
	var req loginUserRequest
	// 使用一個login user request並且使用c.ShouldBindJSON將request放到此架構中
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var user User
	// 再使用一個user負責裝database table中的資料，使用server.database.Where().First(&user)抓取在users中第一筆符合的資料並且放入user中。
	// 切記，使用gorm時，必須注意struct的名稱以及element的名稱要與table與index名稱相符。
	if err := server.database.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
	}
	// 檢查兩個password是否相同 -> 驗證密碼
	if err := util.CheckPassword(user.HashedPassword, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.Token.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.Token.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// Todo -> 加入 Session到rsp中
	rsp := loginUserResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  user,
	}

	c.JSON(http.StatusOK, rsp)
}
