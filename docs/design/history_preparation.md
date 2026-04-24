# History Preparation

## Goal

This document describes how the current codebase is prepared for a future
user history feature without implementing persistence yet.

## Current Direction

The intended ownership is:

- `frontend`: display current result and later render history
- `api-gateway`: authenticate, orchestrate requests, shape assessment data
- `ml-api`: prediction only
- `data-svc`: future ownership of assessment persistence and history queries

## Prepared Flow in `api-gateway`

The risk flow now has an explicit preparation step for future history:

1. receive risk features from the browser
2. forward features to `ml-api`
3. decode the prediction into a typed response
4. extract the authenticated user from request context
5. assemble an `assessmentCandidate`
6. pass the candidate to an `assessmentRecorder`
7. return the prediction to the browser

Today the recorder is a no-op implementation. Later it can be replaced by:

- a client calling `data-svc`
- a temporary PostgreSQL-backed implementation in `api-gateway`

## Why This Helps

- `main.go` no longer needs to absorb all history logic
- the shape of data needed for history is already defined
- the future integration point is isolated behind `assessmentRecorder`
- `data-svc` can be introduced later without rewriting the prediction flow

## Next Implementation Step

The next logical step for history is to add a real recorder implementation
that sends `assessmentCandidate` to `data-svc`, then expose read endpoints
for user assessments through `api-gateway`.
