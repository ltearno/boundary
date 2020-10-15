#!/usr/bin/env bats

load _auth
load _scopes
load _helpers

export NEW_PROJECT='test_project'
export NEW_ORG='test_org'

@test "boundary/login: can login as default user" {
  run login $DEFAULT_LOGIN
  [ "$status" -eq 0 ]
}

@test "boundary/scopes: can create $NEW_ORG organization level scope" {
	run create_scope $DEFAULT_GLOBAL $NEW_ORG
  echo "$output"
	[ "$status" -eq 0 ]
}

@test "boundary/scopes: can read $NEW_ORG organization level scope" {
  local sid=$(scope_id $NEW_ORG $DEFAULT_GLOBAL)
	run read_scope $sid
  echo "$output"
	[ "$status" -eq 0 ]
}

@test "boundary/scopes: can create $NEW_PROJECT project level scope" {
  local parent=$(scope_id $NEW_ORG $DEFAULT_GLOBAL)
	run create_scope $parent $NEW_PROJECT
  echo "$output"
	[ "$status" -eq 0 ]
}

@test "boundary/scopes: can read $NEW_PROJECT project level scope" {
  local parent=$(scope_id $NEW_ORG $DEFAULT_GLOBAL)
  local sid=$(scope_id $NEW_PROJECT $parent)
	run read_scope $sid
  echo "$output"
	[ "$status" -eq 0 ]
}

@test "boundary/scopes: can delete $NEW_PROJECT project level scope" {
  local parent=$(scope_id $NEW_ORG $DEFAULT_GLOBAL)
  local sid=$(scope_id $NEW_PROJECT $parent)
	run delete_scope $sid
  echo "$output"
	[ "$status" -eq 0 ]
}

@test "boundary/scopes: can not read deleted $NEW_PROJECT project level scope" {
  local parent=$(scope_id $NEW_ORG $DEFAULT_GLOBAL)
  local sid=$(scope_id $NEW_PROJECT $parent)
	run read_scope $sid
  echo "$output"
	[ "$status" -eq 1 ]
}

@test "boundary/scopes: can delete $NEW_ORG organization level scope" {
  local sid=$(scope_id $NEW_ORG $DEFAULT_GLOBAL)
	run delete_scope $sid
  echo "$output"
	[ "$status" -eq 0 ]
}

@test "boundary/scopes: can not read deleted $NEW_ORG organization level scope" {
  local sid=$(scope_id $NEW_PROJECT $DEFAULT_GLOBAL)
	run read_scope $sid
  echo "$output"
	[ "$status" -eq 1 ]
}
