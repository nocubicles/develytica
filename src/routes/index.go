package routes

import (
	"net/http"
	"time"

	"github.com/nocubicles/skillbase.io/src/constants"
	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/utils"
)

type AdData struct {
	ID            uint
	Title         string
	Thumbnail     string
	FrameSize     string
	FrameSizeDesc string
	Location      string
	Added         time.Time
	AdType        string
	Direction     string
	BikeType      string
	Price         uint
}

func getAdsData(ads *[]models.Ad) []AdData {
	adsData := []AdData{}
	adDirections := constants.GetAdDirections()
	adLocations := constants.GetAdLocations()
	adTypes := constants.GetAdTypes()
	bikeTypes := constants.GetBikeTypes()
	frameSizes := constants.GetFrameSizes()

	for _, ad := range *ads {
		adImageUrls := utils.GetAdImageUrls(ad.ID)
		thumbNail := ""
		if len(adImageUrls) > 0 {
			thumbNail = adImageUrls[0]
		}
		adData := AdData{
			ID:            ad.ID,
			Title:         ad.Title,
			Price:         ad.Price,
			Thumbnail:     thumbNail,
			FrameSize:     ad.GetAdValueById(frameSizes, ad.ID),
			Added:         ad.CreatedAt,
			FrameSizeDesc: ad.FrameSizeDescription,
			Location:      ad.GetAdValueById(adLocations, ad.ID),
			AdType:        ad.GetAdValueById(adTypes, ad.ID),
			Direction:     ad.GetAdValueById(adDirections, ad.ID),
			BikeType:      ad.GetAdValueById(bikeTypes, ad.ID),
		}
		adsData = append(adsData, adData)
	}
	return adsData
}

func RenderHome(w http.ResponseWriter, r *http.Request) (int, error) {

	ads := &[]models.Ad{}

	db := utils.DbConnection()

	result := db.Limit(12).Select("title", "ID", "CreatedAt", "price").Find(&ads)
	data := []AdData{}
	if result.RowsAffected > 0 {
		data = getAdsData(ads)
		return http.StatusOK, utils.Render(w, "index.html", data)
	}
	return http.StatusOK, utils.Render(w, "index.html", data)
}
