package fasttext

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/ekzhu/go-fasttext"
	_ "github.com/mattn/go-sqlite3"
)

// GetSimilarity returns the cosine similarity between two vector based on fasttext
func GetSimilarity(str1, str2 string) float64 {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return 0.0
	}
	db := fmt.Sprintf("%v/%v", dir, "devsup.db")
	_, err = os.Open(db)
	if err != nil {
		f, err := os.Create(db)
		if err != nil {
			return 0.0
		}
		f.Close()
	}

	ft := fasttext.NewFastText(db)
	path := filepath.Join(os.Getenv("GOPATH"), "src/github.com/mkhoi1998/devsup/model/dev_sup.vec")
	file, err := os.Open(path)
	err = ft.BuildDB(file)
	vect1, err := ft.GetEmb(str1)
	if err != nil {
		return 0.0
	}
	vect2, err := ft.GetEmb(str2)
	if err != nil {
		return 0.0
	}

	dotProduct := 0.0
	for k, v := range vect1 {
		dotProduct += float64(v) * float64(vect2[k])
	}
	sum1 := 0.0
	for _, v := range vect1 {
		sum1 += math.Pow(float64(v), 2)
	}
	sum2 := 0.0
	for _, v := range vect2 {
		sum2 += math.Pow(float64(v), 2)
	}
	magnitude := math.Sqrt(sum1) * math.Sqrt(sum2)
	if magnitude == 0 {
		return 0.0
	}
	return float64(dotProduct) / float64(magnitude)
}
