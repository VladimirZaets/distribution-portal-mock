package metadata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type LocalMetadata struct {
	metadataType string
	list         MetadataList
	file         string
	filecreated  bool
}

func NewLocalMetadata() *LocalMetadata {
	return &LocalMetadata{
		metadataType: "local",
		file:         "metadata.json",
	}
}

func (lm *LocalMetadata) Set(mtd *Metadata) error {
	metadataList, err := lm.getMetadata()
	if err != nil {
		return err
	}
	//fmt.Println("Name:" + mtd.Name)
	fmt.Println(metadataList)
	metadataList[mtd.Name] = mtd
	err = lm.saveMetadataList(metadataList)
	if err != nil {
		return err
	}
	return nil
}

func (lm *LocalMetadata) Get(m *Metadata) (*Metadata, error) {
	return nil, nil
}

func (lm *LocalMetadata) GetList() (MetadataList, error) {
	return nil, nil
}

func (ln *LocalMetadata) Update(m *Metadata) error {
	return fmt.Errorf("Error")
}

func (ln *LocalMetadata) Delete(m *Metadata) error {
	return fmt.Errorf("Error")
}

func (ln *LocalMetadata) GetType() string {
	return ln.metadataType
}

func (ln *LocalMetadata) createMetadataFile() error {
	f, err := os.Create(ln.file)

	if err != nil {
		log.Fatal(err)
		return err
	}

	defer f.Close()

	_, err = f.WriteString("{}")

	if err != nil {
		log.Fatal(err)
		return err
	}

	ln.filecreated = true
	return nil
}

func (lm *LocalMetadata) getMetadata() (MetadataList, error) {
	if lm.filecreated == false {
		fmt.Println("ENTER")
		err := lm.createMetadataFile()
		if err != nil {
			return nil, err
		}
		lm.filecreated = true
	}
	file, err := ioutil.ReadFile(lm.file)
	fmt.Println(string(file))
	if err != nil {
		return nil, err
	}
	fileJson := map[string]*Metadata{}
	json.Unmarshal(file, &fileJson)
	fmt.Println("fileJson", fileJson)
	return fileJson, nil
}

func (lm *LocalMetadata) saveMetadataList(metadataList MetadataList) error {
	fileBuf, err := json.Marshal(metadataList)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(lm.file, fileBuf, 0644)
	if err != nil {
		return err
	}
	return nil
}
