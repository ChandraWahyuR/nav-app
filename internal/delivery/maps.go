package delivery

import (
	"context"
	"io"
	"net/http"
	"proyek1/constant"
	"proyek1/internal/delivery/middleware"
	"proyek1/internal/model"
	"proyek1/utils"
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
	GetDetailTempat(c *gin.Context)
}

type MapsUsecaseInterface interface {
	InsertTempat(ctx context.Context, placeId string) error
	GetTempatPagination(ctx context.Context, name string, limit, page int) ([]model.GetAllTempat, int, error)
	RouteDestination(ctx context.Context, req model.RequestRouteMaps, placeID string) (*model.ResponseRouteMaps, error)
	GetDetailTempat(ctx context.Context, id string) (model.GetDetailTempat, error)
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
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler(constant.StatusFail, "Unauthorize", nil))
		return
	}

	if !crypto.IsUser(dataToken.Role) {
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler(constant.StatusFail, "error unknown token data", nil))
		return
	}

	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler(constant.StatusFail, "query tidak boleh kosong", nil))
		return
	}

	// Panggil langsung service/fungsi GmapsAdd
	results, err := h.gmaps.GmapsSearchObject(query)
	if err != nil {
		c.JSON(utils.ConverResponse(err), utils.ResponseHandler(constant.StatusFail, err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.ResponseHandler(constant.StatusSuccess, "Berhasil mendapatkan data", results))
}

func (h *MapsHandler) GmapsSearchbyList(c *gin.Context) {
	dataToken, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler(constant.StatusFail, "Unauthorize", nil))
		return
	}

	if !crypto.IsUser(dataToken.Role) {
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler(constant.StatusFail, "error unknown token data", nil))
		return
	}

	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler(constant.StatusFail, "query tidak boleh kosong", nil))
		return
	}

	results, err := h.gmaps.GmapsSearchList(query)
	if err != nil {
		c.JSON(utils.ConverResponse(err), utils.ResponseHandler(constant.StatusFail, err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.ResponseHandler(constant.StatusSuccess, "Berhasil mendapatkan data", results))
}

func (h *MapsHandler) GmapsSearchbyPlaceID(c *gin.Context) {
	dataToken, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler(constant.StatusFail, "Unauthorize", nil))
		return
	}

	if !crypto.IsUser(dataToken.Role) {
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler(constant.StatusFail, "error unknown token data", nil))
		return
	}

	placeId := c.Param("id")
	if placeId == "" {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler(constant.StatusFail, "id kosong", nil))
		return
	}

	results, err := h.gmaps.GmapsSearchByPlaceID(placeId)
	if err != nil {
		c.JSON(utils.ConverResponse(err), utils.ResponseHandler(constant.StatusFail, err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.ResponseHandler(constant.StatusSuccess, "Berhasil mendapatkan data", results))
}

func (h *MapsHandler) InsertData(c *gin.Context) {
	dataToken, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler(constant.StatusFail, "Unauthorize", nil))
		return
	}

	if !crypto.IsAdmin(dataToken.Role) {
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler(constant.StatusFail, "error unknown token data", nil))
		return
	}

	placeId := c.Param("id")
	if placeId == "" {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler(constant.StatusFail, "id kosong", nil))
		return
	}

	ctx := c.Request.Context()
	err := h.us.InsertTempat(ctx, placeId)
	if err != nil {
		c.JSON(utils.ConverResponse(err), utils.ResponseHandler(constant.StatusFail, err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.ResponseHandler(constant.StatusSuccess, "Berhasil menambahkan data", nil))
}

func (h *MapsHandler) GetTempatPagination(c *gin.Context) {
	dataToken, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler(constant.StatusFail, "Unauthorize", nil))
		return
	}

	if !crypto.IsUser(dataToken.Role) {
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler(constant.StatusFail, "error unknown token data", nil))
		return
	}

	p := c.Query("page")
	if p == "" {
		p = "1"
	}
	page, err := strconv.Atoi(p)
	if err != nil || page <= 0 {
		page = 1
	}

	n := c.Query("search")
	ctx := c.Request.Context()
	res, pageTotal, err := h.us.GetTempatPagination(ctx, n, 5, int(page))
	if err != nil {
		c.JSON(utils.ConverResponse(err), utils.ResponseHandler(constant.StatusFail, err.Error(), nil))
		return
	}

	metadata := map[string]int{
		"totalPage": pageTotal,
		"page":      page,
	}
	c.JSON(http.StatusOK, utils.MetadataFormatResponse(constant.StatusSuccess, "Berhasil mendapatkan data", metadata, res))
}

func (h *MapsHandler) GetDetailTempat(c *gin.Context) {
	dataToken, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler(constant.StatusFail, "Unauthorize", nil))
		return
	}

	if !crypto.IsUser(dataToken.Role) {
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler(constant.StatusFail, "error unknown token data", nil))
		return
	}
	ctx := c.Request.Context()

	id := c.Param("id")
	data, err := h.us.GetDetailTempat(ctx, id)
	if err != nil {
		c.JSON(utils.ConverResponse(err), utils.ResponseHandler(constant.StatusFail, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseHandler(constant.StatusSuccess, "Berhasil mendapatkan data", data))
}

// Proxy
func (h *MapsHandler) ProxyPhotoHandler(c *gin.Context) {
	photoRef := c.Query("ref")
	if photoRef == "" {
		c.JSON(http.StatusBadRequest, utils.ResponseHandler(constant.StatusFail, "photo ref kosong", nil))
		return
	}

	photoURL, err := h.gmaps.PhotoReference(photoRef)
	if err != nil {
		c.JSON(utils.ConverResponse(err), utils.ResponseHandler(constant.StatusFail, err.Error(), nil))
		return
	}

	client := &http.Client{Timeout: 10 * time.Second} // timeout
	resp, err := client.Get(photoURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, utils.ResponseHandler(constant.StatusFail, "error terjadi kesalahan mengambil gambar", nil))
		return
	}
	defer resp.Body.Close()

	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	c.Status(http.StatusOK)

	io.Copy(c.Writer, resp.Body)
}

func (h *MapsHandler) RouteDestination(c *gin.Context) {
	dataToken, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler(constant.StatusFail, "Unauthorize", nil))
		return
	}

	if !crypto.IsUser(dataToken.Role) {
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler(constant.StatusFail, "error unknown token data", nil))
		return
	}

	placeId := c.Param("id")
	if placeId == "" {
		c.JSON(http.StatusUnauthorized, utils.ResponseHandler(constant.StatusFail, "parameter id tempat tidak boleh kosong", nil))
		return
	}

	ctx := c.Request.Context()

	var req model.RequestRouteMaps
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.ResponseHandler(constant.StatusFail, "error memproses data", nil))
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
		c.JSON(utils.ConverResponse(err), utils.ResponseHandler(constant.StatusFail, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseHandler(constant.StatusSuccess, "Berhasil mendapatkan data", data))
}
