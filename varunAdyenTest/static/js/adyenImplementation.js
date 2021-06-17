const clientKey = document.getElementById("clientKey").innerHTML;
const type = document.getElementById("type").innerHTML;

async function initCheckout() {
  try {
    const paymentMethodsResponse = await callServer("/api/getPaymentMethods", {});
    const configuration = {
      paymentMethodsResponse: filterUnimplemented(JSON.parse(paymentMethodsResponse)),
      clientKey,
      locale: "en_US",
      environment: "test",
      showPayButton: true,
      paymentMethodsConfiguration: {
        ideal: {
          showImage: true,
        },
        card: {
          hasHolderName: true,
          holderNameRequired: true,
          name: "Credit or debit card",
          amount: {
            value: 3000,
            currency: "EUR",
          },
        },
      },
      onSubmit: (state, component) => {
        if (state.isValid) {
          handleSubmission(state, component, "/api/initiatePayment");
        }
      },
      onAdditionalDetails: (state, component) => {
        handleSubmission(state, component, "/api/handlePaymentAddtionaldetails");
      },
    };
    console.info("before adyencheckout")

    const checkout = new AdyenCheckout(configuration);
    checkout.create(type).mount(document.getElementById(type));

    //const dropin = checkout.create('dropin').mount('#dropin-container');

  } catch (error) {
    console.error(error);
    alert("Error occurred. Look at console for details");
  }
}

function filterUnimplemented(pm) {
  console.info("here in the filter")
  pm.paymentMethods = pm.paymentMethods.filter((it) =>
    ["ach", "scheme", "dotpay", "giropay", "ideal", "directEbanking", "klarna_paynow", "klarna", "klarna_account"].includes(it.type)
  );
  return pm;
}

// Event handlers called when the shopper selects the pay button,
// or when additional information is required to complete the payment
async function handleSubmission(state, component, url) {
  try {
    const res = await callServer(url, state.data);
    const res2 = JSON.parse(res);
    handleServerResponse(res2, component);
  } catch (error) {
    console.error(error);
    alert("Error occurred. Look at console for details");
  }
}

// Calls your server endpoints
async function callServer(url, data) {
  const res = await fetch(url, {
    method: "POST",
    body: data ? JSON.stringify(data) : "",
    headers: {
      "Content-Type": "application/json",
    },
  });

  return await res.json();
}

// Handles responses sent from your server to the client
function handleServerResponse(res2, component) {
  if (res2.action) {
    component.handleAction(res2.action);
  } else {
    switch (res2.resultCode) {
      case "Authorised":
        window.location.href = `/end/success?pspReference=${res2.pspReference}`;
        break;
      case "Pending":
      case "Received":
        window.location.href = "/end/pending";
        break;
      case "Refused":
        window.location.href = "/end/failed";
        break;
      default:
        window.location.href = `/end/error?reason=${res2.resultCode}`;
        break;
    }
  }
}

initCheckout();
