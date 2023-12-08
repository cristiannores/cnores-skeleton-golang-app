package service

import (
	"encoding/json"
	"cnores-skeleton-golang-app/app/infrastructure/mongo_client/models"
)

func GetServiceFromJSON() models.ServiceModel {
	var user models.ServiceModel
	_ = json.Unmarshal([]byte(GetServiceBsonJSON()), &user)
	return user
}

func GetServiceBsonJSON() string {
	return `{
  "_id": {
    "$oid": "654c78fd42e70dd56f0521de"
  },
  "name": "Terapia Ocupacional",
  "description": "Terapia para pacientes con TEA",
  "image": "http://example.com/image.jpg",
  "isActive": true,
  "createdBy": {
    "$oid": "654b9b7df05303afe250debd"
  },
  "updatedBy": {
    "$oid": "654b9b805e0559cde1157f85"
  },
  "createdAt": {
    "$date": {
      "$numberLong": "1672531200000"
    }
  },
  "updatedAt": {
    "$date": {
      "$numberLong": "1672617600000"
    }
  }
}`
}
