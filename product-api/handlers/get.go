package handlers

import (
	"context"
	"net/http"

	"github.com/McFlanky/microservices-fullstack-example/api/data"
	protos "github.com/McFlanky/microservices-fullstack-example/currency/protos/currency"
)

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

// ListAll handles GET requests and returns all current products
func (p *Products) ListAll(w http.ResponseWriter, r *http.Request) {
	p.l.Println("[DEBUG] get all records")
	w.Header().Add("Content-Type", "application/json")

	prods := data.GetProducts()

	err := data.ToJSON(prods, w)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Println("[ERROR] serializing product", err)
	}
}

// swagger:route GET /products/{id} products listSingleProduct
// Return a list of products from the database
// responses:
//	200: productResponse
//	404: errorResponse

// ListSingle handles GET requests
func (p *Products) ListSingle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	id := getProductID(r)

	p.l.Println("[DEBUG] get record id", id)

	prod, err := data.GetProductByID(id)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Println("[ERROR] fetching product", err)

		w.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	default:
		p.l.Println("[ERROR] fetching product", err)

		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	// get exchange rate
	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value["GBP"]),
	}

	resp, err := p.cc.GetRate(context.Background(), rr)
	if err != nil {
		p.l.Println("[ERROR] error getting new rate", err)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	p.l.Printf("Resp %#v", resp)

	prod.Price = prod.Price * resp.Rate

	err = data.ToJSON(prod, w)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Println("[ERROR] serializing product", err)
	}
}
