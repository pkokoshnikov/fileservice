package fs

import (
    "net/http"
	"log"
	"io"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"container/list"
	"strconv"
)

type DAO struct {
	Url string 
}

type File struct {
	Id bson.ObjectId `bson:"_id"`
	Filename string `bson:"filename"`
}

func NewDAO() *DAO {
	dao := new(DAO)
	dao.Url = os.Getenv("FSMONGO_URL")		
	
	return dao
}

func (dao *DAO) DeleteFile(filename string) {
	session := dao.getMongoSession()
	defer session.Close()
	
	count, err := session.DB("fileservice").GridFS("fs").Find(bson.M{"filename": filename}).Count();
	if err != nil {
		panic(err)
	}
	
	if(count != 0) {
		if err := session.DB("fileservice").GridFS("fs").Remove(filename); err !=nil {
			log.Print(err)
		}
	} else {
		panic("File not found")
	}		
}

func (dao *DAO) ShowAllFiles() *list.List{
	session := dao.getMongoSession()
	defer session.Close()
	
	iter := session.DB("fileservice").GridFS("fs").Find(nil).Iter()
	defer iter.Close()
	
	list := list.New()
	var result File
	for iter.Next(&result){
		list.PushBack(result)		
	}
	
	return list	
}

func (dao *DAO) UploadFile(r *http.Request) {
	session := dao.getMongoSession()
	defer session.Close()
	
	file, header, err := r.FormFile("fileToUpload")
	if err != nil {
		panic("File cannot be opend from request")		
	}	
	defer file.Close()
	
	gfsfile, err := session.DB("fileservice").GridFS("fs").Create(header.Filename)
	if err != nil {
		panic("File cannot be opend from databse")
	}
	defer gfsfile.Close()
	
	written, err := io.Copy(gfsfile, file)
	if err != nil {
		panic("File cannot be copied to databse")
	}
	log.Printf("File was copied, %s bytes", strconv.FormatInt(written, 10))
}

func (dao *DAO) getMongoSession() (*mgo.Session) {
	session, err := mgo.Dial(dao.Url)
	if err != nil {
		panic(err)
	}
	return session
}