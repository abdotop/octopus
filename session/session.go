package session

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"

	octopus "github.com/abdotop/Octopus"

	// octopus "octopus/Ctx"
	"time"

	"github.com/gofrs/uuid"
)

var Notif sync.Map

type Config struct {
	CookieName string
	Value      string

	Path       string    // optional
	Domain     string    // optional
	Expires    time.Time // optional
	RawExpires string    // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
	Raw      string
	Unparsed []string // Max-Age attribute present and given in seconds
}

type starter struct {
	session *session
	Ctx     *octopus.Ctx
}

type session struct {
	Config      *Config
	database    *sql.DB
	data        *sync.Map
	mu          sync.Mutex
	SessionName string
}

type storage struct {
	cookie *http.Cookie
	id     uuid.UUID
}

func New(c *Config) *session {
	if c == nil {
		c = new(Config)
	}
	if c.CookieName == "" {
		c.CookieName = "mycookie"
	}
	if c.MaxAge == 0 {
		c.MaxAge = int(31536000)
	}
	if c.Expires.IsZero() {
		c.Expires = time.Now().Add(time.Second * time.Duration(c.MaxAge))
	}
	var sS http.SameSite
	if c.SameSite == sS {
		c.SameSite = http.SameSiteNoneMode
	}
	s := &session{Config: c, SessionName: c.CookieName, database: nil, data: &sync.Map{}}
	return s
}

func (s *session) tmp() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					s.mu.Lock()
					if s.database != nil {
						_, err := s.database.Exec(fmt.Sprintf(`DELETE FROM %s WHERE datetime(expiration_date) <= datetime('now')`, s.SessionName))
						if err != nil {
							fmt.Println(err)
							s.mu.Unlock()
							return
						}
					}
					s.data.Range(func(key, value any) bool {
						val, ok := value.(map[string]interface{})
						userKey := val["key"].(uuid.UUID)
						Storage := val["cookie"].(*storage)
						if ok {
							expirationDate := Storage.cookie.Expires
							// Vérifiez si la session a expiré
							if time.Now().After(expirationDate) {
								s.data.Delete(key)
								Notif.Store(userKey, false)
								return true
							}
							//  else {
							// 	Notif.Store(userKey, true)
							// 	return true
							// }
						}
						return true
					})
					s.mu.Unlock()

				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()

		// Attendez que l'application se termine.
		time.Sleep(30 * time.Minute)
		close(quit)
	}()
}

