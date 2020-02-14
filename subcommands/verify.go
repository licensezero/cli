package subcommands

import (
	"errors"
	"licensezero.com/licensezero/api"
	"licensezero.com/licensezero/user"
)

const verifyDescription = "Verify receipts."

var verifyUsage = verifyDescription + "\n\n" +
	"Usage:\n" +
	"  licensezero verify\n"

// Verify saves receipts to disk.
var Verify = &Subcommand{
	Tag:         "buyer",
	Description: verifyDescription,
	Handler: func(env Environment) int {
		identity, err := user.ReadIdentity()
		if err != nil {
			env.Stderr.WriteString("Error reading identity.")
			return 1
		}
		receipts, receiptErrors, err := user.ReadReceipts()
		if err != nil {
			env.Stderr.WriteString("Error reading receipts: " + err.Error())
			return 1
		}
		foundError := false
		for _, receiptError := range receiptErrors {
			foundError = true
			env.Stderr.WriteString(receiptError.Error() + "\n")
		}
		servers := make(map[string]*api.BrokerServer)
		registers := make(map[string]*api.Register)
		for _, receipt := range receipts {
			err = receipt.Validate()
			if err != nil {
				foundError = true
				env.Stderr.WriteString(
					"Receipt for " +
						receipt.License.Values.Server + "/orders/" + receipt.License.Values.OfferID +
						"is not a valid receipt.\n",
				)
			}
			err = receipt.VerifySignature()
			if err != nil {
				foundError = true
				env.Stderr.WriteString(
					"Signature for " +
						receipt.License.Values.Server + "/orders/" + receipt.License.Values.OfferID +
						"is not valid.\n",
				)
			}
			brokerServer := receipt.License.Values.Server
			server, ok := servers[brokerServer]
			if !ok {
				server = &api.BrokerServer{Client: env.Client, Base: brokerServer}
				servers[brokerServer] = server
			}
			register, ok := registers[brokerServer]
			if !ok {
				register, err := server.Register()
				if err != nil {
					foundError = true
					env.Stderr.WriteString(
						"Could not fetch key register for " +
							brokerServer + ".\n",
					)
				}
				registers[brokerServer] = register
			}
			if err = register.ValidateEffectiveDate(receipt); err != nil {
				foundError = true
				env.Stderr.WriteString(
					"Signature for " +
						receipt.License.Values.Server + "/orders/" + receipt.License.Values.OfferID +
						"does not match time frame for the broker's signing key.\n",
				)
			}
			if errs := identity.ValidateReceipt(receipt); errs != nil {
				foundError = true
				uri := receipt.License.Values.Server + "/orders/" + receipt.License.Values.OfferID
				for _, err := range errs {
					switch {
					case errors.Is(err, user.ErrNameMismatch):
						env.Stderr.WriteString("Name on " + uri + "does not match your identity.\n")
					case errors.Is(err, user.ErrJurisdictionMismatch):
						env.Stderr.WriteString("Name on " + uri + "does not match your identity.\n")
					case errors.Is(err, user.ErrEMailMismatch):
						env.Stderr.WriteString("Name on " + uri + "does not match your identity.\n")
					}
				}
			}
		}
		if foundError {
			return 1
		}
		return 0
	},
}
