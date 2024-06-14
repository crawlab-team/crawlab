package vcs

import "sync"

var GitMemStorages = sync.Map{}
var GitMemFileSystem = sync.Map{}
