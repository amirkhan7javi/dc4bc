package client

import (
	"encoding/json"
	"fmt"
	"github.com/depools/dc4bc/client/types"
	"image"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/depools/dc4bc/qr"
	"github.com/depools/dc4bc/storage"
)

type Response struct {
	ErrorMessage string `json:"error_message,omitempty"`
	Result interface{} `json:"result"`
}

func rawResponse(w http.ResponseWriter, response []byte) {
	if _, err := w.Write(response); err != nil {
		panic(fmt.Sprintf("failed to write response: %v", err))
	}
}

func errorResponse(w http.ResponseWriter, statusCode int, error string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	resp := Response{ErrorMessage: error}
	respBz, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Failed to marshal response: %v\n", err)
		return
	}
	if _, err := w.Write(respBz); err != nil {
		panic(fmt.Sprintf("failed to write response: %v", err))
	}
}

func successResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	resp := Response{Result: response}
	respBz, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Failed to marshal response: %v\n", err)
		return
	}
	if _, err := w.Write(respBz); err != nil {
		panic(fmt.Sprintf("failed to write response: %v", err))
	}
}

func (c *Client) StartHTTPServer(listenAddr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/sendMessage", c.sendMessageHandler)
	mux.HandleFunc("/getOperations", c.getOperationsHandler)
	mux.HandleFunc("/getOperationQRPath", c.getOperationQRPathHandler)
	mux.HandleFunc("/readProcessedOperationFromCamera", c.readProcessedOperationFromCameraHandler)

	mux.HandleFunc("/readProcessedOperation", c.readProcessedOperationFromBodyHandler)
	mux.HandleFunc("/getOperationQR", c.getOperationQRToBodyHandler)

	return http.ListenAndServe(listenAddr, mux)
}

func (c *Client) sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusBadRequest, "Wrong HTTP method")
		return
	}
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to read request body: %v", err))
		return
	}

	var msg storage.Message
	if err = json.Unmarshal(reqBytes, &msg); err != nil {
		errorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to unmarshal message: %v", err))
		return
	}

	if err = c.SendMessage(msg); err != nil {
		errorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to send message to the storage: %v", err))
		return
	}

	successResponse(w, "ok")
}

func (c *Client) getOperationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusBadRequest, "Wrong HTTP method")
		return
	}

	operations, err := c.GetOperations()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get operations: %v", err))
		return
	}

	successResponse(w, operations)
}

func (c *Client) getOperationQRPathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusBadRequest, "Wrong HTTP method")
		return
	}
	operationID := r.URL.Query().Get("operationID")

	qrPath, err := c.GetOperationQRPath(operationID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get operation QR path: %v", err))
		return
	}

	successResponse(w, qrPath)
}

func (c *Client) getOperationQRToBodyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusBadRequest, "Wrong HTTP method")
		return
	}
	operationID := r.URL.Query().Get("operationID")

	operationJSON, err := c.getOperationJSON(operationID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get operation in JSON: %v", err))
		return
	}

	encodedData, err := qr.EncodeQR(operationJSON)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode operation: %v", err))
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(encodedData)))
	rawResponse(w, encodedData)
}

func (c *Client) readProcessedOperationFromCameraHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusBadRequest, "Wrong HTTP method")
		return
	}

	if err := c.ReadProcessedOperation(); err != nil {
		errorResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to handle processed operation from camera path: %v", err))
		return
	}

	successResponse(w, "ok")
}

func (c *Client) readProcessedOperationFromBodyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusBadRequest, "Wrong HTTP method")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		errorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse multipat form: %v", err))
		return
	}

	file, _, err := r.FormFile("qr")
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to retrieve a file: %v", err))
		return
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to decode an image: %v", err))
		return
	}

	qrData, err := qr.ReadDataFromQR(img)
	if err != nil {
		return
	}

	var operation types.Operation
	if err = json.Unmarshal(qrData, &operation); err != nil {
		errorResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to unmarshal processed operation: %v", err))
		return
	}
	if err := c.handleProcessedOperation(operation); err != nil {
		errorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to handle an operation: %v", err))
		return
	}

	successResponse(w, "ok")
}