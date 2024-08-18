package handlers

import (
	"net/http"

	"github.com/McFlanky/microservices-fullstack-example/api/data"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// Update handles PUT requests to update products
func (p *Products) UpdateProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	// fetch the product from the context
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Debug("Updating record id", prod.ID)

	err := p.productDB.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		p.l.Error("Product not found", err)

		w.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found in database"}, w)
		return
	}

	// write the no content success header
	w.WriteHeader(http.StatusNoContent)
}
