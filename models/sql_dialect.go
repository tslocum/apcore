// apcore is a server framework for implementing an ActivityPub application.
// Copyright (C) 2019 Cory Slep
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package models

// SqlDialect is a SQL dialect provider.
//
// Note that the order for inputs and outputs listed matter.
type SqlDialect interface {
	/* Table Creation Statements */

	// CreateUsersTable for the User model.
	CreateUsersTable() string
	// CreateFedDataTable for the FedData model.
	CreateFedDataTable() string
	// CreateLocalDataTable for the LocalData model.
	CreateLocalDataTable() string
	// CreateInboxesTable for the Inboxes model.
	CreateInboxesTable() string
	// CreateOutboxesTable for the Outboxes model.
	CreateOutboxesTable() string
	// CreateDeliveryAttemptsTable for the DeliveryAttempts model.
	CreateDeliveryAttemptsTable() string
	// CreatePrivateKeysTable for the PrivateKeys model.
	CreatePrivateKeysTable() string
	// CreateClientInfosTable for the ClientInfos model.
	CreateClientInfosTable() string
	// CreateTokenInfosTable for the TokenInfos model.
	CreateTokenInfosTable() string

	/* Queries */

	// InsertUser:
	//  Params
	//   Email       string
	//   Hashpass    []byte
	//   Salt        []byte
	//   Actor       []byte
	//   Privileges  []byte
	//   Preferences []byte
	//  Returns
	//   ID          string
	InsertUser() string
	// SensitiveUserByEmail:
	//  Params
	//   Email       string
	//  Returns
	//   ID          string
	//   Hashpass    []byte
	//   Salt        []byte
	SensitiveUserByEmail() string
	// UserByID:
	//  Params
	//   ID          string
	//  Returns
	//   ID          string
	//   Email       string
	//   Actor       []byte
	//   Privileges  []byte
	//   Preferences []byte
	UserByID() string
	// ActorIDForOutbox:
	//  Params
	//   OutboxID    string
	//  Returns
	//   ActorID     string
	ActorIDForOutbox() string
	// ActorIDForInbox:
	//  Params
	//   InboxID     string
	//  Returns
	//   ActorID     string
	ActorIDForInbox() string

	// FedCreate:
	//  Params
	//   Payload     []byte
	//  Returns
	FedCreate() string
	// FedUpdate:
	//  Params
	//   ID          string
	//   Payload     []byte
	//  Returns
	FedUpdate() string
	// FedDelete:
	//  Params
	//   ID          string
	//  Returns
	FedDelete() string

	// LocalCreate:
	//  Params
	//   Payload     []byte
	//  Returns
	LocalCreate() string
	// LocalUpdate:
	//  Params
	//   ID          string
	//   Payload     []byte
	//  Returns
	LocalUpdate() string
	// LocalDelete:
	//  Params
	//   ID          string
	//  Returns
	LocalDelete() string

	// InsertInbox:
	//  Params
	//   ActorID     string
	//   Inbox      []byte
	//  Returns
	InsertInbox() string
	// InboxContainsForActor:
	//  Params
	//   ActorID     string
	//   Item        string
	//  Returns
	//   Bool        bool
	InboxContainsForActor() string
	// InboxContains:
	//  Params
	//   InboxID     string
	//   Item        string
	//  Returns
	//   Bool        bool
	InboxContains() string
	// GetInbox:
	//  Params
	//   Inbox       string
	//   Min         int
	//   Max         int
	//  Returns
	//   Page        []byte
	GetInbox() string
	// GetPublicInbox:
	//  Params
	//   Inbox       string
	//   Min         int
	//   Max         int
	//  Returns
	//   Page        []byte
	GetPublicInbox() string
	// GetInboxLastPage:
	//  Params
	//   Inbox       string
	//   N           int
	//  Returns
	//   Page        []byte
	//   StartIndex  int
	GetInboxLastPage() string
	// GetPublicInboxLastPage:
	//  Params
	//   Inbox       string
	//   N           int
	//  Returns
	//   Page        []byte
	//   StartIndex  int
	GetPublicInboxLastPage() string
	// PrependInboxItem:
	//  Params
	//   Inbox       string
	//   Item        string
	//  Returns
	PrependInboxItem() string
	// DeleteInboxItem:
	//  Params
	//   Inbox       string
	//   Item        string
	//  Returns
	DeleteInboxItem() string

	// InsertOutbox:
	//  Params
	//   ActorID     string
	//   Outbox      []byte
	//  Returns
	InsertOutbox() string
	// OutboxContainsForActor:
	//  Params
	//   ActorID     string
	//   Item        string
	//  Returns
	//   Bool        bool
	OutboxContainsForActor() string
	// OutboxContains:
	//  Params
	//   OutboxID    string
	//   Item        string
	//  Returns
	//   Bool        bool
	OutboxContains() string
	// GetOutbox:
	//  Params
	//   Outbox      string
	//   Min         int
	//   Max         int
	//  Returns
	//   Page        []byte
	GetOutbox() string
	// GetPublicOutbox:
	//  Params
	//   Outbox      string
	//   Min         int
	//   Max         int
	//  Returns
	//   Page        []byte
	GetPublicOutbox() string
	// GetOutboxLastPage:
	//  Params
	//   Outbox      string
	//   N           int
	//  Returns
	//   Page        []byte
	//   StartIndex  int
	GetOutboxLastPage() string
	// GetPublicOutboxLastPage:
	//  Params
	//   Outbox      string
	//   N           int
	//  Returns
	//   Page        []byte
	//   StartIndex  int
	GetPublicOutboxLastPage() string
	// PrependOutboxItem:
	//  Params
	//   Outbox      string
	//   Item        string
	//  Returns
	PrependOutboxItem() string
	// DeleteOutboxItem:
	//  Params
	//   Outbox      string
	//   Item        string
	//  Returns
	DeleteOutboxItem() string
	// OutboxForInbox:
	//  Params
	//   Inbox       string
	//  Returns
	//   Outbox      string
	OutboxForInbox() string

	// InsertAttempt:
	//  Params
	//   FromID      string
	//   ToActor     string
	//   Payload     []byte
	//   State       string
	//  Returns
	//   ID          string
	InsertAttempt() string
	// MarkSuccessfulAttempt:
	//  Params
	//   ID          string
	//  Returns
	MarkSuccessfulAttempt() string
	// MarkFailedAttempt:
	//  Params
	//   ID          string
	//  Returns
	MarkFailedAttempt() string

	// CreatePrivateKey:
	//  Params
	//   UserID      string
	//   Purpose     string
	//   PrivKey     []byte
	//  Returns
	CreatePrivateKey() string
	// GetPrivateKeyByUserID:
	//  Params
	//   UserID      string
	//   Purpose     string
	//  Returns
	//   PrivKey     []byte
	GetPrivateKeyByUserID() string

	// CreateClientInfo:
	//  Params
	//   Secret      string
	//   Domain      string
	//   UserID      string
	//  Returns
	//   ID          string
	CreateClientInfo() string
	// GetClientInfoByID:
	//  Params
	//   ID          string
	//  Returns
	//   ID          string
	//   Secret      string
	//   Domain      string
	//   UserID      string
	GetClientInfoByID() string

	// CreateTokenInfo:
	//  Params
	//   ClientID    string
	//   UserID      string
	//   RedirURI    string
	//   Scope       string
	//   Code        string
	//   CodeCreated time.Time
	//   CodeExpires time.Duration
	//   Access      string
	//   AccessCtd   time.Time
	//   AccessExp   time.Duration
	//   Refresh     string
	//   RefrCreated time.Time
	//   RefrExpires time.Duration
	//  Returns
	CreateTokenInfo() string
	// RemoveTokenInfoByCode:
	//  Params
	//   Code        string
	//  Returns
	RemoveTokenInfoByCode() string
	// RemoveTokenInfoByAccess:
	//  Params
	//   Access      string
	//  Returns
	RemoveTokenInfoByAccess() string
	// RemoveTokenInfoByRefresh:
	//  Params
	//   Refresh     string
	//  Returns
	RemoveTokenInfoByRefresh() string
	// GetTokenInfoByCode:
	//  Params
	//   Code        string
	//  Returns
	//   ClientID    string
	//   UserID      string
	//   RedirURI    string
	//   Scope       string
	//   Code        string
	//   CodeCreated time.Time
	//   CodeExpires time.Duration
	//   Access      string
	//   AccessCtd   time.Time
	//   AccessExp   time.Duration
	//   Refresh     string
	//   RefrCreated time.Time
	//   RefrExpires time.Duration
	GetTokenInfoByCode() string
	// GetTokenInfoByAccess:
	//  Params
	//   Access      string
	//  Returns
	//   ClientID    string
	//   UserID      string
	//   RedirURI    string
	//   Scope       string
	//   Code        string
	//   CodeCreated time.Time
	//   CodeExpires time.Duration
	//   Access      string
	//   AccessCtd   time.Time
	//   AccessExp   time.Duration
	//   Refresh     string
	//   RefrCreated time.Time
	//   RefrExpires time.Duration
	GetTokenInfoByAccess() string
	// GetTokenInfoByRefresh:
	//  Params
	//   Refresh     string
	//  Returns
	//   ClientID    string
	//   UserID      string
	//   RedirURI    string
	//   Scope       string
	//   Code        string
	//   CodeCreated time.Time
	//   CodeExpires time.Duration
	//   Access      string
	//   AccessCtd   time.Time
	//   AccessExp   time.Duration
	//   Refresh     string
	//   RefrCreated time.Time
	//   RefrExpires time.Duration
	GetTokenInfoByRefresh() string
}