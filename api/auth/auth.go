package auth

import (
    "encoding/json"
    "net/http"
    "errors"
    "time"
    "fmt"

    "gorm.io/gorm"
    "github.com/plankiton/ServMan/api/util"
    "github.com/plankiton/ServMan/api/user"
)


type Token struct {
    ID          string    `json:"token,omitempty" gorm:"primaryKey,unique"`
    PersonId    string    `json:"person_id,omitempty"`

    CreateTime  time.Time `json:"created_at,omitempty"`
    LastLogTime time.Time `json:"last_login,omitempty"`
}

func (self *Token) Populate(p user.Person) {
    self.PersonId = p.ID

    self.CreateTime = time.Now()
    self.LastLogTime = time.Now()

    str := fmt.Sprintf("%s%s%s",
        p.ID,
        p.PassHash,
        self.CreateTime.String(),
    )
    self.ID = util.ToHash(str)
}

var database *gorm.DB
func PopulateDB(db *gorm.DB) {
    database = db
    database.AutoMigrate(&Token{})
}

func CreateToken(w http.ResponseWriter, r *http.Request) {
    var body util.Request
    json.NewDecoder(r.Body).Decode(&body)
    if CheckToken(body.Token, &user.Person{}) == nil {
        w.WriteHeader(401)
        json.NewEncoder(w).Encode(util.Response{
            Message: "You already logged!",
            Code: "PermissionDenied",
            Type: "error",
            Data: nil,
        })
        return
    }

    person := user.Person{}
    res := database.Where("doc_value = ?", body.Data["document"]).First(&person)
    if errors.Is(res.Error, gorm.ErrRecordNotFound) {
        w.WriteHeader(404)

        json.NewEncoder(w).Encode(util.Response{
            Message: "The person not found!",
            Code: "NotFound",
            Type: "error",
            Data: nil,
        })

        return
    }

    if person.CheckPass(body.Data["password"]) {

        token := Token{}
        token.Populate(person)

        database.Create(token)
        database.Commit()

        json.NewEncoder(w).Encode(util.Response {
            Code: "Login",
            Type: "sucess",
            Data: token,
        })
        return
    }

    w.WriteHeader(401)
    json.NewEncoder(w).Encode(util.Response{
        Message: "You can't get the token!",
        Code: "PermissionDenied",
        Type: "error",
        Data: nil,
    })
}

func TestToken(w http.ResponseWriter, r *http.Request) {
    var body util.Request
    json.NewDecoder(r.Body).Decode(&body)

    person := user.Person{}
    if CheckToken(body.Token, &person) == nil {
        json.NewEncoder(w).Encode(util.Response {
            Code: "Login",
            Type: "sucess",
            Data: person,
        })
        return
    }

    w.WriteHeader(401)
    json.NewEncoder(w).Encode(util.Response{
        Message: "Invalid token!",
        Code: "PermissionDenied",
        Type: "error",
        Data: nil,
    })
}

func DeleteToken(w http.ResponseWriter, r *http.Request) {
    var body util.Request
    json.NewDecoder(r.Body).Decode(&body)

    person := user.Person{}
    if CheckToken(body.Token, &person) == nil {
        token := Token{}
        database.Where("id = ?", body.Token).First(&token)
        print(token.ID, " - ", person.Name, "\n")
        database.Delete(token)
        database.Commit()
        json.NewEncoder(w).Encode(util.Response {
            Code: "Logout",
            Type: "sucess",
        })
        return
    }

    w.WriteHeader(401)
    json.NewEncoder(w).Encode(util.Response{
        Message: "Invalid token!",
        Code: "PermissionDenied",
        Type: "error",
        Data: nil,
    })
}

func CheckToken(token_id string, p *user.Person) error {
    token := Token{}
    res := database.Where("id = ?", token_id).First(&token)
    if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
        res = database.Where("id = ?", token.PersonId).First(&p)
        if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
            token.LastLogTime = time.Now()
            database.Save(&token)
            database.Commit()
            return nil
        }
    }

    return errors.New("Permission denied")
}

func RefuseTokens() {
    for {
        time.Sleep(time.Hour)
        tokens := []Token{}
        database.Find(&tokens)
        for _, token := range(tokens) {
            hours := time.Now().Sub(token.LastLogTime).Hours()
            five_hours := (time.Hour*24*5).Hours()

            if hours >= five_hours {
                database.Delete(token)
                database.Commit()
            }
        }
    }
}
