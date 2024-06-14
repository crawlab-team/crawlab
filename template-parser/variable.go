package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"strings"
)

var userRegexp, _ = regexp.Compile("^user(?:\\[(\\w+)\\])?$")
var extRegexp, _ = regexp.Compile("^(\\w+)?:(\\w+)")

type Variable struct {
	root   interface{}
	tokens []string
	doc    bson.M
}

func (v *Variable) GetValue() (value interface{}, err error) {
	return v.getNodeByIndex(len(v.tokens) - 1)
}

func (v *Variable) getNodeByIndex(index int) (result interface{}, err error) {
	node := v.doc
	for i := 0; i < index && i < len(v.tokens); i++ {
		nextIndex := i + 1
		if nextIndex < len(v.tokens)-1 {
			// root or intermediate node
		} else {
			// value
			return v.getNextValue(node, i), nil
		}

		// next node
		node, err = v.getNextNode(node, i)
		if err != nil {
			return nil, err
		}
		if node == nil {
			return nil, nil
		}
	}
	return node, nil
}

func (v *Variable) getNextNode(currentNode bson.M, currentIndex int) (nextNode bson.M, err error) {
	// next index and token
	nextIndex := currentIndex + 1
	nextToken := v.tokens[nextIndex]

	// attempt to get attribute in current node
	nextNodeRes, ok := currentNode[nextToken]
	if ok {
		nextNode, ok = nextNodeRes.(bson.M)
		if ok {
			return nextNode, nil
		}
	}

	// if next token is user or user[<action>]
	if userRegexp.MatchString(nextToken) {
		return v.getNextNodeUserAction(currentNode, nextIndex, nextToken)
	}

	// if next token is <model>:<ext>
	if extRegexp.MatchString(nextToken) {
		return v.getNextNodeExt(currentNode, nextIndex, nextToken)
	}

	// model
	model := nextToken

	// next id
	nextId, err := v._getNodeModelId(currentNode, model, nextIndex)
	if err != nil {
		return nil, err
	}

	// mongo collection name
	colName := fmt.Sprintf("%ss", model)

	// get next node from mongo collection
	if err := mongo.GetMongoCol(colName).FindId(nextId).One(&nextNode); err != nil {
		return nil, err
	}

	return nextNode, nil
}

func (v *Variable) getNextNodeUserAction(currentNode bson.M, nextIndex int, nextToken string) (nextNode bson.M, err error) {
	// action
	matches := userRegexp.FindStringSubmatch(nextToken)
	action := "create"
	if len(matches) > 1 && matches[1] != "" {
		action = matches[1]
	}

	// id
	id, err := v._getCurrentNodeId(currentNode, nextIndex)
	if err != nil {
		return nil, err
	}

	// artifact
	var artifact bson.M
	if err := mongo.GetMongoCol("artifacts").FindId(id).One(&artifact); err != nil {
		return nil, err
	}

	// artifact._sys
	sysRes, ok := artifact["_sys"]
	if !ok {
		return nil, errors.New(fmt.Sprintf("_sys not exists in artifact of %s", strings.Join(v.tokens[:nextIndex], ".")))
	}
	sys, ok := sysRes.(bson.M)
	if !ok {
		return nil, errors.New(fmt.Sprintf("_sys is invalid in artifact of %s", strings.Join(v.tokens[:nextIndex], ".")))
	}

	// <action>_uid
	uidRes, _ := sys[action+"_uid"]
	uid, _ := uidRes.(primitive.ObjectID)
	if err := mongo.GetMongoCol("users").FindId(uid).One(&nextNode); err != nil {
		return nil, err
	}

	return nextNode, nil
}

func (v *Variable) getNextNodeExt(currentNode bson.M, nextIndex int, nextToken string) (nextNode bson.M, err error) {
	matches := extRegexp.FindStringSubmatch(nextToken)
	var model, ext string
	var mode int
	if matches[1] == "" {
		// :<model>_<ext>
		mode = 0
		arr := strings.Split(matches[2], "_")
		model = arr[0]
		ext = arr[1]
	} else {
		// <model>:<ext>
		mode = 1
		model = matches[1]
		ext = matches[2]
	}

	// id
	var id primitive.ObjectID
	switch mode {
	case 0:
		// id = currentNode._id
		id, err = v._getCurrentNodeId(currentNode, nextIndex)
		if err != nil {
			return nil, err
		}
	case 1:
		// id = currentNode.<model>_id
		id, err = v._getNodeModelId(currentNode, model, nextIndex)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("invalid mode: %d", mode))
	}

	// ext
	colName := fmt.Sprintf("%s_%ss", model, ext)
	if err := mongo.GetMongoCol(colName).FindId(id).One(&nextNode); err != nil {
		return nil, err
	}

	return nextNode, nil
}

func (v *Variable) getNextValue(currentNode bson.M, currentIndex int) (nextValue interface{}) {
	// next index and token
	nextIndex := currentIndex + 1
	nextToken := v.tokens[nextIndex]

	// next value
	nextValue, _ = currentNode[nextToken]

	return nextValue
}

func (v *Variable) _getCurrentNodeId(currentNode bson.M, nextIndex int) (id primitive.ObjectID, err error) {
	idRes, ok := currentNode["_id"]
	if !ok {
		return id, errors.New(fmt.Sprintf("%s is not available in %s", "_id", strings.Join(v.tokens[:nextIndex], ".")))
	}
	switch idRes.(type) {
	case string:
		idStr, ok := idRes.(string)
		if !ok {
			return id, errors.New(fmt.Sprintf("%s is not ObjectId in %s", "_id", strings.Join(v.tokens[:nextIndex], ".")))
		}
		id, err = primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return id, errors.New(fmt.Sprintf("%s is not ObjectId in %s", "_id", strings.Join(v.tokens[:nextIndex], ".")))
		}
		return id, nil
	case primitive.ObjectID:
		return idRes.(primitive.ObjectID), nil
	default:
		return id, errors.New(fmt.Sprintf("%s is not ObjectId in %s", "_id", strings.Join(v.tokens[:nextIndex], ".")))
	}
}

func (v *Variable) _getNodeModelId(currentNode bson.M, model string, nextIndex int) (id primitive.ObjectID, err error) {
	nextIdKey := fmt.Sprintf("%s_id", model)
	nextIdRes, ok := currentNode[nextIdKey]
	if !ok {
		return id, errors.New(fmt.Sprintf("%s is not available in %s", nextIdKey, strings.Join(v.tokens[:nextIndex], ".")))
	}
	nextId, ok := nextIdRes.(primitive.ObjectID)
	if !ok {
		nextIdStr, ok := nextIdRes.(string)
		if !ok {
			return id, errors.New(fmt.Sprintf("%s is not ObjectId in %s", nextIdKey, strings.Join(v.tokens[:nextIndex], ".")))
		}
		nextId, err = primitive.ObjectIDFromHex(nextIdStr)
		if err != nil {
			return id, err
		}
	}
	if nextId.IsZero() {
		return id, nil
	}

	return nextId, nil
}

func NewVariable(root interface{}, placeholder string) (v *Variable, err error) {
	// validate
	if placeholder == "" {
		return nil, errors.New("empty placeholder")
	}
	if !strings.HasPrefix(placeholder, "$") {
		return nil, errors.New("not start with $")
	}

	// tokens
	tokens := strings.Split(placeholder, ".")

	// document
	data, err := json.Marshal(root)
	if err != nil {
		return nil, err
	}
	var doc bson.M
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}

	v = &Variable{
		root:   root,
		tokens: tokens,
		doc:    doc,
	}

	return v, nil
}
