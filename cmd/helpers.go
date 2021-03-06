package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Define an envelope type.
type envelope map[string]interface{}

func (app *application) readIDParam(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	key := vars["id"]
	id, err := strconv.Atoi(key)

	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) readIDStringParam(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	key := vars["id"]
	return key, nil
}

// The readString() helper returns a string value from the query string, or the provided
// default value if no matching key could be found.
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

// The readCSV() helper reads a string value from the query string and then splits it
// into a slice on the comma character. If no matching key could be found, it returns
// the provided default value.
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return defaultValue
	}
	return strings.Split(csv, ",")
}

// The readInt() helper reads a string value from the query string and converts it to an
// integer before returning. If no matching key could be found it returns the provided
// default value. If the value couldn't be converted to an integer, then we record an
// error message in the provided Validator instance.
func (app *application) readInt(qs url.Values, key string, defaultValue int) int {
	// TODO: Add validator
	//func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {

	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		//v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return i
}

// The readInt() helper reads a string value from the query string and converts it to an
// integer before returning. If no matching key could be found it returns the provided
// default value. If the value couldn't be converted to an integer, then we record an
// error message in the provided Validator instance.
func (app *application) readFloat(qs url.Values, key string, defaultValue float64) float64 {
	// TODO: Add validator
	//func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {

	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		//v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return i
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1MB.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Initialize the json.Decoder, and call the DisallowUnknownFields() method on it
	// before decoding. This means that if the JSON from the client now includes any
	// field which cannot be mapped to the target destination, the decoder will return
	// an error instead of just ignoring the field.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode the request body to the destination.
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// If the JSON contains a field which cannot be mapped to the target destination
		// then Decode() will now return an error message in the format "json: unknown
		// field "<name>"". We check for this, extract the field name from the error,
		// and interpolate it into our custom error message. Note that there's an open
		// issue at https://github.com/golang/go/issues/29035 regarding turning this
		// into a distinct error type in the future.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// If the request body exceeds 1MB in size the decode will now fail with the
		// error "http: request body too large". There is an open issue about turning
		// this into a distinct error type at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return our own custom error message.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// Define a writeJSON() helper for sending responses. This takes the destination
// http.ResponseWriter, the HTTP status code to send, the data to encode to JSON, and a
// header map containing any additional HTTP headers we want to include in the response.
func (app *application) writeJSON(w http.ResponseWriter, status int, tojsonfy envelope, headers http.Header) error {
	// Encode the data to JSON, returning the error if there was one.
	js, err := json.MarshalIndent(tojsonfy, "", "\t")
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	// Add Headers
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Add the "Content-Type: application/json" header, then write the status code and
	// JSON response.
	w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")

	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) DescribeContractsBufferRenew() error {
	commodities := []string{"", "power", "gas"}
	for _, c := range commodities {
		d, err := app.contracts.Describe(c)
		if err != nil {
			return err
		}
		app.describeBuffer[c] = d
	}
	return nil
}
