package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"time"
)

// UserEntitlementDueUpdate represents a row from 'trackit.user_entitlement_due_update'.
type UserEntitlementDueUpdate struct {
	ID                     int            `json:"id"`                       // id
	Email                  string         `json:"email"`                    // email
	Auth                   string         `json:"auth"`                     // auth
	NextExternal           sql.NullString `json:"next_external"`            // next_external
	ParentUserID           sql.NullInt64  `json:"parent_user_id"`           // parent_user_id
	AwsCustomerIdentifier  string         `json:"aws_customer_identifier"`  // aws_customer_identifier
	AwsCustomerEntitlement bool           `json:"aws_customer_entitlement"` // aws_customer_entitlement
	NextUpdateEntitlement  time.Time      `json:"next_update_entitlement"`  // next_update_entitlement
}
