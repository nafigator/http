// Package status provides structures for most used HTTP-statuses.
package status

import "net/http"

type Ok struct{}

func (s *Ok) Status() int { return http.StatusOK }

type Created struct{}

func (s *Created) Status() int { return http.StatusCreated }

type Accepted struct{}

func (s *Accepted) Status() int { return http.StatusAccepted }

type NoContent struct{}

func (s *NoContent) Status() int { return http.StatusNoContent }

type ResetContent struct{}

func (s *ResetContent) Status() int { return http.StatusResetContent }
