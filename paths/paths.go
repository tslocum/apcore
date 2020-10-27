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

package paths

import (
	"fmt"
	"net/url"
)

const (
	pubKeyFragmentPrefix = "key-"
	usersRoute           = "/users"
)

type Paths struct {
	scheme string
	host   string
}

func NewPaths(scheme, host string) *Paths {
	return &Paths{
		scheme: scheme,
		host:   host,
	}
}

func (p *Paths) getBase() string {
	return fmt.Sprintf("%s://%s", p.scheme, p.host)
}

func (p *Paths) UsersPath(userUUID string) (u *url.URL, err error) {
	u, err = url.Parse(p.getBase() + usersRoute)
	return
}

func (p *Paths) PublicKeyPath(userUUID, keyUUID string) (u *url.URL, err error) {
	u, err = p.UsersPath(userUUID)
	if err != nil {
		return
	}
	u.Fragment = fmt.Sprintf("%s%s", pubKeyFragmentPrefix, keyUUID)
	return
}