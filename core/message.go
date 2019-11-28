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
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/crypto/ethereum"
	"github.com/pdupub/go-pdu/crypto/pdu"
)

// Message is valid msg in pdu
type Message struct {
	SenderID  common.Hash       `json:"senderID"`
	Reference []*MsgReference   `json:"reference"`
	Value     *MsgValue         `json:"value"`
	Signature *crypto.Signature `json:"signature"`
}

// MsgReference is the msg before current msg
type MsgReference struct {
	SenderID common.Hash `json:"senderID"`
	MsgID    common.Hash `json:"msgID"`
}

// CreateMsg used to create a new msg by user in universe
func CreateMsg(user *User, value *MsgValue, priKey *crypto.PrivateKey, refs ...*MsgReference) (*Message, error) {

	v := &MsgValue{
		ContentType: value.ContentType,
		Content:     value.Content,
	}
	var rs []*MsgReference
	for _, r := range refs {
		rs = append(rs, &MsgReference{SenderID: r.SenderID, MsgID: r.MsgID})
	}

	msg := &Message{
		SenderID:  user.ID(),
		Reference: rs,
		Value:     v,
		Signature: nil,
	}
	switch priKey.Source {
	case pdu.SourceName:
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			return nil, err
		}
		sig, err := pdu.Sign(jsonMsg, priKey)
		if err != nil {
			return nil, err
		}
		sig.PubKey = nil
		msg.Signature = sig
		return msg, nil

	case ethereum.SourceName:
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			return nil, err
		}
		sig, err := ethereum.Sign(jsonMsg, priKey)
		if err != nil {
			return nil, err
		}
		sig.PubKey = nil
		msg.Signature = sig
		return msg, nil
	}

	return nil, errors.New("not support right now")
}

// VerifyMsg is used to valid the msg and the user
func VerifyMsg(msg Message) (bool, error) {
	signature := msg.Signature
	msg.Signature = nil
	switch signature.Source {
	case pdu.SourceName:
		jsonMsg, err := json.Marshal(&msg)
		if err != nil {
			return false, err
		}
		return pdu.Verify(jsonMsg, signature)

	case ethereum.SourceName:
		jsonMsg, err := json.Marshal(&msg)
		if err != nil {
			return false, err
		}
		return ethereum.Verify(jsonMsg, signature)
	}
	return false, errors.New("not support right now")

}

// ID is the id of msg based on content and author info
func (msg Message) ID() common.Hash {
	hash := sha256.New()
	hash.Reset()
	var ref string
	for _, r := range msg.Reference {
		ref += fmt.Sprintf("%v%v", r.SenderID, r.MsgID)
	}
	val := fmt.Sprintf("%v", msg.Value)
	hash.Write(append(append(msg.SenderID[:], ref...), val...))
	return common.Bytes2Hash(hash.Sum(nil))
}

// ParentsID return the parents id
// Parents are the message referenced by this Message
func (msg Message) ParentsID() []common.Hash {
	var parentsID []common.Hash
	for _, ref := range msg.Reference {
		parentsID = append(parentsID, ref.MsgID)
	}
	return parentsID
}

// ChildrenID always return nil for msg (useless in pdu)
func (msg Message) ChildrenID() []common.Hash {
	return nil
}
