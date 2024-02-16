package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Meal struct {
	price    string
	MealName string
	mealType string
}

type MealDay struct {
	date string
	menu []*Meal
}

func mensaGrapThisWeek() []*MealDay {
	URL := "https://www.studentenwerk-wuerzburg.de/bamberg/essen-trinken/speiseplaene/mensa-austrasse-bamberg.html"

	var week []*MealDay

	c := colly.NewCollector()

	c.OnHTML("div.week.currentweek", func(e *colly.HTMLElement) {

		e.ForEach("div.day", func(_ int, day *colly.HTMLElement) {

			date := day.ChildText("h5")

			var dailyMenu []*Meal

			day.ForEach("article.menu", func(_ int, mealT *colly.HTMLElement) {

				var mealType string

				mealT.ForEach("div.theicon", func(_ int, d *colly.HTMLElement) {
					mealType = d.Attr("title")
					if mealType == "Fleischlos" {
						mealType = "Vegetarisch"
					}
				})

				meal := mealT.ChildText("div.title")

				price := mealT.ChildText("div.price")

				adjustedPrice := strings.TrimSpace(price)

				if adjustedPrice == "" {
					adjustedPrice = "-"
				}

				dailyMenu = append(dailyMenu, &Meal{price: adjustedPrice, MealName: meal, mealType: mealType})

			})
			week = append(week, &MealDay{date: date, menu: dailyMenu})
		})
	})

	c.Visit(URL)

	return week
}

func main() {

	allFlag := flag.Bool("a", false, "Get the whole week")

	flag.Parse()

	week := mensaGrapThisWeek()

	today := time.Now().Format("02.01.")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"Essen", "Preis", "Typ"})

	if *allFlag {
		for _, mealDay := range week {

			t.AppendRow(table.Row{mealDay.date})

			for _, meal := range mealDay.menu {
				t.AppendRow(table.Row{meal.MealName, meal.price, meal.mealType})
			}
			t.AppendSeparator()
		}
	} else {
		for _, mealDay := range week {
			if strings.Contains(mealDay.date, today) {

				t.AppendRow(table.Row{mealDay.date})

				for _, meal := range mealDay.menu {
					t.AppendRow(table.Row{meal.MealName, meal.price, meal.mealType})
				}
			}
		}
	}

	t.Render()
}
