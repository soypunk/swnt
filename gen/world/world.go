package world

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/nboughton/rollt"
	"github.com/nboughton/swnt/gen/culture"
	"github.com/nboughton/swnt/gen/name"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// TagsTable represents the collection of Tags
type TagsTable []Tag

// Roll selects a random Tag
func (t TagsTable) Roll() string {
	return fmt.Sprint(t[rand.Intn(len(t))])
}

// Random selects a random tag (used in Adventure seed generation)
func (t TagsTable) Random() string {
	return Tags[rand.Intn(len(Tags))].Name
}

// Find returns the tag specified. The search is case insensitive for convenience
func (t TagsTable) Find(name string) (Tag, error) {
	for _, tag := range t {
		if strings.ToLower(tag.Name) == strings.ToLower(name) {
			return tag, nil
		}
	}

	return Tag{}, fmt.Errorf("no tag with name \"%s\"", name)
}

func selectTags(exclude []string) (Tag, Tag) {
	var t TagsTable
	for _, tag := range Tags {
		if !tag.match(exclude) {
			t = append(t, tag)
		}
	}

	t1Idx, t2Idx := rand.Intn(len(t)), rand.Intn(len(t))
	for t1Idx == t2Idx { // Ensure the same tag isn't selected twice
		t2Idx = rand.Intn(len(t))
	}

	return t[t1Idx], t[t2Idx]
}

func (t Tag) match(s []string) bool {
	for _, x := range s {
		if strings.ToLower(t.Name) == strings.ToLower(x) {
			return true
		}
	}

	return false
}

// Tag represents a complete World Tag structure as extracted from the Stars Without Number core book.
type Tag struct {
	Name          string
	Desc          string
	Enemies       rollt.List
	Friends       rollt.List
	Complications rollt.List
	Things        rollt.List
	Places        rollt.List
}

func (t Tag) String() string {
	return fmt.Sprintf(
		"Name\t:\t%s\nDesc\t:\t%s\nEnemies\t:\t%s\nFriends\t:\t%s\nComplications\t:\t%s\nThings\t:\t%s\nPlaces\t:\t%s\n",
		t.Name, t.Desc, t.Enemies, t.Friends, t.Complications, t.Things, t.Places,
	)
}

// World represents a generated world
type World struct {
	Primary      bool
	Name         string
	Culture      culture.Culture
	Tags         [2]Tag
	Atmosphere   string
	Temperature  string
	Population   string
	Biosphere    string
	TechLevel    string
	Origin       string
	Relationship string
	Contact      string
}

// New World, set culture to culture.Any for a random culture and primary to false
// to include relationship information
func New(c culture.Culture, primary bool, exclude []string) World {
	t1, t2 := selectTags(exclude)

	w := World{
		Primary:     primary,
		Name:        name.Names.ByCulture(c).Place.Roll(),
		Culture:     c,
		Tags:        [2]Tag{t1, t2},
		Atmosphere:  Atmosphere.Roll(),
		Temperature: Temperature.Roll(),
		Population:  Population.Roll(),
		Biosphere:   Biosphere.Roll(),
		TechLevel:   TechLevel.Roll(),
	}

	if !w.Primary {
		w.Origin = Other.Origin.Roll()
		w.Relationship = Other.Relationship.Roll()
		w.Contact = Other.Contact.Roll()
	}

	return w
}

func (w World) String() string {
	var buf = new(bytes.Buffer)

	fmt.Fprintf(buf, "Name\t:\t%s\n", w.Name)
	fmt.Fprintf(buf, "Culture\t:\t%s\n", w.Culture)
	fmt.Fprintf(buf, "Atmosphere\t:\t%s\n", w.Atmosphere)
	fmt.Fprintf(buf, "Temperature\t:\t%s\n", w.Temperature)
	fmt.Fprintf(buf, "Biosphere\t:\t%s\n", w.Biosphere)
	fmt.Fprintf(buf, "Population\t:\t%s\n", w.Population)
	fmt.Fprintf(buf, "Tech Level\t:\t%s\n", w.TechLevel)
	fmt.Fprintln(buf, "\t\nTags\t")
	fmt.Fprintf(buf, w.Tags[0].String())
	fmt.Fprintln(buf, "\t")
	fmt.Fprintf(buf, w.Tags[1].String())
	fmt.Fprintln(buf, "\t")

	if !w.Primary {
		fmt.Fprintf(buf, "%s\t:\t%s\n", Other.Origin.Name, w.Origin)
		fmt.Fprintf(buf, "%s\t:\t%s\n", Other.Relationship.Name, w.Relationship)
		fmt.Fprintf(buf, "%s\t:\t%s\n", Other.Contact.Name, w.Contact)
	}

	return buf.String()
}

// Markdown returns world as a Markdown table
func (w World) Markdown() string {
	var buf = new(bytes.Buffer)

	fmt.Fprintf(buf, "| %s | |\n| --- | --- |\n", w.Name)
	fmt.Fprintf(buf, "| Atmosphere | %s |\n", w.Atmosphere)
	fmt.Fprintf(buf, "| Temperature | %s |\n", w.Temperature)
	fmt.Fprintf(buf, "| Biosphere | %s |\n", w.Biosphere)
	fmt.Fprintf(buf, "| Population | %s |\n", w.Population)
	fmt.Fprintf(buf, "| Culture | %s |\n", w.Culture)
	fmt.Fprintf(buf, "| Tech Level | %s |\n", w.TechLevel)
	fmt.Fprintf(buf, "| **Tags** | |\n")
	fmt.Fprintf(buf, "| %s | %s |\n", w.Tags[0].Name, w.Tags[0].Desc)
	fmt.Fprintf(buf, "| %s | %s |\n", w.Tags[1].Name, w.Tags[1].Desc)

	if !w.Primary {
		fmt.Fprintf(buf, "| **Origins** | |\n")
		fmt.Fprintf(buf, "| %s | %s |\n", Other.Origin.Name, w.Origin)
		fmt.Fprintf(buf, "| %s | %s |\n", Other.Relationship.Name, w.Relationship)
		fmt.Fprintf(buf, "| %s | %s |\n", Other.Contact.Name, w.Contact)
	}

	fmt.Fprintln(buf)

	return buf.String()
}
