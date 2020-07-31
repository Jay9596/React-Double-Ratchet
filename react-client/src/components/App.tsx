import * as React from 'react'
import * as axios from "axios"
import * as CryptoJS from "crypto-js"

class TestButton extends React.Component {

    constructor(props) {
        super(props)
        this.state = {}

        this.apiStart = this.apiStart.bind(this)
    }


    async apiStart(event: MouseEvent) {
        event.preventDefault()
        
        // Generate the elliptic Curve object and generate the keys for the client.
        let EC = require('elliptic').ec;
        let ec = new EC('curve25519');
        let user = ec.genKeyPair();
        
        // Call the /start method on server to receive the Public key for the server
        try {
            let res = await axios.get(`http://localhost:8080/start`)

            // Append the public key to Component state
            this.setState({ key: res.data.key })
            
            
            
        } catch (err) {
            console.log('Error in getting Public key')
            console.error(err)
            return
        }

        // Generate the Elliptic Key from the server Public Key
        let pubKey = ec.keyFromPublic (this.state.key) 

        // Use the server key to generate a shared Key
        let sharedKey1 = user.derive(pubKey.getPublic());
        let keyString = sharedKey1.toString(16)
        // console.log("SharedKey:", keyString) // Test sharedKey

        // Encrypt a simple message using Crypto AES with the generated shared key.
        let cipherText = CryptoJS.AES.encrypt('hello world!', keyString).toString();
        // console.log(cipherText) // Test cipher Text
        
        
        // Send the cipherText to the server along with the sharedKey
        try {
            let res = await axios.post(`http://localhost:8080/recv`, {
                "key": sharedKey1,
                "message": cipherText
            })
            // console.log(res.data)// Test Received data from server
            // res.data == "Received"
            
        } catch (error) {
            console.log("Error in sending Cipher Text")
            console.error(error)
            return
        }
        
        
    }
    
    
    render() {
        return (<button onClick={this.apiStart}>Test</button>)
    }
}

// Basic react Application
export const App = () => (
    <div>
        <h1>React Client</h1>
        <TestButton />
    </div>
);
