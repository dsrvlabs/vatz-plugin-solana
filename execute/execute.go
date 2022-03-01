package execute

import (
	"fmt"
	"log"
	"github.com/go-resty/resty/v2"
)

/* RPC Call Payload
type Payload struct {
	Jsonrpc string 'json:"jsonrpc"'
	ID int 'json:"id"'
	Method string 'json:"method"'
}
*/

// Solana Cluster : Mainnet, Testnet, Devnet
type Cluster struct {
	name string
	pubkey string
	URL string
	Port int
}


var (
	Mainnet = Cluster{name: "Mainnet", pubkey: "xxxxx", URL:"221.149.114.193" ,Port: 8899}
	LidoSolana = Cluster{name: "Lido-Solana", pubkey: "xxxxx", URL:"142.132.201.181", Port: 8899}
	Testnet = Cluster{name: "Testnet", pubkey: "xxxxx", URL:"162.55.243.228" ,Port: 8899}
)

func GetHealth(network Cluster) (string, error) {

	client := resty.New()
	data := fmt.Sprint("http://", network.URL, ":", network.Port, "/health")
	log.Println("Request Info: ", network.name)
	log.Println("  data: ", data)
	log.Println()

	resp, err := client.R().Get(data)
	if err != nil {
		log.Fatalf("failed to get response: %v", err)
		return "", err
	}
	log.Println("Response Info: ", resp)

	return resp.String(), err
}

/*
func getHealth(network Cluster) (result string) {

	data := Payload{
		Jsonrpc: "2.0",
		ID: 1,
		Method: "getHealth"
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("failed to marshalling payload: %v", err)
		return err
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "162.55.243.228:8899", body)
	if err != nil {
		log.Fatalf("failed to create new request:  %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to get response: %v", err)
		return err
	}
	defer resp.Body.Close()

	return result
}
*/
