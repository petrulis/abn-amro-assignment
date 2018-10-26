package main

import "github.com/petrulis/abn-amro-assignment/abnerr"

// Errors
var (
	errNotFound   = abnerr.New("NotFound", "Requested resource couldn't be found", 404)
	errInternal   = abnerr.New("Internal", "Unexpected error occured. Please try again later", 500)
	errBadRequest = abnerr.New("BadRequest", "Request payload is not valid", 400)
)
