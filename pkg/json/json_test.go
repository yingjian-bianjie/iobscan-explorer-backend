package json

import (
	"encoding/json"
	"testing"
)

func TestValidate(t *testing.T) {
	schema := `{
  "type": "object",
  "required": ["name","description","image","attributes"],
  "properties": {
    "name": {
      "type": "string"
    },
    "description": {
      "type": "string"
    },
    "external_url": {
      "type": "string"
    },
    "image": {
      "type": "string"
    },
    "image_data": {
      "type": "string"
    },
    "animation_url": {
      "type": "string"
    },
    "background_color": {
      "type": "string"
    },
    "youtube_url": {
      "type": "string"
    },
    "attributes": {
      "type": "array",
      "items": {
        "type": "object",
        "required": [],
        "properties": {
          "trait_type": {
            "type": "string"
          },
          "value": {
            "type": "number"
          }
        }
      }
    }
  }
}`
	var m interface{}
	json.Unmarshal([]byte(schema), &m)
	bytes, _ := json.Marshal(m)
	t.Log(string(bytes))
	jsonstr := `{
    "name": "Dave Starbelly",
    "description": "Friendly OpenSea Creature that enjoys long swims in the ocean.",
    "external_url": "https://openseacreatures.io/3",
    "image": "https://storage.googleapis.com/opensea-prod.appspot.com/puffs/3.png",
    "image_data": "https://storage.googleapis.com/opensea-prod.appspot.com/puffs/3.png",
    "animation_url": "",
    "background_color": "",
    "youtube_url": "",
    "attributes": [
        {
            "trait_type": "Base",
            "value": "Starfish"
        },
        {
            "trait_type": "Eyes",
            "value": "Big"
        },
        {
            "trait_type": "Mouth",
            "value": "Surprised"
        },
        {
            "trait_type": "Level",
            "value": 5
        },
        {
            "trait_type": "Stamina",
            "value": 1.4
        },
        {
            "trait_type": "Personality",
            "value": "Sad"
        },
        {
            "display_type": "boost_number",
            "trait_type": "Aqua Power",
            "value": 40
        },
        {
            "display_type": "boost_percentage",
            "trait_type": "Stamina Increase",
            "value": 10
        },
        {
            "display_type": "number",
            "trait_type": "Generation",
            "value": 2
        }
    ]
}`
	//jsonstr = `{"name":"","description":"","image":"","attributes":[]}`
	result, err := Validate(schema, jsonstr)
	if err != nil {
		t.Fatal(err)
	}
	if result.Valid() {
		t.Log("The document is valid")
	} else {
		t.Log("The document is invalid. see errors :")
		for _, desc := range result.Errors() {
			t.Log(desc)
		}
	}
}
