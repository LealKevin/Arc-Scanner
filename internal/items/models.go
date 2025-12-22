// Package items provides item data models and operations for arc-scanner.
package items

// Item represents a game item with its properties and relationships.
type Item struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	Value             int             `json:"value"`
	Icon              string          `json:"icon"`
	RecycleComponents *[]RecycleEntry `json:"recycle_components"`
	UsedIn            *[]UsedInEntry  `json:"used_in"`
}

// RecycleEntry represents a component obtained from recycling an item.
type RecycleEntry struct {
	Quantity  int       `json:"quantity"`
	Component Component `json:"component"`
}

// UsedInEntry represents a crafting recipe that uses this item.
type UsedInEntry struct {
	Quantity int  `json:"quantity"`
	Item     Item `json:"item"`
}

// Component represents a base component used in crafting.
type Component struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Response represents the API response containing items.
type Response struct {
	Data []Item `json:"data"`
}

// ItemMap provides fast item lookup by ID.
type ItemMap map[string]Item

// BuildIndex creates an ItemMap from a slice of items for fast lookup.
func BuildIndex(items []Item) ItemMap {
	itemMap := make(ItemMap, len(items))
	for _, item := range items {
		itemMap[item.ID] = item
	}
	return itemMap
}

// Get retrieves an item by ID. Returns the item and true if found,
// or an empty Item and false if not found.
func (m ItemMap) Get(id string) (Item, bool) {
	item, ok := m[id]
	return item, ok
}