func (s *session) UseDB(db *sql.DB) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.database = db
	_, err := s.database.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (id UUID PRIMARY KEY,user_id UUID NOT NULL,expiration_date DATETIME NOT NULL);`, s.SessionName))
	if err != nil {
		log.Fatal(err)
	}
}

func (s *session) Start(c *octopus.Ctx) *starter {
	s.tmp()
	return &starter{session: s, Ctx: c}
}

func (s *starter) Set(value uuid.UUID) error {
	s.session.mu.Lock()
	defer s.session.mu.Unlock()

	session := s.session
	c := session.Config
	db := session.database
	tmpdata := session.data
	// get new id for the session
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	if db != nil {
		var user uuid.UUID
		err = db.QueryRow(fmt.Sprintf("SELECT id FROM %s WHERE user_id = $1", session.SessionName), value).Scan(&user)
		if err != sql.ErrNoRows {

			// // Préparez une instruction SQL pour supprimer la session
			stmt, err := db.Prepare(fmt.Sprintf("DELETE FROM %s WHERE id=$1", session.SessionName))
			if err != nil {
				return err
			}

			// // Exécutez l'instruction SQL avec l'UUID de la session
			_, err = stmt.Exec(user)
			if err != nil {
				return err
			}
			tmpdata.Delete(user.String())
		}
		// fmt.Println(user)
		stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (id, user_id, expiration_date) VALUES($1, $2, $3)", session.SessionName))
		if err != nil {
			return err
		}
		// Exécutez l'instruction SQL avec le nouvel UUID et les autres valeurs
		_, err = stmt.Exec(id, value, time.Now().Add(time.Second*time.Duration(c.MaxAge)))
		if err != nil {
			return err
		}

	}
	// Stockez la session dans un cookie
	cookie := &http.Cookie{
		Name:     c.CookieName,
		Value:    id.String(),
		Secure:   c.Secure,
		Expires:  time.Now().Add(time.Second * time.Duration(c.MaxAge)),
		MaxAge:   c.MaxAge,
		Path:     c.Path,
		Domain:   c.Domain,
		HttpOnly: c.HttpOnly,
		SameSite: c.SameSite,
	}

	http.SetCookie(s.Ctx.Response, cookie)
	tmpdata.Store(id.String(), map[string]interface{}{
		"key":    value,
		"cookie": &storage{cookie: cookie, id: value},
	})
	Notif.Store(value, true)
	return nil
}

func (s *starter) Get() (uuid.UUID, error) {
	s.session.mu.Lock()
	defer s.session.mu.Unlock()

	session := s.session
	c := session.Config
	db := session.database
	tmpdata := session.data

	// Récupérez le cookie
	cookie, err := s.Ctx.Request.Cookie(c.CookieName)
	if err != nil {
		return uuid.Nil, fmt.Errorf("erreur lors de la récupération du cookie : %v", err)
	}
	// Récupérez l'ID de session à partir du cookie
	sessionID := cookie.Value
	value, ok := tmpdata.Load(sessionID)
	if ok {
		val, ok := value.(map[string]interface{})
		userKey := val["key"].(uuid.UUID)
		Storage := val["cookie"].(*storage)
		if ok {
			expirationDate := Storage.cookie.Expires
			if time.Now().After(expirationDate) {
				tmpdata.Delete(sessionID)
				Notif.Store(userKey, false)
				return uuid.Nil, fmt.Errorf("la session a expiré")
			} else {
				Notif.Store(userKey, true)
				return Storage.id, nil
			}
		}
	}
	if db != nil {
		// Récupérez la session de la base de données
		var userID uuid.UUID
		var expirationDate time.Time
		err = db.QueryRow(fmt.Sprintf("SELECT user_id, expiration_date FROM %s WHERE id = $1", session.SessionName), sessionID).Scan(&userID, &expirationDate)
		if err != nil {
			return uuid.Nil, fmt.Errorf("erreur lors de la récupération de la session : %v", err)
		}

		// Vérifiez si la session a expiré
		if time.Now().After(expirationDate) {
			Notif.Store(userID, false)
			return uuid.Nil, fmt.Errorf("la session a expiré")
		}
		Notif.Store(userID, true)
		return userID, nil
	}
	// Retournez l'ID de l'utilisateur et nil
	return uuid.Nil, fmt.Errorf("erreur lors de la récupération du cookie")
}

func (s *starter) Valid() bool {
	s.session.mu.Lock()
	defer s.session.mu.Unlock()

	session := s.session
	c := session.Config
	db := session.database
	tmpdata := session.data
	cookie, err := s.Ctx.Request.Cookie(c.CookieName)
	if err != nil {
		// Le cookie n'existe pas
		return false
	}
	value, ok := tmpdata.Load(cookie.Value)
	if ok {
		val, ok := value.(map[string]interface{})
		userKey := val["key"].(uuid.UUID)
		Storage := val["cookie"].(*storage)
		if ok {
			expirationDate := Storage.cookie.Expires
			if time.Now().After(expirationDate) {
				tmpdata.Delete(cookie.Value)
				Notif.Store(userKey, false)
				return false
			} else {
				Notif.Store(userKey, true)
				return true
			}
		}
	}
	if db != nil {
		// Récupérez la session de la base de données
		var expirationDate time.Time
		var userID uuid.UUID
		err = db.QueryRow(fmt.Sprintf("SELECT user_id, expiration_date FROM %s WHERE id = ?", session.SessionName), cookie.Value).Scan(&userID, &expirationDate)
		if err != nil {
			// La session n'existe pas dans la base de données
			return false
		}

		// Vérifiez si la session a expiré
		if time.Now().After(expirationDate) {
			// La session a expiré
			Notif.Store(userID, false)
			return false
		}

		// Le cookie existe et la session est valide
		Notif.Store(userID, true)
		return true
	}
	return false
}

func (s *starter) Delete() error {
	s.session.mu.Lock()
	defer s.session.mu.Unlock()

	session := s.session
	c := session.Config
	db := session.database
	tmpdata := session.data
	Ctx := s.Ctx
	cookie, err := Ctx.Request.Cookie(c.CookieName)
	if err != nil {
		// Le cookie n'existe pas
		return err
	}
	if db != nil {
		// Préparez une instruction SQL pour supprimer la session
		stmt, err := db.Prepare(fmt.Sprintf("DELETE FROM %s WHERE id=$1", session.SessionName))
		if err != nil {
			return err
		}

		// Exécutez l'instruction SQL avec l'UUID de la session
		_, err = stmt.Exec(cookie.Value)
		if err != nil {
			return err
		}
	}

	value, ok := tmpdata.LoadAndDelete(cookie.Value)
	if ok {
		val, ok := value.(map[string]interface{})
		if ok {
			userKey := val["key"].(uuid.UUID)
			Notif.Store(userKey, false)
		}
	}
	// Supprimez le cookie de la session
	http.SetCookie(Ctx.Response, &http.Cookie{
		Name:    c.CookieName,
		Value:   "",
		Secure:  c.Secure,
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	})

	return nil
}
