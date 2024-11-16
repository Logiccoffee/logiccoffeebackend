package controller

import (
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/report"
	"github.com/gocroot/helper/whatsauth"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
)

func GetYesterdayDistincWAGroup(respw http.ResponseWriter, req *http.Request) {
	var resp model.Response
	filter := bson.M{"_id": report.YesterdayFilter()}
	wagroupidlist, err := atdb.GetAllDistinctDoc(config.Mongoconn, filter, "project.wagroupid", "pushrepo")
	if err != nil {
		resp.Info = "Gagal Query Distincs project.wagroupid"
		resp.Response = err.Error()
		at.WriteJSON(respw, http.StatusUnauthorized, resp)
		return
	}
	for _, wagroupid := range wagroupidlist {
		// Type assertion to convert any to string
		groupID, ok := wagroupid.(string)
		if !ok {
			resp.Info = "wagroupid is not a string"
			resp.Response = "wagroupid is not a string"
			at.WriteJSON(respw, http.StatusUnauthorized, resp)
			return
		}
		//kirim report ke group
		dt := &whatsauth.TextMessage{
			To:       groupID,
			IsGroup:  true,
			Messages: report.GetDataRepoMasukHariIni(config.Mongoconn, groupID) + "\n" + report.GetDataLaporanMasukHariini(config.Mongoconn, groupID),
		}
		_, resp, err := atapi.PostStructWithToken[model.Response]("Token", config.WAAPIToken, dt, config.WAAPIMessage)
		if err != nil {
			resp.Info = "Tidak berhak"
			resp.Response = err.Error()
			at.WriteJSON(respw, http.StatusUnauthorized, resp)
			return
		}
	}
	at.WriteJSON(respw, http.StatusOK, resp)
}

func GetReportHariIni(respw http.ResponseWriter, req *http.Request) {
	var resp model.Response
	//kirim report ke group
	dt := &whatsauth.TextMessage{
		To:       "6281313112053-1492882006",
		IsGroup:  true,
		Messages: report.GetDataRepoMasukHarian(config.Mongoconn) + "\n" + report.GetDataLaporanMasukHarian(config.Mongoconn),
	}
	_, resp, err := atapi.PostStructWithToken[model.Response]("Token", config.WAAPIToken, dt, config.WAAPIMessage)
	if err != nil {
		resp.Info = "Tidak berhak"
		resp.Response = err.Error()
		at.WriteJSON(respw, http.StatusUnauthorized, resp)
		return
	}
	at.WriteJSON(respw, http.StatusOK, resp)
}
