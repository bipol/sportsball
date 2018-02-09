package models

//FieldPosition represents the position of the player on the pitch
type FieldPosition int

const (
	//RB = Right back
	RB FieldPosition = 0
	//RCM = Right Center Midfielder
	RCM FieldPosition = 1
	//LB = Left back
	LB FieldPosition = 2
	//LCM = Left Center Midfielder
	LCM FieldPosition = 3
	//CB = Center back
	CB FieldPosition = 4
	//CM = Center Midfielder
	CM FieldPosition = 5
	//F = Forward
	F FieldPosition = 6
)

//Player represents who is on a team
type Player struct {
	Position FieldPosition `json:"field_position"`
	FullName string        `json:"full_name"`
	ID       int           `json:"id"`
	Team     Team          `json:"team"`
}

//New will instantiate a new player
func (Player) New(fp FieldPosition, name string, team Team) (Player, error){
	return Player{fp, name, 0, team}, nil
}
