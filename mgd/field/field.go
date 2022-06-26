package field

import "go.mongodb.org/mongo-driver/bson/primitive"

const ID = "_id"

func StringToObject(id string) primitive.ObjectID {
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return hex
}

const (
	GroupBy    = "groupBy"
	Boundaries = "boundaries"
	Default    = "default"
	Output     = "output"
)

// $bucketAuto
const (
	// GroupBy     = "groupBy" // Declared
	Buckets = "buckets"
	// Output      = "output" // Declared
	Granularity = "granularity"
)

// $collStats
const (
	LatencyStats = "latencyStats"
	StorageStats = "storageStats"
	Count        = "count"
)

// $currentOp
const (
	AllUsers        = "allUsers"
	IdleConnections = "idleConnections"
	IdleCursors     = "idleCursors"
	IdleSessions    = "idleSessions"
	LocalOps        = "localOps"
)

// $geoNear
const (
	Near               = "near"
	DistanceField      = "distanceField"
	Spherical          = "spherical"
	MaxDistance        = "maxDistance"
	Query              = "query"
	DistanceMultiplier = "distanceMultiplier"
	IncludeLocs        = "includeLocs"
	UniqueDocs         = "uniqueDocs"
	MinDistance        = "minDistance"
	Key                = "key"
)

// $graphLookup
const (
	From                    = "from"
	StartWith               = "startWith"
	ConnectFromField        = "connectFromField"
	ConnectToField          = "connectToField"
	As                      = "as"
	MaxDepth                = "maxDepth"
	DepthField              = "depthField"
	RestrictSearchWithMatch = "restrictSearchWithMatch"
)

// $group
const (
// ID="_id" // Declared
)

// $listLocalSessions
const (
// AllUsers = "allUsers" // Declared
)

const (
	// From         = "from" // Declared
	LocalField   = "localField"
	ForeignField = "foreignField"
	// As           = "as" // Declared

	Let      = "let"
	Pipeline = "pipeline"
)

// $merge
const (
	Into = "into"
	On   = "on"
	// Let            = "let" // Declared
	WhenMatched    = "whenMatched"
	WhenNotMatched = "whenNotMatched"
)

// $replaceRoot
const (
	NewRoot = "newRoot"
)

// $sample
const (
	Size = "size"
)

// $unwind
const (
	Path                       = "path"
	IncludeArrayIndex          = "includeArrayIndex"
	PreserveNullAndEmptyArrays = "preserveNullAndEmptyArrays"
)
