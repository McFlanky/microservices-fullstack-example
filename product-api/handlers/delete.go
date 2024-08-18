package handlers

import (
	"net/http"

	"github.com/McFlanky/microservices-fullstack-example/api/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Return a list of products
// responses:
//	201: noContentResponse
// 	404: errorResponse
// 	501: errorResponse

// Delete handles DELETE requests and removes items from the database
func (p *Products) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := getProductID(r)

	p.l.Debug("Deleting record id", id)

	err := p.productDB.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Error("Deleting record id does not exist")

		w.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	if err != nil {
		p.l.Error("Deleting record", err)

		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
