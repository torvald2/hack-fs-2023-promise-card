package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/torvald2/hack-fs-2023-promise-card/usecases"
)

type UserController struct {
	key           string
	url           string
	namespace     string
	tlUrl         string
	tlHash        string
	pinataKey     string
	ensRootOwner  string
	ensPrivateKey string
	mainEns       string
	rpcUrl        string
	resolver      string
}

type CreateUserRequest struct {
	Nick          string `json:"nick"`
	AvalibleAfter int    `json:"valible_after_hours"`
	Avatar        string `json:"avatar"`
}
type CreateUserResponse struct {
	PublicKey           string `json:"public_key"`
	PrivateKeyEncrypted string `json:"private_key_encrypted"`
}

func (u UserController) CreateUser(c *gin.Context) {
	var body CreateUserRequest
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"err": err.Error()})
		return
	}
	duration := time.Duration(body.AvalibleAfter) * time.Hour
	us := usecases.NewCreateUserUseCase(u.key, u.url, u.namespace, u.tlUrl, u.tlHash, u.pinataKey, u.ensRootOwner, u.ensPrivateKey, u.mainEns, u.rpcUrl, u.resolver)
	err := us.Execute(body.Nick, duration, body.Avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, CreateUserResponse{
		PublicKey:           us.GetAddress(),
		PrivateKeyEncrypted: us.GetPrivateKey(),
	})

}

func NewRouter(key, url, namespace, tlUrl, tlHash, pinataKey, ensRootOwner, ensPrivateKey, mainEns, rpcUrl, resolver string) *gin.Engine {
	usrController := UserController{
		key:           key,
		url:           url,
		namespace:     namespace,
		tlUrl:         tlUrl,
		tlHash:        tlHash,
		pinataKey:     pinataKey,
		ensRootOwner:  ensRootOwner,
		ensPrivateKey: ensPrivateKey,
		mainEns:       mainEns,
		rpcUrl:        rpcUrl,
		resolver:      resolver,
	}

	r := gin.New()

	r.POST("/users", usrController.CreateUser)
	return r
}
