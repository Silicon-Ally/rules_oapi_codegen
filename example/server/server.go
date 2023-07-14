// Package server implements the API interface, api.StrictServerInterface,
// which is auto-generated from the OpenAPI 3 spec, and describes the simple,
// canonical 'pet store' example.
package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Silicon-Ally/rules_oapi_codegen/example/api"
)

type Server struct {
	pets []api.Pet
	idx  int
}

// Returns all pets
// (GET /pets)
func (s *Server) FindPets(ctx context.Context, req api.FindPetsRequestObject) (api.FindPetsResponseObject, error) {
	out := []api.Pet{}
	for _, p := range s.pets {
		if matchesTag(req.Params.Tags, p.Tag) {
			out = append(out, p)
		}
	}
	if req.Params.Limit != nil && *req.Params.Limit > 0 && int(*req.Params.Limit) < len(out) {
		out = out[:*req.Params.Limit]
	}

	return api.FindPets200JSONResponse(out), nil
}

func matchesTag(tags *[]string, tag *string) bool {
	if tags == nil || tag == nil || len(*tags) == 0 {
		return true
	}
	for _, t := range *tags {
		if t == *tag {
			return true
		}
	}
	return false
}

// Creates a new pet
// (POST /pets)
func (s *Server) AddPet(ctx context.Context, req api.AddPetRequestObject) (api.AddPetResponseObject, error) {
	s.idx += 1
	id := int64(s.idx)
	s.pets = append(s.pets, api.Pet{
		Id:   id,
		Name: req.Body.Name,
		Tag:  req.Body.Tag,
	})
	return api.AddPet200JSONResponse{
		Id:   id,
		Name: req.Body.Name,
		Tag:  req.Body.Tag,
	}, nil
}

// Deletes a pet by ID
// (DELETE /pets/{id})
func (s *Server) DeletePet(ctx context.Context, req api.DeletePetRequestObject) (api.DeletePetResponseObject, error) {
	for i, p := range s.pets {
		if p.Id == req.Id {
			s.pets = append(s.pets[:i], s.pets[i+1:]...)
			return api.DeletePet204Response{}, nil
		}
	}
	return api.DeletePetdefaultJSONResponse{
		Body: api.Error{
			Code:    1,
			Message: fmt.Sprintf("no pet with id %d", req.Id),
		},
		StatusCode: http.StatusNotFound,
	}, nil
}

// Returns a pet by ID
// (GET /pets/{id})
func (s *Server) FindPetByID(ctx context.Context, req api.FindPetByIDRequestObject) (api.FindPetByIDResponseObject, error) {
	for _, p := range s.pets {
		if p.Id == req.Id {
			return api.FindPetByID200JSONResponse{
				Id:   p.Id,
				Name: p.Name,
				Tag:  p.Tag,
			}, nil
		}
	}

	return api.FindPetByIDdefaultJSONResponse{
		Body: api.Error{
			Code:    2,
			Message: fmt.Sprintf("no pet with id %d", req.Id),
		},
		StatusCode: http.StatusNotFound,
	}, nil
}
