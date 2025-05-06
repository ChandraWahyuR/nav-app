package delivery

import (
	"io"
	"net/http"
	"proyek1/internal/delivery/middleware"
	jwt "proyek1/utils"
	"proyek1/utils/gmaps"
	"time"

	"github.com/gin-gonic/gin"
)

type MapsHandlerInterface interface {
	GmapsSearchbyObject(c *gin.Context)
	GmapsSearchbyList(c *gin.Context)
	GmapsSearchbyPlaceID(c *gin.Context)
}

type MapsHandler struct {
	jwt   jwt.JWTInterface
	gmaps gmaps.GmapsInterface
}

func NewMapsHandler(jwt jwt.JWTInterface, gmaps gmaps.GmapsInterface) MapsHandler {
	return MapsHandler{
		jwt:   jwt,
		gmaps: gmaps,
	}
}

func (h *MapsHandler) GmapsSearchbyObject(c *gin.Context) {
	_, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "query tidak boleh kosong"})
		return
	}

	// Panggil langsung service/fungsi GmapsAdd
	results, err := h.gmaps.GmapsSearchObject(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal memproses maps", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "berhasil",
		"data":    results,
	})
}

func (h *MapsHandler) GmapsSearchbyList(c *gin.Context) {
	_, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "query tidak boleh kosong"})
		return
	}

	results, err := h.gmaps.GmapsSearchList(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal memproses maps", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "berhasil",
		"data":    results,
	})
}

func (h *MapsHandler) GmapsSearchbyPlaceID(c *gin.Context) {
	_, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	placeId := c.Param("id")
	if placeId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "placeId tidak boleh kosong"})
		return
	}

	results, err := h.gmaps.GmapsSearchByPlaceID(placeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal memproses maps", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "berhasil",
		"data":    results,
	})
}

// Proxy
func (c *MapsHandler) ProxyPhotoHandler(ctx *gin.Context) {
	photoRef := ctx.Query("ref")
	if photoRef == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "photo_reference is required"})
		return
	}

	photoURL, err := c.gmaps.PhotoReference(photoRef)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := &http.Client{Timeout: 10 * time.Second} // timeout
	resp, err := client.Get(photoURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch photo from Google"})
		return
	}
	defer resp.Body.Close()

	ctx.Header("Content-Type", resp.Header.Get("Content-Type"))
	ctx.Status(http.StatusOK)

	io.Copy(ctx.Writer, resp.Body)
}
