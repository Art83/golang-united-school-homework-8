package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func AddItem(item, filename string) error {
	var newItem User
	err := json.Unmarshal([]byte(item), &newItem)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	inFileJson, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	var inFile []User
	if len(inFileJson) != 0 {
		err := json.Unmarshal([]byte(inFileJson), &inFile)
		if err != nil {
			return err
		}
		for _, val := range inFile {
			if val.Id == newItem.Id {
				return fmt.Errorf("such ID already exists")
			}
		}

	}
	inFile = append(inFile, newItem)

	dataIn, err := json.Marshal(inFile)

	if err != nil {
		return err
	}
	err = os.WriteFile(filename, []byte(dataIn), 0666)
	if err != nil {
		return err
	}

	return nil
}

func FindById(id, filename string, writer io.Writer) error {
	fmt.Println(id)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	inFileJson, err := io.ReadAll(file)
	var inFile []User
	var given_Id User
	var found bool
	if len(inFileJson) != 0 {
		err := json.Unmarshal([]byte(string(inFileJson)), &inFile)
		if err != nil {
			return err
		}
		for _, val := range inFile {
			fmt.Println(val.Id)
			if val.Id == id {
				given_Id = val
				found = true
				break
			}
		}
	}
	if found == false {
		_, err := fmt.Fprintf(writer, "")
		if err != nil {
			return err
		}
	}

	ItemToReturn, err := json.Marshal(given_Id)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(writer, string(ItemToReturn))
	if err != nil {
		return err
	}
	return nil

}

func ListItems(filename string) (string, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()
	inFileJson, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	dataFromFile := string(inFileJson)
	return dataFromFile, nil

}

func RemoveItem(item_num, filename string, writer io.Writer) error {
	var found bool
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	inFileJson, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	var inFile []User
	if len(inFileJson) != 0 {
		err := json.Unmarshal([]byte(inFileJson), &inFile)
		if err != nil {
			return err
		}
		for i, val := range inFile {
			if val.Id == item_num {
				inFile = append(inFile[:i], inFile[i+1:]...)
				found = true
				break
			}
		}
		if found == false {
			_, err := fmt.Fprintf(writer, "Item with id %s doesn't exist", item_num)
			if err != nil {
				return err
			}
		}

	}

	dataIn, err := json.Marshal(inFile)

	if err != nil {
		return err
	}
	err = os.WriteFile(filename, []byte(dataIn), 0666)
	if err != nil {
		return err
	}

	return nil
}

func Perform(args Arguments, writer io.Writer) error {
	if args["operation"] == "" {
		return fmt.Errorf("there's no operation flag")
	}
	if args["fileName"] == "" {
		return fmt.Errorf("there's no file name")
	}
	if args["operation"] == "add" {
		err := AddItem(args["item"], args["fileName"])
		if err != nil {
			return err
		}
	} else if args["operation"] == "remove" {
		err := RemoveItem(args["item"], args["fileName"], writer)
		if err != nil {
			return err
		}
	} else if args["operation"] == "list" {
		data, err := ListItems(args["fileName"])
		if err != nil {
			return err
		}
		fmt.Println(data)
	} else if args["operation"] == "findById" {
		err := FindById(args["item"], args["fileName"], writer)
		if err != nil {
			return err
		}
	}
	return nil

}

func parseArgs() Arguments {
	arg := make(Arguments)
	id := flag.String("id", "", "id of participant")
	op := flag.String("operation", "", "type of operation")
	fn := flag.String("fileName", "", "title of filename")
	item := flag.String("item", "", "item for storing")
	flag.Parse()
	arg["id"] = *id
	arg["operation"] = *op
	arg["item"] = *item
	arg["fileName"] = *fn
	return arg
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
