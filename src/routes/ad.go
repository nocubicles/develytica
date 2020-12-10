package routes

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/nocubicles/skillbase.io/src/constants"
	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/utils"
)

type SelectionValue struct {
	Value      string
	ID         uint
	IsSelected bool
}

type PostedAd struct {
	Direction          string
	Location           string
	LocationDetailDesc string
	AdType             string
	BikeType           string
	FrameMaterial      string
	Title              string
	Description        string
	FrameSize          string
	FrameSizeDesc      string
	Phone              string
	Images             []string
	Price              uint
	ValidUntil         time.Time
	PostedWhen         time.Time
}

type SelectionValues struct {
	Directions       []SelectionValue
	Locations        []SelectionValue
	AdTypes          []SelectionValue
	BikeTypes        []SelectionValue
	FrameSizes       []SelectionValue
	FrameMaterials   []SelectionValue
	Ad               models.Ad
	Title            string
	Description      string
	FrameSizeDesc    string
	Price            uint
	LocationDesc     string
	Phone            string
	ImageUploadError string
}

func GetAdFormData(ad *models.Ad) SelectionValues {
	adDirections := constants.GetAdDirections()
	adLocations := constants.GetAdLocations()
	adTypes := constants.GetAdTypes()
	bikeTypes := constants.GetBikeTypes()
	frameSizes := constants.GetFrameSizes()
	frameMaterials := constants.GetFrameMaterials()

	var data = SelectionValues{

		Directions:     getSelectionValues(adDirections, ad.AdDirectionID),
		Locations:      getSelectionValues(adLocations, ad.LocationID),
		AdTypes:        getSelectionValues(adTypes, ad.AdTypeID),
		BikeTypes:      getSelectionValues(bikeTypes, ad.BikeTypeID),
		FrameSizes:     getSelectionValues(frameSizes, ad.FrameSizeID),
		FrameMaterials: getSelectionValues(frameMaterials, ad.FrameMaterialID),
		Ad:             *ad,
		Title:          *&ad.Title,
		Description:    *&ad.Description,
		FrameSizeDesc:  *&ad.FrameSizeDescription,
		Price:          *&ad.Price,
		LocationDesc:   *&ad.LocationDescription,
		Phone:          *&ad.PhoneNo,
	}

	return data
}

func RenderAdForm(w http.ResponseWriter, r *http.Request) {
	emptyAd := &models.Ad{
		AdDirectionID:        0,
		AdTypeID:             0,
		BikeTypeID:           0,
		Description:          "",
		FrameMaterialID:      0,
		FrameSizeDescription: "",
		FrameSizeID:          0,
		LocationID:           0,
		LocationDescription:  "",
		Price:                0,
		PhoneNo:              "",
		Title:                "",
		Weight:               0,
		Condition:            false,
	}
	var data = GetAdFormData(emptyAd)

	utils.Render(w, "ad.html", data)
}

func getSelectionValues(input map[uint]string, userSelection uint) []SelectionValue {
	var selectionValues []SelectionValue

	for id, description := range input {
		isSelected := false
		if id == userSelection {
			isSelected = true
		}
		selectionValue := SelectionValue{
			Value:      description,
			ID:         id,
			IsSelected: isSelected,
		}
		selectionValues = append(selectionValues, selectionValue)
	}

	sort.Slice(selectionValues, func(i, j int) bool {
		return selectionValues[i].ID < selectionValues[j].ID
	})

	return selectionValues
}

func RenderAd(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	db := utils.DbConnection()

	ad := &models.Ad{}
	result := db.First(&ad, vars["adID"])

	if result.RowsAffected > 0 {
		adID := ad.ID
		adDirections := constants.GetAdDirections()
		adLocations := constants.GetAdLocations()
		adTypes := constants.GetAdTypes()
		bikeTypes := constants.GetBikeTypes()
		frameSizes := constants.GetFrameSizes()
		frameMaterials := constants.GetFrameMaterials()

		postedAd := &PostedAd{
			AdType:             ad.GetAdValueById(adTypes, ad.AdTypeID),
			BikeType:           ad.GetAdValueById(bikeTypes, ad.BikeTypeID),
			Direction:          ad.GetAdValueById(adDirections, ad.AdDirectionID),
			FrameMaterial:      ad.GetAdValueById(frameMaterials, ad.FrameMaterialID),
			FrameSize:          ad.GetAdValueById(frameSizes, ad.FrameSizeID),
			Location:           ad.GetAdValueById(adLocations, ad.LocationID),
			Description:        ad.Description,
			FrameSizeDesc:      ad.FrameSizeDescription,
			LocationDetailDesc: ad.LocationDescription,
			Phone:              ad.PhoneNo,
			PostedWhen:         ad.CreatedAt,
			ValidUntil:         ad.OpenUntil,
			Price:              ad.Price,
			Title:              ad.Title,
			Images:             utils.GetAdImageUrls(adID),
		}

		utils.Render(w, "ad[adID].html", postedAd)
	} else {
		utils.Render(w, "notfound.html", "")
	}

}

func ReceiveAdForm(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)

	ad := &models.Ad{
		AdDirectionID:        convertStringToUint(r.PostFormValue("direction")),
		AdTypeID:             convertStringToUint(r.PostFormValue("adType")),
		BikeTypeID:           convertStringToUint(r.PostFormValue("bikeType")),
		FrameSizeID:          convertStringToUint(r.PostFormValue("frameSize")),
		FrameSizeDescription: r.PostFormValue("frameSizeDesc"),
		FrameMaterialID:      convertStringToUint(r.PostFormValue("frameMaterial")),
		Title:                r.PostFormValue("adTitle"),
		Price:                convertStringToUint(r.PostFormValue("price")),
		LocationID:           convertStringToUint(r.PostFormValue("location")),
		LocationDescription:  r.PostFormValue("locationDesc"),
		Description:          r.PostFormValue("adDescription"),
		PhoneNo:              r.PostFormValue("phoneNo"),
		Open:                 true,
		OpenUntil:            getOpenUntil(),
		UserID:               userID,
	}

	if ad.Validate() == false {
		var data = GetAdFormData(ad)
		utils.Render(w, "ad.html", data)
	}

	db := utils.DbConnection()
	result := db.Create(&ad)

	if result.RowsAffected > 0 {
		createdAdID := ad.ID
		formData := r.MultipartForm
		files := formData.File["images"]
		if len(files) > 5 {
			var data = GetAdFormData(ad)
			data.ImageUploadError = "Lubatud on laadida maksimum 5 pilti"
			utils.Render(w, "ad.html", data)
			return
		}
		for i := range files {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				fmt.Println(err)
				return
			}

			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err)
				return
			}

			contentType := http.DetectContentType(fileBytes)
			fileName := files[i].Filename
			fileSize := files[i].Size

			if !strings.Contains(contentType, "image") {
				var data = GetAdFormData(ad)
				data.ImageUploadError = "Palun laadige ainult pilte"
				utils.Render(w, "ad.html", data)
				return
			}

			reader := bytes.NewReader(fileBytes)
			utils.UploadImage(createdAdID, reader, contentType, fileName, fileSize)
		}
		utils.GetAdImageUrls(ad.ID)
	}

	RenderHome(w, r)
}

func getOpenUntil() time.Time {
	return time.Now().Add(time.Hour * 14 * 24)
}
