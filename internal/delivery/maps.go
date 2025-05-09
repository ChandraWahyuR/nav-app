package delivery

import (
	"context"
	"io"
	"net/http"
	"proyek1/internal/delivery/middleware"
	"proyek1/internal/model"
	crypto "proyek1/utils"
	jwt "proyek1/utils"
	"proyek1/utils/gmaps"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MapsHandlerInterface interface {
	GmapsSearchbyObject(c *gin.Context)
	GmapsSearchbyList(c *gin.Context)
	GmapsSearchbyPlaceID(c *gin.Context)

	InsertData(c *gin.Context)
	GetTempatPagination(c *gin.Context)
	ProxyPhotoHandler(c *gin.Context)
	RouteDestination(c *gin.Context)
}

type MapsUsecaseInterface interface {
	InsertTempat(ctx context.Context, placeId string) error
	GetTempatPagination(ctx context.Context, limit, page int) ([]model.GetAllTempat, int, error)
	RouteDestination(ctx context.Context, req model.RequestRouteMaps, placeID string) (*model.ResponseRouteMaps, error)
}
type MapsHandler struct {
	jwt   jwt.JWTInterface
	gmaps gmaps.GmapsInterface
	us    MapsUsecaseInterface
}

func NewMapsHandler(jwt jwt.JWTInterface, gmaps gmaps.GmapsInterface, us MapsUsecaseInterface) MapsHandler {
	return MapsHandler{
		jwt:   jwt,
		gmaps: gmaps,
		us:    us,
	}
}

func (h *MapsHandler) GmapsSearchbyObject(c *gin.Context) {
	dataToken, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	if crypto.IsUser(dataToken.Role) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized, can't identity user"})
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
	dataToken, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	if crypto.IsUser(dataToken.Role) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized, can't identity user"})
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
	dataToken, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	if crypto.IsUser(dataToken.Role) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized, can't identity user"})
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

func (h *MapsHandler) InsertData(c *gin.Context) {
	dataToken, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	if !crypto.IsAdmin(dataToken.Role) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized, hanya admin yang boleh akses halaman ini"})
		return
	}

	placeId := c.Param("id")
	if placeId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "placeId tidak boleh kosong"})
		return
	}

	ctx := c.Request.Context()
	err := h.us.InsertTempat(ctx, placeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal memproses maps", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "berhasil",
	})
}

func (h *MapsHandler) GetTempatPagination(c *gin.Context) {
	dataToken, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	if crypto.IsUser(dataToken.Role) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized, can't identity user"})
		return
	}

	p := c.Param("page")
	if p == "" {
		p = "1"
	}
	page, _ := strconv.Atoi(p)

	ctx := c.Request.Context()
	res, pageTotal, err := h.us.GetTempatPagination(ctx, 5, int(page))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal memproses maps", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "berhasil",
		"total_pages": pageTotal,
		"data":        res,
	})
}

// Proxy
func (h *MapsHandler) ProxyPhotoHandler(c *gin.Context) {
	photoRef := c.Query("ref")
	if photoRef == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo_reference is required"})
		return
	}

	photoURL, err := h.gmaps.PhotoReference(photoRef)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := &http.Client{Timeout: 10 * time.Second} // timeout
	resp, err := client.Get(photoURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch photo from Google"})
		return
	}
	defer resp.Body.Close()

	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	c.Status(http.StatusOK)

	io.Copy(c.Writer, resp.Body)
}

func (h *MapsHandler) RouteDestination(c *gin.Context) {
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

	ctx := c.Request.Context()

	var req model.RequestRouteMaps
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	modelReq := model.RequestRouteMaps{
		Origin: model.Waypoint{
			Location: model.LocationReq{
				LatLng: model.LatLng{
					Latitude:  req.Origin.Location.LatLng.Latitude,
					Longitude: req.Origin.Location.LatLng.Longitude,
				},
			},
		}, Destination: model.Waypoint{
			Location: model.LocationReq{
				LatLng: model.LatLng{
					Latitude:  req.Destination.Location.LatLng.Latitude,
					Longitude: req.Destination.Location.LatLng.Longitude,
				},
			},
		},
		TravelMode: req.TravelMode,
	}
	data, err := h.us.RouteDestination(ctx, modelReq, placeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "berhasil", "data": data})
}
