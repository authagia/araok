package main

import (
	"araok/dam"
	"araok/joysound"
	"fmt"
)

func main() {
	keyword := "NEVER"

	listInfo, err := joysound.Search(keyword)
	if err != nil {
		panic(err)
	}
	for i, song := range listInfo.Items {
		fmt.Printf(
			"%s - %s [%s]\n",
			song.Title,
			song.Date,
			song.Href,
		)

		if i > 50 {
			fmt.Printf("---skip----(left %d items)\n", len(listInfo.Items)-i)
			break
		}
	}

	fmt.Println("----------")

	result, err := dam.Search(keyword)
	if err != nil {
		panic(err)
	}

	// fmt.Println("API Status:", result.Result.StatusCode)
	// fmt.Println("Message:", result.Result.Message)
	// fmt.Println("Keyword:", result.Data.Keyword)
	// fmt.Println("Total:", result.Data.TotalCount)

	for i, song := range result {
		fmt.Printf(
			"%s - %s [%s]\n",
			song.Title,
			song.Artist,
			song.RequestNo,
		)
		if i > 50 {
			fmt.Printf("---skip----(left %d items)\n", len(result)-i)
			break
		}
	}
}
