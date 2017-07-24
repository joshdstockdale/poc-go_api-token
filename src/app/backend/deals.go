// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backend

import (
	"app/handler"
	"net/http"

	"github.com/gorilla/mux"
)

type DealData struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	OriginalPrice float32 `json:"originalPrice"`
	SalePrice     float32 `json:"salePrice"`
}

type DealResponse struct {
	Subscriber bool     `json:"subscriber"`
	Access     bool     `json:"access"`
	Deals      DealData `json:"deals"`
}

type PrivateDealsResponse interface {
	CreatePrivateDealsResponse() DealResponse
}
type PublicDealsResponse interface {
	CreatePublicDealsResponse() DealResponse
}

const (
	DEALS_PATH = "/deals/"
)

func InitDeals(router mux.Router) {
	router.HandleFunc(DEALS_PATH+"private", handlePrivateDeals)
	router.HandleFunc(DEALS_PATH+"public", handlePublicDeals)
}

func handlePrivateDeals(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(ACCESS_COOKIE)
	if err != nil {
		handler.HandleJsonResponse(w, r, err)
		return
	}
	handler.HandleJsonResponse(w, r, &DealData{
		Id:            1234,
		Name:          "Private Deal of the Day",
		Description:   "This is just as the title says. We marked it up bc it is so private.",
		OriginalPrice: 9.99,
		SalePrice:     12.99,
	})
}
func handlePublicDeals(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(ACCESS_COOKIE)
	if err != nil {
		handler.HandleJsonResponse(w, r, err)
		return
	}
	handler.HandleJsonResponse(w, r, &DealData{
		Id:            1235,
		Name:          "Public Deal of the Day",
		Description:   "We marked it down bc it is so public.",
		OriginalPrice: 9.99,
		SalePrice:     1.99,
	})
}
