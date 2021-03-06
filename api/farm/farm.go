package farm

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "net/http"
    "time"
    "fmt"
    "errors"

    "gorm.io/gorm"
    "github.com/plankiton/ServMan/api/util"
    "github.com/plankiton/ServMan/api/user"
)

type Addr struct {
    ID           string `json:"id,omitempty" gorm:"index"`
    Street       string `json:"street,omitempty"`
    State        string `json:"state,omitempty"`
    Number       string `json:"number,omitempty"`
    Code         string `json:"cep,omitempty"`
    City         string `json:"city,omitempty"`
    Neightbourn  string `json:"neighborhood,omitempty"`
}

type Farm struct {
    ID        string `json:"id,omitempty" gorm:"primaryKey"`
    PersonId  string `json:"person,omitempty"`
    AddressId string `json:"address,omitempty" gorm:"uniqueIndex"`
    Name      string `json:"name,omitempty" gorm:"index"`

    CreateTime time.Time `json:"created_at,omitempty" gorm:"index"`
    UpdateTime time.Time `json:"updated_at,omitempty" gorm:"index"`
}


var database *gorm.DB
func PopulateDB(db *gorm.DB) {
    database = db
    database.AutoMigrate(&Farm{})
    database.AutoMigrate(&Addr{})
}

func GetAddr(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    addr := Addr{}

    farm := Farm{}
    res := database.Where("id = ?", params["id"]).First(&farm)
    if errors.Is(res.Error, gorm.ErrRecordNotFound) {
        w.WriteHeader(404)

        json.NewEncoder(w).Encode(util.Response{
            Message: "The farm not found!",
            Code: "NotFound",
            Type: "error",
            Data: nil,
        })

        return
    }

    res = database.Where("id = ?", farm.AddressId).First(&addr)
    if errors.Is(res.Error, gorm.ErrRecordNotFound) {
        w.WriteHeader(404)

        json.NewEncoder(w).Encode(util.Response{
            Message: "The farm dont have address!",
            Code: "NotFound",
            Type: "error",
            Data: nil,
        })

        return
    }

    // TODO: sentence for validate logged user
    json.NewEncoder(w).Encode(util.Response{
            Code: "GetAddr",
            Type: "sucess",
            Data: addr,
        })
}

// GetPeople mostra todos os contatos da variável farm
func GetAllFarms(w http.ResponseWriter, r *http.Request) {
    farm := []Farm{}
    database.Find(&farm)

    // TODO: sentence for validate logged user

    json.NewEncoder(w).Encode(util.Response{
            Code: "GetFarms",
            Type: "sucess",
            Data: farm,
        })
}
func GetFarms(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    farm := []Farm{}
    res := database.Where("person_id = ?", params["id"]).Find(&farm)
    if errors.Is(res.Error, gorm.ErrRecordNotFound) {
        w.WriteHeader(404)

        json.NewEncoder(w).Encode(util.Response{
            Message: "The farm not found!",
            Code: "NotFound",
            Type: "error",
            Data: nil,
        })

        return
    }

    // TODO: sentence for validate logged user
    json.NewEncoder(w).Encode(util.Response{
            Code: "GetFarms",
            Type: "sucess",
            Data: farm,
        })
}

// GetFarm mostra apenas um contato
func GetFarm(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    farm := Farm{}

    res := database.Where("id = ?", params["id"]).First(&farm)
    if errors.Is(res.Error, gorm.ErrRecordNotFound) {
        w.WriteHeader(404)

        json.NewEncoder(w).Encode(util.Response{
            Message: "The farm not found!",
            Code: "NotFound",
            Type: "error",
            Data: nil,
        })

        return
    }

    // TODO: sentence for validate logged user
    json.NewEncoder(w).Encode(util.Response{
            Code: "GetFarm",
            Type: "sucess",
            Data: farm,
        })
}

