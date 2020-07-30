package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/status-im/doubleratchet"
)

type State struct {
	sk      []byte
	user    doubleratchet.Session
	keyPair doubleratchet.DHPair
}

type MessageBody struct {
	Text string `json:"message"`
}

type KeyExchange struct {
	Key string `json:"key"`
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Ready")
}

func (s *State) startComm(w http.ResponseWriter, r *http.Request) {
	k := KeyExchange{
		Key: s.keyPair.PublicKey().String(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(k); err != nil {
		log.Fatal(err)
	}
	fmt.Println(k)

}

func (s *State) receive(w http.ResponseWriter, r *http.Request) {
	var m MessageBody

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("Received"))

	// if s.user.RatchetDecrypt([]byte(m.Text), nil) == "hello world" {
	// 	w.Write([]byte("Experiment Success"))
	// } else {
	// 	w.Write([]byte("Need to try Again"))
	// }

}

func main() {
	shared := []byte{
		0xeb, 0x8, 0x10, 0x7c, 0x33, 0x54, 0x0, 0x20,
		0xe9, 0x4f, 0x6c, 0x84, 0xe4, 0x39, 0x50, 0x5a,
		0x2f, 0x60, 0xbe, 0x81, 0xa, 0x78, 0x8b, 0xeb,
		0x1e, 0x2c, 0x9, 0x8d, 0x4b, 0x4d, 0xc1, 0x40,
	}

	keyPair, err := doubleratchet.DefaultCrypto{}.GenerateDH()
	if err != nil {
		log.Fatal(err)
	}

	bob, err := doubleratchet.New([]byte("bob-session-id"), shared, keyPair, nil)
	if err != nil {
		log.Fatal(err)
	}

	st := &State{
		sk:      shared,
		keyPair: keyPair,
		user:    bob,
	}

	// // Alice MUST be created with the shared secret and Bob's public key.
	// alice, err := doubleratchet.NewWithRemoteKey([]byte("alice-session-id"), sk, keyPair.PublicKey(), nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Alice can now encrypt messages under the Double Ratchet session.
	// m, err := alice.RatchetEncrypt([]byte("Hi Bob!"), nil)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Which Bob can decrypt.
	// plaintext, err := bob.RatchetDecrypt(m, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(plaintext))

	router := mux.NewRouter()

	router.HandleFunc("/ping", ping)
	router.HandleFunc("/start", st.startComm)
	router.HandleFunc("/recv", st.receive)
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(router)))
}
