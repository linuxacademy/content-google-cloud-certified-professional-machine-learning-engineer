// Package cxwh contains an example Dialogflow CX webhook
package cxwh

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type fulfillmentInfo struct {
	Tag string `json:"tag"`
}

type sessionInfo struct {
	Session    string                 `json:"session"`
	Parameters map[string]interface{} `json:"parameters"`
}

type text struct {
	Text []string `json:"text"`
}

type responseMessage struct {
	Text text `json:"text"`
}

type fulfillmentResponse struct {
	Messages []responseMessage `json:"messages"`
}

// webhookRequest is used to unmarshal a WebhookRequest JSON object. Note that
// not all members need to be defined--just those that you need to process.
// As an alternative, you could use the types provided by the Dialogflow protocol buffers:
// https://pkg.go.dev/google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3#WebhookRequest
type webhookRequest struct {
	FulfillmentInfo fulfillmentInfo `json:"fulfillmentInfo"`
	SessionInfo     sessionInfo     `json:"sessionInfo"`
}

// webhookResponse is used to marshal a WebhookResponse JSON object. Note that
// not all members need to be defined--just those that you need to process.
// As an alternative, you could use the types provided by the Dialogflow protocol buffers:
// https://pkg.go.dev/google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3#WebhookResponse
type webhookResponse struct {
	FulfillmentResponse fulfillmentResponse `json:"fulfillmentResponse"`
	SessionInfo         sessionInfo         `json:"sessionInfo"`
}

// confirm handles webhook calls using the "confirm" tag.
func confirm(request webhookRequest) (webhookResponse, error) {
	// Create a text message that utilizes the "size" and "color"
	// parameters provided by the end-user.
	// This text message is used in the response below.
	t := fmt.Sprintf("You can pick up your order for a %s %s hoodie in 5 days.",
		request.SessionInfo.Parameters["size"],
		request.SessionInfo.Parameters["color"])

	// Create session parameters that are populated in the response.
	// The "cancel-period" parameter is referenced by the agent.
	// This example hard codes the value 2, but a real system
	// might look up this value in a database.
	p := map[string]interface{}{"cancel-period": "2"}

	// Build and return the response.
	response := webhookResponse{
		FulfillmentResponse: fulfillmentResponse{
			Messages: []responseMessage{
				{
					Text: text{
						Text: []string{t},
					},
				},
			},
		},
		SessionInfo: sessionInfo{
			Parameters: p,
		},
	}
	return response, nil
}

// handleError handles internal errors.
func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "ERROR: %v", err)
}

// HandleWebhookRequest handles WebhookRequest and sends the WebhookResponse.
func HandleWebhookRequest(w http.ResponseWriter, r *http.Request) {
	var request webhookRequest
	var response webhookResponse
	var err error

	// Read input JSON
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, err)
		return
	}
	log.Printf("Request: %+v", request)

	// Get the tag from the request, and call the corresponding
	// function that handles that tag.
	// This example only has one possible tag,
	// but most agents would have many.
	switch tag := request.FulfillmentInfo.Tag; tag {
	case "confirm":
		response, err = confirm(request)
	default:
		err = fmt.Errorf("Unknown tag: %s", tag)
	}
	if err != nil {
		handleError(w, err)
		return
	}
	log.Printf("Response: %+v", response)

	// Send response
	if err = json.NewEncoder(w).Encode(&response); err != nil {
		handleError(w, err)
		return
	}
}
