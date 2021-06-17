package api

import (
	"bytes"
	"encoding/json"
	fmt "fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)


type RedirectReq struct {
	MD             string
	PaRes          string
}

/* Cache the paymentdata to use Payments/details API */
var CachepaymentData = map[string]string{}

func PaymentMethodHandler(c *gin.Context) {

	fmt.Println("Calling Adyen Payment Methods API")

	amount := map[string]interface{}{
		"currency": "EUR",
		"value": 3000,
	}

	paymenMethodreq := map[string]interface{}{
		"merchantAccount": os.Getenv("MERCHANT_ACCOUNT"),
		"countryCode": "NL",
		"shopperLocale": "nl-NL",
		"amount": amount ,
	}

	bpaymentreqbody,_ := json.Marshal(paymenMethodreq)

	print("Request Body for Payment Method :",bytes.NewBuffer(bpaymentreqbody).String())

	b := bytes.NewBuffer(bpaymentreqbody) // byte array to byte.buffer for http request which needs io reader.

	//postRequest("https://checkout-test.adyen.com/v67/paymentMethods", b.Bytes())

	responsebody := postRequest("https://checkout-test.adyen.com/v67/paymentMethods", b.Bytes())

     fmt.Println("Response Body for Payment Method:",string(responsebody))

	c.JSON(http.StatusOK, responsebody)

	return

}

func PaymentsHandler(c *gin.Context) {

	fmt.Println("Calling Adyen Payments API")

	//data, _ := ioutil.ReadAll(c.Request.Body)

	//cardvalue := c.Request.FormValue("encryptedCardNumber")

	//gg := string(data)

	var fg map[string]interface{}

	err := json.NewDecoder(c.Request.Body).Decode(&fg)

	if  err == nil {
     fmt.Println("Error while processing payment request")
	}

	paymentMethod := fg["paymentMethod"]

	//orderreference := uuid.Must(uuid.NewRandom()).String();

	orderreference := "varun_checkoutChallenge"

	amount := map[string]interface{}{
		"currency": "EUR",
		"value": 3000,
	}

/* used for native 3ds3
	additionaldata := map[string]interface{}{
		"allow3DS2": true,
	}
*/
	var k = c.GetHeader("User-Agent")

	BrowserInfo :=  map[string]interface{}{
		"userAgent": k ,
		"acceptHeader": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	}

	paymentreqs := map[string]interface{}{
		"amount": amount,
		"reference" : orderreference,
		"paymentMethod" : paymentMethod,
		"merchantAccount": os.Getenv("MERCHANT_ACCOUNT"),
		"shopperIP": c.ClientIP(),
		"channel": "Web",
		"origin": "http://localhost:8080",
		"returnUrl" : fmt.Sprintf("http://localhost:8080/api/handleUIRedirect?orderRef=%s", orderreference),
		"browserInfo" : BrowserInfo,
	}

	bpaymentreqbody,_ := json.Marshal(paymentreqs)

	//print("Request Body for Payment :",bytes.NewBuffer(bpaymentreqbody).String())

	b := bytes.NewBuffer(bpaymentreqbody) // byte array to byte.buffer for http request which needs io reader.

	responsebody := postRequestreturnByte("https://checkout-test.adyen.com/v66/payments", b.Bytes())


    var actionvalues map [string] interface { }

	action := map[string]interface{}{
		"action": actionvalues ,
	}

	err = json.NewDecoder(bytes.NewReader(responsebody)).Decode(&action)

	if action["action"] != nil {


		var paymentData = fmt.Sprintf("%v", action["paymentData"])

		CachepaymentData[orderreference] = paymentData

		c.JSON(http.StatusOK, string(responsebody))

		return
	}

	rg := map[string] string {

	}

	err = json.NewDecoder(bytes.NewReader(responsebody)).Decode(&rg)

	if err != nil {
	 }
	pspReference := rg["pspReference"]
	resultCode := rg["resultCode"]
	refusalReason := rg["refusalReason"]

	fmt.Println(resultCode)

	c.JSON(
		http.StatusOK,
		map[string] string {
		"pspReference":  pspReference,
		"resultCode": resultCode,
		"refusalReason": refusalReason,
	})
	return
}


func RedirectUIHandler(c *gin.Context) {

	fmt.Println("Calling Redirect UI Handler")

	var rd RedirectReq

	if err := c.ShouldBind(&rd); err != nil {
		log.Print(" Redirect failed - Unable to bind with request")
		return
	}

	var details map[string]interface{}

	redirectResult := c.Query("redirectResult");

	if redirectResult != ""  {
		 details = map[string]interface{}{
			"redirectResult":  redirectResult,
		}
	} else {
		  details  = map[string]interface{}{
			"MD":    rd.MD,
			"PaRes": rd.PaRes,
		}
	}


	paymentData := CachepaymentData[c.Query("orderRef")]

	paymentdetailsreqs := map[string]interface{}{
		"paymentData":  paymentData,
		"details":     details,
	}

	bpaymentreqbody, _ := json.Marshal(paymentdetailsreqs)

	b := bytes.NewBuffer(bpaymentreqbody) // byte array to byte.buffer for http request which needs io reader.

	responsebody := postRequestreturnByte("https://checkout-test.adyen.com/v66/payments/details", b.Bytes())

	rg := map[string]string{}

	json.NewDecoder(bytes.NewReader(responsebody)).Decode(&rg)

	//pspReference := rg["pspReference"]
	resultCode := rg["resultCode"]
	//refusalReason := rg["refusalReason"]

	if "Authorised"  == resultCode {
		c.Redirect(http.StatusFound, "/end/success")
	} else if "Refused" == resultCode {
		c.Redirect(http.StatusFound, "/end/failed")
	}else if "Authorised" == c.Query("resultCode") {   /* for Ideal*/
		c.Redirect(http.StatusFound, "/end/success")
	} else if  "Refused" == c.Query("resultCode") {
		c.Redirect(http.StatusFound, "/end/success")
	}

}



func postRequest(apiurl string, byts [] byte) string {

	req, err := http.NewRequest("POST", apiurl , bytes.NewBuffer(byts))

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set( "x-API-Key", os.Getenv("API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}


func postRequestreturnByte(apiurl string, byts []byte) []byte {

	req, err := http.NewRequest("POST", apiurl , bytes.NewBuffer(byts))

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set( "x-API-Key", os.Getenv("API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body
}

