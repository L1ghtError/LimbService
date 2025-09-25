package user

type UserSchema struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Email       string `json:"email" bson:"email"`
	UserName    string `json:"username" bson:"username"`
	Password    []byte `json:"password,omitempty" bson:"password"`
	IsActivated bool   `json:"isActivated" bson:"isActivated"`
	// https://github.com/golang/go/issues/45669 omitzero is Likely Accept in Proposals
	// Do not expose it!
	ActivationLink [16]byte `bson:"activationLink"`
	Fullname       string   `json:"fullname"`
	Images         []string `json:"images,omitempty" bson:"images,omitempty"`
}
