package internal


type Room struct {
	Id string
	IsStart, IsEnd bool
	Links map[string]*Room
	Visited bool
	Parent *Room
	ParentInbfs *Room
	Usedinpath bool
	X, Y int
	AllowToJump bool
	CameFromBacktraking bool
    Forclinks  []string
	Allreadypathfound bool
	AntInRoom *Ant
	AntThatwilluseThisRoom []*Ant
	IndexOfAntForNextTime int
	Active bool
}

type Graph struct {
	Rooms map[string]*Room
	StartRoom, EndRoom *Room 
	Paths []*Path
	AntCount int
	AntCountInStartRoom int
	AntCountInEndRoom int
	File []string
}

type Path struct {
	Rooms []*Room
	NumberOfAntsToHold int
	CountAntWalking int
	CountAntWaiting int
	CountAntAllreadyEnterThePath int
	CountOfAntReachTheEnd int
	AntsInHome []*Ant
	NextAntIdToEnterThePath *Ant
	TheAntInTheFronT *Ant
	Len int
	IndexxOfreachedRoom  int

}

type Ant struct  {
	ID       int
	Path     []*Room
	Step int 
	ActualRomm *Room
	NextRoomToWalktIn *Room
	walking bool
	Finished bool
}

type Linkstodelet struct {
	Linksdelt [][2]string
}


