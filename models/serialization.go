package models

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/go-fed/activity/pub"
	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
)

// marshal takes any ActivityStreams type and serializes it to JSON.
func marshal(v vocab.Type) (b []byte, err error) {
	var m map[string]interface{}
	m, err = streams.Serialize(v)
	if err != nil {
		return
	}
	b, err = json.Marshal(m)
	if err != nil {
		return
	}
	return
}

// unmarhsal attempts to deserialize JSON bytes into a value.
func unmarshal(maybeByte, v interface{}) error {
	b, ok := maybeByte.([]byte)
	if !ok {
		return errors.New("failed to assert scan to []byte type")
	}
	return json.Unmarshal(b, v)
}

// singleRow allows *sql.Rows to be treated as *sql.Row
type singleRow interface {
	Scan(dest ...interface{}) error
}

// enforceOneRow ensures that there is only one row in the *sql.Rows
//
// Normally, the SQL operations that assume a single row is being returned by
// the database take only the first row and then discard the rest of the rows
// silently. I would rather we know when our expectations are being violated or
// when the database constraints do not match the expected application logic,
// than silently retrieve an arbitrarily row (since the first one grabbed is
// returned arbitrarily, database-and-driver-dependent).
func enforceOneRow(r *sql.Rows, debugname string, fn func(r singleRow) error) error {
	var n int
	for r.Next() {
		if n > 0 {
			return fmt.Errorf("%s: multiple database rows retrieved when enforcing one row", debugname)
		}
		err := fn(singleRow(r))
		if err != nil {
			return err
		}
		n++
	}
	return r.Err()
}

var _ driver.Valuer = OnFollowBehavior(0)
var _ sql.Scanner = (*OnFollowBehavior)(nil)

// OnFollowBehavior is a wrapper around pub.OnFollowBehavior type that also
// knows how to serialize and deserialize itself for SQL database drivers in a
// more readable manner.
type OnFollowBehavior pub.OnFollowBehavior

const (
	onFollowAlwaysAccept = "ALWAYS_ACCEPT"
	onFollowAlwaysReject = "ALWAYS_REJECT"
	onFollowManual       = "MANUAL"
)

func (o OnFollowBehavior) Value() (driver.Value, error) {
	switch pub.OnFollowBehavior(o) {
	case pub.OnFollowAutomaticallyAccept:
		return onFollowAlwaysAccept, nil
	case pub.OnFollowAutomaticallyReject:
		return onFollowAlwaysReject, nil
	case pub.OnFollowDoNothing:
		fallthrough
	default:
		return onFollowManual, nil
	}
}

func (o *OnFollowBehavior) Scan(src interface{}) error {
	s, ok := src.(string)
	if !ok {
		return errors.New("failed to assert scan to string type")
	}
	switch s {
	case onFollowAlwaysAccept:
		*o = OnFollowBehavior(pub.OnFollowAutomaticallyAccept)
	case onFollowAlwaysReject:
		*o = OnFollowBehavior(pub.OnFollowAutomaticallyReject)
	case onFollowManual:
		fallthrough
	default:
		*o = OnFollowBehavior(pub.OnFollowDoNothing)
	}
	return nil
}

var _ driver.Valuer = ActivityStreams{nil}
var _ sql.Scanner = &ActivityStreams{nil}

// ActivityStreams is a wrapper around any ActivityStreams type that also
// knows how to serialize and deserialize itself for SQL database drivers.
type ActivityStreams struct {
	vocab.Type
}

func (a ActivityStreams) Value() (driver.Value, error) {
	return marshal(a)
}

func (a *ActivityStreams) Scan(src interface{}) error {
	var m map[string]interface{}
	if err := unmarshal(src, &m); err != nil {
		return err
	}
	var err error
	a.Type, err = streams.ToType(context.Background(), m)
	return err
}

var _ driver.Valuer = ActivityStreamsPerson{nil}
var _ sql.Scanner = &ActivityStreamsPerson{nil}

