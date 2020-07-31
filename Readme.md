# Double Ratchet React Test

This repository is a simple test to connect a React Application to a Go Server and have both communicate using Double ratchet encryption.  

The Go server can be found inside [_'go-server'_](./go-server) directory and the react client inside the [_'react-client'_](./react-client) directory  

This is just a demo application using the prebuilt libraries available for both the languages/frameworks.

## Communication:
The main communication inside the application is:
- On start, the server generates a Diffie–Hellman key pair and a user for the double ratchet communication and stores it inside a serve state manager.
- On button press in the react application, it sends a GET request to the server to get a Public Key for the server. _(/start route)_
- The server responds with a Public key, then the server generates a shared key for communication using the elliptic curve (22519).
```json
{
    "key": "Public-key from server"
}
```
- It encrypts the plaintext and sends a POST request with the message and the key. _(/recv route)_
```json
{
     "key": "sharedKey",
     "message": "cipherText"
```
- The server responds with "Received" to confirm the communication and then tries to decrypt the message using the doubleratchet library.


## Libraries
#### GO
  - github.com/gorilla/mux 
      - The main router library for the server
  - github.com/gorilla/handlers
      - For handling CORS inside the router
  - github.com/status-im/doubleratchet
      - For the double ratchet encryption communication
#### React
  - Parcel
      - Build system and bundler
  - React/React DOM
      - Build the front-end application
  - Crypto-JS
      - For encryption/decryption
  - Elliptic
      - For generating the keys and users on the client side
  - Axios
      - Used for API communication

## Usage
#### Server
##### Local Machine
- git clone <project>
- cd go-server
- go build 
- .\go-server\go run main.go
#####  Docker
Build Go Server:
- docker build -t go-server -f .\go-server\dockerfile .\go-server
Run Go Server:
- docker run -it -p 8080:8080 --rm go-server

#### Client
- cd react-server
- npm install
- npm start


## Limitations:
- The Go library is properly implemented and there are no issues besides some lack of documentation.
- I could not find any good libraries for the React (browser) clients and thus theres some incorrect keys being shared.

##### Disclaimer
This is just a test application never meant for production environment. This does not work and all the things are only for development level trials.
