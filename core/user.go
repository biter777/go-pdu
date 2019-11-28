// Copyright 2019 The PDU Authors
// This file is part of the PDU library.
//
// The PDU library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PDU library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PDU library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/core/rule"
	"github.com/pdupub/go-pdu/crypto"
	"math/big"
	"strconv"
)

const (
	rootMName     = "Adam"
	rootMDOBExtra = "Hello World!"
	rootFName     = "Eve"
	rootFDOBExtra = ";-)"
	male          = true
	female        = false
)

// User is the author of any msg in pdu
type User struct {
	Name     string   `json:"name"`
	DOBExtra string   `json:"extra"`
	Auth     *Auth    `json:"auth"`
	DOBMsg   *Message `json:"dobMsg"`
	LifeTime uint64   `json:"lifeTime"`
}

var (
	errContentTypeNotDOB = errors.New("content type is not TypeDOB")
)

// CreateRootUsers try to create two root users by public key
// One Male user and one female user,
func CreateRootUsers(key crypto.PublicKey) ([2]*User, error) {
	rootUsers := [2]*User{nil, nil}
	rootFUser := User{Name: rootFName, DOBExtra: rootFDOBExtra, Auth: &Auth{key}, DOBMsg: nil, LifeTime: rule.MaxLifeTime}
	if rootFUser.Gender() == female {
		rootUsers[0] = &rootFUser
	}
	rootMUser := User{Name: rootMName, DOBExtra: rootMDOBExtra, Auth: &Auth{key}, DOBMsg: nil, LifeTime: rule.MaxLifeTime}
	if rootMUser.Gender() == male {
		rootUsers[1] = &rootMUser
	}
	return rootUsers, nil
}

// CreateNewUser create new user by cosign message
// The msg must be signed by user in local user dag.
// Both parents must be in the local use dag.
// Both parents fit the nature rules.
// The BOD struct signed by both parents.
func CreateNewUser(universe *Universe, msg *Message) (*User, error) {
	if msg.Value.ContentType != TypeDOB {
		return nil, errContentTypeNotDOB
	}
	var dobContent DOBMsgContent
	if err := json.Unmarshal(msg.Value.Content, &dobContent); err != nil {
		return nil, err
	}
	newUser := dobContent.User
	newUser.DOBMsg = msg
	// calculate the life time of new user
	p0 := universe.userD.GetVertex(dobContent.Parents[0].UserID)
	if p0 == nil {
		return nil, errUserNotExist
	}
	maxParentLifeTime := p0.Value().(*User).LifeTime

	p1 := universe.userD.GetVertex(dobContent.Parents[1].UserID)
	if p1 == nil {
		return nil, errUserNotExist
	}
	if maxParentLifeTime < p1.Value().(*User).LifeTime {
		maxParentLifeTime = p1.Value().(*User).LifeTime
	}
	if maxParentLifeTime == rule.MortalLifetime {
		newUser.LifeTime = rule.MortalLifetime
	} else {
		newUser.LifeTime = maxParentLifeTime / rule.LifetimeReduceRate
	}

	return &newUser, nil
}

// ID return the vertex.id, related to parents and value of the vertex
// ID cloud use as address of user account
func (u User) ID() common.Hash {
	hash := sha256.New()
	hash.Reset()
	auth := fmt.Sprintf("%v", u.Auth)
	lifeTime := fmt.Sprintf("%v", u.LifeTime)
	var dobMsg string
	// todo : add init DOBMsg to rootUser
	// todo : so this condition can be deleted
	if u.DOBMsg != nil {
		dobMsg += fmt.Sprintf("%v", u.DOBMsg.SenderID)
		for _, v := range u.DOBMsg.Reference {
			dobMsg += fmt.Sprintf("%v%v", v.SenderID, v.MsgID)
		}
		dobMsg += fmt.Sprintf("%v%v%v", u.DOBMsg.Signature.Signature, u.DOBMsg.Signature.Source, u.DOBMsg.Signature.SigType)
		dobMsg += fmt.Sprintf("%v%v", u.DOBMsg.Value.Content, u.DOBMsg.Value.ContentType)
	}
	hash.Write(append(append(append(append([]byte(u.Name), u.DOBExtra...), auth...), dobMsg...), lifeTime...))
	return common.Bytes2Hash(hash.Sum(nil))
}

// Gender return the gender of user, true = male = end of ID is odd
func (u User) Gender() bool {
	hashID := u.ID()
	if uid := new(big.Int).SetBytes(hashID[:]); uid.Mod(uid, big.NewInt(2)).Cmp(big.NewInt(1)) == 0 {
		return true
	}
	return false
}

// Value return the vertex.value
func (u User) Value() interface{} {
	return nil
}

// ParentsID return the ID of user parents,
// res[0] should be the female parent (id end by even)
// res[1] should be the male parent (id end by odd)
func (u User) ParentsID() [2]common.Hash {
	var parentsID [2]common.Hash
	if u.DOBMsg != nil {
		// get parents from dobMsg
		var dobContent DOBMsgContent
		if err := json.Unmarshal(u.DOBMsg.Value.Content, &dobContent); err != nil {
			return parentsID
		}
		parentsID[0] = dobContent.Parents[0].UserID
		parentsID[1] = dobContent.Parents[1].UserID
	}
	return parentsID
}

// UnmarshalJSON is used to unmarshal json
func (u *User) UnmarshalJSON(input []byte) error {
	userMap := make(map[string]interface{})
	err := json.Unmarshal(input, &userMap)
	if err != nil {
		return err
	}
	u.Name = userMap["name"].(string)
	u.DOBExtra = userMap["dobExtra"].(string)
	u.LifeTime, err = strconv.ParseUint(userMap["lifeTime"].(string), 0, 64)
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(userMap["dobMsg"].(string)), &u.DOBMsg)
	json.Unmarshal([]byte(userMap["auth"].(string)), &u.Auth)

	return nil
}

// MarshalJSON marshal user to json
func (u *User) MarshalJSON() ([]byte, error) {
	userMap := make(map[string]interface{})
	userMap["name"] = u.Name
	userMap["dobExtra"] = u.DOBExtra
	userMap["lifeTime"] = fmt.Sprintf("%v", u.LifeTime)

	auth, err := json.Marshal(&u.Auth)
	if err != nil {
		return []byte{}, err
	}
	userMap["auth"] = string(auth)
	dobMsg, err := json.Marshal(&u.DOBMsg)
	if err != nil {
		return []byte{}, err
	}
	userMap["dobMsg"] = string(dobMsg)

	return json.Marshal(userMap)
}
