package operator

const (
	Abs      = "$abs"
	Add      = "$add"
	Ceil     = "$ceil"
	Divide   = "$divide"
	Exp      = "$exp"
	Floor    = "$floor"
	Ln       = "$ln"
	Log      = "$log"
	Log10    = "$log10"
	Multiply = "$multiply"
	Pow      = "$pow"
	Round    = "$round"
	Sqrt     = "$sqrt"
	Subtract = "$subtract"
	Trunc    = "$trunc"
)
const (
	ArrayToObject = "$arrayToObject"
	ConcatArrays  = "$concatArrays"
	Filter        = "$filter"
	IndexOfArray  = "$indexOfArray"
	IsArray       = "$isArray"
	Map           = "$map"
	ObjectToArray = "$objectToArray"
	Range         = "$range"
	Reduce        = "$reduce"
	ReverseArray  = "$reverseArray"
	Zip           = "$zip"
)

const (
	Cond   = "$cond"
	IfNull = "$ifNull"
	Switch = "$switch"
)

const (
	DateFromParts  = "$dateFromParts"
	DateFromString = "$dateFromString"
	DateToParts    = "$dateToParts"
	DateToString   = "$dateToString"
	DayOfMonth     = "$dayOfMonth"
	DayOfWeek      = "$dayOfWeek"
	DayOfYear      = "$dayOfYear"
	Hour           = "$hour"
	IsoDayOfWeek   = "$isoDayOfWeek"
	IsoWeek        = "$isoWeek"
	IsoWeekYear    = "$isoWeekYear"
	Millisecond    = "$millisecond"
	Minute         = "$minute"
	Month          = "$month"
	Second         = "$second"
	ToDate         = "$toDate"
	Week           = "$week"
	Year           = "$year"
)

const (
	Literal = "$literal"
)

const (
	MergeObjects = "$mergeObjects"
)

const (
	AllElementsTrue = "$allElementsTrue"
	AnyElementTrue  = "$anyElementTrue"
	SetDifference   = "$setDifference"
	SetEquals       = "$setEquals"
	SetIntersection = "$setIntersection"
	SetIsSubset     = "$setIsSubset"
	SetUnion        = "$setUnion"
)

const (
	Concat       = "$concat"
	IndexOfBytes = "$indexOfBytes"
	IndexOfCP    = "$indexOfCP"
	Ltrim        = "$ltrim"
	RegexFind    = "$regexFind"
	RegexFindAll = "$regexFindAll"
	RegexMatch   = "$regexMatch"
	Rtrim        = "$rtrim"
	Split        = "$split"
	StrLenBytes  = "$strLenBytes"
	StrLenCP     = "$strLenCP"
	Strcasecmp   = "$strcasecmp"
	Substr       = "$substr"
	SubstrBytes  = "$substrBytes"
	SubstrCP     = "$substrCP"
	ToLower      = "$toLower"
	ToString     = "$toString"
	Trim         = "$trim"
	ToUpper      = "$toUpper"
)
const (
	Sin              = "$sin"
	Cos              = "$cos"
	Tan              = "$tan"
	Asin             = "$asin"
	Acos             = "$acos"
	Atan             = "$atan"
	Atan2            = "$atan2"
	Asinh            = "$asinh"
	Acosh            = "$acosh"
	Atanh            = "$atanh"
	DegreesToRadians = "$degreesToRadians"
	RadiansToDegrees = "$radiansToDegrees"
)

const (
	Convert = "$convert"
	ToBool  = "$toBool"
	//ToDate     = "$toDate" // Declared
	ToDecimal  = "$toDecimal"
	ToDouble   = "$toDouble"
	ToInt      = "$toInt"
	ToLong     = "$toLong"
	ToObjectID = "$toObjectId"
	//ToString   = "$toString" // Declared
	//Type       = "$type" // Declared
)
const (
	Avg        = "$avg"
	First      = "$first"
	Last       = "$last"
	StdDevPop  = "$stdDevPop"
	StdDevSamp = "$stdDevSamp"
	Sum        = "$sum"
)

const (
	AddFields      = "$addFields"
	Bucket         = "$bucket"
	BucketAuto     = "$bucketAuto"
	CollStats      = "$collStats"
	Count          = "$count"
	Facet          = "$facet"
	GeoNear        = "$geoNear"
	GraphLookup    = "$graphLookup"
	Group          = "$group"
	IndexStats     = "$indexStats"
	Limit          = "$limit"
	ListSessions   = "$listSessions"
	Lookup         = "$lookup"
	Match          = "$match"
	Merge          = "$merge"
	Out            = "$out"
	PlanCacheStats = "$planCacheStats"
	Project        = "$project"
	Redact         = "$redact"
	ReplaceRoot    = "$replaceRoot"
	ReplaceWith    = "$replaceWith"
	Sample         = "$sample"
	// Set            = "$set" // Declared
	Skip = "$skip"
	// Sort           = "$sort" // Declared
	SortByCount = "$sortByCount"
	// Unset          = "$unset" // Declared
	Unwind = "$unwind"
)

const (
	CurrentOp         = "$currentOp"
	ListLocalSessions = "$listLocalSessions"
)

const (
	Eq  = "$eq"
	Gt  = "$gt"
	Gte = "$gte"
	In  = "$in"
	Lt  = "$lt"
	Lte = "$lte"
	Ne  = "$ne"
	Nin = "$nin"
)

const (
	And = "$and"
	Not = "$not"
	Nor = "$nor"
	Or  = "$or"
)

const (
	Exists = "$exists"
	Type   = "$type"
)

const (
	Expr       = "$expr"
	JSONSchema = "$jsonSchema"
	Mod        = "$mod"
	Regex      = "$regex"
	Text       = "$text"
	Where      = "$where"
)

const (
	GeoIntersects = "$geoIntersects"
	GeoWithin     = "$geoWithin"
	Near          = "$near"
	NearSphere    = "$nearSphere"
)

const (
	All       = "$all"
	ElemMatch = "$elemMatch"
	Size      = "$size"
)

const (
	BitsAllClear = "$bitsAllClear"
	BitsAllSet   = "$bitsAllSet"
	BitsAnyClear = "$bitsAnyClear"
	BitsAnySet   = "$bitsAnySet"
)

const (
	Comment = "$comment"
)

const (
	Dollar = "$"
	Meta   = "$meta"
	Slice  = "$slice"
)

const (
	Explain = "$explain"
	Hint    = "$hint"
	// Max     = "$max" // Declared
	MaxTimeMS = "$maxTimeMS"
	// Min       = "$min" // Declared
	OrderBy     = "$orderby"
	Query       = "$query"
	ReturnKey   = "$returnKey"
	ShowDiskLoc = "$showDiskLoc"
)

const (
	Natural = "$natural"
)
const (
	CurrentDate = "$currentDate"
	Inc         = "$inc"
	Min         = "$min"
	Max         = "$max"
	Mul         = "$mul"
	Rename      = "$rename"
	Set         = "$set"
	SetOnInsert = "$setOnInsert"
	Unset       = "$unset"
)

const (
	AddToSet = "$addToSet"
	Pop      = "$pop"
	Pull     = "$pull"
	Push     = "$push"
	PullAll  = "$pullAll"
)

const (
	Each     = "$each"
	Position = "$position"
	// Slice    = "$slice" // Declared
	Sort = "$sort"
)

const (
	Bit = "$bit"
)
