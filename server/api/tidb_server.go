package api

import (
	"net/http"

	"github.com/coreos/etcd/clientv3"
	log "github.com/pingcap/log"
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

	var items []string

	start := tidbInfoPrefix
	end := clientv3.GetPrefixRangeEnd(tidbInfoPrefix)

	for {
		ret, err := kv.LoadRange(start, end, 10)
		if len(ret) == 0 {
			break
		}
		items = append(items, ret...)
		start = ret[len(ret) - 1].
	}
	log.Info(items)
	if err != nil {
		h.rd.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.rd.JSON(w, http.StatusOK, nil)
}
