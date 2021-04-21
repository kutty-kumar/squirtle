package main

import (
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/kutty-kumar/ho_oh/pikachu_v1"
	"log"
)


func main(){
	user := pikachu_v1.UserDto{
		FirstName:            "Kumar",
		LastName:             "",
		Age:                  27,
		Height:               165,
		Weight:               180,
		Gender:               2,
		Status:               1,
		ExternalId:           "kumard1",
	}

	uBytes, err := proto.Marshal(&user)
	if err != nil {
		log.Fatalf("An error %v occurred while marshalling proto to bytes", err)
	}

	client := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379",
			DB:  0,
		})
	status := client.Set("1", uBytes, 0)

	if status.Err() != nil {
		log.Fatalf("An error %v+ occurred while setting data", status.Err())
	}

	getStatus := client.Get("1")
	if getStatus.Err() != nil {
		log.Fatalf("An error %v occurred while getting data", getStatus.Err())
	}
	var kBytes []byte
	var kUserDto pikachu_v1.UserDto

	getStatus.Scan(&kBytes)
	err = proto.UnmarshalMerge(kBytes, &kUserDto)
	if err != nil {
		log.Fatalf("An error %v occurred while unmarshalling data", err)
	}
	log.Printf("%v", kUserDto)
}