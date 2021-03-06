package content

import (
	"bytes"
	"fmt"
	"math/rand"

	"github.com/nboughton/go-roll"
	"github.com/nboughton/swnt/content/format"
	"github.com/nboughton/swnt/content/table"
	"github.com/nboughton/swnt/dice"
)

func init() {
	table.Registry.Add(religionTable.leadership)
}

// Religion is pretty self explanatory
type Religion struct {
	Evolution       string
	Leadership      string
	OriginTradition string
}

// NewReligion with random characteristics
func NewReligion() Religion {
	r := Religion{
		Evolution:       religionTable.evolution.Roll(),
		Leadership:      religionTable.leadership.Roll(),
		OriginTradition: religionTable.origin.Roll(),
	}
	return r
}

// Format returns the religion formatted as type t
func (r Religion) Format(t format.OutputType) string {
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, format.Table(t, []string{"Religion", ""}, [][]string{
		{religionTable.origin.Name, r.OriginTradition},
		{religionTable.evolution.Name, r.Evolution},
		{religionTable.leadership.Name, r.Leadership},
	}))

	return buf.String()
}

func (r Religion) String() string {
	return r.Format(format.TEXT)
}

/*************** TABLES ***************/

var religionTable = struct {
	evolution  roll.List
	origin     roll.List
	leadership roll.Table
}{
	// Evolution SWN Revised Free Edition p193
	roll.List{
		Name: "Evolution",
		Items: []string{
			"New holy book. Someone in the faith’s past penned or discovered a text that is now taken to be holy writ and the expressed will of the divine.",
			"New prophet. This faith reveres the words and example of a relatively recent prophet, esteeming him or her as the final word on the will of God. The prophet may or may not still be living.",
			"Syncretism. The faith has merged many of its beliefs with another religion. Roll again on the origin tradition table; this faith has reconciled the major elements of both beliefs into its tradition.",
			"Neofundamentalism. The faith is fiercely resistant to perceived innovations and deviations from their beliefs. Even extremely onerous traditions and restrictions will be observed to the letter.",
			"Quietism. The faith shuns the outside world and involvement with the affairs of nonbelievers. They prefer to keep to their own kind and avoid positions of wealth and power.",
			"Sacrifices. The faith finds it necessary to make substantial sacrifices to please God. Some faiths may go so far as to offer human sacrifices, while others insist on huge tithes offered to the building of religious edifices.",
			"Schism. The faith’s beliefs are actually almost identical to those of the majority of its origin tradition, save for a few minor points of vital interest to theologians and no practical difference whatsoever to believers. This does not prevent a burning resentment towards the parent faith.",
			"Holy family. God’s favor has been shown especially to members of a particular lineage. It may be that only men and women of this bloodline are permitted to become clergy, or they may serve as helpless figureheads for the real leaders of the faith",
		},
	},

	// OriginTradition SWN Revised Free Edition p193
	roll.List{
		Name: "Origin",
		Items: []string{
			"Paganism",
			"Roman Catholicism",
			"Eastern Orthodox Christianity",
			"Protestant Christianity",
			"Buddhism",
			"Judaism",
			"Islam",
			"Taoism",
			"Hinduism",
			"Zoroastrianism",
			"Confucianism",
			"Ideology",
		},
	},

	// Leadership SWN Revised Free Edition p193
	roll.Table{
		Name: "Leadership",
		ID:   "religion.Leadership",
		Dice: roll.Dice{N: 1, Die: roll.D6},
		Items: []roll.TableItem{
			{Match: []int{1, 2}, Text: "Patriarch/Matriarch. A single leader determines doctrine for the entire religion, possibly in consultation with other clerics."},
			{Match: []int{3, 4}, Text: "Council. A group of the oldest and most revered clergy determine the course of the faith."},
			{Match: []int{5}, Text: "Democracy. Every member has an equal voice in matters of faith, with doctrine usually decided at regular church- wide councils."},
			{Match: []int{6}, Text: "No universal leadership", Action: func() string {
				tbl, _ := table.Registry.Get("religion.Leadership")
				tbl.Dice = roll.Dice{N: 1, Die: dice.D5}

				if rand.Intn(6)+1 == 6 {
					return ""
				}

				return "Each region governed independently by a " + tbl.Roll()
			}},
		},
	},
}
