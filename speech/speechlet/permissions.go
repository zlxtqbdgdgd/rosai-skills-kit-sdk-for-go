package speechlet

type Skill struct {
	// eg. rosai1.ask.skill.78d08-0e59-4467-bdec-c119914974
	SkillId     string `json:"skillId"`
	AccessToken string `json:"accessToken,omitempty"`
}

func NewSkill(id string) *Skill {
	return &Skill{SkillId: id}
}

func (skill *Skill) WithAccessToken(token string) *Skill {
	skill.AccessToken = token
	return skill
}

type User struct {
	// UserId eg. rosai1.ask.account.M7ICKYFWA26TR4TUJ77JNT5BNSOIX5GGU4TJZ3TLE7SZL
	// HKB3OC4RBT34EINMGBMHL7POQQDPGTTOJVFUEMR4UOSZLGHPAS3MCOVQQO3WS6BAQMY7LV4GTHN
	// LXUJQWJNDI6HTDGVTAHVA6Q6DLM54TBRYAM3XCWVDFHYKPFSXIUN6BWW6X2CUMCTJKLCYZKUMK6
	// LNNIJXP7Q
	UserId      string       `json:"userId"`
	AppId       string       `json:"appId"`
	AccessToken string       `json:"accessToken,omitempty"`
	Permissions *Permissions `json:"permissions,omitempty"`
}

func NewUser(uid, appId string) *User {
	return &User{UserId: uid, AppId: appId}
}

func (u *User) WithAccessToken(token string) *User {
	u.AccessToken = token
	return u
}

type Permissions struct {
	ConsentToken string `json:"consentToken"`
}
