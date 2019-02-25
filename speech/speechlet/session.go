package speechlet

import (
	"bytes"
	"encoding/base32"
	"encoding/gob"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/slu"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/util"

	"github.com/garyburd/redigo/redis"
)

const (
	SSK_UPDATED_INTENT string = "updatedIntent"
)

func init() {
	gob.Register(slu.Intent{})
	gob.Register(map[string]*slu.Intent{})
}

type Session struct {
	New bool `json:"new`
	// eg. rosai1.pudding-api.session.302948ed-125e-472d-9df4-84ff69085
	ID         string                 `json:"id"`
	Attributes map[string]interface{} `json:"attributes"`
	Skill      *Skill                 `json:"skill,omitempty"`
	User       *User                  `json:"user,omitempty"`
	Device     *Device                `json:"device,omitempty"`
}

func NewSession(userId, appId, deviceId, skillId string) *Session {
	return &Session{
		New:        true,
		ID:         GenSessionId(userId, appId, deviceId, skillId),
		Attributes: make(map[string]interface{}),
		Skill:      NewSkill(skillId),
		User:       NewUser(userId, appId),
		Device:     NewDevice(deviceId),
	}
}

func (ss *Session) SetNew(b bool) {
	ss.New = b
}

func (ss *Session) WithAttr(name string, value interface{}) *Session {
	if ss.Attributes == nil {
		ss.Attributes = make(map[string]interface{})
	}
	ss.Attributes[name] = value
	return ss
}

func (ss *Session) GetAttrStringValue(name string) string {
	if len(ss.Attributes) == 0 {
		return ""
	}
	if v, ok := ss.Attributes[name]; ok {
		if vv, ok := v.(string); ok {
			return vv
		}
	}
	return ""
}

func (ss *Session) GetAttrIntValue(name string) int {
	if len(ss.Attributes) == 0 {
		return -1
	}
	if v, ok := ss.Attributes[name]; ok {
		if vv, ok := v.(int); ok {
			return vv
		} else if vv, ok := v.(int64); ok {
			return int(vv)
		}
	}
	return -1
}

func (ss *Session) GetAttrFloatValue(name string) float64 {
	if len(ss.Attributes) == 0 {
		return 0.0
	}
	if v, ok := ss.Attributes[name]; ok {
		if vv, ok := v.(float64); ok {
			return vv
		}
	}
	return 0.0
}

func (ss *Session) WithUpdatedIntent(intent *slu.Intent) *Session {
	if intent == nil {
		return ss
	}
	if ss.Attributes == nil {
		ss.Attributes = make(map[string]interface{})
	}
	intents, ok := ss.Attributes[SSK_UPDATED_INTENT]
	if !ok || intents == nil {
		ss.Attributes[SSK_UPDATED_INTENT] = make(map[string]*slu.Intent)
	}
	v, ok := ss.Attributes[SSK_UPDATED_INTENT].(map[string]*slu.Intent)
	if ok {
		v[intent.Name] = intent
		ss.Attributes[SSK_UPDATED_INTENT] = v
	}
	return ss
}

func (ss *Session) GetUpdatedIntent(name string) *slu.Intent {
	if v, ok := ss.Attributes[SSK_UPDATED_INTENT]; ok {
		if v, ok := v.(map[string]*slu.Intent); ok {
			if v, ok := v[name]; ok {
				return v
			}
		}
	}
	return nil
}

func (ss *Session) MergeIntent(obj *slu.Intent) *Session {
	if obj == nil {
		return ss
	}
	origIntent := ss.GetUpdatedIntent(obj.Name)
	if origIntent == nil {
		origIntent = slu.NewIntent(obj.Name).WithSubName(obj.SubName)
		ss.WithUpdatedIntent(origIntent)
	}
	if !origIntent.Merge(obj) {
		log.Println("session merged failed")
	}
	return ss
}

func (ss *Session) ClearAllIntents() {
	ss.Attributes[SSK_UPDATED_INTENT] = nil
}

func FetchSessionFromHistory(userId, appId, deviceId, skillId string) (*Session, error) {
	ssStore := GetRediSession()
	ss, err := ssStore.Fetch(userId, appId, deviceId, skillId)
	if err != nil || ss == nil {
		log.Printf("fetch session[userId: %s, appId: %s, deviceId: %s, skillId: %s] error: %s",
			userId, appId, deviceId, skillId, err)
	}
	return ss, nil
}

func PushSessionToCache(ss *Session) error {
	ssStore := GetRediSession()
	return ssStore.Save(ss)
}

// not safe for concurrent use
type RediSession struct {
	Pool       *redis.Pool
	maxAge     int // default Redis TTL for a maxAge == 0 session
	maxLength  int
	keyPrefix  string
	serializer SessionSerializer
}

var (
	redisPool       *redis.Pool
	redisPoolOnce   sync.Once
	rediSession     *RediSession
	rediSessionOnce sync.Once
)

func GetRedisPool() *redis.Pool {
	redisPoolOnce.Do(func() {
		addr, passwd, db, err := util.GetRedisConf()
		if err != nil {
			log.Println(err)
		}
		redisPool = &redis.Pool{
			MaxIdle:     20,
			IdleTimeout: 300 * time.Second,
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
			Dial: func() (redis.Conn, error) {
				return dialWithDB(addr, passwd, db)
			},
		}
	})
	return redisPool
}

