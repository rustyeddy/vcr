package main

// Settings structure
type Settings struct {
	values   map[string]string
	defaults map[string]string
	allowNew bool
}

// NewSettings accepts a map of strings as the configuration
func NewSettings(defaults map[string]string) (set *Settings) {
	set = &Settings{defaults: defaults}

	// Extract the allowed fields from the structure
	set.values = make(map[string]string)
	for idx, val := range defaults {
		set.values[idx] = val
	}
	return set
}

// Get gets the value otherwise returns false
func (set *Settings) Get(k string) string {
	s, _ := set.values[k]
	return s
}

// Exists determines if the key exists in this map, it
// also reflects an illegal configuration item.
func (set *Settings) Exists(k string) (e bool) {
	_, e = set.values[k]
	return e
}

// Add will add a non existant element to the list. It will return
// false if the element already exists, true if everything went
// as planned.
//
// XXX: This should not be used for our config.
func (set *Settings) Add(k string, v string) bool {
	// Settingss
	if set.allowNew == false {
		return false
	}
	_, e := set.values[k]
	if e == false {
		set.values[k] = v
	}
	return e
}

// Update will update an existing item only. Update will return
// false if the key does not already exist in the settings. True
// will be returned if all is well.
func (set *Settings) Update(k string, val string) bool {
	_, e := set.values[k]
	if e == true {
		set.values[k] = val
	}
	return e
}

// Delete removes the value from the list
func (set *Settings) Delete(k string) bool {
	_, e := set.values[k]
	if e == true {
		delete(set.values, k)
	}
	return e
}
