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

// State handles the state data of the application
type State struct {
	sk      []byte
	user    doubleratchet.Session
	keyPair doubleratchet.DHPair
}

// MessageBody is a single message exchanged by the parties.
type MessageBody struct {
	SenderKey string `json:"key"`
	Text      string `json:"message"`
}

// KeyExchange is a single time public Key exchange by the parties
type KeyExchange struct {
	Key string `json:"key"`
}

// Ping returns a 200(OK) status and returns with "Ready"
func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Ready")
}

// startComm sends the server's public key and to start the encrypted communication
func (s *State) startComm(w http.ResponseWriter, r *http.Request) {
	k := KeyExchange{
		Key: s.keyPair.PublicKey().String(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(k); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Key Exchanged", k)

}

func (s *State) receive(w http.ResponseWriter, r *http.Request) {
	var m MessageBody

	// Recieve the cipher text and key and read into a MessageBody struct
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte("Received"))

	// fmt.Println("Received Text:", m)

	// TODO: Need to generate proper DH key from received key
	messHeader := doubleratchet.MessageHeader{
		DH: []byte(m.SenderKey),
		N:  1,
		PN: 0,
	}
	mess := doubleratchet.Message{
		Header:     messHeader,
		Ciphertext: []byte(m.Text),
	}

	plaintext, err := s.user.RatchetDecrypt(mess, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(plaintext)) // Test plaintext received

	// if s.user.RatchetDecrypt(mss, nil) == "hello world!" {
	// 	w.Write([]byte("Experiment Success"))
	// } else {
	// 	w.Write([]byte("Need to try Again"))
	// }

}

func main() {
	// Preshared secret key
	// Copied from the library
	shared := []byte{
		0xeb, 0x8, 0x10, 0x7c, 0x33, 0x54, 0x0, 0x20,
		0xe9, 0x4f, 0x6c, 0x84, 0xe4, 0x39, 0x50, 0x5a,
		0x2f, 0x60, 0xbe, 0x81, 0xa, 0x78, 0x8b, 0xeb,
		0x1e, 0x2c, 0x9, 0x8d, 0x4b, 0x4d, 0xc1, 0x40,
	}

	// Generate Key pair for the server
	keyPair, err := doubleratchet.DefaultCrypto{}.GenerateDH()
	if err != nil {
		log.Fatal(err)
	}

	// Generate a New user with given pre shared key and generated Key pair
	bob, err := doubleratchet.New([]byte("bob-session-id"), shared, keyPair, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Generate a common state for the server
	st := &State{
		sk:      shared,
		keyPair: keyPair,
		user:    bob,
	}

	// Example from Go double ratchet library
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

	router := mux.NewRouter()

	// Header for CORS handler
	headersOK := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOK := handlers.AllowedOrigins([]string{"*"})
	methodsOK := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"})

	router.HandleFunc("/ping", ping).Methods("GET")
	router.HandleFunc("/start", st.startComm).Methods("GET")
	router.HandleFunc("/recv", st.receive).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headersOK, originsOK, methodsOK)(router)))
}
