# How to start

* [Install](http://solidity.readthedocs.io/en/develop/installing-solidity.html)
Solidity compiler, required to generate Go bindings for .sol files.
**IMPORTANT**: DO NOT use npm to install solc. Use apt-get/Homebrew instead.
* [Install 'dep'](https://github.com/golang/dep)
* 'dep ensure' all the dependencies
* [Install 'abigen' tool.](https://github.com/ethereum/go-ethereum/wiki/Native-DApps:-Go-bindings-to-Ethereum-contracts)
 **IMPORTANT**: run 'go install' from vendor folder, not from %GOPATH/src
* Generate the bindings using 'abigen' by running 'go generate' command
* Fill config.yaml file with according info
* Start the server by running 'go run main.go'
* Alternatively run 'go build main.go' and run produced executable

# How to use

* Interaction with server is performed by utilizing the API
* The endpoint is:
  * **GET** _/exchange/:ethereum_address_ - Send the request to server with ethereum address to which tokens will be sent.
  The response is Bitcoin address to which bitcoins must be sent to buy tokens.

# Transaction statuses

-1 = An error occured. Please look at the 'error' column for details<br/>
0 = A purchase was requested, but the funds haven't arrived yet<br/>
1 = User has successfully purchased tokens using BTC
