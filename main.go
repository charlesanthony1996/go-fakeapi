package main

import (
	"encoding/json"
	"go-fakeapi/types"
	"log"
	"net/http"
	"sync"
	"time"
)

// import (
// 	"/.typesUserDto.go"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"sync"
// 	"time"
// )

// define dto here
type Role struct {
	Permissions []string `json:"permissions"`
}

type Group struct {
	Users []string `json:"users"`
}

type AclEntry struct {
	Action string   `json:"action"`
	Src    []string `json:"src"`
	Dst    []string `json:"dst"`
}

type AclData struct {
	Roles  map[string]Role  `json:"roles"`
	Groups map[string]Group `json:"groups"`
	Acl    []AclEntry       `json:"acl"`
}

type AccessControlListDto struct {
	ID        string  `json:"id"`
	CreatedAt string  `json:"createdAt"`
	CreatedBy string  `json:"createdBy"`
	Data      AclData `json:"data"`
}

// in memory storage for acl data
var aclData AccessControlListDto
var aclMutex sync.Mutex

// initialize some test data
func init() {
	aclData = AccessControlListDto{
		ID:        "1",
		CreatedAt: time.Now().Format(time.RFC3339),
		CreatedBy: "system",
		Data: AclData{
			Roles: map[string]Role{
				"admin": {Permissions: []string{"read", "write", "delete"}},
				"user":  {Permissions: []string{"read"}},
			},
			Groups: map[string]Group{
				"developers": {Users: []string{"john", "jane"}},
			},
			Acl: []AclEntry{
				{Action: "allow", Src: []string{"admin"}, Dst: []string{"/admin"}},
				{Action: "deny", Src: []string{"user"}, Dst: []string{"/admin"}},
			},
		},
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// fetch the latest definitions
func getLatestAclDefinition(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	aclMutex.Lock()
	defer aclMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(aclData)
}

// save a new acl definition
func saveAclDefinition(w http.ResponseWriter, r *http.Request) {
	var newAcl AccessControlListDto
	if err := json.NewDecoder(r.Body).Decode(&newAcl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// enableCors(&w)
	aclMutex.Lock()
	defer aclMutex.Unlock()

	// simulate updating the acl with the new data
	newAcl.ID = aclData.ID
	newAcl.CreatedAt = time.Now().Format(time.RFC3339)
	newAcl.CreatedBy = "user"
	aclData = newAcl

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newAcl)
}

// revert to the people acl definition (simulated by just sending the same one)
func revertAclDefinition(w http.ResponseWriter, r *http.Request) {
	aclMutex.Lock()
	defer aclMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(aclData)
}

// user store to handle viewing permissions
func getUser(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	user := types.UserDto{
		Email:       "emberger@nu-it.at",
		Name:        "Michael Emberger",
		Pic:         "https://lh3.googleusercontent.com/a/ACg8ocKVc9-IpolPYVPs676OyRHg9f9QC4APpV0bCJKMg_4oKuNnYwA=s96-c",
		Permissions: []string{"dashboard.write", "profile.edit"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func main() {
	http.HandleFunc("/nunode/auth/acl/definitions/latest", getLatestAclDefinition)
	http.HandleFunc("/nunode/auth/acl/definitions", saveAclDefinition)
	http.HandleFunc("/nunode/auth/acl/definitions/revert", revertAclDefinition)

	// user permission
	http.HandleFunc("/user", getUser)

	log.Println("Starting on server :5693...")
	log.Fatal(http.ListenAndServe(":5693", nil))
}
