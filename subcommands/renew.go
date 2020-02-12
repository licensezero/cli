package subcommands

import (
	"io"
	"licensezero.com/licensezero/api"
	"licensezero.com/licensezero/user"
	"net/http"
)

const renewDescription = "Renew receipts."

var renewUsage = renewDescription + "\n\n" +
	"Usage:\n" +
	"  licensezero renew\n"

// Renew downloads the lastest receipts for recurring licenses.
var Renew = &Subcommand{
	Tag:         "buyer",
	Description: renewDescription,
	Handler: func(args []string, stdin InputDevice, stdout, stderr io.StringWriter, client *http.Client) int {
		receipts, _, err := user.ReadReceipts()
		if err != nil {
			stderr.WriteString("Error reading receipts: " + err.Error())
			return 1
		}
		foundError := false
		for _, receipt := range receipts {
			if !receipt.License.Values.Recurring {
				continue
			}
			brokerServer := api.BrokerServer{
				Client: client,
				Base:   receipt.License.Values.Server,
			}
			latest, err := brokerServer.Latest(receipt.License.Values.OrderID)
			if err != nil {
				foundError = true
				stderr.WriteString(
					receipt.License.Values.Server +
						" did not return a new receipt for offer " +
						receipt.License.Values.OfferID + "\n",
				)
				continue
			}
			if err = latest.Validate(); err != nil {
				foundError = true
				stderr.WriteString(
					receipt.License.Values.Server +
						" returned an invalid receipt for offer " +
						receipt.License.Values.OfferID + "\n",
				)
				continue
			}
			if err = latest.VerifySignature(); err != nil {
				foundError = true
				stderr.WriteString(
					receipt.License.Values.Server +
						" returned a receipt with an invalid signature for offer " +
						receipt.License.Values.OfferID + "\n",
				)
				continue
			}
			if err = user.SaveReceipt(latest); err != nil {
				foundError = true
				stderr.WriteString(
					"Error saving new receipt for " +
						receipt.License.Values.Server + "/orders/" +
						receipt.License.Values.OfferID + "\n",
				)
				continue
			}
			stdout.WriteString(
				"Saved new receipt for " +
					receipt.License.Values.Server + "/orders/" +
					receipt.License.Values.OfferID + "\n",
			)
		}
		if foundError {
			return 1
		}
		return 0
	},
}
