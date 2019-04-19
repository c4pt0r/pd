package api

import (
	"net/http"

	"github.com/coreos/etcd/clientv3"
	"github.com/ngaut/log"
	"github.com/pingcap/pd/server"
	"github.com/unrolled/render"
)

type TiDBServerInfo struct {
	Version     string `json:"version,omitempty"`
	GitHash     string `json:"git_hash,omitempty"`
	DDLID       string `json:"ddl_id,omitempty"`
	AdvertiseIP string `json:"ip,omitempty"`
	Port        int    `json:"listening_port,omitempty"`
	StatusPort  int    `json:"status_port,omitempty"`
	DDLLease    string `json:"lease,omitempty"`
}

const (
	tidbInfoPrefix string = "/tidb/server/info"
)

type TiDBServerInfoHandler struct {
	svr *server.Server
	rd  *render.Render
}

func newTiDBServerInfoHandler(svr *server.Server, rd *render.Render) *TiDBServerInfoHandler {
	return &TiDBServerInfoHandler{
		svr: svr,
		rd:  rd,
	}
}

func (h *TiDBServerInfoHandler) Get(w http.ResponseWriter, r *http.Request) {
	kv := h.svr.GetStorage()

	var items, keys, vals []string
	var err error

	start := tidbInfoPrefix
	end := clientv3.GetPrefixRangeEnd(tidbInfoPrefix)

	batchSize := 10
	for {
		keys, vals, err = kv.LoadRange(start, end, batchSize)
		if len(vals) > 0 {
			items = append(items, vals...)
		}
		if len(vals) < batchSize {
			break
		}
		start = keys[len(keys)-1]
	}
	if err != nil {
		h.rd.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	// REMOVE ME
	for _, v := range items {
		log.Info(v)
	}

	h.rd.JSON(w, http.StatusOK, nil)
}