// ActivityStreamsPerson is a wrapper around the ActivityStreams type that also
// knows how to serialize and deserialize itself for SQL database drivers.
type ActivityStreamsPerson struct {
	vocab.ActivityStreamsPerson
}

func (a ActivityStreamsPerson) Value() (driver.Value, error) {
	return marshal(a)
}

func (a *ActivityStreamsPerson) Scan(src interface{}) error {
	var m map[string]interface{}
	if err := unmarshal(src, &m); err != nil {
		return err
	}
	res, err := streams.NewJSONResolver(func(ctx context.Context, p vocab.ActivityStreamsPerson) error {
		a.ActivityStreamsPerson = p
		return nil
	})
	if err != nil {
		return err
	}
	return res.Resolve(context.Background(), m)
}

var _ driver.Valuer = ActivityStreamsOrderedCollection{nil}
var _ sql.Scanner = &ActivityStreamsOrderedCollection{nil}

// ActivityStreamsOrderedCollection is a wrapper around the ActivityStreams type
// that also knows how to serialize and deserialize itself for SQL database
// drivers.
type ActivityStreamsOrderedCollection struct {
	vocab.ActivityStreamsOrderedCollection
}

func (a ActivityStreamsOrderedCollection) Value() (driver.Value, error) {
	return marshal(a)
}

func (a *ActivityStreamsOrderedCollection) Scan(src interface{}) error {
	var m map[string]interface{}
	if err := unmarshal(src, &m); err != nil {
		return err
	}
	res, err := streams.NewJSONResolver(func(ctx context.Context, oc vocab.ActivityStreamsOrderedCollection) error {
		a.ActivityStreamsOrderedCollection = oc
		return nil
	})
	if err != nil {
		return err
	}
	return res.Resolve(context.Background(), m)
}

var _ driver.Valuer = ActivityStreamsOrderedCollectionPage{nil}
var _ sql.Scanner = &ActivityStreamsOrderedCollectionPage{nil}

// ActivityStreamsOrderedCollectionPage is a wrapper around the ActivityStreams
// type that also knows how to serialize and deserialize itself for SQL database
// drivers.
type ActivityStreamsOrderedCollectionPage struct {
	vocab.ActivityStreamsOrderedCollectionPage
}

func (a ActivityStreamsOrderedCollectionPage) Value() (driver.Value, error) {
	return marshal(a)
}

func (a *ActivityStreamsOrderedCollectionPage) Scan(src interface{}) error {
	var m map[string]interface{}
	if err := unmarshal(src, &m); err != nil {
		return err
	}
	res, err := streams.NewJSONResolver(func(ctx context.Context, oc vocab.ActivityStreamsOrderedCollectionPage) error {
		a.ActivityStreamsOrderedCollectionPage = oc
		return nil
	})
	if err != nil {
		return err
	}
	return res.Resolve(context.Background(), m)
}

var _ driver.Valuer = NullDuration{}
var _ sql.Scanner = &NullDuration{}

// NullDuration can handle nullable time.Duration values in the database.
type NullDuration struct {
	Duration time.Duration
	Valid    bool
}

func (n NullDuration) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Duration, nil
}

func (n *NullDuration) Scan(src interface{}) error {
	if src == nil {
		n.Duration, n.Valid = 0, false
		return nil
	}
	t, ok := src.(int64)
	if !ok {
		return errors.New("failed to assert scan to int64 type")
	}
	n.Duration, n.Valid = time.Duration(t), true
	return nil
}

var _ driver.Valuer = URL{}
var _ sql.Scanner = &URL{}

// URL handles serializing/deserializing *url.URL types into databases.
type URL struct {
	*url.URL
}

func (u URL) Value() (driver.Value, error) {
	return u.URL.String(), nil
}

func (u *URL) Scan(src interface{}) error {
	s, ok := src.(string)
	if !ok {
		return errors.New("failed to assert scan to string type")
	}
	var err error
	u.URL, err = url.Parse(s)
	return err
}