// NewRediSession returns a new RediSession.
// size: maximum number of idle connections.
func GetRediSession() *RediSession {
	rediSessionOnce.Do(func() {
		rediSession = &RediSession{
			// http://godoc.org/github.com/garyburd/redigo/redis#Pool
			Pool:       GetRedisPool(),
			maxAge:     300,
			maxLength:  65536,
			keyPrefix:  "rosai.sdk.session.",
			serializer: GobSerializer{},
		}
		_, err := rediSession.ping()
		if err != nil {
			log.Println(err)
		}
	})
	return rediSession
}

func GenSessionId(userId, appId, deviceId, skillId string) string {
	return strings.TrimRight(base32.StdEncoding.EncodeToString(
		[]byte(userId+appId+deviceId+skillId)), "=")
}

// Get returns a session for the userId name
func (s *RediSession) New(userId, appId, deviceId, skillId string) *Session {
	return NewSession(userId, appId, deviceId, skillId)
}

// Get returns a session for the given userId after adding it to redis.
func (s *RediSession) Fetch(userId, appId, deviceId, skillId string) (*Session, error) {
	session := NewSession(userId, appId, deviceId, skillId)
	ok, err := s.load(session)
	session.New = !(err == nil && ok) // not new if no error and data available
	return session, err
}

func (s *RediSession) Save(ss *Session) error {
	return s.save(ss)
}

func (s *RediSession) Drop(userId, appId, deviceId, skillId string) error {
	id := GenSessionId(userId, appId, deviceId, skillId)
	return s.drop(id)
}

// Close closes the underlying *redis.Pool
func (s *RediSession) Close() error {
	return s.Pool.Close()
}

// SetMaxLength sets RediSession.maxLength if the `l` argument is greater or equal 0
// maxLength restricts the maximum length of new sessions to l.
// If l is 0 there is no limit to the size of a session, use with caution.
// The default for a new RediSession is 65536. Redis allows for max.
// value sizes of up to 512MB (http://redis.io/topics/data-types)
// Default: 65536
func (s *RediSession) SetMaxLength(l int) {
	if l >= 0 {
		s.maxLength = l
	}
}

func (s *RediSession) SetMaxAge(sec int) {
	if sec >= 0 {
		s.maxAge = sec
	}
}

// SetKeyPrefix set the prefix
func (s *RediSession) SetKeyPrefix(p string) {
	s.keyPrefix = p
}

// SetSerializer sets the serializer
func (s *RediSession) SetSerializer(ss SessionSerializer) {
	s.serializer = ss
}

func dial(address, password string) (redis.Conn, error) {
	c, err := redis.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, err
}

func dialWithDB(address, password, DB string) (redis.Conn, error) {
	c, err := dial(address, password)
	if err != nil {
		return nil, err
	}
	if DB == "" {
		return c, nil
	}
	if _, err := c.Do("SELECT", DB); err != nil {
		c.Close()
		return nil, err
	}
	return c, err
}

// ping does an internal ping against a server to check if it is alive.
func (s *RediSession) ping() (bool, error) {
	conn := s.Pool.Get()
	defer conn.Close()
	data, err := conn.Do("PING")
	if err != nil || data == nil {
		return false, err
	}
	return (data == "PONG"), nil
}

// save stores the session in redis.
func (s *RediSession) save(session *Session) error {
	b, err := s.serializer.Serialize(session)
	if err != nil {
		return err
	}
	if s.maxLength != 0 && len(b) > s.maxLength {
		return errors.New("SessionStore: the value to store is too big")
	}
	conn := s.Pool.Get()
	defer conn.Close()
	if err = conn.Err(); err != nil {
		return err
	}
	_, err = conn.Do("SETEX", s.keyPrefix+session.ID, s.maxAge, b)
	return err
}

// load reads the session from redis.
// returns true if there is a sessoin data in DB
func (s *RediSession) load(session *Session) (bool, error) {
	conn := s.Pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return false, err
	}
	data, err := conn.Do("GET", s.keyPrefix+session.ID)
	if err != nil {
		return false, err
	}
	if data == nil {
		return false, nil // no data was associated with this key
	}
	b, err := redis.Bytes(data, err)
	if err != nil {
		return false, err
	}
	return true, s.serializer.Deserialize(b, session)
}

// delete keys from redis if maxAge<0
func (s *RediSession) drop(k string) error {
	conn := s.Pool.Get()
	defer conn.Close()
	if _, err := conn.Do("DEL", s.keyPrefix+k); err != nil {
		return err
	}
	return nil
}

// SessionSerializer provides an interface hook for alternative serializers
type SessionSerializer interface {
	Deserialize(d []byte, ss *Session) error
	Serialize(ss *Session) ([]byte, error)
}

// GobSerializer uses gob package to encode the session map
type GobSerializer struct{}

// Serialize using gob
func (s GobSerializer) Serialize(ss *Session) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(ss.Attributes); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Deserialize back to map[string]interface{}
func (s GobSerializer) Deserialize(d []byte, ss *Session) error {
	dec := gob.NewDecoder(bytes.NewBuffer(d))
	return dec.Decode(&ss.Attributes)
}
