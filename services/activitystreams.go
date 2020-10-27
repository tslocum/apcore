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

package services

import (
	"net/url"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/go-fed/apcore/app"
	"github.com/go-fed/apcore/paths"
)

func toPersonActor(a app.Application,
	scheme, host, username, preferredUsername, summary string,
	pubKey string) (vocab.ActivityStreamsPerson, *url.URL) {
	p := streams.NewActivityStreamsPerson()
	// id
	idProp := streams.NewJSONLDIdProperty()
	idIRI := paths.UserIRIFor(scheme, host, paths.UserPathKey, username)
	idProp.SetIRI(idIRI)
	p.SetJSONLDId(idProp)

	// inbox
	inboxProp := streams.NewActivityStreamsInboxProperty()
	inboxIRI := paths.UserIRIFor(scheme, host, paths.InboxPathKey, username)
	inboxProp.SetIRI(inboxIRI)
	p.SetActivityStreamsInbox(inboxProp)

	// outbox
	outboxProp := streams.NewActivityStreamsOutboxProperty()
	outboxIRI := paths.UserIRIFor(scheme, host, paths.OutboxPathKey, username)
	outboxProp.SetIRI(outboxIRI)
	p.SetActivityStreamsOutbox(outboxProp)

	// followers
	followersProp := streams.NewActivityStreamsFollowersProperty()
	followersIRI := paths.UserIRIFor(scheme, host, paths.FollowersPathKey, username)
	followersProp.SetIRI(followersIRI)
	p.SetActivityStreamsFollowers(followersProp)

	// following
	followingProp := streams.NewActivityStreamsFollowingProperty()
	followingIRI := paths.UserIRIFor(scheme, host, paths.FollowingPathKey, username)
	followingProp.SetIRI(followingIRI)
	p.SetActivityStreamsFollowing(followingProp)

	// liked
	likedProp := streams.NewActivityStreamsLikedProperty()
	likedIRI := paths.UserIRIFor(scheme, host, paths.LikedPathKey, username)
	likedProp.SetIRI(likedIRI)
	p.SetActivityStreamsLiked(likedProp)

	// name
	nameProp := streams.NewActivityStreamsNameProperty()
	nameProp.AppendXMLSchemaString(username)
	p.SetActivityStreamsName(nameProp)

	// preferredUsername
	preferredUsernameProp := streams.NewActivityStreamsPreferredUsernameProperty()
	preferredUsernameProp.SetXMLSchemaString(preferredUsername)
	p.SetActivityStreamsPreferredUsername(preferredUsernameProp)

	// url
	urlProp := streams.NewActivityStreamsUrlProperty()
	urlProp.AppendIRI(idIRI)
	p.SetActivityStreamsUrl(urlProp)

	// summary
	summaryProp := streams.NewActivityStreamsSummaryProperty()
	summaryProp.AppendXMLSchemaString(summary)
	p.SetActivityStreamsSummary(summaryProp)

	// publicKey property
	publicKeyProp := streams.NewW3IDSecurityV1PublicKeyProperty()

	// publicKey type
	publicKeyType := streams.NewW3IDSecurityV1PublicKey()

	// publicKey id
	pubKeyIdProp := streams.NewJSONLDIdProperty()
	pubKeyIRI := paths.UserIRIFor(scheme, host, paths.HttpSigPubKeyKey, username)
	pubKeyIdProp.SetIRI(pubKeyIRI)
	publicKeyType.SetJSONLDId(pubKeyIdProp)

	// publicKey owner
	ownerProp := streams.NewW3IDSecurityV1OwnerProperty()
	ownerProp.SetIRI(idIRI)
	publicKeyType.SetW3IDSecurityV1Owner(ownerProp)

	// publicKey publicKeyPem
	publicKeyPemProp := streams.NewW3IDSecurityV1PublicKeyPemProperty()
	publicKeyPemProp.Set(pubKey)
	publicKeyType.SetW3IDSecurityV1PublicKeyPem(publicKeyPemProp)

	publicKeyProp.AppendW3IDSecurityV1PublicKey(publicKeyType)
	p.SetW3IDSecurityV1PublicKey(publicKeyProp)
	return p, idIRI
}

func emptyInbox(actorID *url.URL) (vocab.ActivityStreamsOrderedCollection, error) {
	id, err := paths.IRIForActorID(paths.InboxPathKey, actorID)
	if err != nil {
		return nil, err
	}
	first, err := paths.IRIForActorID(paths.InboxFirstPathKey, actorID)
	if err != nil {
		return nil, err
	}
	last, err := paths.IRIForActorID(paths.InboxLastPathKey, actorID)
	if err != nil {
		return nil, err
	}
	return emptyOrderedCollection(id, first, last), nil
}

func emptyOutbox(actorID *url.URL) (vocab.ActivityStreamsOrderedCollection, error) {
	id, err := paths.IRIForActorID(paths.OutboxPathKey, actorID)
	if err != nil {
		return nil, err
	}
	first, err := paths.IRIForActorID(paths.OutboxFirstPathKey, actorID)
	if err != nil {
		return nil, err
	}
	last, err := paths.IRIForActorID(paths.OutboxLastPathKey, actorID)
	if err != nil {
		return nil, err
	}
	return emptyOrderedCollection(id, first, last), nil
}

func emptyOrderedCollection(id, first, last *url.URL) vocab.ActivityStreamsOrderedCollection {
	oc := streams.NewActivityStreamsOrderedCollection()
	// id
	idProp := streams.NewJSONLDIdProperty()
	idProp.SetIRI(id)
	oc.SetJSONLDId(idProp)

	// totalItems
	tiProp := streams.NewActivityStreamsTotalItemsProperty()
	tiProp.Set(0)
	oc.SetActivityStreamsTotalItems(tiProp)

	// orderedItems
	oiProp := streams.NewActivityStreamsOrderedItemsProperty()
	oc.SetActivityStreamsOrderedItems(oiProp)

	// first
	firstProp := streams.NewActivityStreamsFirstProperty()
	firstProp.SetIRI(first)
	oc.SetActivityStreamsFirst(firstProp)

	// last
	lastProp := streams.NewActivityStreamsLastProperty()
	lastProp.SetIRI(last)
	oc.SetActivityStreamsLast(lastProp)

	return oc
}