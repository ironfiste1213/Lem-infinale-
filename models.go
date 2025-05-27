package internal


type Room struct {
	Id string
	IsStart, IsEnd bool
	Links map[string]*Room
	Visited bool
	Parent *Room
	Usedinpath bool
	X, Y int
	AllowToJump bool
	CameFromBacktraking bool
	Blocked bool
    Forclinks  []string
}

type Graph struct {
	Rooms map[string]*Room
	StartRoom, EndRoom *Room 
	Paths []*Path
	AntCount int
}

type Path struct {
	Rooms []*Room
	IsFollowingPathBackward bool
	Linksdelet []*Room
	Len int
}

type Ant struct  {
	Id int 
	Path *Path
}

type Linkstodelet struct {
	Linksdelt [][2]string
}


