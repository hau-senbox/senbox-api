package usecase

import (
	"fmt"
	"io/ioutil"
	"sen-global-api/helper"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func GetLocationListDao(c *gin.Context, req request.LocationSheetRequest) ([]response.ListLocationResponse, error) {
	data, err := ioutil.ReadFile(helper.CREDENTIALS)
	var resp []response.ListLocationResponse
	if err != nil {
		return resp, err
	}
	conf, err := google.JWTConfigFromJSON(data, helper.Scope)
	if err != nil {
		return resp, err
	}
	client := conf.Client(c)

	sheetsService, err := sheets.New(client)
	if err != nil {
		return resp, err
	}
	resps, err := sheetsService.Spreadsheets.Values.Get(req.SheetId, "sheet1").ValueRenderOption("FORMATTED_VALUE").Context(c).Do()
	if err != nil {
		return resp, err
	}
	var id []int64
	var locationName []string
	for _, row := range resps.Values {
		t, _ := strconv.Atoi(fmt.Sprintf("%s", row[0]))
		id = append(id, int64(t))
		locationName = append(locationName, fmt.Sprintf("%s", row[1]))
	}
	for _, m := range id {
		for _, n := range locationName {
			if !CheckLocationData(resp, m, n) {
				resp = append(resp, response.ListLocationResponse{
					LocationId:   m,
					LocationName: n,
				})
			}
		}
	}

	return resp, nil
}
func CheckLocationData(a []response.ListLocationResponse, id int64, location string) bool {
	for _, n := range a {
		if n.LocationId == id || n.LocationName == location {
			return true
		}
	}
	return false
}
