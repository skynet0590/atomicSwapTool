package admin

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/skynet0590/atomicSwapTool/dex/msgjson"
	"github.com/skynet0590/atomicSwapTool/server/account"
	"io/ioutil"
	"net/http"
)

const (
	pongStr   = "pong"
	maxUInt16 = int(^uint16(0))
)

// writeJSON marshals the provided interface and writes the bytes to the
// ResponseWriter. The response code is assumed to be StatusOK.
func writeJSON(w http.ResponseWriter, thing interface{}) {
	writeJSONWithStatus(w, thing, http.StatusOK)
}

// writeJSON marshals the provided interface and writes the bytes to the
// ResponseWriter with the specified response code.
func writeJSONWithStatus(w http.ResponseWriter, thing interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, err := json.MarshalIndent(thing, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("JSON encode error: %v", err)
		return
	}
	w.WriteHeader(code)
	_, err = w.Write(append(b, byte('\n')))
	if err != nil {
		log.Errorf("Write error: %v", err)
	}
}

// apiPing is the handler for the '/ping' API request.
func apiPing(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, pongStr)
}

// apiConfig is the handler for the '/config' API request.
func (s *Server) apiConfig(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, s.core.ConfigMsg())
}

// apiAccounts is the handler for the '/accounts' API request.
func (s *Server) apiAccounts(w http.ResponseWriter, _ *http.Request) {
	accts, err := s.core.Accounts()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve accounts: %v", err), http.StatusInternalServerError)
		return
	}
	writeJSON(w, accts)
}

// apiAccountInfo is the handler for the '/account/{account id}' API request.
func (s *Server) apiAccountInfo(w http.ResponseWriter, r *http.Request) {
	acctIDStr := chi.URLParam(r, accountIDKey)
	acctIDSlice, err := hex.DecodeString(acctIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not decode accout id: %v", err), http.StatusBadRequest)
		return
	}
	if len(acctIDSlice) != account.HashSize {
		http.Error(w, "account id has incorrect length", http.StatusBadRequest)
		return
	}
	var acctID account.AccountID
	copy(acctID[:], acctIDSlice)
	acctInfo, err := s.core.AccountInfo(acctID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve account: %v", err), http.StatusInternalServerError)
		return
	}
	writeJSON(w, acctInfo)
}

// apiNotify is the handler for the '/account/{accountID}/notify' API request.
func (s *Server) apiNotify(w http.ResponseWriter, r *http.Request) {
	acctIDStr := chi.URLParam(r, accountIDKey)
	acctID, err := decodeAcctID(acctIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msg, errCode, err := toNote(r)
	if err != nil {
		http.Error(w, err.Error(), errCode)
		return
	}
	s.core.Notify(acctID, msg)
	w.WriteHeader(http.StatusOK)
}

// apiNotifyAll is the handler for the '/notifyall' API request.
func (s *Server) apiNotifyAll(w http.ResponseWriter, r *http.Request) {
	msg, errCode, err := toNote(r)
	if err != nil {
		http.Error(w, err.Error(), errCode)
		return
	}
	s.core.NotifyAll(msg)
	w.WriteHeader(http.StatusOK)
}

// decodeAcctID checks a string as being both hex and the right length and
// returns its bytes encoded as an account.AccountID.
func decodeAcctID(acctIDStr string) (account.AccountID, error) {
	var acctID account.AccountID
	if len(acctIDStr) != account.HashSize*2 {
		return acctID, errors.New("account id has incorrect length")
	}
	if _, err := hex.Decode(acctID[:], []byte(acctIDStr)); err != nil {
		return acctID, fmt.Errorf("could not decode account id: %w", err)
	}
	return acctID, nil
}

func toNote(r *http.Request) (*msgjson.Message, int, error) {
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("unable to read request body: %w", err)
	}
	if len(body) == 0 {
		return nil, http.StatusBadRequest, errors.New("no message to broadcast")
	}
	// Remove trailing newline if present. A newline is added by the curl
	// command when sending from file.
	if body[len(body)-1] == '\n' {
		body = body[:len(body)-1]
	}
	if len(body) > maxUInt16 {
		return nil, http.StatusBadRequest, fmt.Errorf("cannot send messages larger than %d bytes", maxUInt16)
	}
	msg, err := msgjson.NewNotification(msgjson.NotifyRoute, string(body))
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("unable to create notification: %w", err)
	}
	return msg, 0, nil
}
