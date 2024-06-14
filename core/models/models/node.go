package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Node struct {
	Id               primitive.ObjectID `json:"_id" bson:"_id"`
	Key              string             `json:"key" bson:"key"`
	Name             string             `json:"name" bson:"name"`
	Ip               string             `json:"ip" bson:"ip"`
	Port             string             `json:"port" bson:"port"`
	Mac              string             `json:"mac" bson:"mac"`
	Hostname         string             `json:"hostname" bson:"hostname"`
	Description      string             `json:"description" bson:"description"`
	IsMaster         bool               `json:"is_master" bson:"is_master"`
	Status           string             `json:"status" bson:"status"`
	Enabled          bool               `json:"enabled" bson:"enabled"`
	Active           bool               `json:"active" bson:"active"`
	ActiveTs         time.Time          `json:"active_ts" bson:"active_ts"`
	AvailableRunners int                `json:"available_runners" bson:"available_runners"`
	MaxRunners       int                `json:"max_runners" bson:"max_runners"`
}

func (n *Node) GetId() (id primitive.ObjectID) {
	return n.Id
}

func (n *Node) SetId(id primitive.ObjectID) {
	n.Id = id
}

func (n *Node) GetName() (name string) {
	return n.Name
}

func (n *Node) SetName(name string) {
	n.Name = name
}

func (n *Node) GetDescription() (description string) {
	return n.Description
}

func (n *Node) SetDescription(description string) {
	n.Description = description
}

func (n *Node) GetKey() (key string) {
	return n.Key
}

func (n *Node) GetIsMaster() (ok bool) {
	return n.IsMaster
}

func (n *Node) GetActive() (active bool) {
	return n.Active
}

func (n *Node) SetActive(active bool) {
	n.Active = active
}

func (n *Node) SetActiveTs(activeTs time.Time) {
	n.ActiveTs = activeTs
}

func (n *Node) GetStatus() (status string) {
	return n.Status
}

func (n *Node) SetStatus(status string) {
	n.Status = status
}

func (n *Node) GetEnabled() (enabled bool) {
	return n.Enabled
}

func (n *Node) SetEnabled(enabled bool) {
	n.Enabled = enabled
}

func (n *Node) GetAvailableRunners() (runners int) {
	return n.AvailableRunners
}

func (n *Node) SetAvailableRunners(runners int) {
	n.AvailableRunners = runners
}

func (n *Node) GetMaxRunners() (runners int) {
	return n.MaxRunners
}

func (n *Node) SetMaxRunners(runners int) {
	n.MaxRunners = runners
}

func (n *Node) IncrementAvailableRunners() {
	n.AvailableRunners++
}

func (n *Node) DecrementAvailableRunners() {
	n.AvailableRunners--
}

type NodeList []Node

func (l *NodeList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
