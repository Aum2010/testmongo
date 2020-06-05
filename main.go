package main

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"

	"time"

	
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Map struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Lat float32 `json:"lat,omitempty" bson:"lat,omitempty"`
	Lng float32 `json:"lng,omitempty" bson:"lng,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	Alias []string `json:"alias,omitempty" bson:"alias,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}

var client *mongo.Client

func CreateMap(w http.ResponseWriter ,r *http.Request) {
	w.Header().Add("Content-Type","application/json")
	var maps Map
	maps.CreatedAt = time.Now()
	maps.UpdatedAt = time.Now()
	json.NewDecoder(r.Body).Decode(&maps)
	collection := client.Database("Employee").Collection("employees")
	ctx , _ := context.WithTimeout(context.Background(),10*time.Second)
	result , _ := collection.InsertOne(ctx , maps)
	json.NewEncoder(w).Encode(result)
}

func GetMap(w http.ResponseWriter,r *http.Request) {
	w.Header().Add("Content-Type","application/json")
	var maps []Map
	collection := client.Database("Employee").Collection("employees")
	ctx , _ := context.WithTimeout(context.Background(),10*time.Second)
	cursor , err := collection.Find(ctx,bson.M{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"msg":"`+err.Error()+`"}`))
			return
		}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var _map Map
		cursor.Decode(&_map)
		maps = append(maps , _map)
	}
	err = cursor.Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"msg":"`+err.Error()+`"}`))
			return
	}	
	json.NewEncoder(w).Encode(maps)
}

func GetMapById(w http.ResponseWriter ,r *http.Request){
	w.Header().Add("Content-Type","application/json")
	params := mux.Vars(r)
	id , _ := primitive.ObjectIDFromHex(params["id"])
	var _map Map
	collection := client.Database("Employee").Collection("employees")
	ctx , _ := context.WithTimeout(context.Background(),10*time.Second)
	err := collection.FindOne(ctx , Map{ID:id}).Decode(&_map)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"msg":"`+err.Error()+`"}`))
		return
	}
	json.NewEncoder(w).Encode(_map)

}

func UpdateMap(w http.ResponseWriter ,r *http.Request){
	w.Header().Add("Content-Type","application/json")
	params := mux.Vars(r)
	id , _ := primitive.ObjectIDFromHex(params["id"])
	// var _map Map
	var maps Map
	collection := client.Database("Employee").Collection("employees")
	ctx , _ := context.WithTimeout(context.Background(),10*time.Second)
	// collection.FindOne(ctx , Map{ID:id}).Decode(&_map)
	json.NewDecoder(r.Body).Decode(&maps)
	maps.UpdatedAt = time.Now()
	
	update := bson.D{{"$set", maps }}
	//fmt.Println(maps)
	res , err := collection.UpdateOne(ctx , Map{ID:id} , update)
	// res ,err := collection.UpdateOne(ctx , Map{ID:id,UpdatedAt: time.Now()} , _map)
	// fmt.Println(id , _map)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"msg":"`+err.Error()+`"}`))
		return
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", res.MatchedCount, res.ModifiedCount)
	
	
}

func DeleteMap(w http.ResponseWriter ,r *http.Request){
	w.Header().Add("Content-Type","application/json")
	params := mux.Vars(r)
	id , _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database("Employee").Collection("employees")
	ctx , _ := context.WithTimeout(context.Background(),10*time.Second)
	res , err := collection.DeleteOne(ctx , Map{ID:id})
		if (err != nil ){
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"msg":"`+err.Error()+`"}`))
			return
		}
	fmt.Printf("Found multiple documents (array of pointers): %+v\n", res)
}

func main() {
	ctx , _ := context.WithTimeout(context.Background(),10*time.Second)
	client , _ = mongo.Connect(ctx,options.Client().ApplyURI("mongodb://localhost:27017"))
		
	router := mux.NewRouter()
	router.HandleFunc("/createmap" , CreateMap).Methods("POST")
	router.HandleFunc("/getmap" , GetMap).Methods("GET")
	router.HandleFunc("/getmap/{id}",GetMapById).Methods("GET")
	router.HandleFunc("/updatemap/{id}",UpdateMap).Methods("PUT")
	router.HandleFunc("/deletemap/{id}",DeleteMap).Methods("DELETE")
	http.ListenAndServe(":8080" , router)
}