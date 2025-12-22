package items

type Item struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	Value             int             `json:"value"`
	Icon              string          `json:"icon"`
	RecycleComponents *[]RecycleEntry `json:"recycle_components"`
	UsedIn            *[]UsedInEntry  `json:"used_in"`
}

type RecycleEntry struct {
	Quantity  int       `json:"quantity"`
	Component Component `json:"component"`
}

type UsedInEntry struct {
	Quantity int  `json:"quantity"`
	Item     Item `json:"item"`
}

type Component struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	Data []Item `json:"data"`
}

type ItemMap map[string]Item

func BuildIndex(items []Item) ItemMap {
	itemMap := make(ItemMap, len(items))
	for _, item := range items {
		itemMap[item.ID] = item
	}
	return itemMap
}

func (m ItemMap) Get(id string) (Item, bool) {
	item, ok := m[id]
	return item, ok
}
