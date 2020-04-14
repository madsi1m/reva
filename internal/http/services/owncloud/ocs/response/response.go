// Copyright 2018-2020 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package response

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"reflect"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/cs3org/reva/pkg/appctx"
)

// Response is the top level response structure
type Response struct {
	OCS *Payload `json:"ocs"`
}

// Payload combines response metadata and data
type Payload struct {
	XMLName struct{}    `json:"-" xml:"ocs"`
	Meta    *Meta       `json:"meta" xml:"meta"`
	Data    interface{} `json:"data,omitempty" xml:"data,omitempty"`
}

var (
	elementStartElement = xml.StartElement{Name: xml.Name{Local: "element"}}
	metaStartElement    = xml.StartElement{Name: xml.Name{Local: "meta"}}
	ocsName             = xml.Name{Local: "ocs"}
	dataName            = xml.Name{Local: "data"}
)

// MarshalXML handles ocs specific wrapping of array members in 'element' tags for the data
func (p Payload) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	// first the easy part
	// use ocs as the surrounding tag
	start.Name = ocsName
	if err = e.EncodeToken(start); err != nil {
		return
	}

	// encode the meta tag
	if err = e.EncodeElement(p.Meta, metaStartElement); err != nil {
		return
	}

	// we need to use reflection to determine if p.Data is an array or a slice
	rt := reflect.TypeOf(p.Data)
	if rt != nil && (rt.Kind() == reflect.Array || rt.Kind() == reflect.Slice) {
		// this is how to wrap the data elements in their own <element> tag
		v := reflect.ValueOf(p.Data)
		if err = e.EncodeToken(xml.StartElement{Name: dataName}); err != nil {
			return
		}
		for i := 0; i < v.Len(); i++ {
			if err = e.EncodeElement(v.Index(i).Interface(), elementStartElement); err != nil {
				return
			}
		}
		if err = e.EncodeToken(xml.EndElement{Name: dataName}); err != nil {
			return
		}
	} else if err = e.EncodeElement(p.Data, xml.StartElement{Name: dataName}); err != nil {
		return
	}

	// write the closing <ocs> tag
	if err = e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return
	}
	return
}

// Meta holds response metadata
type Meta struct {
	Status       string `json:"status" xml:"status"`
	StatusCode   int    `json:"statuscode" xml:"statuscode"`
	Message      string `json:"message" xml:"message"`
	TotalItems   string `json:"totalitems,omitempty" xml:"totalitems,omitempty"`
	ItemsPerPage string `json:"itemsperpage,omitempty" xml:"itemsperpage,omitempty"`
}

// MetaOK is the default ok response
var MetaOK = &Meta{Status: "ok", StatusCode: 100, Message: "OK"}

// MetaBadRequest is used for unknown errers
var MetaBadRequest = &Meta{Status: "error", StatusCode: 400, Message: "Bad Request"}

// MetaServerError is returned on server errors
var MetaServerError = &Meta{Status: "error", StatusCode: 996, Message: "Server Error"}

// MetaUnauthorized is returned on unauthorized requests
var MetaUnauthorized = &Meta{Status: "error", StatusCode: 997, Message: "Unauthorised"}

// MetaNotFound is returned when trying to access not existing resources
var MetaNotFound = &Meta{Status: "error", StatusCode: 998, Message: "Not Found"}

// MetaUnknownError is used for unknown errers
var MetaUnknownError = &Meta{Status: "error", StatusCode: 999, Message: "Unknown Error"}

// WriteOCSSuccess handles writing successful ocs response data
func WriteOCSSuccess(w http.ResponseWriter, r *http.Request, d interface{}) {
	WriteOCSData(w, r, MetaOK, d, nil)
}

// WriteOCSError handles writing error ocs responses
func WriteOCSError(w http.ResponseWriter, r *http.Request, c int, m string, err error) {
	WriteOCSData(w, r, &Meta{Status: "error", StatusCode: c, Message: m}, nil, err)
}

// WriteOCSData handles writing ocs data in json and xml
func WriteOCSData(w http.ResponseWriter, r *http.Request, m *Meta, d interface{}, err error) {
	WriteOCSResponse(w, r, &Response{
		OCS: &Payload{
			Meta: m,
			Data: d,
		},
	}, err)
}

// WriteOCSResponse handles writing ocs responses in json and xml
func WriteOCSResponse(w http.ResponseWriter, r *http.Request, res *Response, err error) {
	var encoded []byte

	if err != nil {
		appctx.GetLogger(r.Context()).Error().Err(err).Msg(res.OCS.Meta.Message)
	}

	if r.URL.Query().Get("format") == "json" {
		w.Header().Set("Content-Type", "application/json")
		encoded, err = json.Marshal(res)
	} else {
		w.Header().Set("Content-Type", "application/xml")
		_, err = w.Write([]byte(xml.Header))
		if err != nil {
			appctx.GetLogger(r.Context()).Error().Err(err).Msg("error writing xml header")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		encoded, err = xml.Marshal(res.OCS)
	}
	if err != nil {
		appctx.GetLogger(r.Context()).Error().Err(err).Msg("error encoding ocs response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO map error for v2 only?
	// see https://github.com/owncloud/core/commit/bacf1603ffd53b7a5f73854d1d0ceb4ae545ce9f#diff-262cbf0df26b45bad0cf00d947345d9c
	switch res.OCS.Meta.StatusCode {
	case MetaNotFound.StatusCode:
		w.WriteHeader(http.StatusNotFound)
	case MetaServerError.StatusCode:
		w.WriteHeader(http.StatusInternalServerError)
	case MetaUnknownError.StatusCode:
		w.WriteHeader(http.StatusInternalServerError)
	case MetaUnauthorized.StatusCode:
		w.WriteHeader(http.StatusUnauthorized)
	case 100:
		w.WriteHeader(http.StatusOK)
	case 104:
		w.WriteHeader(http.StatusForbidden)
	default:
		// any 2xx, 4xx and 5xx will be used as is
		if res.OCS.Meta.StatusCode >= 200 && res.OCS.Meta.StatusCode < 600 {
			w.WriteHeader(res.OCS.Meta.StatusCode)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	_, err = w.Write(encoded)
	if err != nil {
		appctx.GetLogger(r.Context()).Error().Err(err).Msg("error writing ocs response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// UserIDToString returns a userid string with an optional idp separated by @: "<opaque id>[@<idp>]"
func UserIDToString(userID *user.UserId) string {
	if userID == nil || userID.OpaqueId == "" {
		return ""
	}
	if userID.Idp == "" {
		return userID.OpaqueId
	}
	return userID.OpaqueId + "@" + userID.Idp
}