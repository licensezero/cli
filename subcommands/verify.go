package subcommands

import (
	"io"
	"licensezero.com/licensezero/api"
	"licensezero.com/licensezero/user"
	"net/http"
)

const verifyDescription = "Verify receipts."

var verifyUsage = verifyDescription + "\n\n" +
	"Usage:\n" +
	"  licensezero verify\n"

// Verify saves receipts to disk.
var Verify = &Subcommand{
	Tag:         "buyer",
	Description: verifyDescription,
	Handler: func(args []string, stdin InputDevice, stdout, stderr io.StringWriter, client *http.Client) int {
		receipts, receiptErrors, err := user.ReadReceipts()
		if err != nil {
			stderr.WriteString("Error reading receipts: " + err.Error())
			return 1
		}
		foundError := false
		for _, receiptError := range receiptErrors {
			foundError = true
			stderr.WriteString(receiptError.Error() + "\n")
		}
		servers := make(map[string]*api.BrokerServer)
		registers := make(map[string]*api.Register)
		for _, receipt := range receipts {
			err = receipt.Validate()
			if err != nil {
				foundError = true
				stderr.WriteString(
					"Receipt for " +
						receipt.License.Values.API + "/orders/" + receipt.License.Values.OfferID +
						"is not a valid receipt.\n",
				)
			}
			err = receipt.VerifySignature()
			if err != nil {
				foundError = true
				stderr.WriteString(
					"Signature for " +
						receipt.License.Values.API + "/orders/" + receipt.License.Values.OfferID +
						"is not valid.\n",
				)
			}
			brokerAPI := receipt.License.Values.API
			server, ok := servers[brokerAPI]
			if !ok {
				server = &api.BrokerServer{Client: client, Base: brokerAPI}
				servers[brokerAPI] = server
			}
			register, ok := registers[brokerAPI]
			if !ok {
				register, err := server.Register()
				if err != nil {
					foundError = true
					stderr.WriteString(
						"Could not fetch key register for " +
							brokerAPI + ".\n",
					)
				}
				registers[brokerAPI] = register
			}
			if err = register.ValidateEffectiveDate(receipt); err != nil {
				foundError = true
				stderr.WriteString(
					"Signature for " +
						receipt.License.Values.API + "/orders/" + receipt.License.Values.OfferID +
						"does not match time frame for the broker's signing key.\n",
				)
			}
		}
		if foundError {
			return 1
		}
		return 0
	},
}