// CreateFarm cria um novo contato
func GetAddrFromCep(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    address := Addr{}
    r_addr, err := http.Get(fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", params["cep"]))

    if err != nil || r_addr.StatusCode != 200 {
        w.WriteHeader(404)

        json.NewEncoder(w).Encode(util.Response{
            Message: "The address not found!",
            Code: "NotFound",
            Type: "error",
        })

        return
    }

    json.NewDecoder(r_addr.Body).Decode(&address)
    json.NewEncoder(w).Encode(util.Response {
        Type:    "sucess",
        Code:    "GetAddressFromCep",
        Data:    address,
    })
}

func CreateFarm(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    var body util.Request
    json.NewDecoder(r.Body).Decode(&body)

    if len(body.Data) == 0 {
        w.WriteHeader(400)

        json.NewEncoder(w).Encode(util.Response{
            Message: "The data sent is invalid!"+
                     `(must be {"data": "..."})`,
            Code: "BadRequest",
            Type: "error",
            Data: nil,
        })

        return
    }

    person := user.Person{}
    res := database.Where("doc_value = ? OR id = ?",
        params["id"], params["id"]).First(&person)
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


    farm := Farm {
        PersonId: person.ID,
    }

    address := Addr{}
    r_addr, err := http.Get(fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", body.Data["cep"]))
    if err == nil {
        json.NewDecoder(r_addr.Body).Decode(&address)
    }

    for i, prop := range(body.Data) {
        switch i {
        case "name":
            farm.Name = prop
        case "cep":
            if address.Code == "" {
                address.Code = prop
            }
        case "street":
            if address.Street == "" {
                address.Street = prop
            }
        case "number":
            if address.Number == "" {
                address.Number = prop
            }
        case "neighborhood":
            if address.Neightbourn == "" {
                address.Neightbourn = prop
            }
        case "state":
            if address.State == "" {
                address.State = prop
            }
        case "city":
            if address.City == "" {
                address.City = prop
            }
        }
    }

    farm.CreateTime = time.Now()
    farm.UpdateTime = time.Now()

    address.ID = util.ToHash(address.Code+address.Street+address.Number)
    {
        addr := Addr{}
        res = database.Where("id = ?", address.ID).First(&addr)
        if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
            address = addr
        } else {
            database.Create(&address)
            database.Commit()
        }
    }
    farm.AddressId = address.ID
    farm.ID = util.ToHash(farm.Name+
        farm.PersonId+
        farm.AddressId+
        farm.CreateTime.Format("%Y%m%d%H%M%S"))

    res = database.Where("address_id = ?", farm.AddressId).First(&farm)
    if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
        w.WriteHeader(403)

        json.NewEncoder(w).Encode(util.Response{
            Message: "The farm already exists!",
            Code: "AlreadyExists",
            Type: "error",
            Data: nil,
        })

        return
    }

    // TODO: sentence for validate logged user
    // Sending all to database
    database.Create(&farm)
    database.Commit()
    json.NewEncoder(w).Encode(util.Response {
        Type:    "sucess",
        Code:    "CreatedFarm",
        Data:    farm,
    })
}

func UpdateFarm(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    var body util.Request
    json.NewDecoder(r.Body).Decode(&body)

    farm, address := Farm{}, Addr{}

    res := database.Where("id = ?", params["id"]).First(&farm)
    if errors.Is(res.Error, gorm.ErrRecordNotFound) {
        w.WriteHeader(404)

        json.NewEncoder(w).Encode(util.Response{
            Message: "The farm not found!",
            Code: "NotFound",
            Type: "error",
        })

        return
    }
    res = database.Where("id = ?", farm.AddressId).First(&address)

    for i, prop := range(body.Data) {
        switch i {
        case "name":
            farm.Name = prop
        case "cep":
            address.Code = prop

            r_addr, err := http.Get(fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", body.Data["cep"]))
            if err == nil {
                json.NewDecoder(r_addr.Body).Decode(&address)
            }
        case "street":
            address.Street = prop
        case "number":
            address.Number = prop
        case "neighbornhood":
            address.Neightbourn = prop
        case "state":
            address.State = prop
        case "city":
            address.City = prop
        }
    }

    {
        address.ID = util.ToHash(
            address.Code+
            address.State+
            address.City+
            address.Street+
            address.Neightbourn+
            address.Number)

        fmt.Printf("%v\n", address);
        database.Delete(&address)
        database.Create(&address)
        database.Commit()
    }

    farm.UpdateTime = time.Now()

    farm.AddressId = address.ID

    // TODO: sentence for validate logged user

    database.Save(&farm)
    database.Commit()
    json.NewEncoder(w).Encode(util.Response{
        Message: fmt.Sprintf("Farm %s did updated!", farm.Name),
        Code: "UpdatedFarm",
        Type: "sucess",
        Data: farm,
    })


}

// DeleteFarm deleta um contato
func DeleteFarm(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    farm := Farm{}
    res := database.Where("id = ?", params["id"]).First(&farm)
    if errors.Is(res.Error, gorm.ErrRecordNotFound) {
        w.WriteHeader(404)

        json.NewEncoder(w).Encode(util.Response{
            Message: "The farm not found!",
            Code: "NotFound",
            Type: "error",
            Data: nil,
        })

        return
    }

    // TODO: sentence for validate logged user

    database.Delete(&farm)
    database.Commit()
    json.NewEncoder(w).Encode(util.Response{
        Type: "sucess",
        Code: "DeleteFarm",
        Data: farm,
    })
}
